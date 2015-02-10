#! /bin/bash

. /opt/pkgwrap/bin/setup-build.sh

install_deps() {
    apt-get update
    
    if [ "$BUILD_DEPS" != "" ]; then
        for pkg in $BUILD_DEPS; do 
            apt-get -y install "$pkg" || exit 1;
        done
    fi
}

install_deps;

if [ "$BUILD_TYPE" == "source" ]; then
    # Build source.
    if [ "$BUILD_CMD" != "" ]; then
        su - $BUILD_USER -c "cd $PROJECT_PATH && $BUILD_CMD" || exit 2
        
        su - $BUILD_USER -c "[ -d ~/debuild ] || mkdir ~/debuild"  || exit 3;
        # Copy package data to deb build env
        su - $BUILD_USER -c "cp -a $PROJECT_PATH/build/$PROJECT ~/debuild/" || exit 4
        su - $BUILD_USER -c "cp -a $PROJECT_PATH/build/$PROJECT ~/debuild/$PROJECT.orig" || exit 5
    else
        echo " ** No build command specified! **"
    fi
fi

su - $BUILD_USER -c "cp -a $REPO_LOCAL_PATH/$PKG_DISTRO/debian ~/debuild/$PROJECT/" || exit 6

su - $BUILD_USER -c "cd ~/debuild/$PROJECT && debuild -us -uc" || exit 7
## https://wiki.debian.org/IntroDebianPackaging

find $BUILD_HOME_DIR/debuild/ -name "$PROJECT*.deb" -exec cp -v '{}' $REPO_LOCAL_PATH/$PKG_DISTRO/ \; || exit 7

cat <<EOF
  
  *
  * DEB successfully built!
  *
EOF
echo -n $PKG_RELEASE > "$REPO_LOCAL_PATH/$PKG_DISTRO/RELEASE"