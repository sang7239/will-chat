#!/bin/bash
export TLSKEY=/etc/letsencrypt/live/api.will-hwang.me/privkey.pem
export TLSCERT=/etc/letsencrypt/live/api.will-hwang.me/fullchain.pem
export DBADDR=user-store:27017
export REDISADDR=sessions-store:6379

docker network create --driver bridge api-server-net
docker run -d --name user-store --network api-server-net mongo
docker run -d --name sessions-store --network api-server-net redis

docker run -d --name 344gateway --network api-server-net -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e SESSIONKEY="secretkey" -e REDISADDR=$REDISADDR -e DBADDR=$DBADDR sang7239/will-slack-apiserver
