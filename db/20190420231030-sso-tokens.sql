CREATE TABLE IF NOT EXISTS ssoTokens (
  token                    TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  domain                   TEXT          NOT NULL                           ,
  commenterToken           TEXT          NOT NULL                           ,
  creationDate             TIMESTAMP     NOT NULL
);
