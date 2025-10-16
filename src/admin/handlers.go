package admin

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

// Handler handles admin routes
type Handler struct {
	db        *sql.DB
	templates embed.FS
}

// NewHandler creates admin handler
func NewHandler(db *sql.DB, templates embed.FS) *Handler {
	return &Handler{
		db:        db,
		templates: templates,
	}
}

// DashboardHandler shows admin dashboard
func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "admin/dashboard.html", map[string]interface{}{
		"ServerTitle":       "Zipcodes",
		"ServerDescription": "US Postal Code Lookup API",
		"PageTitle":         "Admin Dashboard",
	})
}

// SettingsHandler shows admin settings
func (h *Handler) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle settings update
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Update settings in database
		for key, values := range r.Form {
			if len(values) > 0 {
				_, err := h.db.Exec("UPDATE settings SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE key = ?", values[0], key)
				if err != nil {
					http.Error(w, "Failed to update settings", http.StatusInternalServerError)
					return
				}
			}
		}

		http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
		return
	}

	// Get settings from database
	settings, err := h.getSettings()
	if err != nil {
		http.Error(w, "Failed to load settings", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "admin/settings.html", map[string]interface{}{
		"ServerTitle":       "Zipcodes",
		"ServerDescription": "US Postal Code Lookup API",
		"PageTitle":         "Server Settings",
		"Settings":          settings,
	})
}

// DatabaseHandler shows database management
func (h *Handler) DatabaseHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "admin/database.html", map[string]interface{}{
		"ServerTitle":       "Zipcodes",
		"ServerDescription": "US Postal Code Lookup API",
		"PageTitle":         "Database Management",
	})
}

// DatabaseTestHandler tests database connection
func (h *Handler) DatabaseTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Test database connection
	err := h.db.Ping()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"Database connection failed"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"message":"Database connection successful"}`))
}

// LogsHandler shows log viewer
func (h *Handler) LogsHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "admin/logs.html", map[string]interface{}{
		"ServerTitle":       "Zipcodes",
		"ServerDescription": "US Postal Code Lookup API",
		"PageTitle":         "Log Viewer",
	})
}

// AuditHandler shows audit log
func (h *Handler) AuditHandler(w http.ResponseWriter, r *http.Request) {
	// Get audit logs from database
	rows, err := h.db.Query("SELECT id, username, action, resource, timestamp FROM audit_log ORDER BY timestamp DESC LIMIT 100")
	if err != nil {
		http.Error(w, "Failed to load audit log", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AuditEntry struct {
		ID        string
		Username  string
		Action    string
		Resource  string
		Timestamp string
	}

	var logs []AuditEntry
	for rows.Next() {
		var entry AuditEntry
		if err := rows.Scan(&entry.ID, &entry.Username, &entry.Action, &entry.Resource, &entry.Timestamp); err != nil {
			continue
		}
		logs = append(logs, entry)
	}

	h.renderTemplate(w, "admin/audit.html", map[string]interface{}{
		"ServerTitle":       "Zipcodes",
		"ServerDescription": "US Postal Code Lookup API",
		"PageTitle":         "Audit Log",
		"Logs":              logs,
	})
}

// getSettings retrieves all settings from database
func (h *Handler) getSettings() (map[string]string, error) {
	rows, err := h.db.Query("SELECT key, value FROM settings ORDER BY category, key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}

	return settings, nil
}

// renderTemplate renders a template with data
func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	tmplData, err := h.templates.ReadFile("templates/" + name)
	if err != nil {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}

	// Parse base template and specific template
	baseTmpl, err := h.templates.ReadFile("templates/base.html")
	if err != nil {
		http.Error(w, "Base template not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("base").Parse(string(baseTmpl))
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.Parse(string(tmplData))
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

// AdminInfoHandler returns admin information (API)
func (h *Handler) AdminInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"data":{"username":"administrator","role":"admin"}}`))
}

// AdminStatsHandler returns server statistics (API)
func (h *Handler) AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Get statistics from database
	var zipcodeCount int
	h.db.QueryRow("SELECT COUNT(*) FROM zipcodes").Scan(&zipcodeCount)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"success":true,"data":{"zipcodes":%d}}`, zipcodeCount)))
}

// ReloadHandler reloads configuration (API)
func (h *Handler) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reload configuration from database
	// In a real implementation, this would reload settings into memory
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"message":"Configuration reloaded"}`))
}
