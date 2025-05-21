# Update single Sass file to use modern syntax
# This script updates a single Sass file to use modern Sass syntax
# including @use imports and namespaced color variables.

param(
    [Parameter(Mandatory=$true)]
    [string]$FilePath
)

# Update @import to @use with appropriate namespace
function Update-Imports {
    param (
        [string]$Content
    )
    
    # Update imports
    $Content = $Content -replace '@import "colors-main.scss";', '@use "sass:map";
@use "colors-main" as colors;'
    $Content = $Content -replace '@import "common-main.scss";', '@use "common-main" as common;'
    $Content = $Content -replace '@import "checkbox.scss";', '@use "checkbox";'
    $Content = $Content -replace '@import "button.scss";', '@use "button";'
    $Content = $Content -replace '@import "source-sans.scss";', '@use "source-sans";'
    
    return $Content
}

# Update color variable references to use namespace
function Update-ColorVariables {
    param (
        [string]$Content
    )
    
    # Basic colors
    $Content = $Content -replace '\$white', 'colors.$white'
    $Content = $Content -replace '\$black', 'colors.$black'
    
    # Numbered color variables (gray-0 through gray-9, etc.)
    foreach ($color in @("gray", "red", "pink", "grape", "violet", "indigo", "blue", "cyan", "teal", "green", "lime", "yellow", "orange")) {
        for ($i = 0; $i -le 9; $i++) {
            $Content = $Content -replace "\`$$color-$i", "colors.`$$color-$i"
        }
    }
    
    return $Content
}

# Convert map-get to map.get
function Update-MapFunctions {
    param (
        [string]$Content
    )
    
    return $Content -replace "map-get\(", "map.get("
}

# Main execution
if (Test-Path $FilePath) {
    Write-Host "Updating $FilePath..."
    
    # Read the file content
    $content = Get-Content -Path $FilePath -Raw
    
    # Update imports, variables and functions
    $content = Update-Imports -Content $content
    $content = Update-ColorVariables -Content $content
    $content = Update-MapFunctions -Content $content
    
    # Write the updated content back to the file
    Set-Content -Path $FilePath -Value $content
    
    Write-Host "File updated successfully!"
} else {
    Write-Host "File not found: $FilePath"
}
