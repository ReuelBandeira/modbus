FILE="utils/version/version_protocol.txt"
COUNTER=$(cat "$FILE")
COUNTER=$((COUNTER+1))
echo "$COUNTER" > "$FILE"