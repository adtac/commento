-- Introduces page attributes

CREATE TABLE IF NOT EXISTS pages (
  domain                   TEXT          NOT NULL                           ,
  path                     TEXT          NOT NULL                           ,
  isLocked                 BOOLEAN       NOT NULL  DEFAULT false
);

CREATE UNIQUE INDEX pagesUniqueIndex ON pages(domain, path);
