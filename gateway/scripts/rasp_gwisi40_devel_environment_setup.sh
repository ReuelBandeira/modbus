#!/bin/bash

PATH=/sbin:/bin:/usr/sbin:/usr/bin
. /lib/lsb/init-functions

homedir="/home/icts"

main(){
 
    message_header "V0.7: 14/09/2023 - Setup GW ISI 40 -Dev Env"
    echo "---------------------"
	echo "change root password "
    echo "---------------------"
	echo "icts@icts:~ $ sudo su"
	echo "root@icts:/home/icts# passwd"
	echo "New password:icts"
	echo "Retype new password:icts"
	echo "passwd: password updated successfully"
	echo "root@icts:/home/icts#exit"
	echo "icts@icts:~ $ ls -l"
	echo "Press [ENTER] to proceed"
	read wait

    install_basic_software
    
}

message_header(){
  echo ""
  echo "-------------------------------------------"
  echo "##### ${1} ######"
  echo "-------------------------------------------"
}


install_basic_software(){

    message_header "STEP 01/18 Raspi OS Lite 64 bits update & upgrade"
    sudo apt update
    sudo apt -y upgrade

    message_header "STEP 02/18 forcing usermod"
	#sudo usermod -a -G sudo icts


    message_header "STEP 03/18 Installing & configuring ufw"
    sudo apt -y install ufw
    sudo ufw allow http
    sudo ufw allow https
    sudo ufw allow ssh
    sudo ufw allow 443
    sudo ufw allow 8080
    sudo ufw allow 8585
    sudo ufw allow 8686
    sudo ufw allow 1882
    sudo ufw allow 1883
    sudo ufw allow 502
    sudo ufw allow 44818
    sudo ufw disable
    

    message_header "STEP 04/18 Installing git"
    sudo apt -y install git
	
    message_header "STEP 05/18 Installing GoLang"
	#sudo apt policy golang
	sudo rm -r /usr/local/go
	sudo wget https://go.dev/dl/go1.21.1.linux-arm64.tar.gz
	sudo tar -C /usr/local -xzf go1.21.1.linux-arm64.tar.gz
    sudo rm go1.21.1.linux-arm64.tar.gz
    echo go version 
	
    #echo "------------------------------------------"
    #echo "1) Copy bellow Two lines into .profile    "
    #echo "------------------------------------------"
    #echo "PATH=$PATH:/usr/local/go/bin"
    #echo "GOPATH=$HOME/golang"
    #echo "------------------------------------------"
    #echo "2) Press CTL-O to save and CTL-X to exit  "
    #echo "------------------------------------------"
	
    if grep -Fxq "usr/local/go/bin" ~/.profile
    then
    # code if found
    else
    # code if not found
    echo "PATH=\"\$PATH:/usr/local/go/bin\"" >> ~/.profile
    fi
	
    #if grep -Fxq "GOPATH=" ~/.profile
    #then
    # code if found
    #else
    # code if not found
    #echo "GOPATH=\"\$HOME/golang\"" >> ~/.profile
    #fi	
	
    #echo cat ~/.profile
	#echo "Press [ENTER] to proceed"
	source ~/.profile
    go version 
    #read wait

    message_header "STEP 06/18 Installing snapd - required for flutter"
    sudo apt -y install snapd
    sudo snap install core
	
    message_header "STEP 07/18 Installing flutter - Logoff/login required"
    sudo snap install flutter --classic
	sudo sudo snap refresh

    message_header "STEP 08/18 Installing can-utils"
    sudo apt -y install can-utils

    message_header "STEP 09/18 setup can Apapter MSP2525"
    setup_can_adapter_msp2525	


    message_header "STEP 10/18 Clonning Projects from repositories"
    setup_gwisi40

    message_header "STEP 11/18 Setup VPNClient"
    setup_vpn

    message_header "STEP 12/18 Setup VPNClient"
    setup_vpn

    message_header "STEP 13/18 create_service_openvpn"
    create_service_openvpn

    message_header "STEP 14/18 create_service_mqtt"
   create_service_mqtt

    message_header "STEP 15/18 create_service_backend"
   create_service_backend

    message_header "STEP 16/18 create_service_gateway"
   create_service_gateway

    message_header "STEP 17/18 install telnet"
    sudo apt -y install telnet

    message_header "STEP 18/18 setup complete"
    echo "----------------------------------------------------------"
    echo "Complete flutter installation after logoff/login and type:"
    echo "----------------------------------------------------------"
	echo "flutter --version"
    echo "flutter upgrade --force"
    echo "flutter channel beta"
    echo "flutter doctor -v"
    echo "flutter --version"
    echo "dart --version"
}

