#!/bin/sh

APK_SOURCE="/shared/area.apk"

APK_DEST="/app/dist/area.apk"

if [ -f "$APK_SOURCE" ]; then
    cp "$APK_SOURCE" "$APK_DEST"
else
    echo "$APK_SOURCE not found. Skipping copy."
fi

exec "$@"