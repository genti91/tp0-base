#!/bin/bash

CONFIG_FILE="server/config.ini"

SERVER_HOST=$(grep -E '^SERVER_IP\s*=' "$CONFIG_FILE" | cut -d'=' -f2 | xargs)
SERVER_PORT=$(grep -E '^SERVER_PORT\s*=' "$CONFIG_FILE" | cut -d'=' -f2 | xargs)

TEST_MESSAGE="hola"

RESPONSE=$(docker run --rm -i --network tp0_testing_net busybox sh -c "echo \"$TEST_MESSAGE\" | nc $SERVER_HOST $SERVER_PORT")


if [ "$RESPONSE" = "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"    
fi
