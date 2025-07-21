# Migration Roadmap: Next.js to Go (Static HTML MVP)

## Phase 1: Planning & Setup

1. **Requirements Analysis**
   - List all current features (authentication, expenses CRUD, search, pagination, etc.).
   - Decide which features are essential for MVP1 (static HTML, no JS interactivity).

2. **Go Project Setup**
   - Initialize a new Go project (Gin, Echo, or net/http).
   - Set up folders: `/cmd`, `/internal`, `/templates`, `/static`, `/handlers`, `/models`.

3. **Tooling**
   - Set up Go modules, linter, and formatter.
   - Prepare `/static` for CSS and future JS (AlpineJS/HTMX).

4. Set up database
   - Create database tables
   - Migrate data

---

## Phase 2: MVP1 – Static HTML

1. **Authentication**
   - Implement login/logout with session cookies.
   - Render login form and error messages with server-side templates.

2. **Expenses List Page**
   - Render expenses list as a static HTML table.
   - Implement search and pagination with standard HTML forms and full page reloads.

3. **Add/Edit Expense**
   - Render forms for adding/editing expenses.
   - Handle form submissions with standard POST requests and redirects.

4. **Navigation & Layout**
   - Create base templates for consistent layout and navigation.

---

## Phase 3: UI & Code Quality

1. **UI Components**
   - Recreate reusable UI elements as HTML partials/templates.
   - Style with CSS or Tailwind (optional).

2. **Testing**
   - Write unit and integration tests for handlers and templates.

3. **Security**
   - Add middleware for authentication, CSRF, and input validation.

---

## Phase 4: Progressive Enhancement

1. **Add AlpineJS/HTMX**
   - Gradually introduce interactivity (dynamic search, AJAX pagination, inline editing) as needed.
   - Enhance forms and tables without major backend changes.

2. **Performance & Accessibility**
   - Optimize template rendering and static asset delivery.
   - Ensure accessibility best practices.

---

## Phase 5: Migration & Launch

1. **Data Migration**
   - Migrate existing data if needed.

2. **Deployment**
   - Set up deployment pipeline and monitoring.

3. **Cutover**
   - Switch to the new Go application and monitor for issues.

---

This plan lets you launch quickly with static HTML and add interactivity later, minimizing risk and rework.

