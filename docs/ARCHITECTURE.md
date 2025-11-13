# Commento Architecture Documentation

## Table of Contents
- [Overview](#overview)
- [Technology Stack](#technology-stack)
- [Database Design](#database-design)
- [Multi-Tenancy Architecture](#multi-tenancy-architecture)
- [Data Flow](#data-flow)
- [Frontend Architecture](#frontend-architecture)
- [Backend Architecture](#backend-architecture)
- [Authentication & Security](#authentication--security)
- [Customization Guide](#customization-guide)
- [Deployment Architecture](#deployment-architecture)

---

## Overview

**Commento** is a lightweight, privacy-focused, self-hosted commenting platform designed as an alternative to Disqus, Facebook Comments, and similar third-party services. It can be embedded into any website and provides a complete commenting system with moderation, OAuth integration, and multi-tenant capabilities.

### Key Features
- Privacy-first design (GDPR compliant, no tracking)
- Multi-tenant architecture (one instance serves multiple websites)
- OAuth integration (Google, GitHub, Twitter, GitLab)
- Markdown support with real-time preview
- Threaded comments with upvoting/downvoting
- Spam detection (Akismet integration)
- Email notifications
- Moderator system
- Import from Disqus

### Architecture Principles
1. **Domain-based isolation**: Each website is a separate tenant identified by domain
2. **Stateless sessions**: Session tokens stored in PostgreSQL
3. **API-first design**: Backend exposes REST API, frontend consumes it
4. **Embedded widget**: Minimal JavaScript bundle loaded on client websites
5. **Privacy by default**: No third-party tracking or analytics

---

## Technology Stack

### Backend
| Component | Technology | Purpose |
|-----------|-----------|---------|
| Language | Go 1.15+ | High-performance backend |
| Database | PostgreSQL 9.6+ | Persistent data storage |
| Router | Gorilla Mux | HTTP routing |
| CORS | Gorilla Handlers | Cross-origin support |
| Password Hashing | bcrypt | Secure password storage |
| Email | SMTP | Notification delivery |
| Markdown | Custom renderer | Comment formatting |

### Frontend
| Component | Technology | Purpose |
|-----------|-----------|---------|
| Build Tool | Gulp.js | Asset compilation |
| Styling | SASS/SCSS | CSS preprocessing |
| Embed Widget | Vanilla JavaScript | Comment widget |
| Dashboard | Vue.js + jQuery | Admin interface |
| Charts | Chartist.js | Analytics visualization |
| Code Highlighting | Highlight.js | Syntax highlighting in comments |

### Infrastructure
- **Docker**: Multi-stage builds for containerization
- **PostgreSQL**: Primary database with connection pooling
- **SMTP Server**: Email delivery
- **CDN Support**: Static asset delivery

---

## Database Design

### Entity Relationship Overview

```
┌─────────────┐
│   Owners    │ 1:N relationship with domains
└──────┬──────┘     (One owner can have multiple domains)
       │
       │
       ▼
┌─────────────┐
│   Domains   │ Multi-tenancy boundary
└──────┬──────┘     (Each domain = one website)
       │
       │ 1:N
       ▼
┌─────────────┐     ┌──────────────┐
│   Comments  │ N:1 │  Commenters  │
└──────┬──────┘     └──────┬───────┘
       │                   │
       │ 1:N               │ N:1
       ▼                   ▼
┌─────────────┐     ┌──────────────┐
│    Votes    │────▶│  Moderators  │
└─────────────┘     └──────────────┘
       │
       │ N:1
       ▼
┌─────────────┐
│    Pages    │ Metadata per URL path
└─────────────┘
```

### Core Tables

#### 1. Owners Table
**Purpose**: Platform administrators who manage one or more domains (websites)

```sql
CREATE TABLE owners (
    ownerHex TEXT PRIMARY KEY,           -- Unique identifier (hex string)
    email TEXT UNIQUE NOT NULL,          -- Login email
    name TEXT NOT NULL,                  -- Display name
    passwordHash TEXT NOT NULL,          -- bcrypt hash
    confirmedEmail BOOLEAN NOT NULL,     -- Email verification status
    joinDate TIMESTAMP NOT NULL          -- Registration date
);
```

**Key Points**:
- One owner can manage multiple domains
- Email must be unique across the platform
- Passwords are hashed with bcrypt (never stored plaintext)

#### 2. Domains Table
**Purpose**: Represents each website using Commento (tenant isolation boundary)

```sql
CREATE TABLE domains (
    domain TEXT PRIMARY KEY,                    -- Website domain (e.g., "example.com")
    ownerHex TEXT NOT NULL,                     -- Owner who controls this domain
    name TEXT NOT NULL,                         -- Friendly name
    creationDate TIMESTAMP NOT NULL,
    state TEXT NOT NULL,                        -- 'frozen' or 'unfrozen'

    -- Moderation Settings
    autoSpamFilter BOOLEAN NOT NULL,
    requireModeration BOOLEAN NOT NULL,          -- All comments need approval
    requireIdentification BOOLEAN NOT NULL,      -- No anonymous comments
    moderateAllAnonymous BOOLEAN NOT NULL,       -- Anonymous comments need approval

    -- Email Settings
    emailNotificationPolicy TEXT NOT NULL,       -- 'all', 'pending', 'none'

    -- OAuth Provider Configuration (per domain)
    commentoProvider BOOLEAN NOT NULL,           -- Native email/password
    googleProvider BOOLEAN NOT NULL,
    twitterProvider BOOLEAN NOT NULL,
    githubProvider BOOLEAN NOT NULL,
    gitlabProvider BOOLEAN NOT NULL,

    -- SSO Configuration
    ssoProvider BOOLEAN NOT NULL,
    ssoSecret TEXT,
    ssoUrl TEXT,

    -- Display Settings
    defaultSortPolicy TEXT NOT NULL,             -- 'score-desc', 'creationDate-desc', etc.

    -- Analytics
    viewsThisMonth INTEGER,

    FOREIGN KEY (ownerHex) REFERENCES owners(ownerHex)
);
```

**Key Points**:
- Each domain is a separate tenant with isolated data
- OAuth providers can be enabled/disabled per domain
- Moderation policies are configured per domain
- Supports custom SSO integration

#### 3. Commenters Table
**Purpose**: Users who write comments (can comment across multiple domains)

```sql
CREATE TABLE commenters (
    commenterHex TEXT PRIMARY KEY,
    email TEXT,                          -- May be null for OAuth users
    name TEXT NOT NULL,
    link TEXT,                           -- Personal website
    photo TEXT,                          -- Avatar URL
    provider TEXT NOT NULL,              -- 'commento', 'google', 'github', etc.
    joinDate TIMESTAMP NOT NULL,
    state TEXT NOT NULL                  -- 'ok', 'banned'
);
```

**Key Points**:
- Commenters are shared across domains (same user can comment on multiple sites)
- Provider indicates authentication method
- Email may be null for OAuth-only users
- State allows for platform-wide bans

#### 4. Comments Table
**Purpose**: The actual comment content

```sql
CREATE TABLE comments (
    commentHex TEXT PRIMARY KEY,
    domain TEXT NOT NULL,                -- Tenant isolation
    path TEXT NOT NULL,                  -- URL path (e.g., "/blog/post-1")
    commenterHex TEXT NOT NULL,
    markdown TEXT NOT NULL,              -- Original markdown
    html TEXT NOT NULL,                  -- Rendered HTML
    parentHex TEXT NOT NULL,             -- 'root' or parent comment ID (threading)
    score INTEGER NOT NULL,              -- Upvotes - downvotes
    state TEXT NOT NULL,                 -- 'approved', 'unapproved', 'flagged'
    creationDate TIMESTAMP NOT NULL,
    deleted BOOLEAN NOT NULL,
    deleterHex TEXT,                     -- Who deleted (null if not deleted)
    deletionDate TIMESTAMP,

    FOREIGN KEY (domain) REFERENCES domains(domain),
    FOREIGN KEY (commenterHex) REFERENCES commenters(commenterHex)
);

-- Critical indexes for performance
CREATE INDEX idx_comments_domain_path ON comments(domain, path);
CREATE INDEX idx_comments_domain_state ON comments(domain, state);
CREATE INDEX idx_comments_parent ON comments(parentHex);
```

**Key Points**:
- Domain + path identifies where comment appears
- Supports threaded discussions via parentHex
- Both markdown and HTML stored (HTML is sanitized)
- Soft delete (marked deleted but not removed)
- Score allows sorting by popularity

#### 5. Votes Table
**Purpose**: Upvote/downvote system

```sql
CREATE TABLE votes (
    commentHex TEXT NOT NULL,
    commenterHex TEXT NOT NULL,
    direction INTEGER NOT NULL,          -- +1 for upvote, -1 for downvote
    voteDate TIMESTAMP NOT NULL,

    PRIMARY KEY (commentHex, commenterHex),
    FOREIGN KEY (commentHex) REFERENCES comments(commentHex),
    FOREIGN KEY (commenterHex) REFERENCES commenters(commenterHex)
);
```

**Key Points**:
- Composite primary key prevents duplicate votes
- Direction allows changing vote (upvote to downvote)
- Triggers update comment.score automatically

#### 6. Pages Table
**Purpose**: Metadata for each page on a website

```sql
CREATE TABLE pages (
    domain TEXT NOT NULL,
    path TEXT NOT NULL,
    isLocked BOOLEAN NOT NULL,           -- Prevent new comments
    commentCount INTEGER NOT NULL,       -- Cached count
    stickyCommentHex TEXT,               -- Pinned comment
    title TEXT,                          -- Page title (auto-detected)

    PRIMARY KEY (domain, path),
    FOREIGN KEY (domain) REFERENCES domains(domain)
);
```

**Key Points**:
- CommentCount is cached and updated by triggers
- Locking prevents new comments on specific pages
- Sticky comments appear at the top

#### 7. Moderators Table
**Purpose**: Per-domain moderators

```sql
CREATE TABLE moderators (
    domain TEXT NOT NULL,
    email TEXT NOT NULL,
    addDate TIMESTAMP NOT NULL,

    PRIMARY KEY (domain, email),
    FOREIGN KEY (domain) REFERENCES domains(domain)
);
```

**Key Points**:
- Moderators have special permissions for specific domains only
- Identified by email (can moderate before creating commenter account)

#### 8. Session Tables
**Purpose**: Authentication token storage

```sql
CREATE TABLE ownerSessions (
    ownerHex TEXT NOT NULL,
    ownerToken TEXT NOT NULL UNIQUE,
    creationDate TIMESTAMP NOT NULL,
    FOREIGN KEY (ownerHex) REFERENCES owners(ownerHex)
);

CREATE TABLE commenterSessions (
    commenterHex TEXT NOT NULL,
    commenterToken TEXT NOT NULL UNIQUE,
    creationDate TIMESTAMP NOT NULL,
    FOREIGN KEY (commenterHex) REFERENCES commenters(commenterHex)
);
```

**Key Points**:
- Tokens are 32-byte hex strings
- Sessions persist until explicitly deleted
- Separate tables for owners vs commenters

### Database Triggers

#### Auto-update Comment Count
```sql
-- Trigger on comment insert
CREATE TRIGGER commentsInsertTrigger
AFTER INSERT ON comments
FOR EACH ROW
EXECUTE PROCEDURE commentsInsertTriggerFunction();

-- Function increments page.commentCount
CREATE FUNCTION commentsInsertTriggerFunction() RETURNS TRIGGER AS $$
BEGIN
    UPDATE pages
    SET commentCount = commentCount + 1
    WHERE domain = NEW.domain AND path = NEW.path;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

#### Auto-update Vote Score
```sql
-- Trigger on vote insert/update/delete
CREATE TRIGGER votesInsertTrigger
AFTER INSERT OR UPDATE OR DELETE ON votes
FOR EACH ROW
EXECUTE PROCEDURE votesUpdateScore();

-- Recalculates comment.score
CREATE FUNCTION votesUpdateScore() RETURNS TRIGGER AS $$
BEGIN
    UPDATE comments
    SET score = (
        SELECT COALESCE(SUM(direction), 0)
        FROM votes
        WHERE commentHex = NEW.commentHex
    )
    WHERE commentHex = NEW.commentHex;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

---

## Multi-Tenancy Architecture

### Isolation Strategy

Commento implements **domain-based multi-tenancy** where each website (domain) is a logically isolated tenant sharing the same database and application instance.

```
┌─────────────────────────────────────────────────────────────┐
│                    Commento Instance                         │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   Domain A   │  │   Domain B   │  │   Domain C   │    │
│  │ example.com  │  │  myblog.org  │  │  shop.net    │    │
│  │              │  │              │  │              │    │
│  │ - Comments   │  │ - Comments   │  │ - Comments   │    │
│  │ - Moderators │  │ - Moderators │  │ - Moderators │    │
│  │ - OAuth      │  │ - OAuth      │  │ - OAuth      │    │
│  │ - Settings   │  │ - Settings   │  │ - Settings   │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                              │
│              Shared: Database, Commenters, Owners            │
└─────────────────────────────────────────────────────────────┘
```

### Implementation Details

#### 1. Domain Registration Flow

**File**: `api/domain_new.go:10`

```go
func domainNew(ownerHex string, name string, domain string) error {
    // Validate domain format
    if !domainIsValid(domain) {
        return errorInvalidDomain
    }

    // Check if domain already registered
    if exists := domainExists(domain); exists {
        return errorDomainAlreadyExists
    }

    // Create new domain with default settings
    statement := `
        INSERT INTO domains (
            ownerHex, name, domain, creationDate, state,
            autoSpamFilter, requireModeration, requireIdentification,
            moderateAllAnonymous, emailNotificationPolicy,
            commentoProvider, googleProvider, twitterProvider,
            githubProvider, gitlabProvider, ssoProvider,
            defaultSortPolicy, viewsThisMonth
        )
        VALUES ($1, $2, $3, $4, 'unfrozen',
                true, false, false,
                false, 'none',
                true, false, false,
                false, false, false,
                'score-desc', 0);
    `

    _, err := db.Exec(statement, ownerHex, name, domain, time.Now().UTC())
    return err
}
```

#### 2. Request Isolation

Every API request includes a `domain` parameter that enforces data isolation:

**File**: `api/comment_list.go:45`

```go
func commentList(domain string, path string, commenterToken string) ([]comment, error) {
    // All queries filtered by domain
    statement := `
        SELECT commentHex, commenterHex, parentHex, markdown, html,
               score, state, creationDate, deleted
        FROM comments
        WHERE domain = $1 AND path = $2
        ORDER BY creationDate ASC;
    `

    rows, err := db.Query(statement, domain, path)
    // ... process results
}
```

**Critical**: Every database query includes `WHERE domain = $1` to prevent cross-tenant data leakage.

#### 3. Domain Verification

When a website embeds the Commento widget, the backend verifies the domain:

**File**: `api/domain.go:78`

```go
func domainFromRequest(r *http.Request) (string, error) {
    // Get domain from request parameter
    domain := r.FormValue("domain")

    // Verify domain exists in database
    if !domainExists(domain) {
        return "", errorNoSuchDomain
    }

    // Verify domain is not frozen
    d := domainGet(domain)
    if d.State == "frozen" {
        return "", errorDomainFrozen
    }

    return domain, nil
}
```

#### 4. Per-Domain Configuration

Each domain has independent configuration:

| Setting | Scope | Example |
|---------|-------|---------|
| OAuth Providers | Per-domain | Domain A uses Google, Domain B uses GitHub |
| Moderation Policy | Per-domain | Domain A requires approval, Domain B doesn't |
| Moderators | Per-domain | Different moderators for each domain |
| Email Notifications | Per-domain | Domain A sends all, Domain B sends none |
| SSO Settings | Per-domain | Custom SSO per domain |

**File**: `api/domain_update.go:15`

```go
func domainUpdate(ownerToken string, domain string, settings map[string]interface{}) error {
    // Verify ownership
    ownerHex, err := ownerGetByOwnerToken(ownerToken)
    if err != nil {
        return err
    }

    // Verify owner owns this domain
    d := domainGet(domain)
    if d.OwnerHex != ownerHex {
        return errorNotAuthorised
    }

    // Update domain-specific settings
    statement := `
        UPDATE domains
        SET autoSpamFilter = $1,
            requireModeration = $2,
            requireIdentification = $3,
            googleProvider = $4,
            githubProvider = $5
        WHERE domain = $6;
    `

    db.Exec(statement, settings["autoSpamFilter"], settings["requireModeration"],
            settings["requireIdentification"], settings["googleProvider"],
            settings["githubProvider"], domain)

    return nil
}
```

### Data Sharing vs. Isolation

| Data Type | Sharing | Reason |
|-----------|---------|--------|
| **Comments** | Isolated per domain | Core tenant data |
| **Pages** | Isolated per domain | URL paths are domain-specific |
| **Moderators** | Isolated per domain | Different staff per website |
| **Domain Settings** | Isolated per domain | Each site has unique config |
| **Commenters** | Shared across domains | Users can comment on multiple sites |
| **Owners** | Shared across domains | Can manage multiple sites |
| **OAuth Tokens** | Shared across domains | Single sign-on experience |

### Security Considerations

1. **SQL Injection Prevention**: All queries use parameterized statements
2. **Authorization Checks**: Owner can only modify their own domains
3. **Domain Validation**: Domains verified on every request
4. **State Checks**: Frozen domains reject new comments
5. **Moderator Verification**: Checked per-domain, not globally

---

## Data Flow

### Flow 1: Creating a New Comment

```
┌─────────────┐
│   Browser   │
│ example.com │
└──────┬──────┘
       │ 1. User types comment and clicks "Submit"
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  JavaScript (commento.js)                                    │
│  - Captures comment text (markdown)                          │
│  - Gets domain from window.location.host                     │
│  - Gets path from window.location.pathname                   │
│  - Includes commenterToken from localStorage                 │
└──────┬──────────────────────────────────────────────────────┘
       │ 2. POST /api/comment/new
       │    {
       │      domain: "example.com",
       │      path: "/blog/my-post",
       │      markdown: "Great article!",
       │      parentHex: "root",
       │      commenterToken: "abc123..."
       │    }
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend (api/comment_new.go)                                │
│                                                              │
│  Step 1: Parse & Validate Request                            │
│  ├─ Unmarshal JSON body                                      │
│  ├─ Validate domain exists                                   │
│  └─ Check domain not frozen                                  │
│                                                              │
│  Step 2: Authenticate Commenter                              │
│  ├─ Query commenterSessions for commenterToken               │
│  ├─ Get commenterHex                                         │
│  └─ Load commenter details                                   │
│                                                              │
│  Step 3: Load Domain Settings                                │
│  ├─ Query domains table                                      │
│  ├─ Get requireModeration, autoSpamFilter flags              │
│  └─ Check if commenter is moderator                          │
│                                                              │
│  Step 4: Determine Comment State                             │
│  ├─ IF user is moderator → state = "approved"                │
│  ├─ ELSE IF requireModeration = true → state = "unapproved"  │
│  ├─ ELSE IF spam detected → state = "flagged"                │
│  └─ ELSE → state = "approved"                                │
│                                                              │
│  Step 5: Process Markdown                                    │
│  ├─ Sanitize markdown                                        │
│  ├─ Render markdown → HTML                                   │
│  ├─ Escape HTML to prevent XSS                               │
│  └─ Detect code blocks for syntax highlighting               │
│                                                              │
│  Step 6: Generate Comment ID                                 │
│  └─ commentHex = randomHex(32)                               │
└──────┬──────────────────────────────────────────────────────┘
       │ 3. INSERT INTO comments
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  PostgreSQL Database                                         │
│                                                              │
│  INSERT INTO comments (                                      │
│    commentHex, domain, path, commenterHex,                   │
│    parentHex, markdown, html, score,                         │
│    state, creationDate, deleted                              │
│  ) VALUES (                                                  │
│    'def456...', 'example.com', '/blog/my-post',              │
│    'abc123...', 'root', 'Great article!',                    │
│    '<p>Great article!</p>', 0,                               │
│    'approved', NOW(), false                                  │
│  );                                                          │
│                                                              │
│  Trigger: commentsInsertTrigger                              │
│  └─ UPDATE pages SET commentCount += 1                       │
│     WHERE domain='example.com' AND path='/blog/my-post'      │
└──────┬──────────────────────────────────────────────────────┘
       │ 4. Transaction committed
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend (api/comment_new.go) - continued                    │
│                                                              │
│  Step 7: Async Email Notifications                           │
│  └─ go emailNotificationNew(...)                             │
│     ├─ Find parent commenter (if reply)                      │
│     ├─ Query domain email settings                           │
│     ├─ Check if recipient opted-in                           │
│     └─ Send SMTP email                                       │
│                                                              │
│  Step 8: Return Response                                     │
│  └─ JSON: {                                                  │
│       "success": true,                                       │
│       "commentHex": "def456...",                             │
│       "state": "approved",                                   │
│       "html": "<p>Great article!</p>"                        │
│     }                                                        │
└──────┬──────────────────────────────────────────────────────┘
       │ 5. HTTP 200 response
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  JavaScript (commento.js)                                    │
│                                                              │
│  Step 1: Handle Response                                     │
│  ├─ IF state = "approved" → Show immediately                 │
│  ├─ ELSE IF state = "unapproved" → Show pending message      │
│  └─ ELSE IF state = "flagged" → Show spam warning            │
│                                                              │
│  Step 2: Update UI                                           │
│  ├─ Create comment card DOM element                          │
│  ├─ Add avatar, name, timestamp                              │
│  ├─ Insert HTML content                                      │
│  ├─ Add vote buttons, reply button                           │
│  └─ Append to parent or root                                 │
│                                                              │
│  Step 3: Clear Form                                          │
│  └─ Reset textarea, hide preview                             │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────┐
│   Browser   │
│ (Comment    │
│  visible)   │
└─────────────┘
```

### Flow 2: Loading Comments on Page Load

```
┌─────────────┐
│   Browser   │
│ Visits      │
│ example.com │
│ /blog/post  │
└──────┬──────┘
       │ 1. HTML loads
       │    <div id="commento"></div>
       │    <script src="https://cdn.commento.io/js/commento.js"></script>
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  JavaScript (commento.js)                                    │
│                                                              │
│  DOMContentLoaded:                                           │
│  ├─ Find #commento div                                       │
│  ├─ Extract domain from window.location                      │
│  ├─ Extract path from window.location                        │
│  ├─ Check localStorage for commenterToken                    │
│  └─ Call main()                                              │
└──────┬──────────────────────────────────────────────────────┘
       │ 2. POST /api/comment/list
       │    {
       │      domain: "example.com",
       │      path: "/blog/post",
       │      commenterToken: "xyz789..." (or "anonymous")
       │    }
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend (api/comment_list.go)                               │
│                                                              │
│  Step 1: Validate Domain                                     │
│  ├─ Check domain exists                                      │
│  └─ Check domain not frozen                                  │
│                                                              │
│  Step 2: Authenticate (if token provided)                    │
│  ├─ Query commenterSessions                                  │
│  ├─ Get commenterHex                                         │
│  └─ Check if user is moderator for this domain               │
│                                                              │
│  Step 3: Query Comments                                      │
│  └─ SELECT * FROM comments                                   │
│     WHERE domain = 'example.com'                             │
│     AND path = '/blog/post'                                  │
│     AND (state = 'approved'                                  │
│          OR commenterHex = $authenticatedUser)               │
│     ORDER BY creationDate ASC                                │
│                                                              │
│  Step 4: Enrich with Commenter Data                          │
│  └─ For each comment:                                        │
│     SELECT name, photo, link                                 │
│     FROM commenters                                          │
│     WHERE commenterHex = $comment.commenterHex               │
│                                                              │
│  Step 5: Load User's Votes (if authenticated)                │
│  └─ SELECT commentHex, direction                             │
│     FROM votes                                               │
│     WHERE commenterHex = $authenticatedUser                  │
│                                                              │
│  Step 6: Load Domain Settings                                │
│  └─ SELECT requireModeration, requireIdentification,         │
│            commentoProvider, googleProvider, etc.            │
│     FROM domains                                             │
│     WHERE domain = 'example.com'                             │
└──────┬──────────────────────────────────────────────────────┘
       │ 3. Return JSON
       │    {
       │      "success": true,
       │      "comments": [
       │        {
       │          "commentHex": "abc123",
       │          "commenterHex": "def456",
       │          "markdown": "Great!",
       │          "html": "<p>Great!</p>",
       │          "parentHex": "root",
       │          "score": 5,
       │          "state": "approved",
       │          "creationDate": "2023-11-13T10:30:00Z"
       │        },
       │        {...}
       │      ],
       │      "commenters": {
       │        "def456": {
       │          "name": "John Doe",
       │          "photo": "https://...",
       │          "link": "https://..."
       │        }
       │      },
       │      "requireModeration": false,
       │      "isModerator": false,
       │      "configuredOauths": ["google", "github"]
       │    }
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  JavaScript (commento.js)                                    │
│                                                              │
│  Step 1: Build Comment Tree                                  │
│  ├─ Create map of commentHex → comment                       │
│  ├─ Group comments by parentHex                              │
│  └─ Recursively build tree structure                         │
│                                                              │
│  Step 2: Render Comments                                     │
│  └─ For each root comment:                                   │
│     ├─ Create card div                                       │
│     ├─ Add avatar (from commenters map)                      │
│     ├─ Add name and timestamp                                │
│     ├─ Insert HTML content                                   │
│     ├─ Add vote buttons (↑ score ↓)                          │
│     ├─ Highlight user's vote if exists                       │
│     ├─ Add reply button                                      │
│     └─ Recursively render children (indent)                  │
│                                                              │
│  Step 3: Setup Event Listeners                               │
│  ├─ Vote buttons → POST /api/comment/vote                    │
│  ├─ Reply buttons → Show reply form                          │
│  ├─ Edit buttons → Load markdown in form                     │
│  └─ Delete buttons → POST /api/comment/delete                │
│                                                              │
│  Step 4: Render Login/Signup UI                              │
│  └─ Show OAuth buttons based on configuredOauths             │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────┐
│   Browser   │
│ (Comments   │
│  rendered)  │
└─────────────┘
```

### Flow 3: OAuth Authentication (Google Example)

```
┌─────────────┐
│   Browser   │
│ User clicks │
│"Login with  │
│  Google"    │
└──────┬──────┘
       │ 1. Click event
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  JavaScript (commento.js)                                    │
│  └─ window.location = "/api/oauth/google/redirect"           │
│                       + "?domain=example.com"                │
└──────┬──────────────────────────────────────────────────────┘
       │ 2. GET /api/oauth/google/redirect?domain=example.com
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend (api/oauth_google.go)                               │
│                                                              │
│  Step 1: Generate OAuth State Token                          │
│  ├─ state = randomHex(32)                                    │
│  ├─ Store in session: states[state] = domain                 │
│  └─ Set cookie: commento_oauth_state = state                 │
│                                                              │
│  Step 2: Build Google OAuth URL                              │
│  └─ https://accounts.google.com/o/oauth2/auth?               │
│     client_id=xxx&                                           │
│     redirect_uri=https://commento.io/api/oauth/google/callback&│
│     response_type=code&                                      │
│     scope=profile+email&                                     │
│     state=abc123...                                          │
│                                                              │
│  Step 3: Redirect                                            │
│  └─ HTTP 302 → Google                                        │
└──────┬──────────────────────────────────────────────────────┘
       │ 3. Redirect to Google
       │
       ▼
┌─────────────────┐
│  Google OAuth   │
│  User approves  │
│  permissions    │
└──────┬──────────┘
       │ 4. Google redirects back
       │    /api/oauth/google/callback?code=xyz&state=abc123
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend (api/oauth_google.go) - Callback Handler            │
│                                                              │
│  Step 1: Verify State Token                                  │
│  ├─ Read cookie: commento_oauth_state                        │
│  ├─ Compare with state parameter                             │
│  └─ Prevent CSRF attacks                                     │
│                                                              │
│  Step 2: Exchange Code for Access Token                      │
│  └─ POST https://oauth2.googleapis.com/token                 │
│     ├─ code = xyz                                            │
│     ├─ client_id = xxx                                       │
│     ├─ client_secret = yyy                                   │
│     └─ grant_type = authorization_code                       │
│                                                              │
│  Step 3: Fetch User Profile                                  │
│  └─ GET https://www.googleapis.com/oauth2/v2/userinfo        │
│     Authorization: Bearer {access_token}                     │
│     Response:                                                │
│     {                                                        │
│       "id": "12345",                                         │
│       "email": "user@gmail.com",                             │
│       "name": "John Doe",                                    │
│       "picture": "https://..."                               │
│     }                                                        │
│                                                              │
│  Step 4: Create or Update Commenter                          │
│  └─ Check if commenter exists:                               │
│     SELECT * FROM commenters                                 │
│     WHERE provider = 'google' AND email = 'user@gmail.com'   │
│     ├─ IF EXISTS → Update photo, name                        │
│     └─ IF NOT EXISTS → INSERT new commenter                  │
│                                                              │
│  Step 5: Create Session                                      │
│  ├─ commenterToken = randomHex(32)                           │
│  └─ INSERT INTO commenterSessions                            │
│     (commenterToken, commenterHex, creationDate)             │
│     VALUES (token, hex, NOW())                               │
│                                                              │
│  Step 6: Redirect Back to Website                            │
│  └─ HTTP 302 → example.com/blog/post                         │
│     + Set cookie: commentoCommenterToken = token             │
└──────┬──────────────────────────────────────────────────────┘
       │ 5. Redirect to original page
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  JavaScript (commento.js)                                    │
│                                                              │
│  Step 1: Detect Cookie                                       │
│  └─ commentoCommenterToken exists                            │
│                                                              │
│  Step 2: Store in localStorage                               │
│  └─ localStorage.setItem("commentoCommenterToken", token)    │
│                                                              │
│  Step 3: Fetch User Info                                     │
│  └─ POST /api/commenter/self                                 │
│     { commenterToken: token }                                │
│                                                              │
│  Step 4: Update UI                                           │
│  ├─ Show "Logged in as John Doe"                             │
│  ├─ Show avatar                                              │
│  ├─ Enable comment form                                      │
│  └─ Show logout button                                       │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────┐
│   Browser   │
│ (User is    │
│  logged in) │
└─────────────┘
```

### Flow 4: Spam Detection & Moderation

```
Comment Submitted
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend Spam Check (api/spam.go)                            │
│                                                              │
│  IF domain.autoSpamFilter = false                            │
│  └─ Skip spam detection                                      │
│                                                              │
│  ELSE:                                                       │
│  ├─ Check against Akismet API                                │
│  │  POST https://[key].rest.akismet.com/1.1/comment-check    │
│  │  {                                                        │
│  │    blog: "https://example.com",                           │
│  │    user_ip: "1.2.3.4",                                    │
│  │    comment_content: "comment text",                       │
│  │    comment_author: "John Doe",                            │
│  │    comment_author_email: "user@example.com"               │
│  │  }                                                        │
│  │                                                           │
│  └─ Response: "true" (spam) or "false" (not spam)            │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Determine Comment State                                     │
│                                                              │
│  Priority Order:                                             │
│  1. IF commenter is moderator                                │
│     └─ state = "approved" (always)                           │
│                                                              │
│  2. ELSE IF domain.requireModeration = true                  │
│     └─ state = "unapproved" (all comments need approval)     │
│                                                              │
│  3. ELSE IF spam detected                                    │
│     └─ state = "flagged"                                     │
│                                                              │
│  4. ELSE IF domain.moderateAllAnonymous = true               │
│     AND commenter is anonymous                               │
│     └─ state = "unapproved"                                  │
│                                                              │
│  5. ELSE                                                     │
│     └─ state = "approved"                                    │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│  Insert Comment with State                                   │
│  INSERT INTO comments (..., state) VALUES (..., $state)      │
└──────┬──────────────────────────────────────────────────────┘
       │
       ├──────────────┬──────────────┬──────────────┐
       │              │              │              │
       ▼              ▼              ▼              ▼
  state =        state =        state =        state =
  "approved"     "unapproved"   "flagged"      "approved"
       │              │              │              │
       │              │              │              │
       ▼              ▼              ▼              ▼
  Visible        Hidden         Hidden         Visible
  immediately    (pending)      (spam)         immediately
       │              │              │              │
       │              │              │              │
       │              ▼              ▼              │
       │         ┌────────────────────────┐        │
       │         │ Moderator Dashboard    │        │
       │         │ - View pending/flagged │        │
       │         │ - Approve/Delete       │        │
       │         └───────┬────────────────┘        │
       │                 │                         │
       │                 ▼                         │
       │         POST /api/comment/approve         │
       │         └─ UPDATE comments                │
       │            SET state = 'approved'         │
       │                 │                         │
       └─────────────────┴─────────────────────────┘
                         │
                         ▼
                   Now visible
                   to all users
```

---

## Frontend Architecture

### Embedded Widget Structure

```
Website (example.com/blog/post)
├── <head>
│   └── (no dependencies required)
├── <body>
│   ├── ... website content ...
│   ├── <div id="commento"></div>        ← Injection point
│   └── <script src="https://commento.io/js/commento.js"></script>
└── (Widget renders inside #commento div)
```

### JavaScript Architecture

**File**: `frontend/js/commento.js` (Main embed script)

#### Initialization Flow

```javascript
(function(global) {
  "use strict";

  // 1. Configuration
  var origin = "https://commento.io"; // API endpoint
  var cdn = "https://cdn.commento.io"; // Static assets

  // 2. Detect attributes from script tag
  var scripts = document.getElementsByTagName("script");
  var script = scripts[scripts.length - 1];
  var cssOverride = script.getAttribute("data-css-override");
  var autoInit = script.getAttribute("data-auto-init");
  var noFonts = script.getAttribute("data-no-fonts");
  var hideDeleted = script.getAttribute("data-hide-deleted");
  var pageId = script.getAttribute("data-page-id");

  // 3. Extract page info
  var domain = window.location.host;
  var path = pageId || window.location.pathname;

  // 4. Main initialization
  function main() {
    loadCss();
    loadComments();
    renderLoginForm();
  }

  // 5. Auto-start on DOMContentLoaded
  if (autoInit !== "false") {
    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", main);
    } else {
      main();
    }
  }

  // 6. Expose global API
  global.commento = main;
})(window);
```

#### Key Functions

**1. Load Comments**
```javascript
function loadComments() {
  var commenterToken = localStorage.getItem("commentoCommenterToken") || "anonymous";

  post(origin + "/api/comment/list", {
    domain: domain,
    path: path,
    commenterToken: commenterToken
  }, function(resp) {
    if (resp.success) {
      renderComments(resp.comments, resp.commenters);
      renderSettings(resp);
    }
  });
}
```

**2. Build Comment Tree**
```javascript
function buildCommentTree(comments) {
  var commentMap = {};
  var rootComments = [];

  // First pass: create map
  comments.forEach(function(c) {
    commentMap[c.commentHex] = c;
    c.children = [];
  });

  // Second pass: build tree
  comments.forEach(function(c) {
    if (c.parentHex === "root") {
      rootComments.push(c);
    } else {
      var parent = commentMap[c.parentHex];
      if (parent) {
        parent.children.push(c);
      }
    }
  });

  return rootComments;
}
```

**3. Render Comment Card**
```javascript
function renderComment(comment, commenter, depth) {
  var card = createElement("div", "commento-comment-card");
  card.setAttribute("data-comment-hex", comment.commentHex);

  // Avatar
  var avatar = createElement("img", "commento-avatar");
  avatar.src = commenter.photo || "/img/default-avatar.png";

  // Header (name, timestamp)
  var header = createElement("div", "commento-comment-header");
  var name = createElement("span", "commento-name");
  name.innerText = commenter.name;
  var timestamp = createElement("span", "commento-timestamp");
  timestamp.innerText = timeAgo(comment.creationDate);
  header.appendChild(name);
  header.appendChild(timestamp);

  // Body (comment content)
  var body = createElement("div", "commento-comment-body");
  body.innerHTML = comment.html; // Already sanitized by backend

  // Actions (vote, reply, edit, delete)
  var actions = createElement("div", "commento-actions");

  // Upvote button
  var upvote = createElement("button", "commento-upvote");
  upvote.innerHTML = "↑";
  upvote.onclick = function() {
    vote(comment.commentHex, 1);
  };

  // Score display
  var score = createElement("span", "commento-score");
  score.innerText = comment.score;

  // Downvote button
  var downvote = createElement("button", "commento-downvote");
  downvote.innerHTML = "↓";
  downvote.onclick = function() {
    vote(comment.commentHex, -1);
  };

  // Reply button
  var reply = createElement("button", "commento-reply");
  reply.innerText = "Reply";
  reply.onclick = function() {
    showReplyForm(comment.commentHex);
  };

  actions.appendChild(upvote);
  actions.appendChild(score);
  actions.appendChild(downvote);
  actions.appendChild(reply);

  // Assemble card
  card.appendChild(avatar);
  card.appendChild(header);
  card.appendChild(body);
  card.appendChild(actions);

  // Recursively render children
  if (comment.children && comment.children.length > 0) {
    var childContainer = createElement("div", "commento-children");
    childContainer.style.marginLeft = (depth * 30) + "px";
    comment.children.forEach(function(child) {
      childContainer.appendChild(renderComment(child, commenters[child.commenterHex], depth + 1));
    });
    card.appendChild(childContainer);
  }

  return card;
}
```

**4. Submit Comment**
```javascript
function submitComment(parentHex, markdown) {
  var commenterToken = localStorage.getItem("commentoCommenterToken");

  if (!commenterToken) {
    alert("Please login to comment");
    return;
  }

  post(origin + "/api/comment/new", {
    domain: domain,
    path: path,
    commenterToken: commenterToken,
    parentHex: parentHex,
    markdown: markdown
  }, function(resp) {
    if (resp.success) {
      if (resp.state === "approved") {
        // Refresh comments to show new one
        loadComments();
      } else if (resp.state === "unapproved") {
        alert("Your comment is pending moderation");
      } else if (resp.state === "flagged") {
        alert("Your comment was flagged as spam");
      }
    } else {
      alert("Error: " + resp.message);
    }
  });
}
```

**5. Vote Handler**
```javascript
function vote(commentHex, direction) {
  var commenterToken = localStorage.getItem("commentoCommenterToken");

  if (!commenterToken) {
    alert("Please login to vote");
    return;
  }

  post(origin + "/api/comment/vote", {
    commentHex: commentHex,
    commenterToken: commenterToken,
    direction: direction
  }, function(resp) {
    if (resp.success) {
      // Update score in UI
      var scoreEl = document.querySelector(
        '[data-comment-hex="' + commentHex + '"] .commento-score'
      );
      scoreEl.innerText = resp.score;

      // Update button states
      updateVoteButtons(commentHex, direction);
    }
  });
}
```

### Dashboard Architecture

**File**: `frontend/js/dashboard.js`

The dashboard is a separate Vue.js application for domain owners and moderators.

```
Dashboard Structure:
├── Login Page (/login)
├── Dashboard Home (/dashboard)
│   ├── Domain Switcher
│   ├── Statistics (views, comments, spam blocked)
│   └── Quick Actions
├── Moderation (/dashboard/moderation)
│   ├── Pending Comments Tab
│   ├── Flagged Comments Tab
│   └── All Comments Tab
├── Settings (/dashboard/settings)
│   ├── General Settings
│   ├── Authentication Providers
│   ├── Moderators Management
│   ├── Email Settings
│   └── Import/Export
└── Domain Management (/dashboard/domains)
    ├── Add New Domain
    ├── Domain List
    └── Domain Settings per site
```

**Vue Component Example**:
```javascript
var moderationApp = new Vue({
  el: '#moderation-app',
  data: {
    pendingComments: [],
    flaggedComments: [],
    loading: true,
    currentDomain: ""
  },
  methods: {
    loadPendingComments: function() {
      this.loading = true;
      post(origin + "/api/comment/list/pending", {
        ownerToken: getOwnerToken(),
        domain: this.currentDomain
      }, function(resp) {
        this.pendingComments = resp.comments;
        this.loading = false;
      }.bind(this));
    },

    approveComment: function(commentHex) {
      post(origin + "/api/comment/approve", {
        ownerToken: getOwnerToken(),
        commentHex: commentHex
      }, function(resp) {
        if (resp.success) {
          // Remove from pending list
          this.pendingComments = this.pendingComments.filter(function(c) {
            return c.commentHex !== commentHex;
          });
        }
      }.bind(this));
    },

    deleteComment: function(commentHex) {
      if (!confirm("Delete this comment?")) return;

      post(origin + "/api/comment/delete", {
        ownerToken: getOwnerToken(),
        commentHex: commentHex
      }, function(resp) {
        if (resp.success) {
          this.loadPendingComments(); // Refresh
        }
      }.bind(this));
    }
  },

  mounted: function() {
    this.currentDomain = getDomainFromUrl();
    this.loadPendingComments();
  }
});
```

### Build Process

**File**: `frontend/gulpfile.js`

```javascript
var gulp = require('gulp');
var sass = require('gulp-sass');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var cleanCSS = require('gulp-clean-css');

// Compile SASS to CSS
gulp.task('sass', function() {
  return gulp.src('css/**/*.scss')
    .pipe(sass())
    .pipe(cleanCSS()) // Minify
    .pipe(gulp.dest('build/css'));
});

// Bundle JavaScript
gulp.task('js-commento', function() {
  return gulp.src([
    'js/utils.js',
    'js/commento.js',
    'js/render.js',
    'js/auth.js'
  ])
  .pipe(concat('commento.js'))
  .pipe(uglify()) // Minify
  .pipe(gulp.dest('build/js'));
});

gulp.task('js-dashboard', function() {
  return gulp.src([
    'js/dashboard-*.js'
  ])
  .pipe(concat('dashboard.js'))
  .pipe(uglify())
  .pipe(gulp.dest('build/js'));
});

// Watch for changes
gulp.task('watch', function() {
  gulp.watch('css/**/*.scss', ['sass']);
  gulp.watch('js/**/*.js', ['js-commento', 'js-dashboard']);
});

gulp.task('default', ['sass', 'js-commento', 'js-dashboard']);
```

### Styling System

**File**: `frontend/css/commento.scss`

```scss
// Variables
$primary-color: #7f8ea3;
$background-color: #fff;
$border-color: #e1e1e1;
$text-color: #333;

// Comment card
.commento-comment-card {
  background: $background-color;
  border: 1px solid $border-color;
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 15px;

  .commento-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    float: left;
    margin-right: 10px;
  }

  .commento-comment-header {
    font-size: 14px;
    color: #666;
    margin-bottom: 8px;

    .commento-name {
      font-weight: bold;
      color: $text-color;
    }

    .commento-timestamp {
      color: #999;
      margin-left: 8px;
    }
  }

  .commento-comment-body {
    clear: both;
    line-height: 1.6;

    p {
      margin: 0 0 10px 0;
    }

    code {
      background: #f5f5f5;
      padding: 2px 4px;
      border-radius: 2px;
    }

    pre {
      background: #f5f5f5;
      padding: 10px;
      overflow-x: auto;
    }
  }

  .commento-actions {
    margin-top: 10px;

    button {
      background: none;
      border: 1px solid $border-color;
      padding: 5px 10px;
      margin-right: 5px;
      cursor: pointer;
      border-radius: 3px;

      &:hover {
        background: #f5f5f5;
      }
    }
  }

  .commento-children {
    margin-top: 15px;
    padding-left: 30px;
    border-left: 2px solid $border-color;
  }
}

// Custom CSS override support
.commento-root[data-css-override] {
  // User can override any styles via data-css-override attribute
}
```

---

## Backend Architecture

### Application Structure

```
api/
├── main.go                 # Entry point, server initialization
├── router_*.go             # Route definitions
│   ├── router_api.go       # API endpoints
│   ├── router_static.go    # Static file serving
│   └── router_oauth.go     # OAuth routes
├── config.go               # Environment variable parsing
├── database.go             # PostgreSQL connection
├── migrate.go              # Database migrations
│
├── owner_*.go              # Owner management
│   ├── owner_new.go
│   ├── owner_login.go
│   ├── owner_get.go
│   └── owner_confirm_hex.go
│
├── domain_*.go             # Domain management
│   ├── domain_new.go
│   ├── domain_get.go
│   ├── domain_update.go
│   ├── domain_moderator_*.go
│   └── domain_import_disqus.go
│
├── commenter_*.go          # Commenter management
│   ├── commenter_new.go
│   ├── commenter_login.go
│   ├── commenter_get.go
│   └── commenter_update.go
│
├── comment_*.go            # Comment operations
│   ├── comment_new.go      # Create comment
│   ├── comment_list.go     # Fetch comments
│   ├── comment_get.go      # Get single comment
│   ├── comment_vote.go     # Upvote/downvote
│   ├── comment_approve.go  # Moderator approval
│   ├── comment_delete.go   # Delete comment
│   └── comment_edit.go     # Edit comment
│
├── oauth_*.go              # OAuth integration
│   ├── oauth_google.go
│   ├── oauth_github.go
│   ├── oauth_twitter.go
│   └── oauth_gitlab.go
│
├── email_*.go              # Email notifications
│   ├── email_notification_new.go
│   ├── email_update.go
│   └── email_templates.go
│
├── spam.go                 # Spam detection (Akismet)
├── markdown.go             # Markdown to HTML conversion
├── utils.go                # Helper functions
├── errors.go               # Error definitions
└── session.go              # Session management
```

### Main Application Bootstrap

**File**: `api/main.go:10`

```go
package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
)

func main() {
    fmt.Println("Commento starting...")

    // 1. Initialize logger
    loggerCreate()
    logger.Infof("starting Commento")

    // 2. Parse configuration from environment
    if err := configParse(); err != nil {
        logger.Fatalf("cannot parse config: %v", err)
    }

    // 3. Connect to PostgreSQL
    if err := dbConnect(5); err != nil { // 5 retries
        logger.Fatalf("cannot connect to database: %v", err)
    }
    defer db.Close()

    // 4. Run migrations
    if err := migrate(); err != nil {
        logger.Fatalf("cannot run migrations: %v", err)
    }

    // 5. Configure SMTP for emails
    if err := smtpConfigure(); err != nil {
        logger.Warningf("cannot configure SMTP: %v", err)
    }

    // 6. Configure OAuth providers
    if err := oauthConfigure(); err != nil {
        logger.Warningf("cannot configure OAuth: %v", err)
    }

    // 7. Initialize markdown renderer
    markdownRendererCreate()

    // 8. Setup HTTP routes
    router := routesServe()

    // 9. Start server
    address := fmt.Sprintf("%s:%d", os.Getenv("BIND_ADDRESS"), os.Getenv("PORT"))
    logger.Infof("starting server on %s", address)

    server := &http.Server{
        Addr:         address,
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    if err := server.ListenAndServe(); err != nil {
        logger.Fatalf("cannot start server: %v", err)
    }
}
```

### Configuration Management

**File**: `api/config.go:15`

```go
package main

import (
    "os"
    "strconv"
)

type config struct {
    // Database
    Postgres string

    // Server
    BindAddress string
    Port        int
    Origin      string
    CDNPrefix   string

    // SMTP
    SMTPHost     string
    SMTPPort     int
    SMTPUsername string
    SMTPPassword string
    SMTPFromAddress string

    // OAuth
    GoogleKey    string
    GoogleSecret string
    GithubKey    string
    GithubSecret string
    TwitterKey   string
    TwitterSecret string
    GitlabKey    string
    GitlabSecret string

    // Features
    ForbidNewOwners bool
    AkismetKey      string
}

var cfg config

func configParse() error {
    // Required
    cfg.Postgres = os.Getenv("COMMENTO_POSTGRES")
    if cfg.Postgres == "" {
        return fmt.Errorf("COMMENTO_POSTGRES not set")
    }

    cfg.Origin = os.Getenv("COMMENTO_ORIGIN")
    if cfg.Origin == "" {
        return fmt.Errorf("COMMENTO_ORIGIN not set")
    }

    // Optional with defaults
    cfg.BindAddress = getEnvOrDefault("COMMENTO_BIND_ADDRESS", "127.0.0.1")
    cfg.Port = getEnvIntOrDefault("COMMENTO_PORT", 8080)
    cfg.CDNPrefix = getEnvOrDefault("COMMENTO_CDN_PREFIX", "")

    // SMTP
    cfg.SMTPHost = os.Getenv("COMMENTO_SMTP_HOST")
    cfg.SMTPPort = getEnvIntOrDefault("COMMENTO_SMTP_PORT", 587)
    cfg.SMTPUsername = os.Getenv("COMMENTO_SMTP_USERNAME")
    cfg.SMTPPassword = os.Getenv("COMMENTO_SMTP_PASSWORD")
    cfg.SMTPFromAddress = os.Getenv("COMMENTO_SMTP_FROM_ADDRESS")

    // OAuth
    cfg.GoogleKey = os.Getenv("COMMENTO_GOOGLE_KEY")
    cfg.GoogleSecret = os.Getenv("COMMENTO_GOOGLE_SECRET")
    cfg.GithubKey = os.Getenv("COMMENTO_GITHUB_KEY")
    cfg.GithubSecret = os.Getenv("COMMENTO_GITHUB_SECRET")

    // Features
    cfg.ForbidNewOwners = os.Getenv("COMMENTO_FORBID_NEW_OWNERS") == "true"
    cfg.AkismetKey = os.Getenv("COMMENTO_AKISMET_KEY")

    return nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

### Database Connection

**File**: `api/database.go:20`

```go
package main

import (
    "database/sql"
    "time"
    _ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

func dbConnect(retries int) error {
    var err error

    for i := 0; i < retries; i++ {
        db, err = sql.Open("postgres", cfg.Postgres)
        if err != nil {
            logger.Warningf("cannot connect to database (attempt %d/%d): %v",
                          i+1, retries, err)
            time.Sleep(5 * time.Second)
            continue
        }

        // Test connection
        if err = db.Ping(); err != nil {
            logger.Warningf("cannot ping database (attempt %d/%d): %v",
                          i+1, retries, err)
            time.Sleep(5 * time.Second)
            continue
        }

        // Connection successful
        db.SetMaxOpenConns(50)
        db.SetMaxIdleConns(10)
        db.SetConnMaxLifetime(time.Hour)

        logger.Infof("connected to database")
        return nil
    }

    return fmt.Errorf("failed to connect after %d attempts: %v", retries, err)
}
```

### API Route Handler Example

**File**: `api/comment_new.go:35`

```go
package main

import (
    "net/http"
    "time"
)

type commentNewRequest struct {
    CommenterToken string `json:"commenterToken"`
    Domain         string `json:"domain"`
    Path           string `json:"path"`
    ParentHex      string `json:"parentHex"`
    Markdown       string `json:"markdown"`
}

type commentNewResponse struct {
    Success    bool   `json:"success"`
    CommentHex string `json:"commentHex"`
    State      string `json:"state"`
    HTML       string `json:"html"`
}

func commentNewHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req commentNewRequest
    if err := bodyUnmarshal(r, &req); err != nil {
        writeError(w, err)
        return
    }

    // 2. Validate domain
    domain, err := domainGet(req.Domain)
    if err != nil {
        writeError(w, err)
        return
    }

    if domain.State == "frozen" {
        writeError(w, errorDomainFrozen)
        return
    }

    // 3. Authenticate commenter
    commenterHex, err := commenterGetByCommenterToken(req.CommenterToken)
    if err != nil {
        writeError(w, errorNotAuthorised)
        return
    }

    commenter, err := commenterGetByHex(commenterHex)
    if err != nil {
        writeError(w, err)
        return
    }

    if commenter.State == "banned" {
        writeError(w, errorUserBanned)
        return
    }

    // 4. Check if moderator
    isModerator, err := moderatorCheck(req.Domain, commenter.Email)
    if err != nil {
        isModerator = false
    }

    // 5. Determine comment state
    state := "approved"

    if !isModerator {
        if domain.RequireModeration {
            state = "unapproved"
        } else if domain.ModerateAllAnonymous && commenter.Provider == "commento" && commenter.Email == "" {
            state = "unapproved"
        } else if domain.AutoSpamFilter {
            isSpam, err := isSpamAkismet(req.Domain, req.Markdown, commenter.Name,
                                        commenter.Email, r.RemoteAddr)
            if err == nil && isSpam {
                state = "flagged"
            }
        }
    }

    // 6. Validate markdown
    if len(req.Markdown) > 10000 {
        writeError(w, errorCommentTooLong)
        return
    }

    if req.Markdown == "" {
        writeError(w, errorCommentEmpty)
        return
    }

    // 7. Sanitize and render markdown
    markdown := sanitizeMarkdown(req.Markdown)
    html, err := markdownToHTML(markdown)
    if err != nil {
        writeError(w, err)
        return
    }

    // 8. Create comment
    commentHex := randomHex(32)
    creationDate := time.Now().UTC()

    statement := `
        INSERT INTO comments (
            commentHex, domain, path, commenterHex,
            parentHex, markdown, html, score,
            state, creationDate, deleted
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, 0, $8, $9, false);
    `

    _, err = db.Exec(statement, commentHex, req.Domain, req.Path, commenterHex,
                     req.ParentHex, markdown, html, state, creationDate)
    if err != nil {
        logger.Errorf("cannot insert comment: %v", err)
        writeError(w, errorInternal)
        return
    }

    // 9. Send email notifications (async)
    go emailNotificationNew(req.Domain, req.Path, commenterHex, commentHex,
                            html, req.ParentHex, state)

    // 10. Return response
    writeJSON(w, commentNewResponse{
        Success:    true,
        CommentHex: commentHex,
        State:      state,
        HTML:       html,
    })
}
```

### Middleware: CORS Handler

**File**: `api/router_api.go:120`

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow all origins (since embeddable widget)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // Handle preflight
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## Authentication & Security

### Password Hashing

**File**: `api/owner_new.go:25`

```go
import "golang.org/x/crypto/bcrypt"

func ownerNew(email, name, password string) (string, error) {
    // Generate bcrypt hash (cost = 12)
    passwordHash, err := bcrypt.GenerateFromPassword(
        []byte(password),
        12,
    )
    if err != nil {
        return "", err
    }

    ownerHex := randomHex(32)

    statement := `
        INSERT INTO owners (ownerHex, email, name, passwordHash,
                           confirmedEmail, joinDate)
        VALUES ($1, $2, $3, $4, false, $5);
    `

    _, err = db.Exec(statement, ownerHex, email, name,
                     string(passwordHash), time.Now().UTC())

    return ownerHex, err
}

func ownerLogin(email, password string) (string, error) {
    var ownerHex string
    var passwordHash string

    // Fetch hash from database
    err := db.QueryRow(`
        SELECT ownerHex, passwordHash
        FROM owners
        WHERE email = $1
    `, email).Scan(&ownerHex, &passwordHash)

    if err != nil {
        return "", errorInvalidLogin
    }

    // Compare password with hash
    err = bcrypt.CompareHashAndPassword(
        []byte(passwordHash),
        []byte(password),
    )

    if err != nil {
        return "", errorInvalidLogin
    }

    // Create session token
    ownerToken := randomHex(32)

    _, err = db.Exec(`
        INSERT INTO ownerSessions (ownerToken, ownerHex, creationDate)
        VALUES ($1, $2, $3)
    `, ownerToken, ownerHex, time.Now().UTC())

    return ownerToken, err
}
```

### SQL Injection Prevention

All queries use parameterized statements:

```go
// ❌ VULNERABLE (string concatenation)
query := "SELECT * FROM comments WHERE domain = '" + domain + "'"
db.Query(query)

// ✅ SECURE (parameterized)
statement := "SELECT * FROM comments WHERE domain = $1"
db.Query(statement, domain)
```

### XSS Prevention

**File**: `api/markdown.go:40`

```go
import "github.com/microcosm-cc/bluemonday"

var htmlPolicy = bluemonday.UGCPolicy()

func markdownToHTML(markdown string) (string, error) {
    // 1. Render markdown
    html := blackfriday.Run([]byte(markdown))

    // 2. Sanitize HTML (remove dangerous tags/attributes)
    safeHTML := htmlPolicy.SanitizeBytes(html)

    return string(safeHTML), nil
}

// Bluemonday policy allows:
// - Safe HTML tags: p, b, i, a, code, pre, ul, ol, li, blockquote
// - Safe attributes: href (for links), class (for code highlighting)
// - Removes: script, iframe, object, embed, style, onclick, onerror, etc.
```

### CSRF Protection

OAuth state tokens prevent CSRF attacks:

```go
func oauthGoogleRedirect(w http.ResponseWriter, r *http.Request) {
    // Generate random state token
    state := randomHex(32)

    // Store in session (server-side)
    oauthStates[state] = time.Now()

    // Also set as HTTP-only cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "commento_oauth_state",
        Value:    state,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    })

    // Redirect to Google with state parameter
    url := fmt.Sprintf(
        "https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=profile+email&state=%s",
        cfg.GoogleKey,
        cfg.Origin+"/api/oauth/google/callback",
        state,
    )

    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
    // Verify state matches
    stateFromURL := r.URL.Query().Get("state")
    stateFromCookie, _ := r.Cookie("commento_oauth_state")

    if stateFromURL != stateFromCookie.Value {
        writeError(w, errorInvalidOAuthState)
        return
    }

    // Delete used state
    delete(oauthStates, stateFromURL)

    // Continue with OAuth flow...
}
```

### Rate Limiting (Recommended Implementation)

```go
// Not currently in codebase, but recommended:
import "golang.org/x/time/rate"

var commentLimiter = rate.NewLimiter(1, 5) // 1 per second, burst of 5

func commentNewHandler(w http.ResponseWriter, r *http.Request) {
    // Rate limit per IP
    if !commentLimiter.Allow() {
        writeError(w, errorRateLimitExceeded)
        return
    }

    // ... rest of handler
}
```

---

## Customization Guide

### 1. Styling Customization

#### Option A: CSS Override Attribute

Add `data-css-override` to the script tag:

```html
<script
  src="https://commento.example.com/js/commento.js"
  data-css-override="/my-custom-commento.css"
></script>
```

Then create `/my-custom-commento.css`:

```css
/* Override default styles */
.commento-root {
  font-family: 'Your Custom Font', sans-serif;
}

.commento-comment-card {
  background: #f9f9f9;
  border-left: 4px solid #4CAF50;
}

.commento-upvote:hover {
  color: #4CAF50;
}

.commento-downvote:hover {
  color: #f44336;
}
```

#### Option B: Modify Source SCSS

**File**: `frontend/css/commento.scss`

```scss
// Change primary color
$primary-color: #4CAF50; // Green instead of default

// Change fonts
$font-family: 'Roboto', sans-serif;

// Rebuild
// Run: cd frontend && gulp sass
```

#### Option C: Disable Default Styles

```html
<script
  src="https://commento.example.com/js/commento.js"
  data-no-fonts="true"
  data-css-override="/completely-custom.css"
></script>
```

### 2. Adding a New OAuth Provider

#### Step 1: Add Configuration

**File**: `api/config.go`

```go
type config struct {
    // ... existing fields ...

    // New provider
    LinkedInKey    string
    LinkedInSecret string
}

func configParse() error {
    // ... existing code ...

    cfg.LinkedInKey = os.Getenv("COMMENTO_LINKEDIN_KEY")
    cfg.LinkedInSecret = os.Getenv("COMMENTO_LINKEDIN_SECRET")
}
```

#### Step 2: Add Database Column

**File**: `db/new_migration.sql`

```sql
ALTER TABLE domains ADD COLUMN linkedinProvider BOOLEAN DEFAULT false;
```

#### Step 3: Implement OAuth Handler

**File**: `api/oauth_linkedin.go`

```go
package main

func oauthLinkedInRedirect(w http.ResponseWriter, r *http.Request) {
    domain := r.URL.Query().Get("domain")
    state := randomHex(32)
    oauthStates[state] = time.Now()

    http.SetCookie(w, &http.Cookie{
        Name:     "commento_oauth_state",
        Value:    state,
        HttpOnly: true,
    })

    url := fmt.Sprintf(
        "https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=r_liteprofile+r_emailaddress",
        cfg.LinkedInKey,
        cfg.Origin+"/api/oauth/linkedin/callback",
        state,
    )

    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthLinkedInCallback(w http.ResponseWriter, r *http.Request) {
    // 1. Verify state
    state := r.URL.Query().Get("state")
    if !verifyOAuthState(r, state) {
        writeError(w, errorInvalidOAuthState)
        return
    }

    // 2. Exchange code for token
    code := r.URL.Query().Get("code")
    tokenResp, err := http.PostForm("https://www.linkedin.com/oauth/v2/accessToken", url.Values{
        "grant_type":    {"authorization_code"},
        "code":          {code},
        "client_id":     {cfg.LinkedInKey},
        "client_secret": {cfg.LinkedInSecret},
        "redirect_uri":  {cfg.Origin + "/api/oauth/linkedin/callback"},
    })

    // 3. Parse token response
    var tokenData struct {
        AccessToken string `json:"access_token"`
    }
    json.NewDecoder(tokenResp.Body).Decode(&tokenData)

    // 4. Fetch user profile
    req, _ := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
    req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
    profileResp, _ := http.DefaultClient.Do(req)

    var profile struct {
        ID        string `json:"id"`
        FirstName struct {
            Localized map[string]string `json:"localized"`
        } `json:"firstName"`
        LastName struct {
            Localized map[string]string `json:"localized"`
        } `json:"lastName"`
    }
    json.NewDecoder(profileResp.Body).Decode(&profile)

    // 5. Get email
    req, _ = http.NewRequest("GET", "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))", nil)
    req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
    emailResp, _ := http.DefaultClient.Do(req)

    var emailData struct {
        Elements []struct {
            Handle struct {
                EmailAddress string `json:"emailAddress"`
            } `json:"handle~"`
        } `json:"elements"`
    }
    json.NewDecoder(emailResp.Body).Decode(&emailData)

    // 6. Create or update commenter
    name := profile.FirstName.Localized["en_US"] + " " + profile.LastName.Localized["en_US"]
    email := emailData.Elements[0].Handle.EmailAddress

    commenterHex, err := commenterGetByProviderID("linkedin", profile.ID)
    if err != nil {
        // Create new commenter
        commenterHex = randomHex(32)
        db.Exec(`
            INSERT INTO commenters (commenterHex, email, name, provider, joinDate, state)
            VALUES ($1, $2, $3, 'linkedin', $4, 'ok')
        `, commenterHex, email, name, time.Now())
    }

    // 7. Create session
    commenterToken := randomHex(32)
    db.Exec(`
        INSERT INTO commenterSessions (commenterToken, commenterHex, creationDate)
        VALUES ($1, $2, $3)
    `, commenterToken, commenterHex, time.Now())

    // 8. Redirect back
    http.SetCookie(w, &http.Cookie{
        Name:  "commentoCommenterToken",
        Value: commenterToken,
    })
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
```

#### Step 4: Add Routes

**File**: `api/router_oauth.go`

```go
func routesServe() *mux.Router {
    // ... existing routes ...

    router.HandleFunc("/api/oauth/linkedin/redirect", oauthLinkedInRedirect).Methods("GET")
    router.HandleFunc("/api/oauth/linkedin/callback", oauthLinkedInCallback).Methods("GET")

    return router
}
```

#### Step 5: Update Frontend

**File**: `frontend/js/commento.js`

```javascript
function renderOAuthButtons(configuredOauths) {
  // ... existing code ...

  if (configuredOauths.indexOf("linkedin") > -1) {
    var linkedinButton = createElement("button", "commento-oauth-button commento-oauth-linkedin");
    linkedinButton.innerText = "Login with LinkedIn";
    linkedinButton.onclick = function() {
      window.location = origin + "/api/oauth/linkedin/redirect?domain=" + domain;
    };
    oauthContainer.appendChild(linkedinButton);
  }
}
```

### 3. Custom Email Templates

**File**: `templates/email-notification-new.txt`

```
Hello {{.CommenterName}},

{{.ParentCommenterName}} replied to your comment on {{.Domain}}{{.Path}}:

---
{{.HTML}}
---

View the full conversation:
{{.Origin}}{{.Path}}#commento-{{.CommentHex}}

---
To unsubscribe from these emails, click here:
{{.Origin}}/api/commenter/email/unsubscribe?token={{.UnsubscribeToken}}

Powered by Commento
```

Modify as needed, then rebuild.

### 4. Changing Moderation Logic

**File**: `api/comment_new.go:85`

```go
// Current logic
if !isModerator {
    if domain.RequireModeration {
        state = "unapproved"
    } else if domain.AutoSpamFilter && isSpamAkismet(...) {
        state = "flagged"
    }
}

// Example: Auto-approve users with > 10 approved comments
if !isModerator {
    approvedCount := getApprovedCommentCount(commenterHex)

    if approvedCount > 10 {
        state = "approved" // Trusted user
    } else if domain.RequireModeration {
        state = "unapproved"
    } else if domain.AutoSpamFilter && isSpamAkismet(...) {
        state = "flagged"
    }
}
```

### 5. Adding Custom API Endpoints

**File**: `api/comment_stats.go` (new file)

```go
package main

import (
    "net/http"
)

type commentStatsResponse struct {
    TotalComments   int `json:"totalComments"`
    ApprovedComments int `json:"approvedComments"`
    UniqueCommenters int `json:"uniqueCommenters"`
}

func commentStatsHandler(w http.ResponseWriter, r *http.Request) {
    domain := r.URL.Query().Get("domain")

    var stats commentStatsResponse

    db.QueryRow(`
        SELECT COUNT(*) FROM comments WHERE domain = $1
    `, domain).Scan(&stats.TotalComments)

    db.QueryRow(`
        SELECT COUNT(*) FROM comments WHERE domain = $1 AND state = 'approved'
    `, domain).Scan(&stats.ApprovedComments)

    db.QueryRow(`
        SELECT COUNT(DISTINCT commenterHex) FROM comments WHERE domain = $1
    `, domain).Scan(&stats.UniqueCommenters)

    writeJSON(w, stats)
}
```

**File**: `api/router_api.go`

```go
router.HandleFunc("/api/comment/stats", commentStatsHandler).Methods("GET")
```

### 6. Database Backup

```bash
# Backup
pg_dump commento > backup.sql

# Restore
psql commento < backup.sql
```

### 7. Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `COMMENTO_POSTGRES` | ✅ | - | PostgreSQL connection string |
| `COMMENTO_ORIGIN` | ✅ | - | Base URL (e.g., https://commento.io) |
| `COMMENTO_PORT` | ❌ | 8080 | Server port |
| `COMMENTO_BIND_ADDRESS` | ❌ | 127.0.0.1 | Bind address |
| `COMMENTO_CDN_PREFIX` | ❌ | - | CDN URL for static assets |
| `COMMENTO_SMTP_HOST` | ❌ | - | SMTP server |
| `COMMENTO_SMTP_PORT` | ❌ | 587 | SMTP port |
| `COMMENTO_SMTP_USERNAME` | ❌ | - | SMTP auth username |
| `COMMENTO_SMTP_PASSWORD` | ❌ | - | SMTP auth password |
| `COMMENTO_SMTP_FROM_ADDRESS` | ❌ | - | Email sender address |
| `COMMENTO_GOOGLE_KEY` | ❌ | - | Google OAuth client ID |
| `COMMENTO_GOOGLE_SECRET` | ❌ | - | Google OAuth secret |
| `COMMENTO_GITHUB_KEY` | ❌ | - | GitHub OAuth client ID |
| `COMMENTO_GITHUB_SECRET` | ❌ | - | GitHub OAuth secret |
| `COMMENTO_TWITTER_KEY` | ❌ | - | Twitter OAuth API key |
| `COMMENTO_TWITTER_SECRET` | ❌ | - | Twitter OAuth secret |
| `COMMENTO_GITLAB_KEY` | ❌ | - | GitLab OAuth app ID |
| `COMMENTO_GITLAB_SECRET` | ❌ | - | GitLab OAuth secret |
| `COMMENTO_AKISMET_KEY` | ❌ | - | Akismet API key for spam detection |
| `COMMENTO_FORBID_NEW_OWNERS` | ❌ | false | Disable new registrations |

---

## Deployment Architecture

### Docker Deployment

**File**: `docker-compose.yml`

```yaml
version: '3.7'

services:
  server:
    image: registry.gitlab.com/commento/commento:latest
    ports:
      - "8080:8080"
    environment:
      COMMENTO_ORIGIN: https://commento.example.com
      COMMENTO_POSTGRES: postgres://postgres:password@db:5432/commento?sslmode=disable
      COMMENTO_PORT: 8080
      COMMENTO_SMTP_HOST: smtp.sendgrid.net
      COMMENTO_SMTP_PORT: 587
      COMMENTO_SMTP_USERNAME: apikey
      COMMENTO_SMTP_PASSWORD: ${SENDGRID_API_KEY}
      COMMENTO_SMTP_FROM_ADDRESS: noreply@example.com
      COMMENTO_GOOGLE_KEY: ${GOOGLE_OAUTH_KEY}
      COMMENTO_GOOGLE_SECRET: ${GOOGLE_OAUTH_SECRET}
      COMMENTO_GITHUB_KEY: ${GITHUB_OAUTH_KEY}
      COMMENTO_GITHUB_SECRET: ${GITHUB_OAUTH_SECRET}
      COMMENTO_AKISMET_KEY: ${AKISMET_KEY}
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:13
    environment:
      POSTGRES_DB: commento
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres-data:
```

### Nginx Reverse Proxy

```nginx
server {
    listen 443 ssl http2;
    server_name commento.example.com;

    ssl_certificate /etc/letsencrypt/live/commento.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/commento.example.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Cache static assets
    location /js/ {
        proxy_pass http://localhost:8080;
        proxy_cache_valid 200 1d;
        add_header Cache-Control "public, immutable";
    }
}
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: commento
spec:
  replicas: 3
  selector:
    matchLabels:
      app: commento
  template:
    metadata:
      labels:
        app: commento
    spec:
      containers:
      - name: commento
        image: registry.gitlab.com/commento/commento:latest
        ports:
        - containerPort: 8080
        env:
        - name: COMMENTO_ORIGIN
          value: "https://commento.example.com"
        - name: COMMENTO_POSTGRES
          valueFrom:
            secretKeyRef:
              name: commento-secrets
              key: postgres-url
        - name: COMMENTO_SMTP_HOST
          value: "smtp.sendgrid.net"
        - name: COMMENTO_SMTP_PASSWORD
          valueFrom:
            secretKeyRef:
              name: commento-secrets
              key: smtp-password
        livenessProbe:
          httpGet:
            path: /api/health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: commento
spec:
  selector:
    app: commento
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Scaling Considerations

1. **Stateless Design**: Multiple Commento instances can run in parallel
2. **Database Connection Pooling**: Configure `SetMaxOpenConns` appropriately
3. **CDN for Static Assets**: Set `COMMENTO_CDN_PREFIX` to offload JS/CSS
4. **PostgreSQL Optimization**:
   - Use connection pooling (PgBouncer)
   - Index optimization (already included)
   - Read replicas for analytics queries

### Monitoring

```bash
# Prometheus metrics (add this handler)
# File: api/metrics.go
func metricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := fmt.Sprintf(`
# HELP commento_comments_total Total comments created
# TYPE commento_comments_total counter
commento_comments_total %d

# HELP commento_active_sessions Active user sessions
# TYPE commento_active_sessions gauge
commento_active_sessions %d
`, getCommentCount(), getActiveSessionCount())

    w.Write([]byte(metrics))
}
```

---

## Conclusion

This architecture document provides a comprehensive overview of Commento's design, implementation, and customization options. Key takeaways:

1. **Multi-tenancy** is achieved through domain-based isolation
2. **Security** is prioritized with bcrypt, parameterized queries, and XSS prevention
3. **Scalability** is built-in with stateless design and connection pooling
4. **Customization** is flexible through CSS overrides, OAuth providers, and API extensions
5. **Privacy** is core with no tracking and GDPR compliance

For further questions or contributions, refer to the source code in the `api/` and `frontend/` directories.
