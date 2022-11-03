#!/bin/bash
while getopts u: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
        port) port=${OPTARG};;
        contname) contname=${OPTARG};;
    esac
done
UserID=$(id -u $username)
if [ -z "$username" ]
then
    echo "Missing argument"
else
    docker run -it -d --name $username-web --user $UserID -v /usr/local/nwrs/web/$username/html:/usr/share/nginx/html -p $port:8080 nginxinc/nginx-unprivileged:latest
fi

