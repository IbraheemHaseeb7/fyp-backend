#!/bin/bash

# Path to your project
WORKDIR=~/work/fyp_b

# Open 4 tabs in GNOME Terminal, each running the same command
gnome-terminal --tab --title="apee-i" -- bash -c "cd $WORKDIR && nvim api.yaml; exec bash"
gnome-terminal --tab --title="docker" -- bash -c "cd $WORKDIR && clear && sudo docker start 9f8a92492793 && sudo docker exec -it 9f8a92492793 bash; exec bash"
gnome-terminal --tab --title="auth" -- bash -c "cd $WORKDIR/auth && clear && air; exec bash"
gnome-terminal --tab --title="db" -- bash -c "cd $WORKDIR/db && clear && air; exec bash"
gnome-terminal --tab --title="img" -- bash -c "cd $WORKDIR/img && clear && air; exec bash"
