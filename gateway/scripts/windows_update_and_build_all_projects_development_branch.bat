  echo "Updating And Building all projects V3.0 21/11/2023\n\n"


  echo "STEP01 updating and emulators\emucanopentcp ...\n\n"
  
  cd emulators\emucanopentcp
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\

  echo "STEP02 updating and emulators\emuethernetip ...\n\n"
  cd emulators\emuethernetip
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\

  echo "STEP03 updating and emulators\emumodbus ...\n\n"
  cd emulators\emumodbus
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\
 
  echo "STEP04 updating and emulators\emunewprotocol ...\n\n"
  cd emulators\emunewprotocol
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\


  echo "STEP05 updating and building protocols\canopen...\n\n"
  cd protocols\canopen
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\

  echo "STEP06 updating and building protocols\ethernetip...\n\n"
  cd protocols\ethernetip
  git reset --hard
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\

  echo "STEP07 updating and building protocols\modbus...\n\n"
  cd protocols\modbus
  git reset --hard
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\

  echo "STEP08 updating and building protocols\newprotocol...\n\n"
  cd protocols\newprotocol
  git pull
  git checkout development
  git pull
  go build
  cd ..\..\

  echo "STEP09 updating and building backend...\n\n"
  cd backend\gwisi40server
  git reset --hard
  git pull
  git checkout development
  git pull
  go build
  del *.db
  cd ..\..\

  echo "STEP010 updating and building gateway...\n\n"
  cd gateway
  git reset --hard
  git pull
  git checkout development
  git pull
  go build
  cd ..

  echo "STEP011 updating and building factory...\n\n"
  cd factory
  git reset --hard
  git pull
  git checkout development
  git pull
  go build
  cd ..

  echo "STEP12 updating and building frontend...\n\n"
  cd frontend\isi_4_0
  git reset --hard
  git pull
  git checkout development
  git reset --hard
  git pull
  echo "after flutter build type this two bello lines:"
  echo "xcopy /E /I /Y "build\web" "..\..\backend\gwisi40server\web""
  echo "cd ..\..\\n\n"
  flutter clean
  flutter build web
  echo "STEP013 copiyng web folder from frontend to backend.\n\n"
  xcopy /E /I /Y "build\web" "..\..\backend\gwisi40server\web"
  cd ..\..\

  echo "All done!!!\n\n"
