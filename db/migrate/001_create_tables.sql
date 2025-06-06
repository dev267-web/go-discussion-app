-- db/migrate/001_create_tables.sql

-- 1) users
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    full_name       VARCHAR(255),
    bio             TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2) discussions
CREATE TABLE IF NOT EXISTS discussions (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    content         TEXT NOT NULL,
    scheduled_at    TIMESTAMPTZ,          -- NULL if “post immediately”
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3) comments
CREATE TABLE IF NOT EXISTS comments (
    id              SERIAL PRIMARY KEY,
    discussion_id   INTEGER NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4) tags
CREATE TABLE IF NOT EXISTS tags (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5) discussion_tags (many-to-many between discussions and tags)
CREATE TABLE IF NOT EXISTS discussion_tags (
    discussion_id   INTEGER NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    tag_id          INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (discussion_id, tag_id)
);

-- 6) subscriptions (email notifications)
CREATE TABLE IF NOT EXISTS subscriptions (
    id              SERIAL PRIMARY KEY,
    discussion_id   INTEGER NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
    email           VARCHAR(255) NOT NULL,
    subscribed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (discussion_id, email)
);

-- 7) (Optional) You can add an index on scheduled_at if you’ll be querying upcoming discussions frequently
CREATE INDEX IF NOT EXISTS idx_discussions_scheduled_at
    ON discussions(scheduled_at);

-- 8) (Optional) Index on created_at for ordering
CREATE INDEX IF NOT EXISTS idx_discussions_created_at
    ON discussions(created_at);
