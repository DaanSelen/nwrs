#!/bin/bash
while getopts u:c:port flag
do
    case "${flag}" in
        u) username=${OPTARG};;
        c) containername=${OPTARG};;
        port) port=${OPTARG};;
    esac
done
UserID=$(id -u $username)
if [ -z "$username" ] || [ -z "$containername" ] || [ -z $port ]
then
    echo $port
    echo "Missing argument"
else
    docker run -it -d --name $containername --user $UserID -v /usr/local/nwrs/web/$username/html:/usr/share/nginx/html -p $port:8080 nginxinc/nginx-unprivileged:latest
fi

