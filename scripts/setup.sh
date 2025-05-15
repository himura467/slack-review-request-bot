#!/usr/bin/env bash

set -e

ROOT_DIR=$(cd $(dirname $0)/..; pwd)
REVIEWER_MAP_FILE='reviewer_map.json'

OP_VAULT_ID=mdwa7hdut7jl67jl5hkcrqjk7m OP_ITEM_ID=r4arfpv6ybzpf7ithlh6ppb7wm op inject -i .env.template -o .env

echo $(op item get nkzbrgx7vgrb3h7lq62blbi24m --format json) | jq '{
  reviews: [.fields[] | select(.label != "notesPlain" and .value != null) | {key: .label, value: .value}] | from_entries
}' | jq '.reviews' > "$REVIEWER_MAP_FILE"
