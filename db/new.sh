#!/usr/bin/env bash

timestamp=$(date +%Y%m%d%H%M%S)

printf "rules:\n"
printf "  * use hyphens to separate words (not spaces, not underscores)\n"
printf "  * keep it as short as possible (add comments inside the file)\n"
printf "  * try to keep each migration idempotent (roughly, the order of application shouldn't matter)\n"
printf "\n"
printf "good example: 20180416164303-init-schema.sql\n\n"
printf "filename: %s-" "${timestamp}"
read filename

filename="${timestamp}-${filename}"
if [[ ! $filename =~ .sql$ ]]; then
  filename="${filename}.sql"
fi

touch "${filename}"
echo "created ${filename}"
