# Capsule Hosting Service — Project Master

## What This Project Is

A free, donation-supported Gemini capsule hosting service. The goal is to reduce the barrier to publishing on the Gemini small web to near zero — no server knowledge required, no command line, just sign up and write.

This is **not** a social network, blogging platform, or growth-oriented product. It is infrastructure. Every decision should be evaluated against that.

**One-line pitch:** Geocities for Gemini, run ethically, forever simple.

---

## Document Index

Read these files when working on the relevant area. Do not assume — if something seems unclear, the answer is probably in one of these documents.

| File | Contents |
|------|----------|
| `docs/PHILOSOPHY.md` | Core commitments, non-negotiables, what this service will never become |
| `docs/ARCHITECTURE.md` | Full technical stack, infrastructure, DNS/TLS setup, directory layout, API surface |
| `docs/EDITOR.md` | Web editor UI spec, gemtext preview, keyboard shortcuts, file management UX |
| `docs/AUTH.md` | Registration flow, authentication, session management, password reset |
| `docs/SECURITY.md` | File isolation, path traversal protection, rate limiting, what is never logged |
| `docs/MODERATION.md` | Content policy, abuse handling, legal framework (Section 230, DMCA, CSAM) |
| `docs/SUSTAINABILITY.md` | Cost structure, donation model, operator continuity plan |
| `docs/LAUNCH.md` | Four-phase launch plan, pre-launch checklist, post-launch maintenance |
| `docs/NON-FEATURES.md` | Permanent list of things that will never be built and why |
| `docs/RISK.md` | Risk register with mitigations for all identified risks |
| `docs/GEMTEXT.md` | Gemtext format reference, syntax rules, rendering expectations |

---

## Key Decisions (Locked)

These are not open questions. Do not re-litigate them.

- **Protocol:** Gemini only
- **URL structure:** `gemini://username.yourdomain.com` (subdomain per user, wildcard TLS)
- **Custom domains:** No
- **Auth:** Email + password, bcrypt, JWT sessions
- **Funding:** Donation-supported, free to use, no tiers
- **Analytics:** None. Ever.
- **Ads:** None. Ever.
- **Comments:** None. Ever.
- **Database:** SQLite (no separate DB server)
- **Language:** Go for backend
- **Gemini server:** Agate
- **Web server:** Caddy
- **OS:** Ubuntu 24.04 LTS

---

## Infrastructure at a Glance

Single VPS — no microservices, no Kubernetes, no complexity.

```
8 vCPU / 24 GB RAM / 400 GB SSD / 600 Mbit/s / Unlimited inbound
Cost: ~$12/month
Total monthly cost including domain + backups: ~$14/month
```

Storage per user cap: 50 MB (gemtext files are kilobytes — this is generous).
Expected capacity: tens of thousands of capsules on this hardware with headroom to spare.

---

## Repository Structure (Three Repos)

```
capsule-service/    — Go backend (AGPL-3.0)
capsule-editor/     — Web editor frontend, vanilla HTML/CSS/JS (AGPL-3.0)
capsule-deploy/     — Config templates, systemd units, setup scripts, docs
```

---

## Core Rules for All Code

1. **No tracking code of any kind.** No analytics, no telemetry, no logging of user reading behavior.
2. **No external CDN dependencies** in the editor frontend. Everything must be self-hosted.
3. **Path traversal protection is non-negotiable.** All file paths must be resolved and validated within the user's capsule directory before any file operation.
4. **Usernames are DNS subdomains.** Validate strictly: lowercase a-z, digits, hyphens only; 3–32 chars; no leading/trailing hyphens; reserved names blocked.
5. **SQLite only.** Do not introduce Postgres, MySQL, or any external database server.
6. **Single binary deployment** for the Go backend. No runtime dependencies beyond the binary and config file.
7. **Boring technology wins.** If there are two ways to do something, prefer the one that is easier to understand and maintain over the one that is clever.

---

## Quick Reference: What Capsule URLs Look Like

```
gemini://alice.yourdomain.com          ← Alice's capsule root
gemini://alice.yourdomain.com/about.gmi
gemini://alice.yourdomain.com/posts/2025-01-01-hello.gmi
```

The web editor is served at:
```
https://yourdomain.com                 ← landing page + editor
https://yourdomain.com/register
https://yourdomain.com/login
https://yourdomain.com/editor
```
