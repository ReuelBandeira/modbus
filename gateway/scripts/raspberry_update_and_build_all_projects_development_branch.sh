#!/bin/bash

PATH=/sbin:/bin:/usr/sbin:/usr/bin
. /lib/lsb/init-functions

homedir="/home/icts"

main(){

    echo "step01-0 Running Autoupdate v0.3 01-09-2023..."

    echo "step01-1 Stopping services gateway, backend, mqtt and openvpn..."
    source ~/.profile
    cd lab
    sudo service gateway stop
    sudo service backend stop
    sudo service mqtt stop
    sudo service openvpn stop

    echo "step01-2 updating and building backend..."
    cd backend/gwisi40server
    git pull
    git checkout development
    git pull
    go build
    sudo rm *.db
    cd ../../

    echo "step01-3 updating and building gateway..."
    cd gateway
    git pull
    git checkout development
    git pull
    go build
    cd ..

    echo "step01-4 updating and building frontend..."
    cd frontend/isi_4_0
    git pull
    git checkout development
    git pull
    flutter build web
    cp -r build/web/ ../../backend/gwisi40server/
    cd ../../

    echo "step01-5 updating and building protocols/canopen..."
    cd protocols/canopen
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "step01-6 updating and building protocols/ethernetip..."
    cd protocols/ethernetip
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "step01-7 updating and building protocols/modbus..."
    cd protocols/modbus
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "step01-8 updating and building protocols/newprotocol..."
    cd protocols/newprotocol
    git pull
    git checkout development
    git pull
    go build
    cd ../../

    echo "step01-9 Starting services openvpn, mqtt, backend & gateway..."
    sudo service openvpn start
    sudo service mqtt start
    sudo service backend start
    sudo service gateway start
}


# Call the main
main "$@"

