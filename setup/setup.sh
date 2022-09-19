#!bin/bash
sudo echo "This needs to be run as admin."
apt-get install docker.io -y
apt-get install vsftpd -y; systemctl enable vsftpd
