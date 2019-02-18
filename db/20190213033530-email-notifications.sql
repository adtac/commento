-- Email notifications
-- There are two kinds of email notifications: those sent to domain moderators
-- and those sent to commenters. Domain owners can choose to subscribe their
-- moderators to all comments, those pending moderation, or no emails. Each
-- moderator can independently opt out of these emails, of course. Commenters,
-- on the other, can choose to opt into reply notifications by email.

-- TODO: daily and weekly digests instead of just batched real-time emails?

-- TODO: more granular options to unsubscribe from emails for particular
-- domains can be provided - add unsubscribedReplyDomains []TEXT and
-- unsubscribedModeratorDomains []TEXT to emails table?

-- Each address has a cooldown period so that emails aren't sent within 10
-- minutes of each other. Why is this a separate table instead of another
-- column on commenters/owners? Because there may be some mods that haven't
-- logged in to create a row in the commenter table.
CREATE TABLE IF NOT EXISTS emails (
  email TEXT NOT NULL UNIQUE PRIMARY KEY,
  unsubscribeSecretHex TEXT NOT NULL UNIQUE,
  lastEmailNotificationDate TIMESTAMP NOT NULL,
  pendingEmails INTEGER NOT NULL DEFAULT 0,
  sendReplyNotifications BOOLEAN NOT NULL DEFAULT false,
  sendModeratorNotifications BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS unsubscribeSecretHexIndex ON emails(unsubscribeSecretHex);

-- Which comments should be sent?
-- Possible values: all, pending-moderation, none
-- Default to pending-moderation because this is critical. If the user forgets
-- to moderate, some comments will never see the light of day.
ALTER TABLE domains
  ADD COLUMN emailNotificationPolicy TEXT DEFAULT 'pending-moderation';

-- Each page now needs to store the title of the page.
ALTER TABLE pages
  ADD COLUMN title TEXT DEFAULT '';
