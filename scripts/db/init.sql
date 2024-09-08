CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    tg_message_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    length INT,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    chat_id BIGINT NOT NULL,
    sender_id TEXT,
    sender_username TEXT
);
CREATE EXTENSION IF NOT EXISTS "vector";
CREATE TABLE IF NOT EXISTS embeddings_3large (
    id SERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL,
    embedding vector(2000), -- 2000 is max for index
    created_at TIMESTAMP DEFAULT NOW()
    );
CREATE INDEX IF NOT EXISTS idx_emb3l_emb ON embeddings_3large USING hnsw (embedding vector_cosine_ops);