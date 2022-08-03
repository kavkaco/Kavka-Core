#!/usr/bin/env bash

RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
GREEN='\033[0;32m'

filesExists=()
filesList=("allowed_origins" "configs.yml" "redis.yml" "mongo.yml" "smtp.yml")

for f in "${filesList[@]}"
do
    if [ -f "app/configs/${f}" ]; then
        filesExists+=("${f}")
    fi
done

if [ ${#filesExists[@]} -eq 0 ]; then
    cat <<EOT >> app/configs/allowed_origins
*
EOT
    cat <<EOT >> app/configs/configs.yml
LISTEN_PORT: 8000
EOT
    cat <<EOT >> app/configs/redis.yml
HOST: 127.0.0.1
USERNAME:
PASSWORD:
PORT: 6379
DB_NAME: 0
EOT
    cat <<EOT >> app/configs/mongo.yml
HOST: 127.0.0.1
USERNAME: mongo
PASSWORD: mongo
PORT: 27017
DB_NAME: kavka
EOT
    cat <<EOT >> app/configs/smtp.yml
HOST: smtp.gmail.com
PORT: 587
MAIL: 
PASSWORD: 
EOT
    echo -e "${GREEN}+ Configs files initialized"
else
    echo -e "${RED}- Error"
    echo -e "${BLUE}    These config files already exists : "
    for f in "${filesExists[@]}"
    do
        echo -e "${YELLOW}        ${f}"
    done
    echo -e "${BLUE}    Remove them and try again"
fi
