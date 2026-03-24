# Explicit Non-Features

These things will never be built. This list exists so that when feature requests arrive — and they will — the answer is documented and consistent. These are not temporary deferrals. They are permanent decisions.

When someone asks for one of these, the response is: "That's intentionally not part of this service." No apology, no "maybe someday."

---

## The List

### Analytics Dashboard
**Why never:** Requires tracking reader behavior. Incompatible with the privacy commitment. Even "anonymous" analytics involve logging requests. The answer to "how many people read my capsule?" is "you don't know, and that's fine."

### Follower / Subscriber System
**Why never:** Creates social pressure dynamics. People start writing for their follower count rather than for themselves. This is how every platform loses its soul.

### Comments on Capsules
**Why never:** Moderation burden, harassment vector, scope creep into social networking. The small web has other mechanisms for conversation (linking, email, Station). This service is for publishing, not discussion.

### "Popular Capsules" or "Trending" Section
**Why never:** Algorithmic promotion creates incentives to write for engagement. Homepages with popularity signals become self-reinforcing — popular things become more popular, everything else disappears. Antithetical to the small web.

### Email Newsletters from Capsules
**Why never:** Becomes a different product entirely. Requires managing subscriptions, unsubscribes, CAN-SPAM compliance, bounce handling, deliverability. Out of scope and out of character.

### Custom Domains
**Why never:** Requires per-user TLS certificate provisioning (ACME for arbitrary domains), DNS verification, significantly more complex infrastructure, and ongoing support when users misconfigure their DNS. The subdomain approach is clean and sufficient.

### Paid Tiers with Extra Features
**Why never:** Creates a two-class system. Users who pay get a better experience; users who don't feel like second-class citizens. The "free tier" becomes the degraded version. The donation model keeps everyone equal.

### Image / Binary File Hosting
**Why never:** Gemini is a text protocol. Images hosted here would be served over HTTP, not Gemini, which is a different service. Binary uploads are a storage abuse vector (someone uploads 50 MB of MP3s). Gemtext can reference external images via links if users want that.

### Built-in Search of Capsule Content
**Why never:** To search content you must index it. Indexing requires reading and storing user content centrally. Privacy risk, storage cost, and scope creep. External search services (geminispace.info etc.) exist for discovery.

### API Access for Third-Party Applications
**Why never:** Expands the security surface dramatically. Requires API key management, rate limiting at scale, versioning, documentation, and support for breaking changes. One application (the editor) is the right scope.

### Native Mobile Application
**Why never:** The web editor works on mobile. A native app means an entirely separate codebase, App Store / Play Store overhead, review processes, and ongoing maintenance burden for one developer. Not justified.

### "Verified" Accounts or Badges
**Why never:** Creates status hierarchy. The small web has no celebrities. Everyone's capsule is equal.

### Capsule Import from WordPress / Substack / etc.
**Why never:** Sounds helpful, is actually a significant engineering project with endless edge cases. Users who want to migrate can convert their content manually or with external tools. Not a core service responsibility.

### Scheduled / Timed Publishing
**Why never:** Adds complexity to the file write pipeline, requires a job scheduler, and solves a problem that barely exists in gemtext publishing. Write when you're ready, publish when you save.

### Multiple Authors / Collaborative Capsules
**Why never:** Requires a permissions model, ownership transfer logic, conflict resolution, and shared session management. A capsule belongs to one person. Collaboration is out of scope.

### Capsule "Themes" or Custom CSS
**Why never:** Gemtext has no CSS. The rendering is entirely in the client, not the server. There is nothing to theme. (If this service ever adds a web-based capsule viewer, that viewer's appearance is the service's responsibility, not the user's.)

---

## How to Handle Feature Requests

When a user requests one of the above:
1. Thank them for the feedback
2. Explain briefly that it's intentionally not part of the service and why (one sentence)
3. Suggest alternatives if they exist (e.g. "for analytics, you might consider self-hosting a capsule where you can enable whatever you want")
4. Do not promise to "consider it" or add it to a backlog

The non-features list should be linked from the service documentation so users can understand the philosophy before requesting changes.
