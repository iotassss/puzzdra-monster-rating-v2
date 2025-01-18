#!/bin/bash

# URL to fetch
URL="https://game8.jp/pazudora/645084"

# Output file
OUTPUT_FILE="index.html"

# Fetch the page and save it to index.html
curl -o "$OUTPUT_FILE" "$URL"

# Confirm the operation
if [ $? -eq 0 ]; then
    echo "Page successfully fetched and saved to $OUTPUT_FILE"
else
    echo "Failed to fetch the page."
fi
