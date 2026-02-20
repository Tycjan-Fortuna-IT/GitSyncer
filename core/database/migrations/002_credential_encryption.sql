-- +goose Up

ALTER TABLE credentials ADD COLUMN salt BLOB NOT NULL DEFAULT x'';

-- +goose Down

ALTER TABLE credentials DROP COLUMN salt;
