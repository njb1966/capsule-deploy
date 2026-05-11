package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

type User struct {
	ID          int64
	Fingerprint string
	Username    string
	CreatedAt   time.Time
	Banned      bool
}

type Board struct {
	ID          int64
	Slug        string
	Name        string
	Description string
}

type Thread struct {
	ID         int64
	BoardID    int64
	BoardSlug  string
	UserID     int64
	Username   string
	Subject    string
	CreatedAt  time.Time
	LastPostAt time.Time
	PostCount  int
}

type Post struct {
	ID        int64
	ThreadID  int64
	UserID    int64
	Username  string
	Body      string
	CreatedAt time.Time
}

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	fingerprint TEXT    UNIQUE NOT NULL,
	username    TEXT    UNIQUE NOT NULL,
	created_at  INTEGER NOT NULL DEFAULT (unixepoch()),
	banned      INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS boards (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	slug        TEXT UNIQUE NOT NULL,
	name        TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS threads (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	board_id     INTEGER NOT NULL REFERENCES boards(id),
	user_id      INTEGER NOT NULL REFERENCES users(id),
	subject      TEXT    NOT NULL,
	created_at   INTEGER NOT NULL DEFAULT (unixepoch()),
	last_post_at INTEGER NOT NULL DEFAULT (unixepoch()),
	post_count   INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS posts (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	thread_id  INTEGER NOT NULL REFERENCES threads(id),
	user_id    INTEGER NOT NULL REFERENCES users(id),
	body       TEXT    NOT NULL,
	created_at INTEGER NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE IF NOT EXISTS drafts (
	fingerprint TEXT    NOT NULL,
	board_id    INTEGER NOT NULL REFERENCES boards(id),
	subject     TEXT    NOT NULL,
	PRIMARY KEY (fingerprint, board_id)
);
`

var seedBoards = []Board{
	{Slug: "introductions", Name: "Introductions", Description: "Introduce yourself and share your capsule"},
	{Slug: "general", Name: "General Discussion", Description: "Anything and everything Gemini"},
	{Slug: "help", Name: "Help & Support", Description: "Questions about GemCities and Gemini"},
	{Slug: "offtopic", Name: "Off-Topic", Description: "Everything else"},
}

func openDB(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)

	if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"); err != nil {
		return nil, err
	}
	if _, err := sqlDB.Exec(schema); err != nil {
		return nil, fmt.Errorf("schema: %w", err)
	}

	db := &DB{sqlDB}
	for _, b := range seedBoards {
		if _, err := db.Exec(
			`INSERT OR IGNORE INTO boards (slug, name, description) VALUES (?, ?, ?)`,
			b.Slug, b.Name, b.Description,
		); err != nil {
			return nil, fmt.Errorf("seed boards: %w", err)
		}
	}
	return db, nil
}

func (db *DB) getUserByFingerprint(fp string) (*User, error) {
	u := &User{}
	var ts int64
	var banned int
	err := db.QueryRow(
		`SELECT id, fingerprint, username, created_at, banned FROM users WHERE fingerprint = ?`, fp,
	).Scan(&u.ID, &u.Fingerprint, &u.Username, &ts, &banned)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt = time.Unix(ts, 0)
	u.Banned = banned != 0
	return u, nil
}

func (db *DB) usernameExists(username string) (bool, error) {
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&n)
	return n > 0, err
}

func (db *DB) createUser(fp, username string) (*User, error) {
	res, err := db.Exec(`INSERT INTO users (fingerprint, username) VALUES (?, ?)`, fp, username)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &User{ID: id, Fingerprint: fp, Username: username, CreatedAt: time.Now()}, nil
}

func (db *DB) getBoards() ([]Board, error) {
	rows, err := db.Query(`SELECT id, slug, name, description FROM boards ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var boards []Board
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.ID, &b.Slug, &b.Name, &b.Description); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

func (db *DB) getBoardBySlug(slug string) (*Board, error) {
	b := &Board{}
	err := db.QueryRow(
		`SELECT id, slug, name, description FROM boards WHERE slug = ?`, slug,
	).Scan(&b.ID, &b.Slug, &b.Name, &b.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return b, err
}

func (db *DB) getThreadsByBoard(boardID int64) ([]Thread, error) {
	rows, err := db.Query(`
		SELECT t.id, t.board_id, b.slug, t.user_id, u.username,
		       t.subject, t.created_at, t.last_post_at, t.post_count
		FROM threads t
		JOIN users  u ON u.id = t.user_id
		JOIN boards b ON b.id = t.board_id
		WHERE t.board_id = ?
		ORDER BY t.last_post_at DESC`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanThreads(rows)
}

func (db *DB) getThread(id int64) (*Thread, error) {
	t := &Thread{}
	var createdAt, lastPostAt int64
	err := db.QueryRow(`
		SELECT t.id, t.board_id, b.slug, t.user_id, u.username,
		       t.subject, t.created_at, t.last_post_at, t.post_count
		FROM threads t
		JOIN users  u ON u.id = t.user_id
		JOIN boards b ON b.id = t.board_id
		WHERE t.id = ?`, id,
	).Scan(&t.ID, &t.BoardID, &t.BoardSlug, &t.UserID, &t.Username,
		&t.Subject, &createdAt, &lastPostAt, &t.PostCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t.CreatedAt = time.Unix(createdAt, 0)
	t.LastPostAt = time.Unix(lastPostAt, 0)
	return t, nil
}

func (db *DB) createThread(boardID, userID int64, subject, body string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT INTO threads (board_id, user_id, subject) VALUES (?, ?, ?)`,
		boardID, userID, subject,
	)
	if err != nil {
		return 0, err
	}
	threadID, _ := res.LastInsertId()

	if _, err = tx.Exec(
		`INSERT INTO posts (thread_id, user_id, body) VALUES (?, ?, ?)`,
		threadID, userID, body,
	); err != nil {
		return 0, err
	}
	return threadID, tx.Commit()
}

func (db *DB) getPostsByThread(threadID int64) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.thread_id, p.user_id, u.username, p.body, p.created_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.thread_id = ?
		ORDER BY p.created_at ASC`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		var ts int64
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.UserID, &p.Username, &p.Body, &ts); err != nil {
			return nil, err
		}
		p.CreatedAt = time.Unix(ts, 0)
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (db *DB) createPost(threadID, userID int64, body string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(
		`INSERT INTO posts (thread_id, user_id, body) VALUES (?, ?, ?)`,
		threadID, userID, body,
	); err != nil {
		return err
	}
	if _, err = tx.Exec(
		`UPDATE threads SET last_post_at = unixepoch(), post_count = post_count + 1 WHERE id = ?`,
		threadID,
	); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) getDraft(fingerprint string, boardID int64) (string, error) {
	var subject string
	err := db.QueryRow(
		`SELECT subject FROM drafts WHERE fingerprint = ? AND board_id = ?`,
		fingerprint, boardID,
	).Scan(&subject)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return subject, err
}

func (db *DB) saveDraft(fingerprint string, boardID int64, subject string) error {
	_, err := db.Exec(
		`INSERT OR REPLACE INTO drafts (fingerprint, board_id, subject) VALUES (?, ?, ?)`,
		fingerprint, boardID, subject,
	)
	return err
}

func (db *DB) deleteDraft(fingerprint string, boardID int64) error {
	_, err := db.Exec(
		`DELETE FROM drafts WHERE fingerprint = ? AND board_id = ?`,
		fingerprint, boardID,
	)
	return err
}

func scanThreads(rows *sql.Rows) ([]Thread, error) {
	var threads []Thread
	for rows.Next() {
		var t Thread
		var createdAt, lastPostAt int64
		if err := rows.Scan(
			&t.ID, &t.BoardID, &t.BoardSlug, &t.UserID, &t.Username,
			&t.Subject, &createdAt, &lastPostAt, &t.PostCount,
		); err != nil {
			return nil, err
		}
		t.CreatedAt = time.Unix(createdAt, 0)
		t.LastPostAt = time.Unix(lastPostAt, 0)
		threads = append(threads, t)
	}
	return threads, rows.Err()
}
