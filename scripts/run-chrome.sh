#!/bin/bash
Xvfb :99 -screen 0 1280x720x16 &> xvfb.log &
export DISPLAY=:99
/usr/bin/google-chrome-stable --disable-gpu --autoplay-policy=no-user-gesture-required --enable-logging=stderr  --no-sandbox --disable-infobars --kiosk --start-maximized --window-position=0,0 --window-size=1280,720 --app="https://www.youtube.com/watch?v=7QvXwp45TD8" &
sleep 30
rm /tmp/.X99-lock
rm xvfb.log
rm -rf ~/.config/google-chrome/Singleton*
exit 0
