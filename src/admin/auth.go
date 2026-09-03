package admin

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

type Session struct {
	ID        string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
	IP        string
}

type AuthManager struct {
	username       string
	passwordHash   []byte
	passwordSalt   []byte
	apiToken       string
	sessionTimeout time.Duration
	sessions       map[string]*Session
	mu             sync.RWMutex
	secureCookie   bool
}

func NewAuthManager(username, password, apiToken string, sessionTimeoutSecs int, sslEnabled bool) *AuthManager {
	am := &AuthManager{
		username:       username,
		apiToken:       apiToken,
		sessionTimeout: time.Duration(sessionTimeoutSecs) * time.Second,
		sessions:       make(map[string]*Session),
		secureCookie:   sslEnabled,
	}

	if password != "" {
		am.passwordSalt = make([]byte, 16)
		rand.Read(am.passwordSalt)
		am.passwordHash = argon2.IDKey([]byte(password), am.passwordSalt, 1, 64*1024, 4, 32)
	}

	go am.cleanupExpiredSessions()

	return am
}

func (am *AuthManager) Authenticate(username, password string) bool {
	if am.username == "" || am.passwordHash == nil {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(username), []byte(am.username)) != 1 {
		return false
	}
	hash := argon2.IDKey([]byte(password), am.passwordSalt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, am.passwordHash) == 1
}

func (am *AuthManager) CreateSession(username, ip string) *Session {
	am.mu.Lock()
	defer am.mu.Unlock()

	idBytes := make([]byte, 32)
	rand.Read(idBytes)
	id := base64.URLEncoding.EncodeToString(idBytes)

	session := &Session{
		ID:        id,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(am.sessionTimeout),
		IP:        ip,
	}
	am.sessions[id] = session
	return session
}

func (am *AuthManager) GetSession(id string) (*Session, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	session, ok := am.sessions[id]
	if !ok || time.Now().After(session.ExpiresAt) {
		return nil, false
	}
	return session, true
}

func (am *AuthManager) RefreshSession(id string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	if session, ok := am.sessions[id]; ok {
		session.ExpiresAt = time.Now().Add(am.sessionTimeout)
	}
}

func (am *AuthManager) DeleteSession(id string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.sessions, id)
}

func (am *AuthManager) ValidateAPIToken(token string) bool {
	if am.apiToken == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(am.apiToken)) == 1
}

func (am *AuthManager) GetSessionFromRequest(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return nil, false
	}
	return am.GetSession(cookie.Value)
}

func (am *AuthManager) SetSessionCookie(w http.ResponseWriter, session *Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    session.ID,
		Path:     "/admin",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   am.secureCookie,
		SameSite: http.SameSiteStrictMode,
	})
}

func (am *AuthManager) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   am.secureCookie,
		SameSite: http.SameSiteStrictMode,
	})
}

func (am *AuthManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		am.mu.Lock()
		now := time.Now()
		for id, session := range am.sessions {
			if now.After(session.ExpiresAt) {
				delete(am.sessions, id)
			}
		}
		am.mu.Unlock()
	}
}

func GetTokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}
