-- Build the comments count column

UPDATE pages
SET commentCount = subquery.commentCount
FROM (
  SELECT COUNT(commentHex) as commentCount
  FROM comments
  WHERE state = 'approved'
  GROUP BY (domain, path)
) as subquery;
