-- Make all login providers optional (but enabled by default)

ALTER TABLE domains
  ADD commentoProvider BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE domains
  ADD googleProvider BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE domains
  ADD twitterProvider BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE domains
  ADD githubProvider BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE domains
  ADD gitlabProvider BOOLEAN NOT NULL DEFAULT true;
