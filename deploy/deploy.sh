#!/bin/bash

eval "$(cat ./deploy/.env <(echo) <(declare -x))"

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
ssh -t ${USER_NAME}@${HOST} -p ${PORT} "sudo systemctl stop archibe.service"
scp -P ${PORT} archibe ${USER_NAME}@${HOST}:/home/${USER_NAME}/
ssh -t ${USER_NAME}@${HOST} -p ${PORT} "sudo systemctl restart archibe.service"
