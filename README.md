# capsule-deploy

Server configuration, systemd units, and setup scripts for [GemCities](https://gemcities.com) — a free Gemini capsule hosting service.

Use this repo to self-host your own Gemini capsule hosting service.

## Contents (Planned)

- `caddy/` — Caddyfile template
- `agate/` — Agate config template (with logging disabled)
- `systemd/` — Unit files for capsule-service, agate, caddy
- `scripts/` — VPS setup, backup (Restic), TLS cert automation
- `config/` — Annotated config.toml template for capsule-service

## Related Repos

- [capsule-service](https://github.com/njb1966/capsule-service) — Go backend API (also contains full project docs)
- [capsule-editor](https://github.com/njb1966/capsule-editor) — Web editor frontend

## License

AGPL-3.0 — see [LICENSE](LICENSE)