# Universal Server Template Implementation TODO

## Progress Overview
- **Database Schema**: ✅ Complete
- **Authentication Core**: ✅ Complete
- **Middleware**: ✅ Complete
- **Setup Wizard**: 🔄 In Progress
- **Auth Routes**: ⏳ Not Started
- **User Interface**: ⏳ Not Started
- **Admin Interface**: ⏳ Not Started
- **UI Components**: ⏳ Not Started

---

## Completed ✅

### Database Layer
- [x] Complete schema (users, sessions, tokens, settings, audit_log, scheduled_tasks)
- [x] Schema initialization with default settings
- [x] Zipcode database integration

### Authentication System
- [x] User, Session, Token models
- [x] Password hashing (bcrypt)
- [x] Token generation (secure random)
- [x] Input validation (username, email, password)
- [x] Auth repository (database operations)
- [x] Session middleware
- [x] Auth middleware (RequireAuth, RequireAdmin)
- [x] Token auth middleware
- [x] Client IP detection

---

## In Progress 🔄

### User Setup Wizard
- [x] `/user/setup` - Welcome screen
- [x] `/user/setup/register` - First user registration form
- [x] `/user/setup/admin` - Administrator creation form
- [x] `/user/setup/complete` - Completion and redirect
- [x] Setup flow handlers (src/web/setup.go)
- [x] Setup detection (redirect if users exist)

---

## Not Started ⏳

### Authentication Routes & Handlers
- [ ] `/auth/login` - Login form (GET)
- [ ] `/auth/login` - Process login (POST)
- [ ] `/auth/logout` - Logout and redirect (GET)
- [ ] `/auth/register` - Registration form (GET) [if enabled]
- [ ] `/auth/register` - Process registration (POST)
- [ ] `/auth/password/reset` - Request reset form (GET)
- [ ] `/auth/password/reset` - Send reset email (POST)
- [ ] `/auth/password/new` - New password form (GET)
- [ ] `/auth/password/new` - Set new password (POST)
- [ ] Login handler (validate credentials, create session)
- [ ] Logout handler (delete session, clear cookie)
- [ ] Password reset handlers

### User Routes & Handlers
- [ ] `/user` - User dashboard
- [ ] `/user/profile` - Profile management (GET/POST)
- [ ] `/user/settings` - User settings (GET/POST)
- [ ] `/user/tokens` - Token management (GET/POST/DELETE)
- [ ] `/user/sessions` - Active sessions (GET/DELETE)
- [ ] `/user/security` - Security settings (GET/POST)
- [ ] `/user/avatar` - Avatar upload/delete
- [ ] User dashboard handler
- [ ] Profile update handler
- [ ] Settings update handler
- [ ] Token CRUD handlers
- [ ] Session management handlers

### Admin Routes & Handlers
- [ ] `/admin` - Admin dashboard
- [ ] `/admin/users` - User management list
- [ ] `/admin/users/:id` - View/edit specific user
- [ ] `/admin/settings` - Server settings
- [ ] `/admin/database` - Database management
- [ ] `/admin/logs` - Log viewer
- [ ] `/admin/audit` - Audit log
- [ ] `/admin/backup` - Backup management
- [ ] Admin dashboard handler (stats, charts)
- [ ] User management handlers (list, view, edit, delete)
- [ ] Settings handlers (get, update by category)
- [ ] Database status handler
- [ ] Log viewing handlers
- [ ] Audit log handler
- [ ] Backup/restore handlers

### HTML Templates (Must follow spec exactly)

#### Base Templates
- [ ] `base.html` - Base template with proper HTML5 structure
  - DOCTYPE, meta tags, viewport
  - Header with logo, nav, profile menu
  - Main content area
  - Footer (always at bottom)
  - Modal container
  - Toast container

