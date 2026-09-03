package admin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	auth      *AuthManager
	version   string
	startTime time.Time
}

func NewHandler(username, password, apiToken string, sessionTimeout int, sslEnabled bool, version string) *Handler {
	return &Handler{
		auth:      NewAuthManager(username, password, apiToken, sessionTimeout, sslEnabled),
		version:   version,
		startTime: time.Now(),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/admin", h.handleAdminLogin)
	r.Post("/admin/login", h.handleAdminLoginPost)
	r.Get("/admin/logout", h.handleAdminLogout)
	r.Get("/admin/dashboard", h.requireSession(h.handleAdminDashboard))
	r.Get("/admin/settings", h.requireSession(h.handleAdminSettings))
	r.Post("/admin/settings", h.requireSession(h.handleAdminSettings))

	r.Get("/api/v1/admin/status", h.requireToken(h.handleAPIStatus))
	r.Get("/api/v1/admin/config", h.requireToken(h.handleAPIConfig))
	r.Put("/api/v1/admin/config", h.requireToken(h.handleAPIConfig))
	r.Post("/api/v1/admin/reload", h.requireToken(h.handleAPIReload))
}

func (h *Handler) requireSession(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := h.auth.GetSessionFromRequest(r)
		if !ok {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		h.auth.RefreshSession(session.ID)
		handler(w, r)
	}
}

func (h *Handler) requireToken(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := GetTokenFromRequest(r)
		if token == "" || !h.auth.ValidateAPIToken(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		handler(w, r)
	}
}

func (h *Handler) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.auth.GetSessionFromRequest(r); ok {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	h.renderLoginPage(w, "")
}

func (h *Handler) handleAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if h.auth.Authenticate(username, password) {
		session := h.auth.CreateSession(username, r.RemoteAddr)
		h.auth.SetSessionCookie(w, session)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	h.renderLoginPage(w, "Invalid username or password")
}

func (h *Handler) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if session, ok := h.auth.GetSessionFromRequest(r); ok {
		h.auth.DeleteSession(session.ID)
	}
	h.auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	h.renderDashboardPage(w)
}

func (h *Handler) handleAdminSettings(w http.ResponseWriter, r *http.Request) {
	message := ""
	if r.Method == http.MethodPost {
		message = "Settings updated successfully"
	}
	h.renderSettingsPage(w, message)
}

func (h *Handler) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": h.version,
		"uptime":  time.Since(h.startTime).String(),
		"memory": map[string]interface{}{
			"alloc": m.Alloc,
			"sys":   m.Sys,
			"numGC": m.NumGC,
		},
		"runtime": map[string]interface{}{
			"goroutines": runtime.NumGoroutine(),
			"cpus":       runtime.NumCPU(),
			"goVersion":  runtime.Version(),
		},
	})
}

func (h *Handler) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"version": h.version})
}

func (h *Handler) handleAPIReload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "reloaded"})
}

func (h *Handler) renderLoginPage(w http.ResponseWriter, errorMsg string) {
	tmpl := template.Must(template.New("login").Parse(loginTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{"Error": errorMsg})
}

func (h *Handler) renderDashboardPage(w http.ResponseWriter) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{
		"Version":    h.version,
		"MemAlloc":   fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		"Goroutines": runtime.NumGoroutine(),
		"Uptime":     time.Since(h.startTime).Round(time.Second).String(),
	})
}

func (h *Handler) renderSettingsPage(w http.ResponseWriter, message string) {
	tmpl := template.Must(template.New("settings").Parse(settingsTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{"Message": message})
}

const loginTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Login - Zipcodes</title>
    <style>
        :root { --bg-color: #282a36; --fg-color: #f8f8f2; --accent: #bd93f9; --red: #ff5555; --input-bg: #44475a; }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg-color); color: var(--fg-color); min-height: 100vh; display: flex; align-items: center; justify-content: center; }
        .login-container { background: var(--input-bg); padding: 2rem; border-radius: 8px; width: 100%; max-width: 400px; }
        h1 { text-align: center; margin-bottom: 1.5rem; color: var(--accent); }
        .error { background: var(--red); color: #fff; padding: 0.75rem; border-radius: 4px; margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; }
        input[type="text"], input[type="password"] { width: 100%; padding: 0.75rem; border: none; border-radius: 4px; background: var(--bg-color); color: var(--fg-color); margin-bottom: 1rem; }
        button { width: 100%; padding: 0.75rem; border: none; border-radius: 4px; background: var(--accent); color: var(--bg-color); font-weight: 600; cursor: pointer; }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>Admin Login</h1>
        {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
        <form method="POST" action="/admin/login">
            <label for="username">Username</label>
            <input type="text" id="username" name="username" required autofocus>
            <label for="password">Password</label>
            <input type="password" id="password" name="password" required>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>`

const dashboardTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Dashboard - Zipcodes</title>
    <style>
        :root { --bg-color: #282a36; --fg-color: #f8f8f2; --accent: #bd93f9; --card-bg: #44475a; }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg-color); color: var(--fg-color); min-height: 100vh; }
        .navbar { background: var(--card-bg); padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; }
        .navbar h1 { color: var(--accent); font-size: 1.5rem; }
        .navbar a { color: var(--fg-color); text-decoration: none; margin-left: 1rem; }
        .container { max-width: 1200px; margin: 2rem auto; padding: 0 1rem; }
        .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; }
        .card { background: var(--card-bg); padding: 1.5rem; border-radius: 8px; }
        .card h3 { color: var(--accent); margin-bottom: 0.5rem; }
        .card p { font-size: 1.5rem; font-weight: bold; }
    </style>
</head>
<body>
    <nav class="navbar"><h1>Zipcodes Admin</h1><div><a href="/admin/dashboard">Dashboard</a><a href="/admin/settings">Settings</a><a href="/admin/logout">Logout</a></div></nav>
    <div class="container">
        <div class="cards">
            <div class="card"><h3>Version</h3><p>{{.Version}}</p></div>
            <div class="card"><h3>Memory</h3><p>{{.MemAlloc}}</p></div>
            <div class="card"><h3>Goroutines</h3><p>{{.Goroutines}}</p></div>
            <div class="card"><h3>Uptime</h3><p>{{.Uptime}}</p></div>
        </div>
    </div>
</body>
</html>`

const settingsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings - Zipcodes</title>
    <style>
        :root { --bg-color: #282a36; --fg-color: #f8f8f2; --accent: #bd93f9; --card-bg: #44475a; --green: #50fa7b; }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg-color); color: var(--fg-color); min-height: 100vh; }
        .navbar { background: var(--card-bg); padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; }
        .navbar h1 { color: var(--accent); font-size: 1.5rem; }
        .navbar a { color: var(--fg-color); text-decoration: none; margin-left: 1rem; }
        .container { max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
        .message { background: var(--green); color: #000; padding: 1rem; border-radius: 4px; margin-bottom: 1rem; }
        .card { background: var(--card-bg); padding: 1.5rem; border-radius: 8px; }
        .card h2 { color: var(--accent); margin-bottom: 1rem; }
    </style>
</head>
<body>
    <nav class="navbar"><h1>Zipcodes Admin</h1><div><a href="/admin/dashboard">Dashboard</a><a href="/admin/settings">Settings</a><a href="/admin/logout">Logout</a></div></nav>
    <div class="container">
        {{if .Message}}<div class="message">{{.Message}}</div>{{end}}
        <div class="card"><h2>Settings</h2><p>Settings configuration coming soon.</p></div>
    </div>
</body>
</html>`
