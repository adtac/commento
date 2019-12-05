-- Default sort policy for each domain

CREATE TYPE sortPolicy AS ENUM (
  'score-desc',
  'creationdate-desc',
  'creationdate-asc'
);

ALTER TABLE domains
  ADD defaultSortPolicy sortPolicy NOT NULL DEFAULT 'score-desc';
