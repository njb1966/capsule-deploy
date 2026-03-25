# Security

## File Isolation & Path Traversal Protection

This is the most critical security concern in the entire application. Every file operation must be validated before execution.

### The Rule

Every file path provided by the user must be:
1. Cleaned (remove `.`, `..`, double slashes)
2. Joined with the user's capsule root directory
3. Resolved to an absolute path
4. Verified to still be within the user's capsule root

**In Go:**
```go
func safeFilePath(capsuleRoot, userProvidedPath string) (string, error) {
    // Clean and join
    joined := filepath.Join(capsuleRoot, filepath.Clean("/"+userProvidedPath))
    // Resolve to absolute
    abs, err := filepath.Abs(joined)
    if err != nil {
        return "", err
    }
    // Verify it is still inside the capsule root
    if !strings.HasPrefix(abs, capsuleRoot+string(filepath.Separator)) {
        return "", errors.New("path traversal detected")
    }
    return abs, nil
}
```

This must be called for **every** file read, write, delete, and rename operation. No exceptions.

---

## File Validation

### Filename Rules (enforced server-side, not just client-side)
- Allowed characters: `a-z`, `A-Z`, `0-9`, `-`, `_`, `.`
- Must end in `.gmi`
- No hidden files (names starting with `.`)
- Maximum filename length: 64 characters
- Maximum path depth: 5 levels of subdirectory

### Directory Name Rules
- Allowed characters: `a-z`, `A-Z`, `0-9`, `-`, `_`
- No dots in directory names
- Maximum name length: 64 characters

### File Size Limits
- Maximum per file: **1 MB** (gemtext files have no legitimate reason to exceed this)
- Maximum total per user: **50 MB**
- Storage usage tracked in `users.storage_used_bytes` column, updated on every write and delete

---

## Storage Enforcement

Before writing any file:
1. Check new file size against per-file limit
2. Calculate `current_storage + new_file_size - old_file_size` (for updates)
3. Check against per-user 50 MB limit
4. Return HTTP 413 with clear message if limit exceeded

Storage drift (DB value vs. actual disk usage) is reconciled by a nightly cron job, not on every login. Recalculating from disk on login would be slow for users with many files and adds latency to every session start.

---

## Authentication Security

See `docs/AUTH.md` for full details. Summary:
- bcrypt cost 12 for passwords
- httpOnly JWT cookies (not localStorage)
- Server-side session invalidation on logout and password change
- Rate limiting on all auth endpoints

---

## What Is Never Logged

Application logs record: errors, service start/stop, and failed auth attempts (for rate limiting only — IP and timestamp, no username on failed logins).

**Explicitly never logged:**
- IP addresses of Gemini capsule visitors
- Which capsules a visitor reads
- How long a visitor spends reading
- Referrer/source information for any request
- Contents of files being edited
- Search queries (if any search is ever added)

The Gemini server (Agate) access logs, if enabled, should be **disabled or piped to /dev/null**. Agate access logs would contain visitor IPs. This is incompatible with the privacy commitment.

```toml
# In Agate config — disable access logging
# Check current Agate docs for exact config key
```

---

## Content Security (Web Editor)

The editor frontend is served as static files. No user-generated content is ever rendered as HTML by the web server.

- Gemtext preview is rendered in a sandboxed `<iframe>` or via a DOM-based renderer that never uses `innerHTML` with unsanitized user content
- File contents are always treated as plain text in the editor textarea
- If using iframe for preview: `sandbox="allow-same-origin"` attribute, no scripts in preview
- CSP header on all editor pages:
  ```
  Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self'; img-src 'none'; connect-src 'self'
  ```

---

## HTTPS Configuration (Caddy)

Caddy's defaults are strong. Verify:
- TLS 1.2 minimum (Caddy default)
- HSTS header enabled
- No mixed content

---

## Dependency Surface

Intentionally minimal:
- **Go standard library** handles most things
- **bcrypt:** `golang.org/x/crypto/bcrypt`
- **JWT:** one well-maintained library (e.g. `golang-jwt/jwt`)
- **SQLite driver:** `modernc.org/sqlite` (pure Go, no CGo dependency)
- **ZIP generation:** Go standard library `archive/zip`

No web framework. Use `net/http` directly. Fewer dependencies = smaller attack surface and easier security auditing.

---

## Abuse / Rate Limiting Infrastructure

Rate limiting implemented in the Go application layer (not via external tools like fail2ban, to keep the stack simple).

In-memory rate limit counters using a sliding window algorithm. On server restart, counters reset — acceptable for this scale.

Failed login attempts: if 10 failed attempts from same IP in 15 minutes, return 429 for subsequent attempts from that IP for the remainder of the window.

---

## Incident Response

If a security issue is discovered:
1. Take the affected component offline if necessary
2. Assess scope
3. If user data may be compromised: notify affected users by email within 72 hours
4. Document what happened, what was done, and what changed in the private incident log
5. If the issue affects other Gemini hosts using Agate: notify upstream responsibly

No public bug bounty program — but a security contact email address (`security@gemcities.com`) should be published.
