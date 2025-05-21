-- Add new features: spoiler, downvote counter, reactions, timestamps
-- May 18, 2025

-- Add spoiler flag to comments table
ALTER TABLE comments
ADD COLUMN spoiler BOOLEAN NOT NULL DEFAULT false;

-- Add reactions table
CREATE TABLE IF NOT EXISTS reactions (
  commentHex               TEXT          NOT NULL                           ,
  commenterHex             TEXT          NOT NULL                           ,
  type                     TEXT          NOT NULL                           ,
  reactionDate             TIMESTAMP     NOT NULL
);

-- Create unique index to ensure each commenter has only one reaction type per comment
CREATE UNIQUE INDEX reactionsUniqueIndex ON reactions(commentHex, commenterHex, type);

-- Create function to count reactions by type
CREATE OR REPLACE FUNCTION get_reaction_count(cHex TEXT, rType TEXT) 
RETURNS INTEGER AS $$
BEGIN
  RETURN (SELECT COUNT(*) FROM reactions WHERE commentHex = cHex AND type = rType);
END;
$$ LANGUAGE plpgsql;
