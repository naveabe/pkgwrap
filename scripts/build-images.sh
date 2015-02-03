#! /bin/bash

#
# Run from project home: i.e scripts/build-images.sh
#

DOCKER_CMD="docker -H 127.0.0.1:5555"

IMAGEFILES_DIR="data/imagefiles"
IMG_PREFIX="buildsys-"

[ -d "$IMAGEFILES_DIR" ] || {
    echo "Image files directory not found: $IMAGEFILES_DIR";
    exit 1;
}

build_images() {
    distro=$1
    cd "$IMAGEFILES_DIR/$distro";
    $DOCKER_CMD build -t ${IMG_PREFIX}${distro} . ;
    cd - ;
    for release in `ls $IMAGEFILES_DIR/${distro} | grep -v Dockerfile`; do
        cd "$IMAGEFILES_DIR/${distro}/$release" ;
        $DOCKER_CMD build -t ${IMG_PREFIX}${distro}:${release} . ;
        cd - ;
    done;
}

DISTRO=$1

case "$DISTRO" in
    ubuntu)
        build_images "ubuntu"
        ;;
    centos)
        build_images "centos"
        ;;
    *)
        build_images "ubuntu"
        build_images "centos"
        ;;
esac
