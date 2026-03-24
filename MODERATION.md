# Content Moderation & Legal Framework

## Why Gemini Is a Cleaner Baseline

Gemini's protocol constraints eliminate entire abuse categories before moderation begins:
- **No images** — no visual CSAM, gore, shock content, or image-based harassment
- **No JavaScript** — no phishing kits, crypto miners, drive-by exploits, malware delivery
- **No comments** — no harassment threads, coordinated pile-ons, or spam vectors
- **No embeds** — no content injection attacks

The realistic threat surface is almost entirely written text. Serious, but manageable.

---

## Terms of Service

Publish a short, plain-language ToS on the landing page. No legalese. Users must agree at registration.

**Core prohibitions:**
- No illegal content of any kind
- No content designed to threaten or harass specific individuals
- No hate speech targeting people based on protected characteristics
- No doxxing — posting private personal information about others without consent
- No copyright infringement

Short and unambiguous. No one can claim they did not see it.

---

## Realistic Threat Assessment

| Content Type | Likelihood | Severity |
|-------------|------------|----------|
| Written hate speech / manifestos | Possible | High |
| Doxxing (posting personal info to harass) | Possible | High |
| Threats against specific individuals | Possible | High |
| CSAM in text form | Low | Critical — federal law applies |
| Copyright infringement (books, etc.) | Moderate | Medium — DMCA process |
| Drug marketplace text | Low | Medium |
| Spam/SEO capsules | Low | Low |

---

## Moderation Approach (Five Layers)

### Layer 1 — Registration Friction
Require email verification before capsule goes live. Single most effective filter. Bad actors using throwaway emails still face friction, and every account has a real email address for follow-up.

### Layer 2 — Terms of Service
Plain-language ToS agreed to at registration. Establishes the rules clearly and provides the legal basis for account suspension.

### Layer 3 — Abuse Reporting
Published abuse email address: `abuse@yourdomain.com`. Linked from every page of the web UI and from the service's own Gemini capsule. Monitored by the operator. Target: review and act within 48 hours of receiving a report.

### Layer 4 — Periodic Spot-Checks
Occasional manual review of newly created capsules — not reading every word, just a brief check of new accounts. Think of this as a librarian walking the stacks, not surveillance. Catches obvious problems without automated monitoring.

### Layer 5 — Incident Log
A private log of every abuse report received and every action taken. Date, nature of report, action taken, outcome. A plain text file or spreadsheet is sufficient. This is evidence of good-faith moderation if legal scrutiny ever arises.

---

## Enforcement Actions

| Situation | Action |
|-----------|--------|
| Single policy-violating file | Remove specific file, email user with explanation |
| Pattern of violations | Suspend account, remove all content, email user |
| CSAM — any amount | Remove immediately, report to NCMEC, preserve evidence, permanent ban |
| Valid DMCA takedown | Remove content within 48 hours, notify user |
| Credible threat of violence | Remove immediately, consider reporting to authorities |
| Doxxing post | Remove immediately, warn user, suspend on repeat |
| Repeat violator | Permanent account suspension |

---

## US Legal Framework

> **Note:** This section reflects general US law. Consult a lawyer before launch for jurisdiction-specific advice. A single 1-hour consultation is strongly recommended.

### Section 230 — Communications Decency Act
As a US-based host, Section 230 provides significant protection: you are generally not liable for content your users post, as long as you act in good faith to remove illegal content when notified. You are an interactive computer service provider, not the author of user capsules. The layered moderation approach above demonstrates good faith.

### DMCA Safe Harbor
Protects from copyright liability for user-posted content, provided:

1. **Register a DMCA designated agent** with the US Copyright Office
   - URL: `copyright.gov/dmca-agent`
   - Cost: $6
   - Time: ~10 minutes
   - **Do this before launch — without it, safe harbor does not apply**

2. **Publish a DMCA takedown procedure** on the site (a page at `/dmca`)

3. **Respond promptly** to valid takedown notices by removing the content

4. **Implement a repeat infringer policy** — account suspension for users with multiple valid takedowns

### CSAM — Mandatory Reporting Law (18 U.S.C. § 2258A)

**There is no safe harbor for CSAM. No exceptions.**

Child Sexual Abuse Material in any form — including written text — triggers mandatory federal reporting obligations for any electronic service provider.

**If you discover CSAM:**
1. **Remove the content immediately**
2. **Report to NCMEC** via CyberTipline at `cybertipline.org` — this is legally required
3. **Preserve the evidence** — do not delete the account or wipe logs before reporting
4. **Permanently ban the account**

Failure to report is itself a federal crime. Act immediately — do not delay for any reason.

Bookmark `cybertipline.org` now and keep it accessible.

---

## Pre-Launch Legal Checklist

- [ ] Register DMCA designated agent at `copyright.gov/dmca-agent` ($6)
- [ ] Publish Terms of Service at `yourdomain.com/terms`
- [ ] Publish Privacy Policy at `yourdomain.com/privacy`
- [ ] Publish DMCA procedure at `yourdomain.com/dmca`
- [ ] Publish abuse contact: `abuse@yourdomain.com`
- [ ] Publish security contact: `security@yourdomain.com`
- [ ] Bookmark NCMEC CyberTipline: `cybertipline.org`
- [ ] Consult a lawyer (1-hour review of ToS and liability exposure)
- [ ] Create private incident log file

---

## Privacy Policy Requirements

The privacy policy should clearly state:
- What data is collected (email, username, password hash, files)
- What data is NOT collected (reader IPs, reading behavior, analytics)
- How data is used (account management and recovery only)
- Data retention (data held until account deletion, then deleted)
- No third-party sharing of any kind
- How to delete your account and data

---

## What Is Not Done

- No automated content scanning or AI moderation
- No keyword filtering
- No hash-matching against known CSAM databases (PhotoDNA etc.) — this requires a formal program, legal agreements, and mental health support for operators; out of scope for a solo operation
- No proactive surveillance of user content

The sole exception: if a credible external tip arrives about specific illegal content, manual review of that specific content is appropriate and necessary.
