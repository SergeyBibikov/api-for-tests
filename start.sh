#! /bin/bash

if [ -z "$DBPass" ]; then
    echo "Please set DBPass env variable to contain the database password"
    exit 1
fi

docker rmi db
docker build -t db --no-cache .
docker run --rm -it -p 5432:5432 --name db -e POSTGRES_PASSWORD=$DBPass -d db 
echo "DB started"

go build -o ./app
echo "App built"

./app > log &
