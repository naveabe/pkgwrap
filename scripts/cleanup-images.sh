#! /bin/bash

#
# Remove all containers with an 'Exited' status
#
# Todo: 
#   The logic needs to be improved to include 

DOCKER_OPTS="-H 127.0.0.1:5555 -D"

for cont in `docker $DOCKER_OPTS ps -a | grep Exited | awk '{print \$1 }'`; do
    docker $DOCKER_OPTS rm $cont 
done
