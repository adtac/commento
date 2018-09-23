ALTER TABLE pages
  ADD commentCount INTEGER NOT NULL DEFAULT 0;

CREATE OR REPLACE FUNCTION commentsInsertTriggerFunction() RETURNS TRIGGER AS $trigger$
BEGIN
  UPDATE pages
  SET commentCount = commentCount + 1
  WHERE domain = new.domain AND path = new.path;

  RETURN NEW;
END;
$trigger$ LANGUAGE plpgsql;

CREATE TRIGGER commentsInsertTrigger AFTER INSERT ON comments
FOR EACH ROW EXECUTE PROCEDURE commentsInsertTriggerFunction();
