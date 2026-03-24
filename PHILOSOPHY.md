# Philosophy & Guiding Principles

## The Core Commitments

These are not aspirations or marketing copy. They are hard constraints that shape what gets built and what gets permanently refused.

**No surveillance.**
No pageview counters, no visitor graphs, no heatmaps, no session recording. Users who read capsules are not tracked in any way. Application logs do not record IP addresses of readers, what capsules are visited, how long a user spends on a page, or referrer information.

**No engagement mechanics.**
No likes, no follower counts, no trending sections, no algorithmic promotion. Nothing that creates incentives to write for an audience rather than for oneself.

**No ads.**
Ever. Not now, not if the service grows, not as a "free tier" trade-off.

**Radical simplicity.**
If a feature could be cut without losing the core purpose, it should be cut. Complexity is a liability, not an asset. When in doubt, do less.

**Easy exit.**
Users can export all their content as raw gemtext files at any time with one click. The service should never feel like a trap. See `docs/EDITOR.md` for export implementation.

**Open source.**
All server software is published under AGPL-3.0. Anyone can run the same stack. If this service shuts down tomorrow, the ecosystem survives.

---

## What This Service Is

A piece of infrastructure for people who want to publish on the Gemini small web without managing a server. That is the whole thing.

It is not a social network.
It is not competing with WordPress or Substack.
It is not trying to grow into something larger.
It is not optimizing for retention, engagement, or growth metrics.

---

## What This Service Is Not Trying to Do

- Convert people to Gemini (users already know what they want)
- Build a community (community emerges from good tools, not from features designed to create it)
- Monetize user attention in any form
- Scale to millions of users (the small web is small on purpose)

---

## The Philosophical Test

When evaluating any new feature request or technical decision, ask:

1. Does this make it easier to publish on Gemini? If no, it probably doesn't belong.
2. Does this require tracking, storing, or analyzing user behavior? If yes, refuse it.
3. Does this create any social pressure or engagement dynamic? If yes, refuse it.
4. Does this add ongoing maintenance burden disproportionate to its value? If yes, cut it.
5. Would this change what kind of person uses the service? Consider carefully.

---

## On Community Pressure

The small web community will request features. Some requests will be well-intentioned and reasonable-sounding. The job is to hold the line on the commitments above even when the argument for adding something seems compelling.

History of every platform that started simple: each feature addition seemed reasonable at the time. The cumulative effect was not.

"We could add X, it would only be small" is how every platform lost its way. Saying no is a feature.
