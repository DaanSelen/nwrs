#!/bin/bash
while getopts u: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
    esac
done
if [ -z "$username" ]
then
    echo "Missing argument"
else
    deluser $username
    rm -rfd /home/$username; rm -rfd /home/web/$username
    sed -i "/DenyUsers $username/d" /etc/ssh/sshd_config
    sed -i "/$username/d" /home/scripts/var/users
fi
