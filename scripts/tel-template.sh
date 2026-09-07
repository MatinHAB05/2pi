#!/bin/sh
set -e

ELASTIC_URL=${ELASTIC_URL:-"http://elasticsearch:9200"}
ELASTIC_PASSWORD=${ELASTIC_PASSWORD:?ELASTIC_PASSWORD is required}

TEMPLATE_NAME="telegram_logs_template"

echo "Waiting for Elasticsearch..."
until curl -s -u "elastic:$ELASTIC_PASSWORD" "$ELASTIC_URL/_cat/health?h=status" | grep -qE 'green|yellow'; do
  sleep 2
done

echo "Elasticsearch is up!"

# Check if index template exists
if curl -s  -u "elastic:$ELASTIC_PASSWORD" "$ELASTIC_URL/_index_template/$TEMPLATE_NAME" > /dev/null; then
  echo "Index template '$TEMPLATE_NAME' already exists. Skipping creation."
  exit 0
fi

echo "Creating index template '$TEMPLATE_NAME'..."

curl -sS -X PUT "$ELASTIC_URL/_index_template/$TEMPLATE_NAME" \
  -u "elastic:$ELASTIC_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
  "index_patterns": ["telegram-logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 1,
      "index.mapping.coerce": false
    },
    "mappings": {
      "properties": {
        "@timestamp": {
          "type": "date"
        },
        "event": {
          "properties": {
            "action": {
              "type": "keyword"
            }
          }
        },
        "telegram_url": {
          "type": "keyword"
        },
        "telegram_payload": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        },
        "telegram": {
          "properties": {
            "payload_json": {
              "type": "object",
              "dynamic": true
            }
          }
        }
      }
    }
  }
}'

echo "\nIndex template '$TEMPLATE_NAME' created successfully."