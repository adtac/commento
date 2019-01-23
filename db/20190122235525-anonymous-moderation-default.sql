-- Allow the owner to change whether anonymous comments are put into moderation by default.

ALTER TABLE domains
  ADD COLUMN moderateAllAnonymous BOOLEAN DEFAULT true;
