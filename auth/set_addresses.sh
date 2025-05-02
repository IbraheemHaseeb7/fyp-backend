#!/bin/bash

if [ -z "$1" ]; then
  intf="enp4s0"
else
  intf="$1"
fi

# finding the current address and the assigned one
BASE=$(cat .env | grep BASE_ADDRESS | cut -d '=' -f2)
REAL=$(ip addr show $intf | grep 'inet ' | awk '{print $2}' | cut -d '/' -f1)

# replacing the texts
sed -i "s/$BASE/$REAL/g" ".env"

