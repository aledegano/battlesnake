#!/bin/bash

timestamp=$(date +%s)
export image=degano/battlesnake-go:${timestamp}

docker build -t $image .
docker push $image

envsubst '$image' < deployment.yaml.tmpl | kubectl apply -f -
