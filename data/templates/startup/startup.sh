#! /bin/sh

#
# {{.Name}}
#

# Used for status stop start. (more acurate than storing a pid)
pgrep="pgrep -f"

NAME={{.Name}}
BIN={{.RunnablePath}}
BIN_ARGS="{{.RunnableArgs}}"

LOGFILE={{.Logfile}}

CMD="${BIN} ${BIN_ARGS}"

RETVAL=0

status() {
    pids=`${pgrep} "${CMD}"`;
    if [ "$pids" == "" ]; then
        echo "Not running!";
    else
        echo "Running: $pids";
    fi
}

start() {
    $CMD > "$LOGFILE" 2>&1 &

    pids=`${pgrep} "${CMD}"`;
    if [ "$pids" == "" ]; then
        echo "Failed"
    else 
        echo "Started: $pids"
    fi
}

stop() {
    pids=`${pgrep} "${CMD}"`;
    for pid in $pids; do 
        kill $pid
    done

    pids=`${pgrep} "${CMD}"`;
    if [ "$pids" == "" ]; then
        echo "Stopped"
    else 
        echo "Failed: $pids"
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        start
        sleep 3
        stop
        ;;
    status)
        status
        ;;
    *)
        echo "Invalid option: '$1'!"
        RETVAL=1
        ;;
esac

exit $RETVAL
