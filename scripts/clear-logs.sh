#!/bin/sh
set -e

rm logs/app/*.log
echo "clear app logs succusfully"

rm logs/elastic/*.log
echo "clear elastic logs succusfully"

rm logs/telegram/*.log
echo "clear telegram logs succusfully"
