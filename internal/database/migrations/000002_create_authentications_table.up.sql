CREATE TABLE IF NOT EXISTS authentications (
    token TEXT NOT NULL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_authentications_user_id ON authentications(user_id);
