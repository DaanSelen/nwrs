#!/bin/bash
while getopts u: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
    esac
done
source /usr/local/nwrs/scripts/var/vars
newNextPort=$((NextPort+1)); echo "NextPort=$newNextPort" > /usr/local/nwrs/scripts/var/vars
UserID=$(id -u $username)
if [ -z "$username" ]
then
    echo "Missing argument"
else
    docker run -it -d --name $username-web --user $UserID -v /usr/local/nwrs/web/$username/html:/usr/share/nginx/html -p $NextPort:8080 nginxinc/nginx-unprivileged:latest
fi
echo $Highest

