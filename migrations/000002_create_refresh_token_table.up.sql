CREATE TABLE refresh_token (
                                user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
                                token_hash TEXT PRIMARY KEY,
                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);