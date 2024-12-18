#!/bin/sh

cp /shared/area.apk /app/dist/

nohup npx vite preview --port 8081 --host
