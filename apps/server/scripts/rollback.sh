#!/bin/bash

DEST=$1

pm2 stop $DEST/ecosystem.config.cjs
cp $DEST/main.bak $DEST/main
pm2 start $DEST/ecosystem.config.cjs
