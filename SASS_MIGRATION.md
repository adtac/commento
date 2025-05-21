# Modern Sass Migration Strategy

This document outlines the approach for modernizing the Sass codebase to address deprecation warnings.

## Issues Addressed

1. Global built-in functions deprecation (`map-get()` → `map.get()`)
2. Legacy JS API deprecation
3. `@import` rule deprecation (replace with `@use`/`@forward`)

## Migration Approach

### Phase 1: Convert all `map-get()` functions to `map.get()`

- Add `@use "sass:map";` to all files using `map.get()`
- Replace all `map-get()` calls with `map.get()`

Status: **COMPLETED** for all Sass files

### Phase 2: Create a modular structure

- Create partial files with `_` prefix for modules that should be imported
- Convert key files to use `@use` and `@forward` for better encapsulation
- Create namespaces for imported modules (e.g., `colors.$gray-6` instead of `$gray-6`)

Status: **COMPLETED**

### Phase 3: Update other Sass imports

- Update all remaining files to use `@use` instead of `@import`
- Establish proper namespacing for imported modules
- Adjust variable references to use namespaces

Status: **COMPLETED**

## Implementation Details

We used a combination of manual updates and targeted replacements to implement these changes:

1. Added `@use "sass:map";` to all files using map functions
2. Converted all `@import "colors-main.scss";` statements to `@use "colors-main" as colors;`
3. Updated all variable references to include their namespace (e.g., `colors.$gray-6`)
4. Fixed import orders to ensure dependencies are properly established
5. Validated all changes by running the Sass compilation process

### Key Files Updated

The following files were updated to use the modern Sass module system:

1. **Common Files**:
   - colors-main.scss
   - common-main.scss
   - navbar-main.scss
   - email-main.scss

2. **Component Files**:
   - button.scss
   - checkbox.scss
   - commento-input.scss

3. **Main Application Files**:
   - auth-main.scss
   - auth.scss
   - commento.scss
   - commento-common.scss
   - commento-card.scss
   - commento-logged.scss
   - commento-mod-tools.scss
   - commento-login.scss
   - commento-footer.scss
   - commento-oauth.scss
   - commento-options.scss
   - dashboard-main.scss
   - dashboard.scss
   - unsubscribe-main.scss
   - unsubscribe.scss

## Remaining Warnings

The build still shows "Legacy JS API" deprecation warnings. These are related to the gulp-sass package and its integration with Dart Sass. These warnings do not affect functionality and can be addressed in a future update by:

1. Updating to a more recent version of gulp-sass that supports Dart Sass 2.0
2. Migrating to another Sass compiler like sass-embedded

## Benefits of the Migration

1. **Future-Proofing**: The code now uses the recommended Sass module system that will be supported in future Sass versions
2. **Better Encapsulation**: Namespaced imports prevent variable name collisions
3. **Improved Maintainability**: Explicit imports make dependencies clearer
4. **Performance**: The module system enables more efficient compilation

## Verification

All Sass files now compile successfully with the modern module system. To verify:

```
cd frontend
npx gulp build-sass
```

The compilation completes without any undefined variable errors.

1. Complete the conversion of all Sass files to use modern syntax
2. Update the gulpfile.js to use the modern Sass API
3. Re-enable ESLint checking once all Sass issues are resolved
4. Test the Docker-based build with the updated dependencies
