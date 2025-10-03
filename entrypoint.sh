#!/bin/bash
CMD=$1
shift

exec /app/"$CMD" $*
