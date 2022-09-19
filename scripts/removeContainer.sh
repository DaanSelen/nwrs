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
    docker stop $username-web
    docker rm $username-web
fi
