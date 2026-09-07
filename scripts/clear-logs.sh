#!/bin/sh
set -e

rm -f logs/app/*.log
echo "clear app logs succusfully"

rm -f logs/elastic/*.log
echo "clear elastic logs succusfully"

rm -f logs/telegram/*.log
echo "clear telegram logs succusfully"


rm -f clean-logs/app/*.log
echo "clear app logs succusfully"

rm -f clean-logs/elastic/*.log
echo "clear elastic logs succusfully"

rm -f clean-logs/telegram/*.log
echo "clear telegram logs succusfully"

