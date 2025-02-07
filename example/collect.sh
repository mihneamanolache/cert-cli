FILE="./example/domains.txt"
OUTPUT_DIR="./example/$(date +%d-%m-%y)"
PROXY="$YOUR_PROXY_SERVER"
ERROR_FILE="$OUTPUT_DIR/error.txt"

mkdir -p "$OUTPUT_DIR"  

collect_certificates() {
  DOMAIN="$1"
  echo "Collecting $DOMAIN"
  
  # Try to collect certificates
  cert-cli -q "$DOMAIN" --match "LIKE" --proxy "$PROXY" -o "$OUTPUT_DIR/$DOMAIN" > /dev/null
  STATUS=$?

  if [ $STATUS -ne 0 ]; then
    echo "Error collecting $DOMAIN, status code: $STATUS"
    echo "$DOMAIN" >> "$ERROR_FILE"
  else
    echo "Done collecting $DOMAIN"
  fi
}

export -f collect_certificates 
export OUTPUT_DIR PROXY ERROR_FILE

cat "$FILE" | xargs -P 25 -I {} bash -c 'collect_certificates "{}"'
