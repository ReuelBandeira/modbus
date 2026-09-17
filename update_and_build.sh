#!/bin/bash

source ~/.profile
echo $PATH

echo "step01-1 updating and building protocols/canopen..."
cd protocols/canopen
git pull
git checkout development
git pull
go build
cd ..

echo "step01-2 updating and building ethernetip..."
cd ethernetip
git pull
git checkout development
git pull
go build
cd ..

echo "step01-3 updating and building modbus..."
cd modbus
git pull
git checkout development
git pull
go build
cd ..

echo "step01-4 updating and building newprotocol..."
cd newprotocol
git pull
git checkout development
git pull
go build
cd ../..

echo "step01-5 updating and building backend..."
cd backend/gwisi40server
git pull
git checkout development
git pull
go build
rm -f *.db
cd ../..

echo "step01-6 updating and building gateway..."
cd gateway
git pull
git checkout development
git pull
go build
cd ..

echo "step01-7 updating and building frontend..."
cd frontend/isi_4_0
git pull
git checkout development
git pull
yarn install
flutter clean
flutter build web
echo "Movendo Build"
cp -r build/web/ ../../backend/gwisi40server/
cd ../..
