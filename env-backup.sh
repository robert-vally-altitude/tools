#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <local_backup_path>"
    exit 1
fi

# Check for config file
CONFIG_FILE=".env-backup-config"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found at $CONFIG_FILE"
    echo "Please create it with the following format:"
    echo "ENV_NAME|USERNAME@SERVER_IP|KEY_FILE_PATH|SOURCE_FILE_PATH|DEST_FILENAME_PATTERN"
    echo "Example:"
    echo "DEV|ubuntu@11.22.333.444|~/keys/DEV-singapore.pem|/var/www/html/server/.env|env-{env}-{timestamp}"
    exit 1
fi

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_PATH="$1"
mkdir -p "$BACKUP_PATH"

# Read config file and build command list
# The || [ -n "$line" ] ensures the last line is processed even without a newline
while IFS='|' read -r ENV_NAME SERVER_SPEC KEY_FILE SOURCE_FILE DEST_PATTERN || [ -n "$ENV_NAME" ]; do
    [[ -z "$ENV_NAME" || "$ENV_NAME" =~ ^# ]] && continue
    
    KEY_FILE="${KEY_FILE/#\~/$HOME}"
    DEST_FILENAME=$(echo "$DEST_PATTERN" | sed "s/{env}/$ENV_NAME/g" | sed "s/{timestamp}/$TIMESTAMP/g")
    
    CMD="scp -q -i \"$KEY_FILE\" $SERVER_SPEC:$SOURCE_FILE \"$BACKUP_PATH/$DEST_FILENAME\""
    echo -n "Copying $ENV_NAME:$SOURCE_FILE to $BACKUP_PATH/$DEST_FILENAME... "
    
    if eval "$CMD" 2>/dev/null; then
        echo "Success"
    else
        echo "Failure"
    fi
done < "$CONFIG_FILE"

echo "All files have been copied to $BACKUP_PATH" 