#!/bin/sh
set -e

ES_URL=${ES_URL:-"http://elasticsearch:9200"}
ES_PASSWORD=${ES_PASSWORD:?ES_PASSWORD is required}
KIBANA_SYSTEM_PASSWORD=${KIBANA_SYSTEM_PASSWORD:?KIBANA_SYSTEM_PASSWORD is required}

INDEX_NAME="articles-v1"
ALIAS_NAME="articles"

echo "Waiting for Elasticsearch..."
until curl -s -f -u "elastic:$ES_PASSWORD" "$ES_URL/_cat/health?h=status" | grep -qE 'green|yellow'; do
  sleep 2
done

echo "Elasticsearch is up!"

# 1. Configure kibana_system password
echo "Setting kibana_system password..."
curl -sS -f -X POST "$ES_URL/_security/user/kibana_system/_password" \
  -u "elastic:$ES_PASSWORD" \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"$KIBANA_SYSTEM_PASSWORD\"}"
echo "kibana_system password configured successfully."
echo ""

# 2. Create Index with Settings & Mappings (if not exists)
# TODO : Synonyms file
if ! curl -s -f -u "elastic:$ES_PASSWORD" "$ES_URL/$INDEX_NAME" > /dev/null; then
  echo "Creating index $INDEX_NAME..."
  curl -sS -f -X PUT "$ES_URL/$INDEX_NAME" \
    -u "elastic:$ES_PASSWORD" \
    -H "Content-Type: application/json" \
    -d '{
      "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0,
        "index.mapping.coerce": false,
        "analysis": {
          "char_filter": {
            "zero_width_spaces": {
              "type": "mapping",
              "mappings": [
                "\\u200C=>\\u0020"
              ]
            }
          },
          "filter": {
            "persian_stop": {
              "type": "stop",
              "stopwords": "_persian_"
            },
            "english_stop": {
              "type": "stop",
              "stopwords": "_english_"
            },
            "english_stemmer": {
              "type": "stemmer",
              "language": "english"
            },
            "english_possessive_stemmer": {
              "type": "stemmer",
              "language": "possessive_english"
            },
            "custom_synonyms": {
              "type": "synonym",
              "synonyms": []
            }
          },
          "analyzer": {
            "fa_en_analyzer": {
              "type": "custom",
              "tokenizer": "standard",
              "char_filter": [
                "zero_width_spaces"
              ],
              "filter": [
                "english_possessive_stemmer",
                "lowercase",
                "decimal_digit",
                "arabic_normalization",
                "persian_normalization",
                "persian_stop",
                "english_stop",
                "custom_synonyms",
                "persian_stem",
                "english_stemmer"
              ]
            },
            "fa_en_fuzzy_analyzer": {
              "type": "custom",
              "tokenizer": "standard",
              "char_filter": [
                "zero_width_spaces"
              ],
              "filter": [
                "lowercase",
                "decimal_digit",
                "arabic_normalization",
                "persian_normalization"
              ]
            }
          }
        }
      },
      "mappings": {
        "dynamic": "strict",
        "properties": {
          "id": {
            "type": "keyword"
          },
          "created_at": {
            "type": "date"
          },
          "deleted_at": {
            "type": "date"
          },
          "updated_at": {
            "type": "date"
          },
          "title": {
            "type": "text",
            "analyzer": "fa_en_analyzer",
            "fields": {
              "fuzzy": {
                "type": "text",
                "analyzer": "fa_en_fuzzy_analyzer"
              },
              "keyword": {
                "type": "keyword",
                "ignore_above": 256
              }
            }
          },
          "description": {
            "type": "text",
            "analyzer": "fa_en_analyzer",
            "fields": {
              "fuzzy": {
                "type": "text",
                "analyzer": "fa_en_fuzzy_analyzer"
              }
            }
          },
          "image_url": {
            "type": "keyword",
            "index": false
          },
          "url": {
            "type": "keyword",
            "index": false
          }
        }
      }
    }'
  echo ""
else
  echo "Index $INDEX_NAME already exists, skipping creation."
fi

# 3. Setup Alias
echo "Ensuring alias $ALIAS_NAME points to $INDEX_NAME..."
curl -sS -f -X POST "$ES_URL/_aliases" \
  -u "elastic:$ES_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
    "actions": [
      {
        "add": {
          "index": "'"$INDEX_NAME"'",
          "alias": "'"$ALIAS_NAME"'"
        }
      }
    ]
  }'
echo ""

echo "Elasticsearch setup completed successfully!"