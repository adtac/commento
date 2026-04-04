-- SaaS multi-user: plans, subscriptions, and usage tracking

CREATE TABLE IF NOT EXISTS plans (
  planId                   TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  name                     TEXT          NOT NULL                           ,
  maxDomains               INTEGER       NOT NULL                           ,
  maxCommentsPerMonth      INTEGER       NOT NULL                           ,
  maxPageViewsPerMonth     INTEGER       NOT NULL                           ,
  isActive                 BOOLEAN       NOT NULL  DEFAULT true
);

INSERT INTO plans (planId, name, maxDomains, maxCommentsPerMonth, maxPageViewsPerMonth)
VALUES
  ('free',     'Free',     1,   500,    10000),
  ('starter',  'Starter',  5,   5000,   100000),
  ('business', 'Business', 50,  50000,  1000000);

CREATE TABLE IF NOT EXISTS subscriptions (
  ownerHex                 TEXT          NOT NULL                           ,
  planId                   TEXT          NOT NULL  DEFAULT 'free'           ,
  status                   TEXT          NOT NULL  DEFAULT 'active'         ,
  currentPeriodStart       TIMESTAMP     NOT NULL                           ,
  currentPeriodEnd         TIMESTAMP     NOT NULL                           ,
  createdDate              TIMESTAMP     NOT NULL                           ,
  PRIMARY KEY (ownerHex)
);

CREATE TABLE IF NOT EXISTS usage (
  ownerHex                 TEXT          NOT NULL                           ,
  month                    TEXT          NOT NULL                           ,
  commentCount             INTEGER       NOT NULL  DEFAULT 0               ,
  pageViewCount            INTEGER       NOT NULL  DEFAULT 0               ,
  PRIMARY KEY (ownerHex, month)
);

ALTER TABLE owners ADD COLUMN IF NOT EXISTS isSuperAdmin BOOLEAN NOT NULL DEFAULT false;

-- Seed: give every existing owner a free subscription
INSERT INTO subscriptions (ownerHex, planId, status, currentPeriodStart, currentPeriodEnd, createdDate)
SELECT ownerHex, 'free', 'active', NOW(), NOW() + INTERVAL '100 years', NOW()
FROM owners
WHERE ownerHex NOT IN (SELECT ownerHex FROM subscriptions);
