#!bin/bash
sed -i '/NextPort=/c\NextPort=10001' var/nextport
echo  -n 'Done'
