-- +goose Up

CREATE TABLE providers (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT    NOT NULL,
    type            TEXT    NOT NULL,
    base_url        TEXT    NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE repositories (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_id     INTEGER NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    name            TEXT    NOT NULL,
    clone_url       TEXT    NOT NULL,
    description     TEXT    NOT NULL DEFAULT '',
    is_mirror       INTEGER NOT NULL DEFAULT 0,
    default_branch  TEXT    NOT NULL DEFAULT 'master',
    last_synced_at  DATETIME,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE credentials (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_id     INTEGER NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    label           TEXT    NOT NULL DEFAULT '',
    auth_type       TEXT    NOT NULL,
    auth_data       TEXT    NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE sync_schedules (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id   INTEGER NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    cron_expr       TEXT    NOT NULL,
    enabled         INTEGER NOT NULL DEFAULT 1,
    last_run_at     DATETIME,
    next_run_at     DATETIME,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE sync_history (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id   INTEGER NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    status          TEXT    NOT NULL,
    started_at      DATETIME NOT NULL DEFAULT (datetime('now')),
    finished_at     DATETIME,
    error_message   TEXT    NOT NULL DEFAULT '',
    details         TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE settings (
    key             TEXT PRIMARY KEY,
    value           TEXT NOT NULL DEFAULT '',
    updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_repositories_provider_id ON repositories(provider_id);
CREATE INDEX idx_credentials_provider_id ON credentials(provider_id);
CREATE INDEX idx_sync_schedules_repository_id ON sync_schedules(repository_id);
CREATE INDEX idx_sync_history_repository_id ON sync_history(repository_id);
CREATE INDEX idx_sync_history_status ON sync_history(status);

-- +goose Down

DROP INDEX IF EXISTS idx_sync_history_status;
DROP INDEX IF EXISTS idx_sync_history_repository_id;
DROP INDEX IF EXISTS idx_sync_schedules_repository_id;
DROP INDEX IF EXISTS idx_credentials_provider_id;
DROP INDEX IF EXISTS idx_repositories_provider_id;

DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS sync_history;
DROP TABLE IF EXISTS sync_schedules;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS repositories;
DROP TABLE IF EXISTS providers;
