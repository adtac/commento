#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Function to update map-get to map.get
function updateMapGet(content) {
  return content.replace(/map-get\(/g, 'map.get(');
}

// Function to add @use 'sass:map'; at the beginning of the file if it contains map.get
function addMapImport(content) {
  if (content.includes('map.get(') && !content.includes('@use "sass:map"')) {
    return '@use "sass:map";\n\n' + content;
  }
  return content;
}

// Process a single scss file
function processFile(filePath) {
  console.log(`Processing ${filePath}...`);
  let content = fs.readFileSync(filePath, 'utf8');
  
  // Update map-get to map.get
  content = updateMapGet(content);
  
  // Add map import if needed
  content = addMapImport(content);
  
  // Write back the updated content
  fs.writeFileSync(filePath, content, 'utf8');
  console.log(`Updated ${filePath}`);
}

// Get all scss files
const sassDir = path.join(__dirname, 'frontend', 'sass');
const scssFiles = fs.readdirSync(sassDir)
  .filter(file => file.endsWith('.scss'))
  .map(file => path.join(sassDir, file));

// Process each file
scssFiles.forEach(processFile);

console.log('All scss files have been updated with map.get syntax!');
