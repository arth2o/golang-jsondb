#!/bin/bash

if [ ! -d ./bin ]; then
    mkdir -p ./bin
fi

if [ ! -d ./data ]; then
    mkdir -p ./data
fi

if [ -f ./bin/server ]; then
    rm ./bin/server
fi

go build -o bin/server cmd/server/main.go

if [ -f .env.development ]; then
    source .env.development
    ./bin/server
else
    echo "No .env.development file found. Please create one."
    ## ./bin/server    
fi

## run with encryption
## ENABLE_ENCRYPTION=true ENCRYPTION_KEY=0123456789abcdef0123456789abcdef ./bin/server
