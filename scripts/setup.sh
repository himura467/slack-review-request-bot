#!/usr/bin/env bash

set -e

REVIEWER_MAP_FILE='reviewer_map.json'

op item get nkzbrgx7vgrb3h7lq62blbi24m --vault "Slack Review Request Bot" --format json | jq '{
  reviews: [.fields[] | select(.label != "notesPlain" and .value != null) | {key: .label, value: .value}] | from_entries
}' | jq '.reviews' > "$REVIEWER_MAP_FILE"
