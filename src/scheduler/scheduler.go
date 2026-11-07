package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Task represents a scheduled task
type Task struct {
	Name     string
	Interval time.Duration
	Func     func() error
	LastRun  time.Time
	NextRun  time.Time
	Enabled  bool
	NextRunFunc func(time.Time) time.Time
}

// Scheduler manages periodic tasks
type Scheduler struct {
	tasks   map[string]*Task
	stop    chan struct{}
	running bool
	inFlight map[string]bool
	mu      sync.RWMutex
}

// New creates a new scheduler
func New() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
		stop:  make(chan struct{}),
		inFlight: make(map[string]bool),
	}
}

// AddTask adds a new task to the scheduler
func (s *Scheduler) AddTask(name string, interval time.Duration, fn func() error) {
	s.AddTaskAt(name, interval, time.Now().Add(interval), fn)
}

// AddTaskAt registers a task with a persisted or otherwise selected first run.
func (s *Scheduler) AddTaskAt(name string, interval time.Duration, nextRun time.Time, fn func() error) {
	s.AddTaskSchedule(name, interval, nextRun, nil, fn)
}

// AddTaskSchedule registers a task with a persisted first run and custom cadence.
func (s *Scheduler) AddTaskSchedule(name string, interval time.Duration, nextRun time.Time, nextRunFunc func(time.Time) time.Time, fn func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[name] = &Task{Name: name, Interval: interval, Func: fn, NextRun: nextRun, Enabled: true, NextRunFunc: nextRunFunc}
	log.Printf("Scheduler: Added task '%s' (interval: %v)", name, interval)
}

// RemoveTask removes a task from the scheduler
func (s *Scheduler) RemoveTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, name)
	log.Printf("Scheduler: Removed task '%s'", name)
}

// EnableTask enables a task
func (s *Scheduler) EnableTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, ok := s.tasks[name]; ok {
		task.Enabled = true
		task.NextRun = time.Now().Add(task.Interval)
	}
}

// DisableTask disables a task
func (s *Scheduler) DisableTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, ok := s.tasks[name]; ok {
		task.Enabled = false
	}
}

// Start begins the scheduler loop
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stop = make(chan struct{})
	s.mu.Unlock()

	log.Printf("Scheduler: Started with %d tasks", len(s.tasks))

	go func() {
		// Poll periodically for tasks that have become due.
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.stop:
				log.Println("Scheduler: Stopped")
				return
			case <-ticker.C:
				s.runDueTasks()
			}
		}
	}()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stop)
	s.running = false
}

// runDueTasks executes due tasks independently and avoids overlapping runs.
func (s *Scheduler) runDueTasks() {
	s.mu.Lock()
	now := time.Now()
	dueTasks := make([]*Task, 0)
	for name, task := range s.tasks {
		if task.Enabled && !s.inFlight[name] && !now.Before(task.NextRun) {
			s.inFlight[name] = true
			dueTasks = append(dueTasks, task)
		}
	}
	s.mu.Unlock()

	for _, task := range dueTasks {
		go s.executeTask(task)
	}
}

// executeTask isolates task failures and always advances the next run time.
func (s *Scheduler) executeTask(task *Task) {
	log.Printf("Scheduler: Running task '%s'", task.Name)
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Printf("Scheduler: Task '%s' panicked: %v", task.Name, recovered)
		}
		s.mu.Lock()
		task.LastRun = time.Now()
		if task.NextRunFunc != nil {
			task.NextRun = task.NextRunFunc(task.LastRun)
		} else {
			task.NextRun = task.LastRun.Add(task.Interval)
		}
		delete(s.inFlight, task.Name)
		s.mu.Unlock()
	}()
	if err := task.Func(); err != nil {
		log.Printf("Scheduler: Task '%s' failed: %v", task.Name, err)
		return
	}
	log.Printf("Scheduler: Task '%s' completed", task.Name)
}

// RunNow immediately runs a task by name
func (s *Scheduler) RunNow(name string) error {
	s.mu.Lock()
	task, ok := s.tasks[name]
	if !ok {
		s.mu.Unlock()
		return nil
	}
	if s.inFlight[name] {
		s.mu.Unlock()
		return fmt.Errorf("scheduled task %q is already running", name)
	}
	s.inFlight[name] = true
	s.mu.Unlock()

	log.Printf("Scheduler: Running task '%s' immediately", name)
	var taskErr error
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				taskErr = fmt.Errorf("scheduled task %q panicked: %v", name, recovered)
				log.Printf("Scheduler: %v", taskErr)
			}
			s.mu.Lock()
			task.LastRun = time.Now()
			if task.NextRunFunc != nil {
				task.NextRun = task.NextRunFunc(task.LastRun)
			} else {
				task.NextRun = task.LastRun.Add(task.Interval)
			}
			delete(s.inFlight, name)
			s.mu.Unlock()
		}()
		taskErr = task.Func()
	}()
	return taskErr
}

// GetTasks returns all registered tasks
func (s *Scheduler) GetTasks() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, *t)
	}
	return tasks
}

// ParseInterval parses interval strings like "hourly", "daily", "weekly"
func ParseInterval(s string) time.Duration {
	switch s {
	case "minutely":
		return time.Minute
	case "hourly":
		return time.Hour
	case "daily":
		return 24 * time.Hour
	case "weekly":
		return 7 * 24 * time.Hour
	case "monthly":
		return 30 * 24 * time.Hour
	default:
		// Try to parse as duration
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
		return 24 * time.Hour // Default to daily
	}
}
