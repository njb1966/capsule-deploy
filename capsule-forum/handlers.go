package main

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	gemini "git.sr.ht/~adnano/go-gemini"
)

func buildMux(db *DB) *gemini.Mux {
	mux := &gemini.Mux{}
	mux.HandleFunc("/", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		homeHandler(w, db)
	})
	mux.HandleFunc("/register", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		registerHandler(w, r, db)
	})
	mux.HandleFunc("/profile", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		profileHandler(w, r, db)
	})
	mux.HandleFunc("/board/", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		boardRouter(w, r, db)
	})
	mux.HandleFunc("/thread/", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		threadRouter(w, r, db)
	})
	return mux
}

// --- helpers ---

func fingerprint(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sum[:])
}

func peerCert(r *gemini.Request) *x509.Certificate {
	tls := r.TLS()
	if tls == nil || len(tls.PeerCertificates) == 0 {
		return nil
	}
	return tls.PeerCertificates[0]
}

func writeHeader(w gemini.ResponseWriter) {
	w.WriteHeader(gemini.StatusSuccess, "text/gemini; charset=utf-8")
}

func redirectTo(w gemini.ResponseWriter, target string) {
	w.WriteHeader(gemini.StatusRedirect, target)
}

func prompt(w gemini.ResponseWriter, p string) {
	w.WriteHeader(gemini.StatusInput, p)
}

func certRequired(w gemini.ResponseWriter) {
	w.WriteHeader(gemini.StatusCertificateRequired,
		"A client certificate is required. Please activate one in your Gemini client.")
}

func notFound(w gemini.ResponseWriter) {
	w.WriteHeader(gemini.StatusNotFound, "Not found")
}

func fail(w gemini.ResponseWriter, msg string) {
	w.WriteHeader(gemini.StatusPermanentFailure, msg)
}

func fmtTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04")
}

func validUsername(s string) bool {
	if len(s) < 3 || len(s) > 20 {
		return false
	}
	if s[0] == '-' || s[len(s)-1] == '-' {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}
	return true
}

func queryString(r *gemini.Request) string {
	if r.URL.RawQuery == "" {
		return ""
	}
	s, _ := url.QueryUnescape(r.URL.RawQuery)
	return strings.TrimSpace(s)
}

// --- handlers ---

func homeHandler(w gemini.ResponseWriter, db *DB) {
	boards, err := db.getBoards()
	if err != nil {
		fail(w, "Internal error")
		return
	}

	writeHeader(w)
	fmt.Fprintln(w, "# GemCities Forum")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "A community space for GemCities capsule authors.")
	fmt.Fprintln(w, "Read freely. To post, activate a client certificate in your Gemini client.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Boards")
	fmt.Fprintln(w)
	for _, b := range boards {
		fmt.Fprintf(w, "=> /board/%s %s — %s\n", b.Slug, b.Name, b.Description)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Account")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "=> /register Register or check registration")
	fmt.Fprintln(w, "=> /profile  Your profile")
}

func registerHandler(w gemini.ResponseWriter, r *gemini.Request, db *DB) {
	cert := peerCert(r)
	if cert == nil {
		certRequired(w)
		return
	}

	fp := fingerprint(cert)
	user, err := db.getUserByFingerprint(fp)
	if err != nil {
		fail(w, "Internal error")
		return
	}
	if user != nil {
		redirectTo(w, "/profile")
		return
	}

	input := queryString(r)
	if input == "" {
		prompt(w, "Choose a username (3–20 chars, letters/digits/hyphens):")
		return
	}

	username := strings.ToLower(input)
	if !validUsername(username) {
		prompt(w, "Invalid username — 3–20 chars, letters, digits, hyphens, no leading/trailing hyphens:")
		return
	}

	exists, err := db.usernameExists(username)
	if err != nil {
		fail(w, "Internal error")
		return
	}
	if exists {
		prompt(w, "That username is taken. Choose another:")
		return
	}

	if _, err := db.createUser(fp, username); err != nil {
		fail(w, "Could not create account")
		return
	}
	redirectTo(w, "/profile")
}

func profileHandler(w gemini.ResponseWriter, r *gemini.Request, db *DB) {
	cert := peerCert(r)
	if cert == nil {
		certRequired(w)
		return
	}

	fp := fingerprint(cert)
	user, err := db.getUserByFingerprint(fp)
	if err != nil {
		fail(w, "Internal error")
		return
	}
	if user == nil {
		redirectTo(w, "/register")
		return
	}

	writeHeader(w)
	fmt.Fprintln(w, "# Your Profile")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Username:    %s\n", user.Username)
	fmt.Fprintf(w, "Member since: %s UTC\n", fmtTime(user.CreatedAt))
	fmt.Fprintln(w)
	fmt.Fprintln(w, "=> / Forum home")
}

func boardRouter(w gemini.ResponseWriter, r *gemini.Request, db *DB) {
	// /board/{slug}
	// /board/{slug}/new
	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 3)
	if len(parts) < 2 || parts[1] == "" {
		notFound(w)
		return
	}

	board, err := db.getBoardBySlug(parts[1])
	if err != nil || board == nil {
		notFound(w)
		return
	}

	if len(parts) == 3 && parts[2] == "new" {
		newThreadHandler(w, r, db, board)
		return
	}

	boardHandler(w, db, board)
}