#### Setup Templates
- [ ] `setup/welcome.html` - Welcome screen
- [ ] `setup/register.html` - First user registration
- [ ] `setup/admin.html` - Admin creation
- [ ] `setup/complete.html` - Setup complete

#### Auth Templates
- [ ] `auth/login.html` - Login form
- [ ] `auth/register.html` - Registration form
- [ ] `auth/password-reset.html` - Password reset request
- [ ] `auth/password-new.html` - New password form

#### User Templates
- [ ] `user/dashboard.html` - User dashboard
- [ ] `user/profile.html` - Profile management
- [ ] `user/settings.html` - User settings
- [ ] `user/tokens.html` - Token management
- [ ] `user/sessions.html` - Active sessions
- [ ] `user/security.html` - Security settings

#### Admin Templates
- [ ] `admin/dashboard.html` - Admin dashboard with stats
- [ ] `admin/users.html` - User management list
- [ ] `admin/user-detail.html` - User detail/edit
- [ ] `admin/settings.html` - Server settings
- [ ] `admin/database.html` - Database management
- [ ] `admin/logs.html` - Log viewer
- [ ] `admin/audit.html` - Audit log

### CSS Requirements (Professional Design System)
- [ ] CSS variables for theming (dark/light)
- [ ] Responsive breakpoints (mobile-first)
- [ ] Typography system
- [ ] Color system
- [ ] Spacing system
- [ ] Component styles:
  - [ ] Buttons (primary, secondary, danger)
  - [ ] Forms (inputs, labels, validation states)
  - [ ] Cards
  - [ ] Tables (sortable, hover, striped)
  - [ ] Navigation (header, sidebar)
  - [ ] Footer
  - [ ] Modals
  - [ ] Alerts/Toasts
  - [ ] Toggles/Switches
  - [ ] Banners
  - [ ] Loading states
  - [ ] Tooltips
  - [ ] Dropdowns
  - [ ] Tabs

### JavaScript Requirements (NO basic popups!)
- [ ] UI class with modal methods (showModal, confirm)
- [ ] Toast notification system
- [ ] Form validation
- [ ] AJAX helpers
- [ ] Timezone conversion
- [ ] Relative time updates
- [ ] Theme toggle
- [ ] Dropdown menus
- [ ] Tab switching
- [ ] Table sorting
- [ ] File upload handling
- [ ] Session keepalive

### UI Components (MANDATORY - No alert/confirm/prompt)
- [ ] Modal component (fade in/out, backdrop, ESC to close)
- [ ] Toggle/Switch component (iOS-style)
- [ ] Alert/Notification component (info, success, warning, error)
- [ ] Banner component (site-wide notifications)
- [ ] Loading states (button loading, page loading, skeleton screens)
- [ ] Enhanced forms (floating labels, validation, hints)
- [ ] Tooltips (not title attribute)
- [ ] Enhanced tables (sortable, searchable, paginated)
- [ ] Cards (consistent styling)
- [ ] Tabs (smooth transitions, keyboard nav)
- [ ] Progress indicators (bars, circular, step indicators)

### Navigation Requirements
- [ ] Profile menu (top right, dropdown)
  - [ ] Avatar display (uploaded/gravatar/initials)
  - [ ] User options when logged in
  - [ ] Login/Register when logged out
  - [ ] Admin panel link for administrators
- [ ] Session persistence (30-day default)
- [ ] "Remember me" checkbox
- [ ] Multi-device support
- [ ] Automatic token refresh
- [ ] Avatar system priority:
  1. User uploaded image
  2. Gravatar (if email)
  3. Generated initials (colored background)
  4. Default icon

### Server Integration
- [ ] Update main.go to initialize auth system
- [ ] Update server.go to add all routes
- [ ] Setup detection on server start
- [ ] First-run redirect logic
- [ ] Session cleanup scheduler
- [ ] Health check integration (show auth status)

