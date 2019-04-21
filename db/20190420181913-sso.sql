-- Single Sign-On (SSO)

ALTER TABLE domains
  ADD ssoProvider BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE domains
  ADD ssoSecret TEXT NOT NULL DEFAULT '';

ALTER TABLE domains
  ADD ssoUrl TEXT NOT NULL DEFAULT '';