func boardHandler(w gemini.ResponseWriter, db *DB, board *Board) {
	threads, err := db.getThreadsByBoard(board.ID)
	if err != nil {
		fail(w, "Internal error")
		return
	}

	writeHeader(w)
	fmt.Fprintf(w, "# %s\n", board.Name)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s\n", board.Description)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "=> /board/%s/new + New Thread\n", board.Slug)
	fmt.Fprintln(w)

	if len(threads) == 0 {
		fmt.Fprintln(w, "No threads yet. Be the first to post.")
	} else {
		fmt.Fprintln(w, "## Threads")
		fmt.Fprintln(w)
		for _, t := range threads {
			fmt.Fprintf(w, "=> /thread/%d [%d] %s — by %s — %s\n",
				t.ID, t.PostCount, t.Subject, t.Username, fmtTime(t.LastPostAt))
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "=> / Forum home")
}

func newThreadHandler(w gemini.ResponseWriter, r *gemini.Request, db *DB, board *Board) {
	cert := peerCert(r)
	if cert == nil {
		certRequired(w)
		return
	}

	fp := fingerprint(cert)
	user, err := db.getUserByFingerprint(fp)
	if err != nil {
		fail(w, "Internal error")
		return
	}
	if user == nil {
		redirectTo(w, "/register")
		return
	}
	if user.Banned {
		fail(w, "Your account has been suspended.")
		return
	}

	input := queryString(r)
	draft, err := db.getDraft(fp, board.ID)
	if err != nil {
		fail(w, "Internal error")
		return
	}

	switch {
	case draft == "" && input == "":
		// Step 1: ask for subject
		prompt(w, fmt.Sprintf("New thread in %q — enter subject:", board.Name))

	case draft == "" && input != "":
		// Step 2: received subject, ask for body
		if len(input) > 200 {
			prompt(w, "Subject too long (max 200 chars). Enter subject:")
			return
		}
		if err := db.saveDraft(fp, board.ID, input); err != nil {
			fail(w, "Internal error")
			return
		}
		prompt(w, fmt.Sprintf("Subject saved. Now enter your post body:"))

	case draft != "" && input == "":
		// User navigated back — re-prompt for body
		prompt(w, fmt.Sprintf("Subject: %q — enter your post body:", draft))

	case draft != "" && input != "":
		// Step 3: received body, create thread
		threadID, err := db.createThread(board.ID, user.ID, draft, input)
		if err != nil {
			fail(w, "Could not create thread")
			return
		}
		_ = db.deleteDraft(fp, board.ID)
		redirectTo(w, fmt.Sprintf("/thread/%d", threadID))
	}
}

func threadRouter(w gemini.ResponseWriter, r *gemini.Request, db *DB) {
	// /thread/{id}
	// /thread/{id}/reply
	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 3)
	if len(parts) < 2 || parts[1] == "" {
		notFound(w)
		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		notFound(w)
		return
	}

	thread, err := db.getThread(id)
	if err != nil || thread == nil {
		notFound(w)
		return
	}

	if len(parts) == 3 && parts[2] == "reply" {
		replyHandler(w, r, db, thread)
		return
	}

	threadHandler(w, db, thread)
}

func threadHandler(w gemini.ResponseWriter, db *DB, thread *Thread) {
	posts, err := db.getPostsByThread(thread.ID)
	if err != nil {
		fail(w, "Internal error")
		return
	}

	writeHeader(w)
	fmt.Fprintf(w, "# %s\n", thread.Subject)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Board: %s | Started by %s | %s UTC\n",
		thread.BoardSlug, thread.Username, fmtTime(thread.CreatedAt))
	fmt.Fprintln(w)
	fmt.Fprintf(w, "=> /thread/%d/reply + Reply\n", thread.ID)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "---")
	fmt.Fprintln(w)

	for i, p := range posts {
		fmt.Fprintf(w, "## #%d — %s — %s UTC\n", i+1, p.Username, fmtTime(p.CreatedAt))
		fmt.Fprintln(w)
		fmt.Fprintln(w, p.Body)
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "=> /board/%s Back to %s\n", thread.BoardSlug, thread.BoardSlug)
	fmt.Fprintln(w, "=> / Forum home")
}

func replyHandler(w gemini.ResponseWriter, r *gemini.Request, db *DB, thread *Thread) {
	cert := peerCert(r)
	if cert == nil {
		certRequired(w)
		return
	}

	fp := fingerprint(cert)
	user, err := db.getUserByFingerprint(fp)
	if err != nil {
		fail(w, "Internal error")
		return
	}
	if user == nil {
		redirectTo(w, "/register")
		return
	}
	if user.Banned {
		fail(w, "Your account has been suspended.")
		return
	}

	input := queryString(r)
	if input == "" {
		prompt(w, fmt.Sprintf("Reply to \"%s\":", thread.Subject))
		return
	}

	if err := db.createPost(thread.ID, user.ID, input); err != nil {
		fail(w, "Could not post reply")
		return
	}
	redirectTo(w, fmt.Sprintf("/thread/%d", thread.ID))
}
