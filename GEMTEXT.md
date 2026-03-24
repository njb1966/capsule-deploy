# Gemtext Format Reference

Gemtext is the native document format of the Gemini protocol. It is intentionally simple. There are seven line types. Every line is interpreted independently based on its prefix — there is no inline formatting, no nesting, no block context (with one exception: preformatted blocks).

---

## Line Types

### Heading Lines
```
# First level heading
## Second level heading
### Third level heading
```
Headings must have the `#` at the very start of the line. A space after `#` is conventional and recommended. Rendered as H1, H2, H3 respectively.

### Link Lines
```
=> gemini://example.com
=> gemini://example.com A human-readable label
=> https://example.com Links to other protocols are valid
```
Links must be on their own line. The URL immediately follows `=>` and a space. An optional label follows the URL separated by whitespace. If no label is provided, clients typically display the URL itself. This is the **only** way to create links in gemtext — there is no inline linking.

### List Item Lines
```
* An unordered list item
* Another item
```
List items begin with `* ` (asterisk space). Consecutive list item lines are typically grouped visually by clients. There are no ordered/numbered lists in gemtext.

### Blockquote Lines
```
> This is a quoted passage.
```
Begins with `> ` (greater-than space).

### Preformatted Toggle Lines
```
```
(content here is rendered as-is, monospace, no line-type interpretation)
```
```
The triple-backtick line toggles preformatted mode on and off. Content between two toggle lines is treated as a preformatted block — no gemtext interpretation occurs inside. An optional alt-text label can follow the opening toggle: ` ```alt text here ` — clients may display this as a caption or use it for accessibility.

### Text Lines
Any line that does not match the above prefixes is a text/paragraph line. Rendered as a paragraph. Blank lines are valid text lines and typically render as paragraph spacing.

---

## Key Rules and Constraints

**One element per line.** There is no inline formatting. You cannot bold a word, italicize a phrase, or create an inline link within a paragraph. Emphasis is achieved through word choice and structure, not markup.

**Links must be on their own line.** This is the most common adjustment for people coming from Markdown or HTML. A paragraph cannot contain a link. Links are always standalone.

**No tables.** Gemtext has no table syntax. Tabular data is typically presented as preformatted text or described in prose.

**No images.** Gemtext cannot embed images. A link line pointing to an image URL is valid, and some clients will display it inline, but this is client behavior, not part of the spec.

**No nested structures.** You cannot nest lists, nest blockquotes, or create definition lists. The format is intentionally flat.

---

## File Conventions

- File extension: `.gmi` (conventional, not strictly required by the protocol)
- Encoding: UTF-8
- Line endings: LF (`\n`) preferred, CRLF tolerated by most clients
- MIME type: `text/gemini`

---

## Example Document

```gemini
# My Capsule

Welcome to my small corner of the internet.

## Recent Posts

=> /posts/2025-03-01-hello.gmi  Hello, Gemini — my first post
=> /posts/2025-03-15-tools.gmi  The tools I use every day

## About Me

I write about technology, woodworking, and whatever else is on my mind.

This is a text paragraph. It can be as long as you like.
Another text paragraph follows.

## Elsewhere

=> https://sr.ht/~username  My sourcehut profile
=> gemini://geminispace.info  Search the Geminispace

> "The best time to plant a tree was twenty years ago.
> The second best time is now."

## Code Sample

```
#!/bin/sh
echo "Hello, world"
```
```

---

## Rendering Notes for the Editor Preview

The live preview in the editor renders gemtext using these rules:

- `#`, `##`, `###` → `<h1>`, `<h2>`, `<h3>`
- `=> URL label` → `<a href="URL">label</a>` (label falls back to URL)
- `* text` → `<li>` items, consecutive `*` lines grouped in `<ul>`
- `> text` → `<blockquote>`
- ` ``` ` toggle → `<pre><code>` block
- Everything else → `<p>`
- Blank lines → paragraph break (rendered as `<p>&nbsp;</p>` or CSS margin)

Links in the preview are **not clickable** — this prevents accidental navigation away from the editor.

---

## Gemtext vs Markdown

| Feature | Gemtext | Markdown |
|---------|---------|----------|
| Headings | Yes (`#`) | Yes (`#`) |
| Bold / italic | No | Yes |
| Inline links | No | Yes |
| Standalone links | Yes (`=>`) | Via inline links |
| Images | No (link only) | Yes (`![]()`) |
| Tables | No | Yes (GFM) |
| Code blocks | Yes (toggle) | Yes (fences/indent) |
| Ordered lists | No | Yes |
| Unordered lists | Yes (`*`) | Yes |
| Nested elements | No | Yes |
| Inline HTML | No | Yes (CommonMark) |

Gemtext is intentionally less capable than Markdown. This is a feature, not a bug.