setup_vpn() {
    sudo apt install openfortivpn
    echo "---------------------------------------------------"
    echo " To setup OpenForthVpn perform the following steps:"
    echo "--------------------------------------------------"
	echo "Edit /etc/openfortivpn/config and replase values:"
    ### config file for openfortivpn, see man openfortivpn(1) ###
    #
    # host = vpn.example.org
    # port = 443
    # username = vpnuser
    # password = VPNpassw0rd    echo "trusted-cert = Replace with First Run Return Error"


    echo "Type Your VPN Username followed by [ENTER]:"
    read VPNUSER

    echo "Type Ypur VPN PASSWORD followed by [ENTER]:"
    read VPNPASSWORD
	
	# generating a OpenForti VPN configuration file 
	sudo sed -i 's/# host = vpn.example.org/host = 186.208.253.18/g' /etc/openfortivpn/config
	sudo sed -i 's/# port = 443/port = 10443/g' /etc/openfortivpn/config
	sudo sed -i 's/# username = vpnuser/username = $VPNUSER/g' /etc/openfortivpn/config
	sudo sed -i 's/# password = VPNpassw0rd/password = $VPNPASSWORD/g' /etc/openfortivpn/config

    sudo openfortivpn
    read wait
    echo "trusted-cert = <ReturnedHashText>"
    echo "<copy>  line with trusted-cert = <ReturnedHashText>"
    echo "<paste>  followed by [ENTER]"
    read TRUSTEDCERT

    #sed -i 's/s/#trusted-cert = /$TRUSTEDCERT/g' /etc/openfortivpn/config
	sudo nano /etc/openfortivpn/config
    sudo openfortivpn &
 	
}

# NOT REQUIRED - it was already setup by WinApp[Raspberry Ip Imager]
config_wifi_ssid(){
    echo "Type your WI-FI [SSID] followed by [ENTER]:"
    read SSID

    echo "Type your WI-FI [PASSWORD] followed by [ENTER]:"
    read PASSWORD

   # Generate a network configuration file
   #cat << EOF > /etc/wpa_supplicant/wpa_supplicant.conf
   #country=US
   #update_config=1
   #ctrl_interface=/var/run/wpa_supplicant
   #
   #network={
   #    ssid="$SSID"
   #    psk="$PASSWORD"
   #}
   #EOF
   
   sudo service wpa_supplicant start 
}

setup_gwisi40(){
    echo "--------------------------------------------"
    echo " Adding a new SSH key to your GitHub account"
    echo "--------------------------------------------"
	echo "1) Type:"
    echo "--------------------------------------------"
	echo " ssh-keygen -t rsa -b 4096 -C \"your_email@grupoicts.com.br\""
    echo "--------------------------------------------"
	echo "   Generating public/private ed25519 key pair."
	echo "   Windows: Enter file in which to save the key (C:\Users\<YourUserName>/.ssh/id_rsa): <keyFileName>"
	echo "   Linux:   Enter file in which to save the key (/home/devel/.ssh/id_ed25519): <keyFileName>"
	echo "   Enter passphrase (empty for no passphrase):"
	echo "   Enter same passphrase again:"
	echo "   Your identification has been saved in  <keyFileName>"
	echo "   Your public key has been saved in  <keyFileName>.pub"
	echo "   The key fingerprint is:"
	echo "   SHA256:kxJZfwmUW/5osDEQWzVB/JHUmM8YkDCGToEXhQeTmmA your_email@grupoicts.com.br"
	echo "   The key's randomart image is:"
	echo "   +--[ED25519 256]--+"
	echo "   |       .B@O**=.= |"
	echo "   |    E .oB*oo+o* .|"
	echo "   |   . .o*.o.+o. * |"
	echo "   |      o...=.. o o|"
	echo "   |      . S  = o   |"
	echo "   |       . .. o .  |"
	echo "   |           .     |"
	echo "   |                 |"
	echo "   |                 |"
	echo "   +----[SHA256]-----+"
	echo "2) Edit <KeyFileName.pub> copy it contents"
	echo "3) Enter in your github account >> User Settings >> SSH Keys"
	echo "   3.1) Paste it data and save"
	echo " press [ENTER] to procceed..."
    read wait
    mkdir lab 
    cd lab
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/gateway.git
    git clone git@lab.grupoicts.com.br:finep-isi/backend.git
    git clone git@lab.grupoicts.com.br:finep-isi/frontend.git
    git clone https://github.com/mochi-co/mqtt.git
    mkdir protocols
	cd protocols
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/modbus.git
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/ethernetip.git
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/canopen.git
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/newprotocol.git

    git config --global credential.helper store
    git config --global credential.helper cache
    
	if [ -d "$modbus" ]; then
      git config pull.rebase false
      git checkout development
      go build	  
    else
      echo "The folder modbus does not exist."
    fi
    
	if [ -d "$ethernetip" ]; then
      git config pull.rebase false   
      git checkout development
      go build	  
    else
      echo "The folder ethernetip does not exist."
    fi
    
	if [ -d "$canopen" ]; then
      git config pull.rebase false   
      git checkout development
      go build	  
    else
      echo "The folder canopen does not exist."
    fi

	if [ -d "$newprotocol" ]; then
      git config pull.rebase false   
      git checkout development
      go build	  
    else
      echo "The folder newprotocol does not exist."
    fi
	
	cd /home/icts/lab
	
	if [ -d "$gateway" ]; then
      git config pull.rebase false   
      git checkout development
      go build	  
    else
      echo "The folder gateway does not exist."
    fi
	
	if [ -d "$backend" ]; then
      git config pull.rebase false   
      git checkout development
	  cd gwisi40server
      go build
      cd ..	  
    else
      echo "The folder backend does not exist."
    fi
	
    if [ -d "$frontend" ]; then
      git config pull.rebase false   
      git checkout development
	  cd isi_4_0
      flutter build web
	  cp -r build/web ../../backend/gwisi40server
      cd ..	  
    else
      echo "The folder frontend does not exist."
    fi
    
}   


