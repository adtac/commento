# Building Commento with New Features

This document explains how to build and run Commento with the new features that have been added.

## New Features Added

1. **Mark comment as spoiler** - Blur out comment content until hovered
2. **Counter to downvotes** - Display the count of downvotes separately
3. **Timestamp hyperlinks** - Convert timestamps to clickable links
4. **Reactions** - Add Funny/Interesting/Upsetting/Sad reactions to comments
5. **Sort by reaction score** - Sort comments by reaction counts
6. **Filter by poster username** - Filter comments to show only those from a specific user

## Files Modified/Added

### Backend (API)
- `api/comment.go` - Updated comment struct with new fields
- `api/comment_reaction.go` - New file for handling reactions
- `api/comment_list.go` - Updated to include new fields when retrieving comments
- `api/router_api.go` - Added new API endpoints

### Database
- `db/20250518000000-new-features.sql` - New migration for added features

### Frontend
- `frontend/js/commento.js` - Updated with new UI elements and functionality
- `frontend/js/commento-reactions.js` - New file for reaction functionality
- `frontend/js/commento-spoiler.js` - New file for spoiler functionality
- `frontend/js/commento-timestamps.js` - New file for timestamp hyperlinks
- `frontend/js/commento-filter.js` - New file for username filtering
- `frontend/sass/commento-features.scss` - New CSS styles for the features
- `frontend/gulpfile.js` - Updated to include new JS files in the build

## Building the Project

### Prerequisites
- Docker and Docker Compose (recommended)
- Node.js v16+ and npm (alternative)
- Python 2.7 (for node-sass)
- Go 1.12+ (if building the backend manually)

### Using Docker (Recommended)
```bash
# In the project root directory
docker-compose up --build
```

### Manual Build (Alternative)
If Docker is not available, you can build the project manually:

1. Build the frontend:
```bash
cd frontend
npm install
npm run prod  # or 'npm run dev' for development build
```

2. Build the backend:
```bash
cd api
go build
```

3. Run the database migrations:
```bash
cd db
./run-migrations.sh
```

4. Run the application:
```bash
./commento
```

## Troubleshooting

### Frontend Build Issues
If you encounter issues with the frontend build:

1. Make sure you're using Node.js 16.x or newer
2. Check the `DEPENDENCY_UPDATES.md` file for details on recent dependency updates
3. If you have previously installed node-sass, you might need to clear your npm cache:
   ```bash
   npm cache clean --force
   ```
3. Try using a Node.js version that's compatible with the node-sass version (Node.js v14 or v16)

### Database Migration Issues
If you have issues with the new database migration:

1. Ensure PostgreSQL is running
2. Run the migrations manually:
```sql
psql -U postgres -d commento -f db/20250518000000-new-features.sql
```

## Using the New Features

See the NEW_FEATURES.md file for detailed information on how to use the new features.
