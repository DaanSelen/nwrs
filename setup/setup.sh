#!bin/bash
sudo echo "This needs to be run as admin."
apt-get install docker.io vsftpd -y; systemctl enable vsftpd
cp ./vsftpd.conf /etc/vsftpd.conf; touch /etc/vsftpd.chroot_list
chmod 444 /etc/vsftpd.chroot_list
chmod 666 /usr/local/nwrs/scripts/var/*

