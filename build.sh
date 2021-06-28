#!/bin/bash
WORKING_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
LINUX_BINARY=linux_sindel_auth
LINUX_BINARY_NEW=linux_sindel_auth.new

REMOTE_BINARY=sindel_auth
REMOTE_BINARY_NEW=sindel_auth.new

NATIVE_BINARY=sindel_auth
NATIVE_BINARY_NEW=sindel_auth.new

REMOTE_CONN=skabox.ru
REMOTE_SERVICE_NAME=sindel_auth
REMOTE_DIR=/var/www/skabox.ru/services/auth

build_linux() {
  docker build -t sindelsecurify:1.0 .
  CONTAINER_ID=`docker run -e GOOS=darwin -e GOARCH=amd64 -id sindelsecurify:1.0`
  docker cp $CONTAINER_ID:/app/$LINUX_BINARY ./build/$LINUX_BINARY
  docker stop $CONTAINER_ID
}

build_native() {
  go build -o $WORKING_PATH/build/$NATIVE_BINARY .
}

remote_deploy() {
  rsync -avz ./build/$LINUX_BINARY $REMOTE_CONN:$REMOTE_DIR/$REMOTE_BINARY_NEW
  ssh $REMOTE_CONN "mv $REMOTE_DIR/$REMOTE_BINARY_NEW $REMOTE_DIR/$REMOTE_BINARY"
  ssh $REMOTE_CONN "pkill $REMOTE_SERVICE_NAME"
}

run_native() {
  SERVICE_BINARY="$WORKING_PATH/build/$NATIVE_BINARY -c $WORKING_PATH/build/config.json"
  if [ ! -f $WORKING_PATH/build/run.pid ] ; then
    daemonize -p $WORKING_PATH/build/run.pid -o $WORKING_PATH/build/run.log -e $WORKING_PATH/build/err.log -c $WORKING_PATH/build $SERVICE_BINARY

    echo "Service started"
  else
    echo "Service already started"
  fi  
}

stop_native() {
  if [ -f $WORKING_PATH/build/run.pid ]; then
    kill $(cat $WORKING_PATH/build/run.pid)
    rm -f $WORKING_PATH/build/run.pid
    echo "Service stopped"
  else
    echo "Service not started"
  fi
}

build_proto() {
  protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative  proto/auth2.proto
}

case "$1" in
  proto)
    build_proto
    ;;
  build_native)
    build_native
    ;;
  build_linux)
    build_linux
    ;;
  remote_deploy)
    remote_deploy
    ;;
  run_native)
    run_native
    ;;
  stop_native)
    stop_native
    ;;  
  restart_native)
    stop_native
    run_native
    ;;  
  all)
    build_linux
    build_native
    remote_deploy
    ;;
  *)
    cat << EOF

 Sindel Auth Service tool

  Usage $0 [command]

    Available commands:
      proto          - reBuild proto3 structures
      build_linux    - Build linux binary 
      build_native   - Build native binary using native golang compiler
      remote_deploy  - Deploy linux binary to the production service
      all            - Run all of {build_linux,build_native,remote_deploy}
      run_native     - Run sindel_auth service on the host machine (daemonize required)
      stop_native    - Stop sindel_auth service on the host machine
      restart_native - Restarts sindel_auth service on the host machine
EOF
    exit 1

esac
