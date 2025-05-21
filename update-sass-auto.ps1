# Enhanced Sass Modernization Script
# This script updates Sass files to use the modern syntax:
# 1. Converts map-get to map.get
# 2. Adds @use "sass:map" import where needed
# 3. Updates @import to @use with namespaces
# 4. Updates color variables to use the namespace

# Define the working directory
$sassDir = "c:\Users\Pishi\Desktop\Commento\frontend\sass"

# Function to replace all color variable references with namespaced versions
function Update-ColorVariableReferences {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Content
    )
    
    # Get all simple color variable references and replace them with colors.$ prefix
    $colorVars = @(
        'white', 'black', 
        'gray-0', 'gray-1', 'gray-2', 'gray-3', 'gray-4', 'gray-5', 'gray-6', 'gray-7', 'gray-8', 'gray-9',
        'red-0', 'red-1', 'red-2', 'red-3', 'red-4', 'red-5', 'red-6', 'red-7', 'red-8', 'red-9',
        'pink-0', 'pink-1', 'pink-2', 'pink-3', 'pink-4', 'pink-5', 'pink-6', 'pink-7', 'pink-8', 'pink-9',
        'grape-0', 'grape-1', 'grape-2', 'grape-3', 'grape-4', 'grape-5', 'grape-6', 'grape-7', 'grape-8', 'grape-9',
        'violet-0', 'violet-1', 'violet-2', 'violet-3', 'violet-4', 'violet-5', 'violet-6', 'violet-7', 'violet-8', 'violet-9',
        'indigo-0', 'indigo-1', 'indigo-2', 'indigo-3', 'indigo-4', 'indigo-5', 'indigo-6', 'indigo-7', 'indigo-8', 'indigo-9',
        'blue-0', 'blue-1', 'blue-2', 'blue-3', 'blue-4', 'blue-5', 'blue-6', 'blue-7', 'blue-8', 'blue-9',
        'cyan-0', 'cyan-1', 'cyan-2', 'cyan-3', 'cyan-4', 'cyan-5', 'cyan-6', 'cyan-7', 'cyan-8', 'cyan-9',
        'teal-0', 'teal-1', 'teal-2', 'teal-3', 'teal-4', 'teal-5', 'teal-6', 'teal-7', 'teal-8', 'teal-9',
        'green-0', 'green-1', 'green-2', 'green-3', 'green-4', 'green-5', 'green-6', 'green-7', 'green-8', 'green-9',
        'lime-0', 'lime-1', 'lime-2', 'lime-3', 'lime-4', 'lime-5', 'lime-6', 'lime-7', 'lime-8', 'lime-9',
        'yellow-0', 'yellow-1', 'yellow-2', 'yellow-3', 'yellow-4', 'yellow-5', 'yellow-6', 'yellow-7', 'yellow-8', 'yellow-9',
        'orange-0', 'orange-1', 'orange-2', 'orange-3', 'orange-4', 'orange-5', 'orange-6', 'orange-7', 'orange-8', 'orange-9'
    )
