#!/usr/bin/bash  

    export KAVKA_ENV=test  

    go run cmd/server/server.go &

	sleep 5

	pid=$(lsof -i :3000 | awk 'NR>1 {print $2}')


    go test ./tests/e2e/*

	sleep 1
	
    kill $pid 

