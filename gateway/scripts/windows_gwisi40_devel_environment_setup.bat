    echo "--------------------------------------------"
    echo " Adding a new SSH key to your GitHub account"
    echo "--------------------------------------------"
	echo "1) Type:"
    echo "--------------------------------------------"
	echo " ssh-keygen -t rsa -b 4096 -C \"your_email@grupoicts.com.br\""
    echo "--------------------------------------------"
	echo "   Generating public/private ed25519 key pair."
	echo "   Windows: Enter file in which to save the key (C:\Users\<YourUserName>/.ssh/id_rsa): <keyFileName>"
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
    d: 
    mkdir wks_cts_finep_gwisi40
    cd wks_cts_finep_gwisi40	
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/gateway.git
    git clone git@lab.grupoicts.com.br:finep-isi/backend.git
    git clone git@lab.grupoicts.com.br:finep-isi/frontend.git
    git clone https://github.com/mochi-co/mqtt.git
	cd protocols
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/modbus.git
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/ethernetip.git
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/canopen.git
    git clone git@lab.grupoicts.com.br:finep-isi/gateway/newprotocol.git
	cd ..
    mkdir emulators
	cd emulators
    git clone git@lab.grupoicts.com.br:finep-isi/emulators/emumodbus.git
    git clone git@lab.grupoicts.com.br:finep-isi/emulators/emuethernetip.git
    git clone git@lab.grupoicts.com.br:finep-isi/emulators/emucanopentcp.git
    git clone git@lab.grupoicts.com.br:finep-isi/emulators/emunewprotocol.git
	cd ..
    mkdir test_tools
	cd test_tools
    git clone git@lab.grupoicts.com.br:finep-isi/emutesttools/canopentcpclient.git  
    git clone git@lab.grupoicts.com.br:finep-isi/emutesttools/eiptcpclient.git
    git clone git@lab.grupoicts.com.br:finep-isi/emutesttools/gw40post.git
    git clone git@lab.grupoicts.com.br:finep-isi/emutesttools/mbtcpclient.git
    git clone git@lab.grupoicts.com.br:finep-isi/emutesttools/mewprotocolclient.git
    cd ..
	git config --global credential.helper store
    git config --global credential.helper cache
    echo "-----------------------------------------------------------------------------------"
    echo " 1. Install Latest GoLang  : https://go.dev/dl/"
    echo " 2. Install Latest Flutter : https://docs.flutter.dev/get-started/install"
    echo " 3. Install Latest MQTTX   : https://mqttx.app/docs/cli/downloading-and-installation"
    echo " 3. Install Latest MQTTBox : https://mqttbox.softonic.com.br/"
    echo "-----------------------------------------------------------------------------------"
