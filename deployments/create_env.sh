#!/bin/sh

if [ ! -f ./deployments/.env ]
then
    cp ./deployments/.env.example ./deployments/.env
fi