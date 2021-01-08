#!/bin/bash

CONTAINER_NAME=$2
CONTAINER_IMAGE="node-cmd-injection"

function start() {
  docker build --rm -t $CONTAINER_IMAGE .
  docker run --rm -v $(pwd)/app:/src/app -p 3000:3000 --name $CONTAINER_NAME $CONTAINER_IMAGE
}

function stop() {
  docker stop $CONTAINER_NAME && docker rm $CONTAINER_IMAGE
}

function restart() {
  stop
  start
}

function build() {
    docker build --rm -t $CONTAINER_IMAGE .
}

function_exists() {
  declare -f -F $1 > /dev/null
  return $?
}

if [ $# -lt 2 ]
then
  echo "Usage : $0 start|stop|restart|build <container name>"
  exit
fi

case "$1" in
  start)    function_exists start && start
          ;;
  stop)  function_exists stop && stop
          ;;
  restart)  function_exists restart && restart
          ;;
  build) function_exists build && build
          ;;
  *)      echo "Invalid command - Valid->start|stop|restart"
          ;;
esac