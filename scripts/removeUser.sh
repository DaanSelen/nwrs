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
    /usr/sbin/deluser $username
    rm -rd /usr/local/nwrs/$username; rm -rd /usr/local/nwrs/web/$username
    sed -i "/DenyUsers $username/d" /etc/ssh/sshd_config
fi
