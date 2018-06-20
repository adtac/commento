ALTER TABLE ownerSessions
RENAME COLUMN session TO ownerToken;

ALTER TABLE commenterSessions
RENAME COLUMN session TO commenterToken
