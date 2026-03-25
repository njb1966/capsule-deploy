# Launch Plan

## Phase 1 — Foundation (Weeks 1–4)

Server and infrastructure setup. Nothing user-facing yet.

- [x] Register domain name — gemcities.com
- [x] Provision VPS, install OS — Debian 13 (trixie) on Contabo
- [x] Initial server hardening: key-only SSH, UFW configured (22/80/443/1965)
- [x] Configure wildcard DNS records (`*.gemcities.com` → 62.146.182.23)
- [x] Install Caddy — v2.11.2, serving gemcities.com
- [x] Wildcard TLS for Gemini — self-signed P-256 cert, TOFU-compatible
- [x] Configure certificate auto-renewal — Caddy handles HTTPS automatically
- [x] Install and configure Agate Gemini server — dynamic virtual hosting via wrapper script
- [x] Verify Agate serves correctly for a test subdomain — gemini://test.gemcities.com confirmed
- [x] Backups — Restic 0.17.3, SFTP via Tailscale, 3am daily cron, restore tested 2026-03-24
- [x] Set up systemd services for Agate and Caddy
- [x] Set up monitoring: disk usage alert at 70%, cert expiry alert at 30 days, backup failure alert
- [x] Register DMCA designated agent at copyright.gov — Registration: DMCA-1070853, Pay.gov tracking ID: 280P42RK, 2026-03-24

**Phase 1 complete when:** A test capsule is live and accessible at `gemini://test.yourdomain.com`, backups are running and tested.

**COMPLETED: 2026-03-24**

---

## Phase 2 — Application Build (Weeks 5–8)

Build and deploy the Go backend and editor frontend.

### Backend (Go)
- [ ] Project scaffold: Go module, directory structure, config loading
- [ ] SQLite setup: schema creation, migrations
- [ ] Username validation (see `docs/AUTH.md` for full rules)
- [ ] Registration endpoint with email verification
- [ ] Login / logout endpoints with JWT session handling
- [ ] Password reset flow
- [ ] File API: list, read, write, delete, rename, mkdir
- [ ] Path traversal protection (see `docs/SECURITY.md` — critical)
- [ ] Storage limit enforcement
- [ ] Export endpoint (ZIP generation)
- [ ] Account deletion endpoint
- [ ] Rate limiting on auth endpoints
- [ ] Systemd service for Go application
- [ ] Caddy reverse proxy configuration for `/api/*`

### Editor Frontend
- [ ] Landing page (plain HTML/CSS — explain the service, registration link)
- [ ] Registration page
- [ ] Login page
- [ ] Password reset pages
- [ ] Editor page: file tree, textarea, preview pane, top bar
- [ ] Gemtext live preview renderer (vanilla JS)
- [ ] Keyboard shortcuts
- [ ] File operations: new file, new folder, save, rename, delete
- [ ] Export button
- [ ] Account settings page
- [ ] Account deletion flow
- [ ] Error states and toast notifications
- [ ] Mobile responsive layout
- [ ] Terms of Service page
- [ ] Privacy Policy page
- [ ] DMCA page
- [ ] Abuse contact page

**Phase 2 complete when:** Full registration → edit → publish → export flow works end to end on the live server.

**COMPLETED: 2026-03-24** — Full flow tested live. Email via admin@gemcities.com (Protonmail SMTP). Legal pages still needed before Phase 3.

---

## Phase 3 — Soft Launch (Weeks 9–10)

Private beta with a small group of real users.

- [ ] Invite 10–20 people from the Gemini community (Station, mailing lists, known capsule owners)
- [ ] Provide a feedback channel (email, or a simple feedback form — not a public forum)
- [ ] Observe editor usability closely — this is the most important thing to get right
- [ ] Fix reported issues: bugs, confusing UX, missing edge cases
- [ ] No new features during this phase — only fixes
- [ ] Create the service's own Gemini capsule at `gemini://yourdomain.com` explaining the project (use your own service)
- [ ] Write the Terms of Service, Privacy Policy, and DMCA pages (have a lawyer review ToS)
- [ ] Test backup restoration with real user data

**Phase 3 complete when:** No critical bugs outstanding, feedback incorporated, legal pages reviewed.

---

## Phase 4 — Public Launch (Week 11+)

Open to everyone.

- [ ] Remove any invite gate from registration
- [ ] Announce on Gemini community spaces: geminispace.info, Station (`gemini://station.martinrue.com`)
- [ ] Post on small-web adjacent communities: tildes (tilde.club, etc.), lobste.rs, sourcehut community spaces
- [ ] Post on the service's own Gemini capsule
- [ ] Monitor server error logs closely for first two weeks
- [ ] Monitor disk usage
- [ ] Be available to respond to abuse reports promptly in the first weeks

---

## Post-Launch Maintenance Schedule

| Frequency | Task |
|-----------|------|
| Daily | Check backup completion |
| Weekly | Review error logs briefly |
| Weekly | Check abuse report inbox |
| Monthly | Apply security updates (OS packages, Go, Agate, Caddy) |
| Monthly | Review disk usage |
| Quarterly | Test backup restoration (actually restore to a test location) |
| Quarterly | Review and update reserved username list |
| Annually | Publish financial transparency post |
| Annually | Renew DMCA agent registration (copyright.gov) |

---

## Pre-Launch Checklist (Final Gate)

Do not open public registration until all of these are complete:

- [ ] Wildcard TLS cert working and auto-renewal tested
- [ ] Path traversal protection tested with adversarial inputs
- [ ] Rate limiting tested on login and registration endpoints
- [ ] Storage limits enforced and tested
- [ ] Export function tested — ZIP downloads correctly with correct directory structure
- [ ] Account deletion tested — files actually removed from disk
- [ ] Password reset flow tested end to end
- [ ] Email verification tested
- [ ] Backup tested — restore performed successfully at least once
- [ ] DMCA agent registered
- [ ] ToS, Privacy Policy, DMCA pages live
- [ ] Abuse email address active and monitored
- [ ] Gemtext preview renders all line types correctly
- [ ] Mobile layout tested on real device
- [ ] Agate access logging disabled (no visitor IP logging)
- [ ] CSP headers verified on editor pages
