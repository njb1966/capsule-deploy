# Technical Architecture

## Stack Summary

| Component | Choice | Rationale |
|-----------|--------|-----------|
| OS | Ubuntu 24.04 LTS | Long support cycle, excellent docs, predictable |
| Gemini server | Agate (Rust) | Actively maintained, simple config, solid TLS/SNI support |
| Web server | Caddy | Auto-TLS via Let's Encrypt, simple config, reverse proxy built in |
| Backend language | Go | Fast, single binary deploy, low memory, easy to maintain solo |
| Database | SQLite | No separate DB server, trivial backup, sufficient for this scale |
| File storage | Flat files on disk | Gemtext files are plain text — no object storage needed |
| TLS (Gemini) | Let's Encrypt wildcard via DNS-01 challenge | Covers *.yourdomain.com with one cert |
| TLS (Web UI) | Caddy auto-TLS | Handles yourdomain.com automatically |
| Process management | systemd | Standard, reliable, no extra tooling |
| Backups | Restic → Backblaze B2 | Cheap, reliable, encrypted, incremental |

---

## Infrastructure

Single VPS. No microservices. No container orchestration.

```
8 vCPU / 24 GB RAM / 400 GB SSD / 600 Mbit/s port / Unlimited inbound
~$12/month
```

This hardware is significantly over-specified for the load. Gemtext files average 2–5 KB. Even 50,000 users at the 50 MB cap each would require 2.5 TB — storage upgrades happen long before CPU or RAM becomes a concern. The generous spec provides headroom and avoids premature scaling decisions.

---

## DNS Configuration (One-Time Setup)

```
A     yourdomain.com          →  <server IP>
A     *.yourdomain.com        →  <server IP>    # wildcard — covers all usernames
A     www.yourdomain.com      →  <server IP>
```

The wildcard A record means every new user's subdomain resolves automatically with no DNS action required per user.

---

## TLS Certificate Setup (One-Time Setup)

Use `acme.sh` or `certbot` with a DNS-01 challenge to obtain a wildcard certificate:

```
*.yourdomain.com
yourdomain.com
```

This requires DNS API access (most registrars and DNS providers support this).
Renewal is automated via cron. The 90-day Let's Encrypt cert window provides ample time to catch renewal failures.

**Alert requirement:** Set up an email alert if cert renewal fails. A lapsed wildcard cert takes down every user's capsule simultaneously.

---

## Agate Configuration

Agate uses SNI (Server Name Indication) to serve different content per subdomain from a single running instance. Each user's capsule maps to a directory:

```toml
# /etc/agate/config.toml (conceptual — check Agate docs for exact format)
# Agate serves from a root directory; subdomain routing handled via --hostname
# or virtual hosting configuration per Agate's current docs
```

Each capsule directory lives at `/srv/capsules/<username>/`.

When Agate receives a request for `gemini://alice.yourdomain.com/index.gmi`, it serves `/srv/capsules/alice/index.gmi`.

---

## Filesystem Layout

```
/srv/capsules/                        # root for all capsule content
  alice/                              # one directory per username
    index.gmi
    about.gmi
    posts/
      2025-01-01-hello.gmi
  bob/
    index.gmi

/srv/capsule-service/                 # Go application
  capsule-service                     # compiled binary
  config.toml                         # app config

/var/lib/capsule-service/             # persistent app data
  users.db                            # SQLite database

/etc/agate/                           # Agate config

/var/log/capsule-service/             # application logs
                                      # NOTE: no user reading behavior logged here
```

---

## Go Backend — API Surface

The Go application handles everything the Gemini server does not: accounts, auth, file operations, and the editor API. It runs as a systemd service and exposes an HTTPS API consumed by the editor frontend via Caddy reverse proxy.

### Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/register` | Create account (username, email, password) |
| POST | `/api/verify-email` | Confirm email verification token |
| POST | `/api/login` | Authenticate, return session token |
| POST | `/api/logout` | Invalidate session |
| POST | `/api/password-reset-request` | Send password reset email |
| POST | `/api/password-reset` | Complete password reset with token |
| GET | `/api/capsule` | List all user's files and directories |
| GET | `/api/file/*path` | Return content of a single .gmi file |
| PUT | `/api/file/*path` | Save/update a file (writes to disk) |
| POST | `/api/file` | Create a new file |
| DELETE | `/api/file/*path` | Delete a file |
| POST | `/api/mkdir` | Create a subdirectory |
| GET | `/api/export` | Return ZIP of all user's .gmi files |
| DELETE | `/api/account` | Delete account and all files permanently |
| GET | `/api/health` | Service health check — returns 200 OK if running |

### API Error Response Format

All error responses use a consistent JSON envelope:

