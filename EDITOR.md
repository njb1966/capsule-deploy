# Web Editor Specification

## Overview

The editor is a clean, minimal single-page application served at `https://yourdomain.com/editor`. It is the primary interface users interact with. It must feel like a text editor, not a CMS.

**Technology:** Vanilla HTML, CSS, and JavaScript. No frameworks, no build steps, no npm dependencies, no external CDN calls. Everything self-hosted. The editor should work in any modern browser without JavaScript bundles exceeding ~50 KB total.

---

## Layout

Three-panel layout on desktop, collapsible on mobile:

```
┌─────────────┬──────────────────────────┬───────────────────┐
│  File Tree  │       Editor Pane        │   Preview Pane    │
│             │                          │                   │
│  index.gmi  │  # My Capsule            │  rendered output  │
│  about.gmi  │                          │                   │
│  posts/     │  This is my space on...  │                   │
│    *.gmi    │                          │                   │
│             │                          │                   │
└─────────────┴──────────────────────────┴───────────────────┘
```

**Top bar (above all three panels):**
- Left: service logo/name
- Center: current filename (editable for rename)
- Right: Save button, New File button, Preview toggle, Export button, Account menu

---

## File Tree (Left Panel)

- Lists all `.gmi` files and subdirectories in the user's capsule
- Directories are expandable/collapsible
- Clicking a file opens it in the editor pane
- Right-click (or long-press on mobile) opens context menu: Rename, Delete
- "New File" button at the top of the panel
- "New Folder" button at the top of the panel
- Current open file is highlighted

---

## Editor Pane (Center Panel)

- Plain `<textarea>` with monospace font (system monospace stack: `Consolas, 'Courier New', monospace`)
- Comfortable line height: 1.6
- Font size: 15px
- Soft wrap enabled
- Tab key inserts 2 spaces (not a tab character — tabs cause issues in gemtext)
- Browser spell-check enabled by default (users may disable it in their browser settings)
- Line numbers: optional, off by default
- Auto-saves to localStorage every 30 seconds as a draft (not to server — prevents accidental overwrites)
- Unsaved changes indicator in the top bar (e.g. a dot next to the filename)

---

## Preview Pane (Right Panel)

Live preview renders gemtext as the user types. Updates on each keystroke with a 300ms debounce to avoid thrashing.

### Gemtext Rendering Rules for Preview

| Gemtext | Rendered as |
|---------|-------------|
| `# text` | `<h1>` |
| `## text` | `<h2>` |
| `### text` | `<h3>` |
| `=> URL label` | `<a href="URL">label</a>` (label falls back to URL if absent) |
| `* text` | `<li>` inside `<ul>` (consecutive `*` lines grouped) |
| `> text` | `<blockquote>` |
| ```` ``` ```` (toggle) | `<pre><code>` block until closing ```` ``` ```` |
| (any other line) | `<p>` |

Blank lines between paragraphs render as paragraph breaks.

The preview is purely visual — links are not navigable in the preview pane (prevent accidental navigation away).

---

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+S` | Save current file to server |
| `Ctrl+1` | Insert `# ` at line start |
| `Ctrl+2` | Insert `## ` at line start |
| `Ctrl+3` | Insert `### ` at line start |
| `Ctrl+L` | Insert `=> gemini://` at cursor |
| `Ctrl+Shift+L` | Insert `* ` at line start |
| `Ctrl+>` | Insert `> ` at line start |
| `` Ctrl+` `` | Insert ```` ``` ```` toggle |
| `Ctrl+Z` | Undo (browser native) |
| `Ctrl+Shift+Z` | Redo (browser native) |
| `Escape` | Dismiss any open dialog |

Shortcuts displayed in a help overlay accessible via `?` key or help icon.

---

## File Operations

### New File
1. Click "New File" button or use file tree button
2. Modal prompt: enter filename
3. Filename validated: alphanumeric, hyphens, underscores; `.gmi` extension appended automatically if not provided
4. File created on server with empty content
5. File opens in editor immediately

### New Folder
1. Click "New Folder" button
2. Modal prompt: enter folder name
3. Name validated: alphanumeric, hyphens, underscores only
4. Directory created on server
5. Appears in file tree

### Save
1. `Ctrl+S` or Save button
2. PUT request to `/api/file/<path>` with file content
3. Unsaved indicator cleared on success
4. Error toast on failure (network error, storage limit exceeded, etc.)

### Rename (file)
1. Right-click file → Rename, or select file and click Rename in top bar
2. Modal prompt: pre-filled with current filename
3. POST request to `/api/files/rename` with old and new path
4. File tree updates
5. Note displayed: "Links to this file in other pages are not automatically updated"

### Rename (folder)
1. Right-click folder → Rename, or select folder and click Rename in top bar
2. Modal prompt: pre-filled with current folder name
3. POST request to `/api/files/rename` with old and new path
4. File tree updates; expanded-folder state updated to new path
5. Warning toast: "Folder renamed. Links to files inside may be broken."

### Delete
1. Right-click file → Delete
2. Confirmation dialog: "Delete `filename.gmi`? This cannot be undone."
3. User must click Confirm (not just press Enter — prevent accidental deletion)
4. DELETE request to `/api/file/<path>`
5. If deleted file was open in editor, editor clears and shows file tree

### Export
1. Click Export button in top bar
2. GET request to `/api/export`
3. Server returns ZIP of all `.gmi` files preserving directory structure
4. Browser downloads `username-capsule-export-YYYY-MM-DD.zip`
5. No confirmation needed — export is always safe

---

## Account Settings Page

Accessible from account menu in top bar. Separate page at `/settings`.

- Change email address (requires re-verification)
- Change password (requires current password)
- Storage usage indicator: "X MB of 50 MB used"
- Export capsule button (same as editor export)
- Delete account section (at bottom, visually separated, destructive styling)

### Account Deletion Flow
1. User clicks "Delete my account"
2. Warning displayed: "This will permanently delete your capsule and all files. This cannot be undone."
3. User must type their username to confirm
4. DELETE request to `/api/account`
5. All files deleted from disk
6. Account record deleted from database
7. Username held for 30 days before being available again
8. User redirected to landing page with message: "Your account has been deleted."

---

## Mobile Behavior

On screens narrower than 768px:
- File tree hidden by default, accessible via hamburger/drawer
- Preview pane hidden by default, toggle button in top bar to show
- Editor pane takes full width
- All keyboard shortcuts still work on hardware keyboards
- Touch-friendly tap targets (minimum 44px)

---

## Error States

| Error | Display |
|-------|---------|
| Save failed (network) | Toast: "Save failed — check your connection. Your changes are preserved locally." |
| Save failed (storage limit) | Toast: "Storage limit reached (50 MB). Delete some files to continue." |
| File not found | Editor shows: "This file no longer exists." with link back to file tree |
| Session expired | Redirect to login with message: "Your session expired. Please log in again." |
| Server error | Toast: "Something went wrong. Please try again." |

---

## Starter Capsule Content

When a new account is created, this default `index.gmi` is written to their capsule directory:

```gemini
# Welcome to my capsule

This is the beginning of my space on the small web.

## About

Tell people a little about yourself here.

## Posts

Your writing will go here.

## Links

=> gemini://geminispace.info  Geminispace — search and discover other capsules
```

This ensures the capsule is valid and non-empty from the first second, and gives the user something to edit rather than a blank page.
