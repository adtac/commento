# Frontend Dependency Updates

## Changes Made

The following updates were made to modernize the frontend dependencies:

1. **Updated `gulp-sass` from 4.0.1 to 5.1.0**
   - The newer version uses Dart Sass instead of Node Sass
   - This eliminates the Python 2.7 dependency that was causing build issues

2. **Updated `sass` from 1.5.1 to 1.69.5**
   - This is the latest version of the Dart Sass implementation
   - Much more maintainable and actively developed than node-sass

3. **Modified the gulpfile.js**
   - Updated to work with the new gulp-sass API which requires explicitly providing the sass compiler

4. **Other dependency updates:**
   - jquery: 3.2.1 → 3.7.1
   - eslint: 5.10.0 → 8.56.0
   - highlightjs: 9.10.0 → 11.8.0
   - html-minifier: 3.5.7 → 4.0.0
   - uglify-js: 3.4.1 → 3.17.4
   - vue: 2.5.16 → 2.7.15

## Benefits

- No longer depends on the deprecated node-sass package
- No longer requires Python 2.7 for building
- Uses modern, actively maintained dependencies
- Better compatibility with newer Node.js versions
- Improved security by updating packages with known vulnerabilities

## How to Build

Run the following commands to build with the updated dependencies:

```bash
cd frontend
npm install
npm run prod
```

This should build successfully without the previous errors related to node-sass and Python 2.7.
