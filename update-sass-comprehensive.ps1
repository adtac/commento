# Comprehensive Sass Modernization Script

# Define the working directory
$sassDir = "c:\Users\Pishi\Desktop\Commento\frontend\sass"

# Function to update map-get to map.get in a file
function Update-MapGetSyntax {
    param (
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    Write-Host "Updating map-get syntax in $FilePath..."
    
    # Read file content
    $content = Get-Content -Path $FilePath -Raw
    
    # Check if the file contains map-get
    if ($content -match "map-get\(") {
        # Replace map-get with map.get
        $updatedContent = $content -replace "map-get\(", "map.get("
        
        # Add @use "sass:map" if it doesn't exist already
        if (-not ($updatedContent -match '@use "sass:map"')) {
            $updatedContent = '@use "sass:map";' + "`r`n`r`n" + $updatedContent
        }
        
        # Write updated content back to file
        Set-Content -Path $FilePath -Value $updatedContent
        
        Write-Host "✓ Updated map-get syntax in $FilePath"
        return $true
    } else {
        Write-Host "- No map-get functions found in $FilePath"
        return $false
    }
}

# Function to update specific color file imports
function Update-MainColorImports {
    $colorsFile = Join-Path $sassDir "colors-main.scss"
    $commonFile = Join-Path $sassDir "common-main.scss"
    
    # Create _index.scss file if it doesn't exist
    $indexFile = Join-Path $sassDir "_index.scss"
    if (-not (Test-Path $indexFile)) {
        $indexContent = @"
@forward "colors-main";
@forward "common-main";
@forward "source-sans";
"@
        Set-Content -Path $indexFile -Value $indexContent
        Write-Host "✓ Created _index.scss"
    }
    
    # Update common-main.scss to use colors-main
    $commonContent = Get-Content -Path $commonFile -Raw
    if ($commonContent -match '@import "colors-main.scss";') {
        $updatedCommonContent = $commonContent -replace '@import "colors-main.scss";', '@use "colors-main" as colors;'
        Set-Content -Path $commonFile -Value $updatedCommonContent
        Write-Host "✓ Updated common-main.scss to use @use for colors-main"
    }
}

# Main execution
Write-Host "Starting Sass Modernization Process..."

# 1. First update colors-main.scss to use map.get
$colorsFile = Join-Path $sassDir "colors-main.scss"
Update-MapGetSyntax -FilePath $colorsFile

# 2. Update all other files with map-get calls
$allScssFiles = Get-ChildItem -Path $sassDir -Filter "*.scss"
foreach ($file in $allScssFiles) {
    if ($file.Name -ne "colors-main.scss") {
        Update-MapGetSyntax -FilePath $file.FullName
    }
}

# 3. Update imports for main files
Update-MainColorImports

Write-Host "Sass Modernization Completed!"
