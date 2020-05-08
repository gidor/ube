# configurazione



sudo setcap CAP_NET_BIND_SERVICE=+eip /opt/ube/ube



   sudo mkdir /opt/ube
   sudo cp ube /opt/ube/
   sudo cp .env /opt/ube/
   sudo cp ube.service /opt/ube/
   sudo cp base.yaml /opt/ube/

   sudo mkdir /etc/ube
   sudo ln -s /opt/ube/base.yaml  /etc/ube/base.yaml
   sudo mkdir /var/log/ube
   sudo useradd  ube -s /usr/sbin/nologin -M
   sudo ln -s /opt/ube/ube.service /lib/systemd/system/ube.service
   sudo setcap CAP_NET_BIND_SERVICE=+eip /opt/ube/ube


   sudo systemctl enable ube.service 
   sudo systemctl status  ube.service 
   sudo systemctl start  ube.service 

   sudo lsof -i -P -n | grep LISTEN

   sudo vi /opt/ube/ube.service 
   sudo vi /opt/ube/.env 

   sudo systemctl daemon-reload 
   sudo systemctl  restart ube
   sudo systemctl  status ube

