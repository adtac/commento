DROP TRIGGER IF EXISTS commentsDeleteTrigger ON comments;

DROP FUNCTION IF EXISTS commentsDeleteTriggerFunction();

ALTER TABLE comments
  ADD deleted BOOLEAN NOT NULL DEFAULT false;
