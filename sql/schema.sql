CREATE TABLE IF NOT EXISTS todos
(
    id
    TEXT
    PRIMARY
    KEY,
    title
    TEXT
    NOT
    NULL,
    completed
    BOOLEAN
    NOT
    NULL
    DEFAULT
    FALSE,
    created_at
    TIMESTAMPTZ
    NOT
    NULL,
    updated_at
    TIMESTAMPTZ
    NOT
    NULL
);

CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos (created_at DESC);
