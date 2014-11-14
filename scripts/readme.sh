#!/bin/sh

BIN=go-readme
PIDFILE=/tmp/readme.pid

case $1 in
    stop)
        if ! [ -f $PIDFILE ]; then
            echo "not running"
            exit 1
        fi
        kill -9 `cat $PIDFILE` && rm -f $PIDFILE
        echo "stopped"
        ;;
    *)
        if [ -f $PIDFILE ]; then
            echo "already running"
            exit 1
        fi
        $BIN &
        echo $! > $PIDFILE
        ;;
esac

exit 0
