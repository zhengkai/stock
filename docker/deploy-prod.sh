#! /usr/bin/env bash

DOCKER_NAME="stock-prod"
# DOCKER_IMAGE="zhengkai/stock:latest"
DOCKER_IMAGE="stock-image"

if [ "${HOSTNAME,,}" == "doll" ]; then
	exit
fi

cd "$(dirname "$(readlink -f "$0")")" || exit 1

# sudo docker pull "$DOCKER_IMAGE"

ENV_FILE="./env-prod.sh"
if [ ! -e "$ENV_FILE" ]; then
	cp env.sh.example "$ENV_FILE"
	mkdir -p static
	echo "export STOCK_DIR=\"$(pwd)/static\"" >> "$ENV_FILE"
fi
. "$ENV_FILE"

sudo docker stop "$DOCKER_NAME" || :
sudo docker rm "$DOCKER_NAME" || :

set -x
sudo docker run -d --name "$DOCKER_NAME" \
	--env "STOCK_VAPID_PUBLIC_KEY=${STOCK_VAPID_PUBLIC_KEY}" \
	--env "STOCK_VAPID_PRIVATE_KEY=${STOCK_VAPID_PRIVATE_KEY}" \
	--env "STOCK_HOST=${HOSTNAME,,}" \
	-p "${STOCK_WEB}:80" \
	--mount "type=bind,source=${STOCK_DIR},target=/static" \
	--restart always \
	"$DOCKER_IMAGE"
