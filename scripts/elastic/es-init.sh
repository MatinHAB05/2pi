#!/bin/sh
set -e

ELASTIC_URL=${ELASTIC_URL:-"http://elasticsearch:9200"}
ELASTIC_PASSWORD=${ELASTIC_PASSWORD:?ELASTIC_PASSWORD is required}
KIBANA_SYSTEM_PASSWORD=${KIBANA_SYSTEM_PASSWORD:?KIBANA_SYSTEM_PASSWORD is required}


echo "Waiting for Elasticsearch..."
until curl -s -f -u "elastic:$ELASTIC_PASSWORD" "$ELASTIC_URL/_cat/health?h=status" | grep -qE 'green|yellow'; do
  sleep 2
done

echo "Elasticsearch is up!"

# 1. Configure kibana_system password
echo "Setting kibana_system password..."
curl -sS -f -X POST "$ELASTIC_URL/_security/user/kibana_system/_password" \
  -u "elastic:$ELASTIC_PASSWORD" \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"$KIBANA_SYSTEM_PASSWORD\"}"
echo "kibana_system password configured successfully."
echo ""
