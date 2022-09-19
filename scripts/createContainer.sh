#!/bin/bash
while getopts u: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
    esac
done
source /home/scripts/var/highestPort
newHighest=$((Highest+1)); echo "Highest=$newHighest" > /home/scripts/var/highestPort
UserID=$(id -u $username)
if [ -z "$username" ]
then
    echo "Missing argument"
else
    docker run -it -d --name $username-web --user $UserID -v /home/web/$username/html:/usr/share/nginx/html -p $Highest:8080 nginxinc/nginx-unprivileged:latest
fi
echo $Highest

