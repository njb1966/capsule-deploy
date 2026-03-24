# Sustainability & Business Model

## Cost Structure

| Item | Monthly Cost |
|------|-------------|
| VPS (8 vCPU, 24 GB RAM, 400 GB SSD) | $12.00 |
| Domain name (annualized) | ~$1.50 |
| Backblaze B2 backup storage (~50 GB) | ~$0.50 |
| Transactional email — password resets (low volume) | ~$0.00 (free tier) |
| **Total** | **~$14/month (~$168/year)** |

This is an exceptionally low cost base. A few dozen modest donations per year cover it entirely. The service does not require growth to remain financially viable.

---

## Donation Model

The service is free to use. Donations are accepted but never required or incentivized.

**How donations work:**
- Accepted via a simple, privacy-respecting processor (Stripe direct, or Ko-fi/Buy Me a Coffee — operator's choice)
- No donation tiers
- No rewards for donating
- No premium features for donors — donating buys nothing except keeping the lights on
- One understated donation link in the footer of the web UI. No popups, no banners, no nag screens.

**Transparency:**
- An annual financial post published on the service's own Gemini capsule
- Contents: total operating costs for the year, total donations received, nothing more
- No individual donor amounts or names published

---

## Why Not a Paid Model

A flat annual fee ($20/year) was considered and rejected for this implementation. Reasons:
- Donation model is more aligned with small web ethos
- Payment processing adds complexity, compliance requirements (PCI), and support burden
- The cost base is low enough that donations from even a small fraction of users covers it
- Keeps the barrier to entry at zero

If donations prove consistently insufficient after 12 months of operation, revisit the fee model. Do not switch to ads.

---

## Operator Continuity Plan

The biggest sustainability risk for a solo-operated service is operator unavailability — burnout, life changes, or unexpected incapacitation. This must be planned for explicitly before launch.

### Succession Document
A private document maintained securely (e.g. in a password manager shared entry, or a sealed physical document) containing:
- VPS provider login credentials
- Domain registrar login credentials
- Backblaze B2 credentials
- SSH key location and passphrase
- DMCA agent registration details
- Any other credentials needed to operate or transfer the service

This document is shared with **one trusted person** who understands what to do if the operator is unavailable.

### If the Service Must Shut Down
1. Post a shutdown notice prominently on the web UI and on the service Gemini capsule
2. Send email to all users: **90 days notice minimum**
3. Send export reminder emails at 90 days, 30 days, and 7 days before shutdown
4. On shutdown date: export function remains available for an additional 30 days (read-only mode — no new registrations, no edits, just export)
5. After final shutdown: server decommissioned

### Open Source as Continuity
All server software is published as open source (AGPL-3.0). Any community member with the technical skills can:
- Stand up a new instance of the same service
- Import user exports from this service
- Continue the project independently

This is the strongest long-term continuity guarantee: the code outlives any individual operator.

---

## Growth Considerations

The service is not designed to scale to millions of users. The small web is intentionally small.

However, if growth significantly exceeds expectations:

**Storage:** At 50 MB/user, the 400 GB SSD holds ~8,000 users at maximum usage. In practice, average usage will be much lower (most capsules are a handful of small files). Storage upgrades are cheap and straightforward on most VPS providers.

**Compute:** Gemini is a simple, low-overhead protocol. 8 vCPU and 24 GB RAM is significantly over-provisioned for this workload. CPU is unlikely to ever be a constraint.

**What to watch:** Disk usage percentage (alert at 70%), backup storage costs, and email sending volume (password resets).

**If the service unexpectedly becomes very popular:** Consider whether that popularity is consistent with the small web ethos before investing in scaling. There is no obligation to grow.
