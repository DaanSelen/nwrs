#!bin/bash
source /usr/local/nwrs/scripts/var/vars
sed -i '/NextPort=/c\NextPort=10001' /usr/local/nwrs/scripts/var/vars
echo  -n "Done" $NextPort
