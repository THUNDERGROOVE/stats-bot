#!/bin/bash

git pull
cd ../census
git pull
cd stats-bot

go build
pkill stats-bot
nohup ./stats-bot&