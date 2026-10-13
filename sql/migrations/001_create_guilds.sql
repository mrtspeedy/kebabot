-- +goose Up
CREATE TABLE guilds (
    id BIGINT PRIMARY KEY,                        -- To store Discord's 64 bit snowflakes
    prefix TEXT NOT NULL DEFAULT '!',             -- The prefix to be used to listen for commands
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()  -- Automatically filled in timestamp
);

-- +goose Down
DROP TABLE guilds;