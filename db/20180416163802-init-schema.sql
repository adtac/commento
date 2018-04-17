-- Initial PostgreSQL database schema

CREATE TABLE IF NOT EXISTS config (
  allowNewOwners           BOOLEAN       NOT NULL
);

INSERT INTO
config (allowNewOwners)
VALUES (true);

CREATE TABLE IF NOT EXISTS owners (
  ownerHex                 TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  email                    TEXT          NOT NULL  UNIQUE                   ,
  name                     TEXT          NOT NULL                           ,
  passwordHash             TEXT          NOT NULL                           ,
  confirmedEmail           TEXT          NOT NULL  DEFAULT false            ,
  joinDate                 TIMESTAMP     NOT NULL
);

CREATE TABLE IF NOT EXISTS ownerSessions (
  session                  TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  ownerHex                 TEXT          NOT NULL                           ,
  loginDate                TIMESTAMP     NOT NULL
);

CREATE TABLE IF NOT EXISTS ownerConfirmHexes (
  confirmHex               TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  ownerHex                 TEXT          NOT NULL                           ,
  sendDate                 TEXT          NOT NULL
);

CREATE TABLE IF NOT EXISTS ownerResetHexes (
  resetHex                 TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  ownerHex                 TEXT          NOT NULL                           ,
  sendDate                 TEXT          NOT NULL
);

CREATE TABLE IF NOT EXISTS domains (
  domain                   TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  ownerHex                 TEXT          NOT NULL                           ,
  name                     TEXT          NOT NULL                           ,
  creationDate             TIMESTAMP     NOT NULL                           ,
  state                    TEXT          NOT NULL  DEFAULT 'unfrozen'       ,
  importedComments         TEXT          NOT NULL  DEFAULT false            ,
  autoSpamFilter           BOOLEAN       NOT NULL  DEFAULT true             ,
  requireModeration        BOOLEAN       NOT NULL  DEFAULT false            ,
  requireIdentification    BOOLEAN       NOT NULL  DEFAULT true             ,
  viewsThisMonth           INTEGER       NOT NULL  DEFAULT 0
);

CREATE TABLE IF NOT EXISTS moderators (
  domain                   TEXT          NOT NULL                           ,
  email                    TEXT          NOT NULL                           ,
  addDate                  TIMESTAMP     NOT NULL                           ,
  PRIMARY KEY (domain, email)
);

CREATE TABLE IF NOT EXISTS commenters (
  commenterHex             TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  email                    TEXT          NOT NULL                           ,
  name                     TEXT          NOT NULL                           ,
  link                     TEXT          NOT NULL                           ,
  photo                    TEXT          NOT NULL                           ,
  provider                 TEXT          NOT NULL                           ,
  joinDate                 TIMESTAMP     NOT NULL                           ,
  state                    TEXT          NOT NULL  DEFAULT 'ok'
);

CREATE TABLE IF NOT EXISTS commenterSessions (
  session                  TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  commenterHex             TEXT          NOT NULL  DEFAULT 'none'           ,
  creationDate             TIMESTAMP     NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
  commentHex               TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  domain                   TEXT          NOT NULL                           ,
  path                     TEXT          NOT NULL                           ,
  commenterHex             TEXT          NOT NULL                           ,
  markdown                 TEXT          NOT NULL                           ,
  html                     TEXT          NOT NULL                           ,
  parentHex                TEXT          NOT NULL                           ,
  score                    INTEGER       NOT NULL  DEFAULT 0                ,
  state                    TEXT          NOT NULL  DEFAULT 'unapproved'     , -- not a BOOLEAN because I expect more states in the future
  creationDate             TIMESTAMP     NOT NULL
);

-- DELETEing a comment should recursively delete all children
CREATE OR REPLACE FUNCTION commentsDeleteTriggerFunction() RETURNS TRIGGER AS $trigger$
BEGIN
  DELETE FROM comments
  WHERE parentHex = old.commentHex;

  RETURN NULL;
END;
$trigger$ LANGUAGE plpgsql;

CREATE TRIGGER commentsDeleteTrigger AFTER DELETE ON comments
FOR EACH ROW EXECUTE PROCEDURE commentsDeleteTriggerFunction();

CREATE TABLE IF NOT EXISTS votes (
  commentHex               TEXT          NOT NULL                           ,
  commenterHex             TEXT          NOT NULL                           ,
  direction                INTEGER       NOT NULL                           ,
  voteDate                 TIMESTAMP     NOT NULL
);

CREATE UNIQUE INDEX votesUniqueIndex ON votes(commentHex, commenterHex);

CREATE OR REPLACE FUNCTION votesInsertTriggerFunction() RETURNS TRIGGER AS $trigger$
BEGIN
  UPDATE comments
  SET score = score + new.direction
  WHERE commentHex = new.commentHex;

  RETURN NEW;
END;
$trigger$ LANGUAGE plpgsql;

CREATE TRIGGER votesInsertTrigger AFTER INSERT ON votes
FOR EACH ROW EXECUTE PROCEDURE votesInsertTriggerFunction();

CREATE OR REPLACE FUNCTION votesUpdateTriggerFunction() RETURNS TRIGGER AS $trigger$
BEGIN
  UPDATE comments
  SET score = score - old.direction + new.direction
  WHERE commentHex = old.commentHex;

  RETURN NEW;
END;
$trigger$ LANGUAGE plpgsql;

CREATE TRIGGER votesUpdateTrigger AFTER UPDATE ON votes
FOR EACH ROW EXECUTE PROCEDURE votesUpdateTriggerFunction();

CREATE TABLE IF NOT EXISTS views (
  domain                   TEXT          NOT NULL                           ,
  commenterHex             TEXT          NOT NULL                           ,
  viewDate                 TIMESTAMP     NOT NULL
);

CREATE INDEX IF NOT EXISTS domainIndex ON views(domain);

CREATE OR REPLACE FUNCTION viewsInsertTriggerFunction() RETURNS TRIGGER AS $trigger$
BEGIN
  UPDATE domains
  SET viewsThisMonth = viewsThisMonth + 1
  WHERE domain = new.domain;

  RETURN NULL;
END;
$trigger$ LANGUAGE plpgsql;

CREATE TRIGGER viewsInsertTrigger AFTER INSERT ON views
FOR EACH ROW EXECUTE PROCEDURE viewsInsertTriggerFunction();
