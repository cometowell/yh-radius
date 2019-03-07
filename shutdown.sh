# !/bin/bash

PID_FILE="$PWD/rad.pid"
echo $PID_FILE

if [ -f "$PID_FILE" ]; then
    PID=$(cat $PID_FILE)
    echo "kill radius $PID"
    kill $PID
fi