create_service_openvpn() {
   # Generate a openvpn Forti Service file
   sudo cat << EOF > openvpn.service
[Unit]
Description=OpenFortiVPN
After=multi-user.target

[Service]
User=root
Group=root
WorkingDirectory=/home/icts
ExecStart=/usr/bin/openfortivpn
Type=simple

[Install]
WantedBy=multi-user.target}
EOF
   sudo mv openvpn.service /etc/systemd/system/openvpn.service
   sudo systemctl daemon-reload  
   sudo service openvpn start 
}

create_service_mqtt() {
   # Generate a mqtt service file
   sudo cat << EOF > mqtt.service
[Unit]
Description=Local MQTT
After=multi-user.target

[Service]
WorkingDirectory=/home/icts/lab/mqtt/cmd
ExecStart=/home/icts/lab/mqtt/cmd/cmd
Type=simple

[Install]
WantedBy=multi-user.target
EOF
   sudo mv mqtt.service /etc/systemd/system/mqtt.service
   sudo systemctl daemon-reload   
   sudo service mqtt start 
}


create_service_backend() {
   # Generate a backend service file
   sudo cat << EOF > backend.service
[Unit]
Description=ISI40 Backend
After=multi-user.target

[Service]
WorkingDirectory=/home/icts/lab/backend/gwisi40server
ExecStart=/home/icts/lab/backend/gwisi40server/gwisi40server
Type=simple

[Install]
WantedBy=multi-user.target
EOF
   sudo mv backend.service /etc/systemd/system/backend.service
   sudo systemctl daemon-reload   
   sudo service backend start 
}

create_service_gateway() {
   # Generate a gateway service file
   sudo cat << EOF > gateway.service
[Unit]
Description=ISI40 Gateway
After=multi-user.target

[Service]
WorkingDirectory=/home/icts/lab/gateway
ExecStart=/home/icts/lab/gateway/gateway
Type=simple

[Install]
WantedBy=multi-user.target
EOF
   sudo mv gateway.service /etc/systemd/system/gateway.service
   sudo systemctl daemon-reload
   sudo service gateway start 
}

create_service_emunewprotocol() {
   # Generate a gateway service file
sudo cat << EOF > emunewprotocol.service
[Unit]
Description=Emulate NewProtocol
After=multi-user.target

[Service]
WorkingDirectory=/home/icts/lab/emulators/emunewprotocol
ExecStart=/home/icts/lab/emulators/emunewprotocol/emunewprotocol
Type=simple

[Install]
WantedBy=multi-user.target
EOF
   sudo mv emunewprotocol.service /etc/systemd/system/emunewprotocol.service
   sudo systemctl daemon-reload
   sudo service emunewprotocol start 
}


setup_can_adapter_msp2525() {
	
    #echo "---------------------------------------------------"
    #echo " To setup Can Adapter MSP2525 steps:"
    #echo "--------------------------------------------------"
	#echo "Edit /boot/config.txt Add Following Lines:"
    #echo "#dtparam=spi=on remove #"
    #echo "dtparam=spi=on"
    #echo "dtoverlay=mcp2515-can0,oscillator=8000000,interrupt=24"
    #echo "dtoverlay=spi-bcm2835-overlay"
	#echo " press [ENTER] to procceed..."
    #read wait
	#sudo nano /boot/config.txt
	sudo sed -i 's/#dtparam=spi=on/dtparam=spi=on\ndtoverlay=mcp2515-can0,oscillator=8000000,interrupt=24\ndtoverlay=spi-bcm2835-overlay/g' /boot/config.txt
	
    echo "---------------------------------------------------"
    echo " Adding Can Adapter MSP2525 to boot"
    echo "--------------------------------------------------"
	echo "Create  /usr/local/bin/config_can.sh Add Following Lines:"
    echo "#!/bin/bash"
    echo "sudo ip link set can0 up type can bitrate 500000"
	echo " press [ENTER] to procceed..."
    read wait
	sudo nano /usr/local/bin/config_can.sh
	sudo chmod +x /usr/local/bin/config_can.sh

    if [ -d "$/sys/bus/spi/devices/spi0.0" ]; then
      if [ -d "$/sys/bus/spi/devices/spi0.0/net" ]; then
    	  sudo ip link set can0 up type can bitrate 500000 
      else
        echo "The folder /sys/bus/spi/devices/spi0.0/net does not exist."
      fi
    else
      echo "The folder /sys/bus/spi/devices/spi0.0 does not exist."
    fi
}

# Call the main
main "$@"
