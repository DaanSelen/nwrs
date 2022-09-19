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
    echo -e "$password\n$password" | adduser $username
    echo $username:$password | /usr/sbin/chpasswd
    echo "DenyUsers $username" >> /etc/ssh/sshd_config
    mkdir /home/web/$username; mkdir /home/web/$username/html
    chown -R $username /home/web/$username/; chmod -R 700 /home/web/$username/
    echo $username >> var/users
fi

