-- Create the resetHexes table

ALTER TABLE ownerResetHexes RENAME TO resetHexes;

ALTER TABLE resetHexes RENAME ownerHex TO hex;

ALTER TABLE resetHexes
  ADD entity TEXT NOT NULL DEFAULT 'owner';
