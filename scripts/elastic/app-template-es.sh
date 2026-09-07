#!/bin/sh
set -e

ELASTIC_URL=${ELASTIC_URL:-"http://elasticsearch:9200"}
ELASTIC_PASSWORD=${ELASTIC_PASSWORD:?ELASTIC_PASSWORD is required}

TEMPLATE_NAME="app_logs_template"

echo "Waiting for Elasticsearch..."
until curl -s -f -u "elastic:$ELASTIC_PASSWORD" "$ELASTIC_URL/_cat/health?h=status" | grep -qE 'green|yellow'; do
  sleep 2
done

echo "Elasticsearch is up!"

# Check if index template exists
if curl -s -f -u "elastic:$ELASTIC_PASSWORD" "$ELASTIC_URL/_index_template/$TEMPLATE_NAME" > /dev/null; then
  echo "Index template '$TEMPLATE_NAME' already exists. Skipping creation."
  exit 0
fi

echo "Creating index template '$TEMPLATE_NAME'..."

curl -sS -f -X PUT "$ELASTIC_URL/_index_template/$TEMPLATE_NAME" \
  -u "elastic:$ELASTIC_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
  "index_patterns": ["app-logs-*"],
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
        "@version": {
          "type": "keyword"
        },
        "app_name": {
          "type": "keyword"
        },
        "log": {
          "properties": {
            "level": {
              "type": "keyword"
            },
            "logger": {
              "type": "keyword"
            },
            "file": {
              "properties": {
                "path": {
                  "type": "keyword"
                }
              }
            }
          }
        },
        "message": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        },
        "category": {
          "type": "keyword"
        },
        "sub_category": {
          "type": "keyword"
        },
        "extra": {
          "type": "flattened"
        },
        "tags": {
          "type": "keyword"
        },
        "ecs": {
          "properties": {
            "version": {
              "type": "keyword"
            }
          }
        },
        "host": {
          "properties": {
            "name": {
              "type": "keyword"
            },
            "hostname": {
              "type": "keyword"
            },
            "ip": {
              "type": "keyword"
            }
          }
        },
        "agent": {
          "properties": {
            "id": {
              "type": "keyword"
            },
            "name": {
              "type": "keyword"
            },
            "type": {
              "type": "keyword"
            },
            "version": {
              "type": "keyword"
            }
          }
        }
      }
    }
  }
}'

echo "\nIndex template '$TEMPLATE_NAME' created successfully."