function Update-MapGetSyntax {
    param (
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    Write-Host "Updating map-get syntax in $FilePath..."
    $fileName = Split-Path $FilePath -Leaf
    
    # Read file content
    $content = Get-Content -Path $FilePath -Raw
    
    # Check if the file contains map-get
    if ($content -match "map-get\(") {
        # Replace map-get with map.get
        $updatedContent = $content -replace "map-get\(", "map.get("
        
            foreach ($var in $colorVars) {
        # Replace variables that aren't already namespaced
        $Content = $Content -replace "([^a-zA-Z0-9\.])\$${var}([^a-zA-Z0-9])", "`$1colors.`$${var}`$2"
    }
    
    return $Content
}

# Function to update @import to @use in a file
function Update-ImportToUse {
    param (
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    Write-Host "Processing $FilePath..."
    
    # Read file content
    $content = Get-Content -Path $FilePath -Raw
    $originalContent = $content
    
    # Keep track if we need to add the colors import
    $needsColorsImport = $false
    $needsMapImport = $false
    
    # Replace @import "colors-main.scss"; with @use "colors-main" as colors;
    if ($content -match '@import "colors-main.scss";') {
        $content = $content -replace '@import "colors-main.scss";', '@use "colors-main" as colors;'
        $needsColorsImport = $false
    } else {
        # Check if we have color variables without the import
        $colorVarPattern = "\$(?:gray|red|pink|grape|violet|indigo|blue|cyan|teal|green|lime|yellow|orange)-[0-9]"
        $whiteBlackPattern = "\$(white|black)"
        
        if (($content -match $colorVarPattern) -or ($content -match $whiteBlackPattern)) {
            $needsColorsImport = $true
        }
    }
    
    # Replace @import "common-main.scss"; with @use "common-main" as common;
    if ($content -match '@import "common-main.scss";') {
        $content = $content -replace '@import "common-main.scss";', '@use "common-main" as common;'
    }
    
    # Replace @import "checkbox.scss"; with @use "checkbox";
    if ($content -match '@import "checkbox.scss";') {
        $content = $content -replace '@import "checkbox.scss";', '@use "checkbox";'
    }
    
    # Replace @import "button.scss"; with @use "button";
    if ($content -match '@import "button.scss";') {
        $content = $content -replace '@import "button.scss";', '@use "button";'
    }
      # Check for map-get usage
    if ($content -match "map-get\(") {
        $content = $content -replace "map-get\(", "map.get("
        $needsMapImport = $true
    }
    
    # Add imports if needed
    $importString = ""
    if ($needsMapImport) {
        $importString += '@use "sass:map";' + "`r`n"
    }
    
    if ($needsColorsImport) {
        $importString += '@use "colors-main" as colors;' + "`r`n"
    }
    
    # Add imports to the beginning of the file
    if ($importString -ne "") {
        $content = $importString + $content
    }
    
    # Update all color variables to use the colors namespace
    if ($content -match "colors\.") {
        # File already has some namespaced colors, update all of them
        $content = Update-ColorVariableReferences -Content $content
    }
    
    # Write updated content back to file if changes were made
    if ($content -ne $originalContent) {
        Set-Content -Path $FilePath -Value $content
        Write-Host "✓ Updated $FilePath"
        return $true
    } else {
        Write-Host "- No changes needed for $FilePath"
        return $false
    }
}

# Main execution
Write-Host "Starting Comprehensive Sass Modernization..."

# Update all scss files
$allScssFiles = Get-ChildItem -Path $sassDir -Filter "*.scss" | Where-Object { -not $_.Name.StartsWith("_") }
$updatedFiles = 0

foreach ($file in $allScssFiles) {
    $success = Update-ImportToUse -FilePath $file.FullName
    if ($success) {
        $updatedFiles++
    }
}

# Create the _index.scss file if it doesn't exist
$indexFile = Join-Path $sassDir "_index.scss"
if (-not (Test-Path $indexFile)) {
    $indexContent = @"
@forward "colors-main";
@forward "common-main";
@forward "source-sans";
"@
    Set-Content -Path $indexFile -Value $indexContent
    Write-Host "Created _index.scss"
}

Write-Host "Sass Modernization Completed! Updated $updatedFiles files."
}

# Function to process all files
function Process-AllFiles {
    $allScssFiles = Get-ChildItem -Path $sassDir -Filter "*.scss"
    
    # First update colors-main.scss
    $colorsMainPath = Join-Path $sassDir "colors-main.scss"
    Update-MapGetSyntax -FilePath $colorsMainPath
    
    # Update all other files
    foreach ($file in $allScssFiles) {
        if ($file.Name -ne "colors-main.scss") {
            $updated = Update-MapGetSyntax -FilePath $file.FullName
            $updated = Update-ColorImports -FilePath $file.FullName
        }
    }
}

# Main execution
Write-Host "Starting Enhanced Sass Modernization Process..."
Process-AllFiles
Write-Host "Sass Modernization Completed!"
