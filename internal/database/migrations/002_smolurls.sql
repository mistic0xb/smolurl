-- Write your migrate up statements here
CREATE TABLE smolurls(
    id BIGINT PRIMARY KEY NOT NULL,
    original_url TEXT NOT NULL,
    smol_url TEXT NOT NULL,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_smolurls_smolurl
ON smolurls (smol_url);

CREATE INDEX idx_smolurls_expirationTime
ON smolurls (expires_at);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
