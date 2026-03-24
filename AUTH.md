# Authentication & Registration

## Registration Flow

1. User visits `yourdomain.com` and clicks "Get a Capsule"
2. Registration form collects: username, email, password
3. Username validated client-side and server-side (see rules below)
4. Password validated: minimum 10 characters, no other restrictions
5. Account created in SQLite
6. Verification email sent with a single-use token link
7. Capsule directory created on disk immediately (capsule is live but shows a placeholder until verified — or goes live immediately, operator choice)
8. User redirected to editor with banner: "Check your email to verify your address"

Email verification is required before the capsule is publicly accessible. This is the primary anti-abuse filter at registration.

---

## Username Rules

Usernames become DNS subdomains and filesystem directory names. Strict validation required.

**Allowed characters:** `a-z`, `0-9`, `-` (hyphen)
- Lowercase only (enforce on input, not just validation)
- No underscores (not valid in DNS labels per RFC 1123)
- No spaces, dots, or special characters

**Length:** 3–32 characters

**Format rules:**
- Cannot start with a hyphen
- Cannot end with a hyphen
- Cannot contain consecutive hyphens (`--`) — reserved by DNS for IDN labels

**Reserved names (blocked at registration):**
```
www, mail, ftp, smtp, pop, imap, api, admin, administrator,
gemini, static, assets, media, cdn, help, about, contact,
support, abuse, postmaster, hostmaster, webmaster, root,
info, status, blog, news, feed, rss, atom, capsule, capsules,
user, users, account, accounts, login, logout, register,
signup, signin, dashboard, editor, settings, profile,
well-known, geminispace
```

Plus a maintained list of common offensive terms (maintain separately, do not hardcode in the application — load from a config file so it can be updated without redeployment).

**Uniqueness:** checked against database at registration time.
**Username hold:** deleted usernames are unavailable for 30 days after deletion.

---

## Password Handling

- Hashed with **bcrypt**, cost factor **12**
- Never stored in plaintext anywhere
- Never logged
- Minimum length: 10 characters
- No maximum length restriction (bcrypt handles long inputs gracefully via pre-hashing if needed)
- No forced complexity rules (length is the strongest single factor)

---

## Session Management

- Sessions use **signed JWT tokens**
- JWT stored in **httpOnly, Secure, SameSite=Strict cookie** — not localStorage
- Session duration: **30 days**
- Token hash stored in `sessions` table for server-side invalidation
- Logout invalidates the session server-side (deletes from sessions table)
- No "remember me" toggle — all sessions are 30 days by default

### JWT Claims

```json
{
  "sub": "<user_id>",
  "username": "<username>",
  "iat": <issued_at_unix>,
  "exp": <expires_at_unix>
}
```

---

## Password Reset Flow

1. User clicks "Forgot password" on login page
2. Enters email address
3. If email exists: single-use reset token generated, hashed and stored in `password_reset_tokens` table, email sent with link
4. If email does not exist: same response shown (do not reveal whether email is registered)
5. Reset link valid for **1 hour**
6. Single-use: token marked `used = 1` immediately on first visit to reset page
7. User sets new password on reset page
8. All existing sessions invalidated on password change
9. User redirected to login

---

## Rate Limiting

| Endpoint | Limit |
|----------|-------|
| POST `/api/login` | 10 attempts per 15 minutes per IP |
| POST `/api/register` | 5 registrations per hour per IP |
| POST `/api/password-reset-request` | 3 requests per hour per IP |
| All other API endpoints | 120 requests per minute per authenticated user |

Rate limit responses return HTTP 429 with a `Retry-After` header.

---

## Email Verification Token

- Token: 32 bytes of cryptographically random data, hex-encoded (64 char string)
- Stored as SHA-256 hash in database (token itself never stored)
- Valid for **72 hours** from generation
- Single-use
- If user requests resend: old token invalidated, new token issued (same 72-hour window)

---

## Email Addresses

- Stored in database
- Used only for: email verification, password reset, critical service notices
- Never shared, never used for marketing, never used for any form of engagement communication
- Users can change email address in account settings (re-verification required)

---

## What Is Never Done With Accounts

- No login history stored beyond what is needed for rate limiting
- No "last seen" timestamp
- No device fingerprinting
- No "suspicious login" detection that would require storing location data
- No third-party auth (no "Login with Google/GitHub") — keeps the dependency surface minimal
