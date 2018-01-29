#!/bin/bash

go build
./mushunew > /dev/null 2>&1 &
PID=$!
node js/test.js
kill $PID
