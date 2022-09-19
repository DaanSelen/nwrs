#!bin/bash
sed -i '/NextPort=/c\NextPort=10001' var/vars
echo  -n 'Done'
