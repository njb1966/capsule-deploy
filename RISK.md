# Risk Register

## How to Use This Document

Review this document quarterly. Update mitigations as circumstances change. Add new risks as they are identified. Mark resolved risks rather than deleting them.

---

## Risk Ratings

**Impact:** High / Medium / Low
**Likelihood:** High / Medium / Low / Unlikely

---

## Active Risks

### R01 — Operator Burnout or Unavailability
**Impact:** High | **Likelihood:** Medium (all solo projects face this)

The service depends on one person. If that person burns out, gets sick, or simply loses interest, the service could go dark with no notice.

**Mitigations:**
- Succession document maintained and shared with one trusted person (see `docs/SUSTAINABILITY.md`)
- All software open source — community can fork and continue
- Shutdown procedure documented: 90-day notice, export reminders, read-only grace period
- Export always available — users are never trapped
- Deliberately low operational burden (simple stack, minimal maintenance requirements)

---

### R02 — Illegal or Seriously Harmful Content Posted
**Impact:** High | **Likelihood:** Medium (will happen eventually at any scale)

A user posts CSAM, credible threats, or other seriously illegal content.

**Mitigations:**
- Email verification at registration (raises barrier)
- Clear Terms of Service agreed to at registration
- Abuse reporting email monitored regularly
- DMCA agent registered before launch
- NCMEC CyberTipline bookmarked and ready (`cybertipline.org`)
- Incident log maintained
- See `docs/MODERATION.md` for full response procedures

---

### R03 — Wildcard TLS Certificate Renewal Failure
**Impact:** High | **Likelihood:** Low (but catastrophic if it occurs)

Let's Encrypt wildcard cert expires. Every user's capsule becomes inaccessible simultaneously with TLS errors.

**Mitigations:**
- Automated renewal via cron with DNS-01 challenge
- Email alert configured for cert expiry within 30 days
- 90-day cert window gives ample response time after any failure
- Manual renewal procedure documented and tested
- Monitor renewal logs weekly

---

### R04 — VPS Provider Outage or Data Loss
**Impact:** High | **Likelihood:** Unlikely (but possible)

VPS provider has an outage, or in a worst case, loses data.

**Mitigations:**
- Daily backups to Backblaze B2 (separate provider, separate infrastructure)
- Backup restoration procedure documented and tested quarterly
- Restic backup includes all capsule files and SQLite database
- Recovery time objective: acceptable to be down for hours (this is not a commercial SLA service)
- Users can re-export their content after restoration

---

### R05 — Storage Exhaustion
**Impact:** Medium | **Likelihood:** Low

Disk fills up, new writes fail, existing capsules may become unreadable.

**Mitigations:**
- 50 MB per-user cap enforced in application
- 400 GB SSD — at 50 MB/user cap, theoretical max ~8,000 users at full capacity
- Real-world average usage far lower (typical capsule is a handful of small files)
- Monitoring alert at 70% disk usage
- VPS storage upgrades available from provider if needed
- Periodic review of inactive accounts (no auto-deletion without notice)

---

### R06 — DDoS Attack Against the Service or a Specific Capsule
**Impact:** Medium | **Likelihood:** Possible

Either the entire service or a specific user's capsule is targeted with high traffic volume.

**Mitigations:**
- Rate limiting at Caddy layer for web UI
- 600 Mbit/s port — significant capacity for a small web service
- Individual capsule can be temporarily suspended if it attracts attack traffic
- VPS provider abuse reporting available for volumetric attacks
- For sustained attacks: upstream provider filtering or temporary block of attacking IP ranges

---

### R07 — Community Backlash (Perceived Centralization)
**Impact:** Medium | **Likelihood:** Possible

The Gemini community views centralized hosting services with suspicion. Someone argues that this service is bad for the small web.

**Mitigations:**
- Open source from day one — transparency about how it works
- Easy export always available — no lock-in
- Clear philosophy documentation (`docs/PHILOSOPHY.md`)
- No features that create lock-in or dependency
- Explicitly not trying to be the only way to publish — self-hosting documentation welcomed
- Service's own capsule explains the rationale honestly

---

### R08 — Legal Demand (DMCA, Subpoena, etc.)
**Impact:** Medium | **Likelihood:** Possible over time

A valid DMCA takedown notice arrives. Or a law enforcement subpoena for user data.

**Mitigations:**
- DMCA agent registered — safe harbor applies
- DMCA response procedure documented in `docs/MODERATION.md`
- Minimal user data collected — limited useful data for law enforcement requests
- No IP logging of capsule visitors — cannot be compelled to produce what does not exist
- Consult a lawyer before launch for jurisdiction-specific guidance

---

### R09 — Username Squatting at Launch
**Impact:** Low | **Likelihood:** High

When registration opens, people register desirable usernames speculatively or maliciously.

**Mitigations:**
- 30-day hold on deleted usernames prevents quick grab-release cycles
- No username trading or transfer mechanism
- Reserved names list pre-populated (`docs/AUTH.md`)
- No way to monetize usernames on this service — reduces squatting incentive

---

### R10 — Email Deliverability Problems
**Impact:** Medium | **Likelihood:** Possible

Password reset and verification emails land in spam or are blocked.

**Mitigations:**
- Use a reputable transactional email service (even free tier of Mailgun, Postmark, etc.)
- Configure SPF, DKIM, and DMARC records for the domain
- From address uses a real domain, not a no-reply on a shared IP
- Provide manual verification fallback (operator can manually verify email if user contacts abuse address)

---

### R11 — Go Application Crashes / Hangs
**Impact:** Medium | **Likelihood:** Low

The Go backend crashes or becomes unresponsive.

**Mitigations:**
- systemd `Restart=on-failure` — automatic restart
- Health check endpoint (`/api/health`) for monitoring
- Structured error logging — crash context captured
- Gemini server (Agate) continues serving existing capsule content even if Go backend is down — only editor and new registrations are affected

---

## Resolved Risks

*(None yet — populate as risks are addressed)*

---

## Risk Review Log

| Date | Reviewer | Changes Made |
|------|----------|-------------|
| (initial) | — | Document created |
