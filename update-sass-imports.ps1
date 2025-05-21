# This script converts the basic Sass files to use @forward instead of @import
# This allows other files to continue importing them with @use instead of @import

# First, let's update colors-main.scss to use @forward for exporting variables
$colorsMainPath = "c:\Users\Pishi\Desktop\Commento\frontend\sass\colors-main.scss"
$colorsMainContent = Get-Content -Path $colorsMainPath -Raw

# Remove the existing @use directive since we'll add it later
$colorsMainContent = $colorsMainContent -replace '@use "sass:map";\s+', ''

# Add both @use and @forward directives
$colorsMainContent = '@use "sass:map";' + "`r`n" + '@forward "sass:map";' + "`r`n`r`n" + $colorsMainContent

# Write the updated content back
Set-Content -Path $colorsMainPath -Value $colorsMainContent
Write-Host "Updated colors-main.scss"

# Now let's update common-main.scss which imports colors-main.scss
$commonMainPath = "c:\Users\Pishi\Desktop\Commento\frontend\sass\common-main.scss"
$commonMainContent = Get-Content -Path $commonMainPath -Raw

# Replace @import with @use
$updatedCommonMain = $commonMainContent -replace '@import "colors-main.scss";', '@use "colors-main" as colors;'
$updatedCommonMain = $updatedCommonMain -replace '@import "source-sans.scss";', '@use "source-sans";'

# Write the updated content back
Set-Content -Path $commonMainPath -Value $updatedCommonMain
Write-Host "Updated common-main.scss"

# Create a forward module for all commonly used files
$forwardModulePath = "c:\Users\Pishi\Desktop\Commento\frontend\sass\_index.scss"
$forwardModuleContent = @"
@forward "colors-main";
@forward "common-main";
@forward "source-sans";
"@

Set-Content -Path $forwardModulePath -Value $forwardModuleContent
Write-Host "Created _index.scss forward module"
