# capsule-forum Admin Reference

## Server locations

| Item | Path |
|---|---|
| Binary | `/usr/local/bin/capsule-forum` |
| Config | `/etc/capsule-forum/config.toml` |
| TLS cert | `/etc/capsule-forum/forum-cert.pem` |
| TLS key | `/etc/capsule-forum/forum-key.pem` |
| Database | `/var/lib/capsule-forum/forum.db` |
| Systemd unit | `/etc/systemd/system/capsule-forum.service` |
| Backups | `debian-mac2:/home/nick/backups/capsule-forum/` |

All commands below run on **contabo1** (`ssh contabo1`) unless noted.

---

## Service management

```bash
systemctl status capsule-forum
systemctl restart capsule-forum
systemctl stop capsule-forum
journalctl -u capsule-forum -n 50 --no-pager
journalctl -u capsule-forum -f          # live log tail
```

---

## Database access

Open an interactive SQLite shell:

```bash
sqlite3 /var/lib/capsule-forum/forum.db
```

Useful read queries:

```sql
-- list all users
SELECT id, username, datetime(created_at, 'unixepoch'), banned FROM users;

-- list all threads with post counts
SELECT t.id, b.slug, u.username, t.subject, t.post_count,
       datetime(t.last_post_at, 'unixepoch')
FROM threads t
JOIN boards b ON b.id = t.board_id
JOIN users  u ON u.id = t.user_id
ORDER BY t.last_post_at DESC;

-- read all posts in a thread (replace 1 with thread id)
SELECT p.id, u.username, datetime(p.created_at, 'unixepoch'), p.body
FROM posts p JOIN users u ON u.id = p.user_id
WHERE p.thread_id = 1
ORDER BY p.created_at;
```

---

## Banning a user

Prevents further posting. Existing posts remain visible.

```bash
sqlite3 /var/lib/capsule-forum/forum.db \
  "UPDATE users SET banned = 1 WHERE username = 'badactor';"
```

Unban:

```bash
sqlite3 /var/lib/capsule-forum/forum.db \
  "UPDATE users SET banned = 0 WHERE username = 'badactor';"
```

---

## Deleting a post

```bash
sqlite3 /var/lib/capsule-forum/forum.db \
  "DELETE FROM posts WHERE id = 42;"
```

If it was the only post in the thread, also delete the thread:

```bash
sqlite3 /var/lib/capsule-forum/forum.db "
  DELETE FROM posts   WHERE thread_id = 7;
  DELETE FROM threads WHERE id = 7;
"
```

After deleting posts manually, fix the post count on the thread:

```bash
sqlite3 /var/lib/capsule-forum/forum.db "
  UPDATE threads SET post_count = (
    SELECT COUNT(*) FROM posts WHERE thread_id = threads.id
  );
"
```

---

## Deleting a user

Only do this if the account is spam/empty. If the user has posts, ban instead.

```bash
sqlite3 /var/lib/capsule-forum/forum.db "
  DELETE FROM drafts WHERE fingerprint =
    (SELECT fingerprint FROM users WHERE username = 'spammer');
  DELETE FROM users WHERE username = 'spammer';
"
```

This will fail with a foreign key error if the user has posts — delete or reassign those first.

---

## Clearing a stuck draft

If a user reports being stuck in the new-thread flow:

```bash
sqlite3 /var/lib/capsule-forum/forum.db \
  "DELETE FROM drafts WHERE fingerprint IN
     (SELECT fingerprint FROM users WHERE username = 'alice');"
```

---

## Deploying an update

Build on the local dev machine and push:

```bash
# in capsule-forum/
GOOS=linux GOARCH=amd64 go build -o /tmp/capsule-forum .
scp /tmp/capsule-forum contabo1:/usr/local/bin/capsule-forum
ssh contabo1 "systemctl restart capsule-forum && systemctl status capsule-forum --no-pager"
```

---

## TLS certificate renewal

The self-signed cert is valid for 10 years (generated 2026-05-11, expires 2036-05-09).
When renewal is needed:

```bash
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 \
  -keyout /tmp/forum-key.pem -out /tmp/forum-cert.pem -days 3650 -nodes \
  -subj "/CN=forum.gemcities.com"

scp /tmp/forum-cert.pem /tmp/forum-key.pem contabo1:/etc/capsule-forum/
ssh contabo1 "
  chown root:capsule-forum /etc/capsule-forum/forum-key.pem /etc/capsule-forum/forum-cert.pem
  chmod 640 /etc/capsule-forum/forum-key.pem /etc/capsule-forum/forum-cert.pem
  systemctl restart capsule-forum
"
```

Note: Gemini clients use TOFU, so a new cert will prompt users to re-accept on first visit after renewal.

---

## Backups

Backups run nightly via cron on **debian-mac2**, stored at:

```
/home/nick/backups/capsule-forum/forum-YYYY-MM-DD.db
```

Kept for 30 days. To run a manual backup immediately:

```bash
# on debian-mac2
rsync -az root@217.77.0.238:/var/lib/capsule-forum/forum.db \
  /home/nick/backups/capsule-forum/forum-$(date +%F).db
```

### Restore from backup

```bash
# on contabo1
systemctl stop capsule-forum
cp /var/lib/capsule-forum/forum.db /var/lib/capsule-forum/forum.db.pre-restore
# copy backup from debian-mac2
rsync -az nick@192.168.0.20:/home/nick/backups/capsule-forum/forum-2026-05-11.db \
  /var/lib/capsule-forum/forum.db
chown capsule-forum:capsule-forum /var/lib/capsule-forum/forum.db
systemctl start capsule-forum
```

---

## Checking cert fingerprints

If a user reports a registration problem or you need to identify which cert belongs to whom:

```bash
sqlite3 /var/lib/capsule-forum/forum.db \
  "SELECT username, fingerprint, datetime(created_at, 'unixepoch') FROM users;"
```

A user's fingerprint can be confirmed in Lagrange under Identities → select the identity → the SHA-256 fingerprint is shown there.
