#!/bin/bash
while getopts u:p: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
        p) password=${OPTARG};;
    esac
done
if [ -z "$username" ] || [ -z "$password" ]
then
    echo "Missing one (or both) argument(s)"
else
    echo -e "$password\n$password" | /usr/sbin/adduser $username
    echo $username:$password | /usr/sbin/chpasswd
    echo "DenyUsers $username" >> /etc/ssh/sshd_config
    mkdir /usr/local/nwrs/web/$username; mkdir /usr/local/nwrs/web/$username/html
    chown -R $username /usr/local/nwrs/web/$username/; chmod -R 700 /usr/local/nwrs/web/$username/
fi

