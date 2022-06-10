#!/usr/bin/sh

# Remove `quotes.json`
rm artefacts/quotes.json
rm cmd/quotes.json
rm internal/quotes.json

# Remove `help.json`
rm core/help.json

# Create `quotes.json`
echo "{\"quotes\":[]}" > artefacts/quotes.json