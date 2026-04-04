-- Session hardening: add expiry and last-seen tracking

ALTER TABLE ownerSessions ADD COLUMN IF NOT EXISTS expiresAt TIMESTAMP;
ALTER TABLE ownerSessions ADD COLUMN IF NOT EXISTS lastSeenDate TIMESTAMP;

-- Set default expiry for existing sessions (30 days from now)
UPDATE ownerSessions SET expiresAt = NOW() + INTERVAL '30 days' WHERE expiresAt IS NULL;
UPDATE ownerSessions SET lastSeenDate = loginDate WHERE lastSeenDate IS NULL;

-- Make expiresAt non-nullable going forward
ALTER TABLE ownerSessions ALTER COLUMN expiresAt SET NOT NULL;
ALTER TABLE ownerSessions ALTER COLUMN expiresAt SET DEFAULT NOW() + INTERVAL '30 days';
