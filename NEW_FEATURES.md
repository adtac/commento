# New Features in Commento

## Recent Feature Additions (May 2025)

### Mark Comment as Spoiler
- Comment owners and moderators can mark comments as spoilers
- Spoiler comments are blurred and revealed on hover
- Great for hiding sensitive information or story spoilers

### Downvote Counter
- Displays the number of downvotes a comment has received
- Shows alongside the existing score counter

### Timestamp Links
- When a user writes a timestamp (e.g., 1:30, 5:45), it's automatically converted to a clickable link
- Clicking the timestamp link triggers a custom event that can be captured to navigate to specific timestamps in media

### Reaction System
- Four reactions available: Funny, Interesting, Upsetting, and Sad
- Multiple users can leave reactions on the same comment
- Reaction counts are displayed

### Sort by Reaction Score
- Sort comments by reactions (Funny, Interesting, Upsetting, Sad)
- Click the new sort options in the comments header

### Filter by Username
- Filter comments to show only those from a specific user
- Great for finding all comments from a particular contributor

## How to Use

### For Users
- Click the reaction buttons to add your reaction to a comment
- Use the sort options to see comments sorted by different reaction types
- Use the username filter to find all comments from a specific user
- Timestamps like 2:30 in comments will automatically become clickable

### For Developers
- Capture timestamp clicks by listening for the 'commentoTimestamp' event
- Example:
  ```javascript
  document.addEventListener('commentoTimestamp', function(e) {
    console.log('Timestamp clicked:', e.detail.timestamp);
    // Navigate video player to this timestamp, for example
  });
  ```

### For Website Owners
- All features work out of the box with no configuration needed
- The spoiler feature helps manage sensitive content
- Reaction sorting provides additional engagement metrics
