# PowerShell script to update Sass files
# Update map-get to map.get and add @use "sass:map" if needed

$sassDir = "c:\Users\Pishi\Desktop\Commento\frontend\sass"
$files = Get-ChildItem -Path $sassDir -Filter "*.scss"

foreach ($file in $files) {
    $filePath = $file.FullName
    Write-Host "Processing $filePath..."
    
    # Read file content
    $content = Get-Content -Path $filePath -Raw
    
    # Replace map-get with map.get
    $updatedContent = $content -replace "map-get\(", "map.get("
    
    # Add @use "sass:map" if the file now uses map.get but doesn't have the import
    if ($updatedContent -match "map\.get\(" -and -not ($updatedContent -match '@use "sass:map"')) {
        $updatedContent = '@use "sass:map";' + "`r`n`r`n" + $updatedContent
    }
    
    # Write updated content back to file
    Set-Content -Path $filePath -Value $updatedContent
    
    Write-Host "Updated $filePath"
}

Write-Host "All Sass files updated successfully!"
