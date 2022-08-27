#! /bin/bash
ps -U sergey | grep app | awk '{ print $1 }' | xargs kill
docker stop db &> /dev/null

echo "App and DB container are stopped"