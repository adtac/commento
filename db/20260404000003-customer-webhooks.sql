-- Customer webhook system

CREATE TABLE IF NOT EXISTS webhooks (
  webhookHex               TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  domain                   TEXT          NOT NULL                           ,
  url                      TEXT          NOT NULL                           ,
  events                   TEXT          NOT NULL                           ,
  secret                   TEXT          NOT NULL                           ,
  createdDate              TIMESTAMP     NOT NULL
);

CREATE INDEX IF NOT EXISTS webhooksDomainIndex ON webhooks(domain);
