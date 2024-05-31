#!/bin/bash
set -e
DEST_PATH=$1

cp $DEST_PATH/main $DEST_PATH/main.bak

if [ -z "$DEST_PATH" ]; then
  echo "Usage: deploy.sh <src> <dest>"
  exit 1
fi

if [ ! -d "$DEST_PATH" ]; then
  echo "Destination path does not exist"
  exit 1
fi

if [ ! -f "$DEST_PATH/main.new" ]; then
  echo "main.new does not exist"
  exit 1
fi

pm2 stop $DEST_PATH/ecosystem.config.cjs
cp $DEST_PATH/main.new $DEST_PATH/main
pm2 start $DEST_PATH/ecosystem.config.cjs

