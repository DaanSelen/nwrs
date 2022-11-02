#!bin/bash
sudo echo "This needs to be run as admin."
apt-get install docker.io vsftpd -y; systemctl enable vsftpd docker
cp ./vsftpd.conf /etc/vsftpd.conf; touch /etc/vsftpd.chroot_list
chmod 440 /etc/vsftpd.chroot_list
chmod 660 /usr/local/nwrs/scripts/var/*
systemctl status vsftpd docker

