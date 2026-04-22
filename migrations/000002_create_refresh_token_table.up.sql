CREATE TABLE refresh_token (
                                user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
                                token_hash CHAR(64) PRIMARY KEY,
                                expires_at TIMESTAMP NOT NULL,
                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                revoked_at TIMESTAMP NULL
);