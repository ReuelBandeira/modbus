#!/bin/bash

PATH=/sbin:/bin:/usr/sbin:/usr/bin
. /lib/lsb/init-functions

homedir="/home/icts"

main(){

    echo "================================================"
    echo "step01-0 Running Linux Build All v0.3 23-11-2023"
    echo "================================================"

    source ~/.profile

    echo "==================================================="
    echo "step01-2 updating and building protocols/canopen..."
    echo "==================================================="
    cd protocols/canopen
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "======================================================"
    echo "step01-3 updating and building protocols/ethernetip..."
    echo "======================================================"
    cd protocols/ethernetip
	git reset --hard
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "=================================================="
    echo "step01-4 updating and building protocols/modbus..."
    echo "=================================================="
    cd protocols/modbus
	git reset --hard
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "======================================================="
    echo "step01-5 updating and building protocols/newprotocol..."
    echo "======================================================="
    cd protocols/newprotocol
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "======================================================="
    echo "step01-6 updating and building emulatorrs/emucanopen..."
    echo "======================================================="
    cd emulators/emucanopentcp
    git pull
    git checkout development
    git pull
    go build
    cd ../../
	
    echo "=========================================================="
    echo "step01-7 updating and building emulatorrs/emuethernetip..."
    echo "=========================================================="
    cd emulators/emuethernetip
    git pull
    git checkout development
    git pull
    go build
    cd ../../


    echo "======================================================"
    echo "step01-8 updating and building emulatorrs/emumodbus..."
    echo "======================================================"
    cd emulators/emumodbus
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "==========================================================="
    echo "step01-9 updating and building emulatorrs/emunewprotocol..."
    echo "==========================================================="
    cd emulators/emunewprotocol
    git pull
    git checkout development
    git pull
    go build
    cd ../../


    echo "=========================================="
    echo "step01-10 updating and building gateway..."
    echo "=========================================="
    cd gateway
	git reset --hard
    git pull
    git checkout development
    git pull
    go build
    cd ..

    echo "=========================================="
    echo "step01-11 updating and building backend..."
    echo "=========================================="
    cd backend/gwisi40server
	git reset --hard
    git pull
    git checkout development
    git pull
    go build
    sudo rm *.db
    cd ../../

 
    echo "=========================================="
    echo "STEP01-12 updating and building factory..."
    echo "=========================================="
    cd factory
    git reset --hard
    git pull
    git checkout development
    git pull
    go build
    cd ..

    echo "==========================================="
    echo "step01-13 updating and building frontend..."
    echo "==========================================="
    cd frontend/isi_4_0
	git reset --hard
    git pull
    git checkout development
    git pull
	flutter clean   
    flutter build web
    cp -r build/web/ ../../backend/gwisi40server/
    cd ../../
    echo " ============"
    echo " All done!!!!"
    echo " ============"
}

# Call the main
main "$@"

