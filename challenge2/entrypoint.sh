#!/bin/bash
while true; do
  java -XX:+ExitOnOutOfMemoryError \
       --add-opens java.base/java.lang=ALL-UNNAMED \
       --add-opens java.base/java.util=ALL-UNNAMED \
       --add-opens java.base/java.io=ALL-UNNAMED \
       --add-opens java.base/java.lang.reflect=ALL-UNNAMED \
       -jar /app/app.jar
  
  EXIT_CODE=$?
  case $EXIT_CODE in
    0)  echo "Clean shutdown."; exit 0 ;;
    42) echo "Manual restart triggered, restarting JVM..." ;;
    3)  echo "OOM detected, restarting JVM..." ;;
    *)  echo "Unexpected exit $EXIT_CODE, stopping."; exit $EXIT_CODE ;;
  esac
done