package scheduler

import (
	"errors"
	"testing"
	"time"
)

func TestTaskLifecycle(t *testing.T) {
	s := New()
	called := false
	s.AddTask("task", time.Hour, func() error { called = true; return nil })
	if got := s.GetTasks(); len(got) != 1 || !got[0].Enabled {
		t.Fatalf("tasks = %#v", got)
	}
	if err := s.RunNow("missing"); err != nil {
		t.Fatalf("missing task: %v", err)
	}
	if err := s.RunNow("task"); err != nil || !called {
		t.Fatalf("RunNow: called=%t err=%v", called, err)
	}
	s.DisableTask("task")
	if s.GetTasks()[0].Enabled {
		t.Fatal("task remained enabled")
	}
	s.EnableTask("task")
	if !s.GetTasks()[0].Enabled {
		t.Fatal("task remained disabled")
	}
	s.RemoveTask("task")
	if len(s.GetTasks()) != 0 {
		t.Fatal("task was not removed")
	}
	s.Start()
	s.Start()
	s.Stop()
	s.Stop()
}

func TestDueTasksAndIntervals(t *testing.T) {
	s := New()
	done := make(chan struct{}, 1)
	s.AddTask("ok", time.Second, func() error { done <- struct{}{}; return nil })
	s.AddTask("bad", time.Second, func() error { return errors.New("expected") })
	s.AddTask("disabled", time.Second, func() error { t.Error("disabled task ran"); return nil })
	s.DisableTask("disabled")
	s.mu.Lock()
	for _, task := range s.tasks {
		task.NextRun = time.Now().Add(-time.Second)
	}
	s.mu.Unlock()
	s.runDueTasks()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("due task did not run")
	}

	for input, want := range map[string]time.Duration{"minutely": time.Minute, "hourly": time.Hour, "daily": 24 * time.Hour, "weekly": 7 * 24 * time.Hour, "monthly": 30 * 24 * time.Hour, "2h": 2 * time.Hour, "invalid": 24 * time.Hour} {
		if got := ParseInterval(input); got != want {
			t.Errorf("ParseInterval(%q) = %s, want %s", input, got, want)
		}
	}
}