```json
{
  "error": "Human-readable message",
  "code": "MACHINE_READABLE_CODE"
}
```

HTTP status codes are used correctly (400 bad request, 401 unauthorized, 403 forbidden, 404 not found, 409 conflict, 413 payload too large, 429 too many requests, 500 internal server error). The `code` field allows the frontend to handle specific errors without string-matching the `error` message.

Example error codes:

| Code | Meaning |
|------|---------|
| `USERNAME_TAKEN` | Registration: username already exists |
| `EMAIL_TAKEN` | Registration: email already in use |
| `INVALID_USERNAME` | Username fails validation rules |
| `INVALID_CREDENTIALS` | Login: wrong email or password |
| `EMAIL_NOT_VERIFIED` | Login: account not yet verified |
| `RATE_LIMITED` | Too many requests from this IP |
| `FILE_TOO_LARGE` | File exceeds 1 MB per-file limit |
| `STORAGE_FULL` | User has reached 50 MB total cap |
| `PATH_INVALID` | Filename or path fails validation |
| `NOT_FOUND` | File or directory does not exist |
| `SESSION_EXPIRED` | JWT is invalid or expired |

### Config File Structure

```toml
# /srv/capsule-service/config.toml
[server]
listen = "127.0.0.1:8080"      # Caddy proxies to this
capsules_root = "/srv/capsules"
db_path = "/var/lib/capsule-service/users.db"
log_path = "/var/log/capsule-service/app.log"

[email]
smtp_host = "smtp.protonmail.ch"
smtp_port = 587                        # STARTTLS; use 465 for SSL/TLS
from_address = "admin@gemcities.com"
# smtp_username and smtp_password are NOT set here — see Secrets below

[limits]
max_file_size_bytes = 1048576   # 1 MB per file
max_total_storage_bytes = 52428800  # 50 MB per user
max_files_per_user = 500

[auth]
jwt_secret = "..."              # generated at setup, never committed
session_duration_days = 30
bcrypt_cost = 12
```

### Secrets Management

Credentials are never stored in `config.toml` (which may be committed to the repo). They live in a root-only systemd environment file:

```
/etc/capsule-service/env      # chmod 600, owned by root
```

Contents:
```
SMTP_USERNAME=admin@gemcities.com
SMTP_PASSWORD=<proton-app-password>   # generated in Proton Settings → Security → App passwords
JWT_SECRET=<random-64-char-string>    # generated once at setup: openssl rand -hex 32
```

The systemd unit file loads this via:
```ini
[Service]
EnvironmentFile=/etc/capsule-service/env
```

The Go backend reads these with `os.Getenv("SMTP_USERNAME")` etc. `config.toml` fields for the same keys are ignored if the environment variable is set. This file must never be committed, backed up to untrusted storage, or logged.

---

## Caddy Configuration (Conceptual)

```
yourdomain.com, www.yourdomain.com {
    root * /srv/capsule-editor/public
    file_server
    reverse_proxy /api/* 127.0.0.1:8080
    encode gzip
}
```

Caddy handles TLS for the web UI automatically. Agate handles TLS for Gemini separately using the wildcard cert.

---

## SQLite Schema (Core Tables)

```sql
CREATE TABLE users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    username    TEXT UNIQUE NOT NULL,
    email       TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,         -- bcrypt
    email_verified INTEGER DEFAULT 0,
    created_at  TEXT NOT NULL,
    storage_used_bytes INTEGER DEFAULT 0
);

CREATE TABLE sessions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    token_hash  TEXT UNIQUE NOT NULL,   -- hash of JWT, for invalidation
    created_at  TEXT NOT NULL,
    expires_at  TEXT NOT NULL
);

CREATE TABLE password_reset_tokens (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    token_hash  TEXT UNIQUE NOT NULL,
    created_at  TEXT NOT NULL,
    expires_at  TEXT NOT NULL,
    used        INTEGER DEFAULT 0
);

CREATE TABLE abuse_log (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    reported_username TEXT,
    report_type TEXT,
    report_details TEXT,
    received_at TEXT NOT NULL,
    resolved_at TEXT,
    action_taken TEXT
);
```

---

## Systemd Services

Three services managed by systemd:

```
capsule-service.service     # Go backend API
agate.service               # Gemini server
caddy.service               # Web server (likely installed as package)
```

All set to `Restart=on-failure` and `WantedBy=multi-user.target`.

---

## Backup Configuration

Restic to Backblaze B2, running daily via cron at 3:00 AM:

```bash
# Backs up capsule files and SQLite database
restic -r b2:bucket-name:/capsules backup /srv/capsules /var/lib/capsule-service/
```

Retention policy: 7 daily, 4 weekly, 6 monthly snapshots.
Backup restoration procedure must be tested quarterly.
Alert on backup failure via email.
