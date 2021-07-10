#!/usr/bin/env bash

SCRIPT_PATH=$( cd "$(dirname "$0")" ; pwd -P )

function clean(){
  cd "${SCRIPT_PATH}"/../ && docker-compose  -f docker-compose.yaml down -v
}

function build(){
  cd "${SCRIPT_PATH}"/../ && docker-compose build
}

function up(){
  cd "${SCRIPT_PATH}"/../ && docker-compose up --force-recreate
}

function down(){
  cd "${SCRIPT_PATH}"/../ && docker-compose down -v --remove-orphans
}

for opt in "$@"
do
    case "$opt" in
        up)
            up
            ;;
        build)
            build
            ;;
        down)
            down
            ;;
        stop)
            down
            ;;
        start)
            up
            ;;
        clean)
            clean
            ;;
        restart)
            down
            clean
            up
            ;;
        *)
            echo $"Usage: $0 {up|down|build|start|stop|clean|restart}"
            exit 1

esac
done