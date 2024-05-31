#!/bin/bash
source .env
echo "GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DATABASE_URI goose -dir deployments/migration/pogresql/ $1"
GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DATABASE_URI goose -dir deployments/migration/pogresql/ $1
