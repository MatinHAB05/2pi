#!/bin/sh
set -e

ELASTIC_PASSWORD=${ELASTIC_PASSWORD:?ELASTIC_PASSWORD is required}
KIB_URL=${KIBANA_URL:-"http://kibana:5601"}


curl -sS -X POST "$KIB_URL/api/saved_objects/_import?overwrite=true" \
  -H "kbn-xsrf: true" \
  -u "elastic:$ELASTIC_PASSWORD" \
  -F file=@/scripts/export.ndjson