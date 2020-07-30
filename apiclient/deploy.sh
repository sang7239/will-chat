#!/bin/bash
docker run -d --name 344client -it -p 80:3000 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro sang7239/will-slack-client