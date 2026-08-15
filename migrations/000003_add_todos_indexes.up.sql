CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos (completed);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos (created_at DESC);