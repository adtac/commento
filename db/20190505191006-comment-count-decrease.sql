-- This trigger is called every time a comment is deleted, so the comment count for the page where the comment belong is updated
CREATE OR REPLACE FUNCTION commentsDeleteTriggerFunction() RETURNS TRIGGER AS $trigger$
BEGIN
  UPDATE pages
  SET commentCount = commentCount - 1
  WHERE domain = old.domain AND path = old.path;

  DELETE FROM comments
  WHERE parentHex = old.commentHex;

  RETURN NEW;
END;
$trigger$ LANGUAGE plpgsql;