### Security Implementation
- [ ] HTTPS detection and enforcement
- [ ] Security headers (all responses)
- [ ] CORS configuration
- [ ] Rate limiting (login, API)
- [ ] Session security (HttpOnly, Secure, SameSite)
- [ ] Token rotation
- [ ] Password reset token expiry
- [ ] Account lockout after failed attempts
- [ ] IP-based rate limiting

### API Endpoints for Auth
- [ ] `POST /api/v1/auth/login` - Login
- [ ] `POST /api/v1/auth/logout` - Logout
- [ ] `POST /api/v1/auth/refresh` - Refresh token
- [ ] `GET /api/v1/auth/status` - Check auth status
- [ ] `GET /api/v1/user` - Get current user info
- [ ] `PUT /api/v1/user` - Update user info
- [ ] `GET /api/v1/user/sessions` - List sessions
- [ ] `DELETE /api/v1/user/sessions/:id` - Kill session
- [ ] `GET /api/v1/admin/users` - List users (admin only)
- [ ] `GET /api/v1/admin/users/:id` - Get user (admin only)
- [ ] `PUT /api/v1/admin/users/:id` - Update user (admin only)
- [ ] `DELETE /api/v1/admin/users/:id` - Delete user (admin only)

### Additional Features (Per Spec)
- [ ] Email verification system
- [ ] 2FA setup and verification
- [ ] Password reset via email
- [ ] Account deletion (with grace period)
- [ ] Data export (GDPR compliance)
- [ ] Audit logging (all admin actions)
- [ ] Scheduled task management UI
- [ ] Backup/restore UI
- [ ] SSL certificate management
- [ ] Database configuration UI
- [ ] Monitoring/metrics configuration

### Testing
- [ ] Unit tests for auth functions
- [ ] Integration tests for auth flow
- [ ] E2E tests for setup wizard
- [ ] E2E tests for login/logout
- [ ] Session management tests
- [ ] API endpoint tests

### Documentation Updates
- [ ] Update README with auth setup
- [ ] Document API authentication
- [ ] Document user management
- [ ] Document admin features
- [ ] Add security best practices

---

## Notes

### Current State
- Zipcode API is fully functional
- Database schema is ready for auth system
- Auth core (models, repository, middleware) is complete
- Need to build entire frontend and integrate

### Critical Path (Minimum for working system)
1. Setup wizard (first run flow)
2. Login/logout pages and handlers
3. Basic user dashboard
4. Basic admin dashboard
5. Base HTML template
6. Minimal CSS/JS for functionality

### Spec Compliance Requirements
- All routes must be scoped correctly (/user/*, /admin/*, /auth/*)
- NO simple popups (alert/confirm/prompt) - use professional UI components
- 30-day session persistence by default
- Administrator account is separate from first user
- Admin can only access /admin/* routes (browse as guest elsewhere)
- All forms must have validation and error states
- Mobile responsive required (98% width < 720px, 90% >= 720px)
- Dark theme as default
- Security by default on everything

### File Count Estimate
- HTML Templates: ~20 files
- Go Handlers: ~15 files
- CSS Files: ~5 files
- JavaScript: ~8 files
- Total new files needed: ~48 files
- Total new lines: ~4000-5000 LOC

---

## Priority Order

### Phase 1: Core Auth (Next)
1. Setup wizard (4 templates + 1 handler file)
2. Login/logout (2 templates + 1 handler file)
3. Basic base template
4. Minimal CSS for functionality
5. Integrate with main.go and server.go

### Phase 2: User Interface
1. User dashboard
2. User profile
3. User settings
4. Session management

### Phase 3: Admin Interface
1. Admin dashboard
2. User management
3. Server settings
4. Logs and audit

### Phase 4: Professional UI
1. Complete CSS design system
2. UI component library
3. JavaScript interactions
4. Mobile responsive polish

### Phase 5: Advanced Features
1. Email verification
2. 2FA
3. Password reset
4. Backup/restore
5. Monitoring
