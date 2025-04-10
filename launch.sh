#!/bin/bash

# starting auth service
echo "Starting auth service..."
cd auth && air > /dev/null 2>&1 &

# obtaining auth service process id
pid=$!

# display success message of starting auth service with process id
echo "Successfully started auth service...Process ID: $pid"

# starting db service
echo "Starting db service..."
cd db && air > /dev/null 2>&1 &

# obtaining db service process id
pid=$!

# display success message of starting db service with process id
echo "Successfully started db service...Process ID: $pid"
