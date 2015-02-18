#! /bin/bash

#
# Remove all containers with an 'Exited' status
#
# Todo: 
#   The logic needs to be improved to include 

DOCKER_OPTS="-H 127.0.0.1:5555 -D"

CLEAN_TYPE=$1

clean_containers() {
    echo -e "\n- Cleaning up containers...\n"
    for cont in `docker $DOCKER_OPTS ps -a | grep Exited | awk '{print \$1 }'`; do
        docker $DOCKER_OPTS rm $cont 
    done
}

clean_images() {
    echo -e "\n- Cleaning up images...\n"
    for img in `docker $DOCKER_OPTS images --no-trunc | grep "<none>" | awk '{print \$3}'`; do
        docker $DOCKER_OPTS rmi $img
done
}

case "$CLEAN_TYPE" in
    images)
        clean_images
        ;;
    containers)
        clean_containers
        ;;
    *)
        clean_images
        clean_containers
        ;;
esac