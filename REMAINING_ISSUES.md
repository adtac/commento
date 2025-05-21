# Remaining Issues in Commento Dependency Updates

## ESLint Errors

During the build, we encountered 120 ESLint errors. These errors are likely due to the update from ESLint 5.x to 8.x, which has many new or stricter rules.

### Possible Solutions:
1. Create an updated `.eslintrc` file with more permissive rules
2. Fix the actual code issues one by one
3. Use `eslint --fix` to automatically resolve formatting issues

## Sass Deprecation Warnings

The build shows numerous deprecation warnings related to Sass:

1. Legacy JS API deprecation:
   ```
   Deprecation Warning [legacy-js-api]: The legacy JS API is deprecated and will be removed in Dart Sass 2.0.0.
   More info: https://sass-lang.com/d/legacy-js-api
   ```

2. @import rule deprecation:
   ```
   Deprecation Warning [import]: Sass @import rules are deprecated and will be removed in Dart Sass 3.0.0.
   More info and automated migrator: https://sass-lang.com/d/import
   ```

3. Global built-in functions:
   ```
   Deprecation Warning [global-builtin]: Global built-in functions are deprecated and will be removed in Dart Sass 3.0.0.
   Use map.get instead.
   ```

### Possible Solutions:
1. Update Sass files to use the `@use` rule instead of `@import`
2. Update global built-in function calls to use the namespace syntax (`map.get` instead of `map-get`)
3. Use the automated migrator tool mentioned in the warnings

## Next Steps

1. Fix the highlightjs version to prevent breaking changes (completed)
2. Create an .eslintrc.json file with updated rules
3. Address the Sass deprecation warnings
4. Test the build with Docker to ensure the updated Dockerfile works
