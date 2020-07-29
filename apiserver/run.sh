#!/bin/bash
export PORT=4000
export HOST=localhost
export TLSCERT=/Users/willhwang/go/src/github.com/will-slack/apiserver/fullchain.pem
export TLSKEY=/Users/willhwang/go/src/github.com/will-slack/apiserver/privkey.pem
export SESSIONKEY="secretkey"
export REDISADDR=localhost:6379
export DBADDR=localhost:27017