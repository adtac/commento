-- add export feature

CREATE TABLE IF NOT EXISTS exports (
  exportHex TEXT NOT NULL UNIQUE PRIMARY KEY,
  binData BYTEA NOT NULL,
  domain TEXT NOT NULL,
  creationDate TIMESTAMP NOT NULL
);
