#!/bin/bash
while getopts u: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
        id) id=${OPTARG};;
        port) port=${OPTARG};;
    esac
done
UserID=$(id -u $username)
if [ -z "$username" ] || [ -z "$id" ]
then
    echo "Missing argument"
else
    docker run -it -d --name $username-web-$id --user $UserID -v /usr/local/nwrs/web/$username/html:/usr/share/nginx/html -p $port:8080 nginxinc/nginx-unprivileged:latest
fi

