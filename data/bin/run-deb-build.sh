#! /bin/bash
#
# https://wiki.debian.org/IntroDebianPackaging
#
. /opt/pkgwrap/bin/setup-build.sh

install_deps() {
    apt-get update
    
    if [ "$BUILD_DEPS" != "" ]; then
        for pkg in $BUILD_DEPS; do 
            apt-get -y install "$pkg" || exit 1;
        done
    fi
}
apt-get -y install tree
install_deps;

su - $BUILD_USER -c "[ -d ~/debuild/$PROJECT-$PKG_VERSION ] || mkdir -p ~/debuild/$PROJECT-$PKG_VERSION" || exit 3;
su - $BUILD_USER -c "cp -a $REPO_LOCAL_PATH/$PKG_DISTRO/debian ~/debuild/$PROJECT-$PKG_VERSION/" || exit 6

su - $BUILD_USER -c "ls -lah ~/debuild/"

if [ "$BUILD_TYPE" == "source" ]; then
    # Build source.
    if [ "$BUILD_CMD" != "" ]; then
        su - $BUILD_USER -c "cp -a $PROJECT_PATH ~/debuild/$PROJECT-$PKG_VERSION.orig" || exit 2
        su - $BUILD_USER -c "cd $PROJECT_PATH && $BUILD_CMD" || exit 2
        
        # Copy package data to deb build env
        su - $BUILD_USER -c "cp -a $PROJECT_PATH/build ~/debuild/$PROJECT-$PKG_VERSION/" || exit 4
        
    else
        echo " ** No build command specified! **"
    fi
fi

su - $BUILD_USER -c "cd ~/debuild/$PROJECT-$PKG_VERSION && debuild -us -uc" || exit 7

find $BUILD_HOME_DIR/debuild/ -name "$PROJECT*.deb" -exec cp -v '{}' $REPO_LOCAL_PATH/$PKG_DISTRO/ \; || exit 7
cat <<EOF
  
  *
  * DEB successfully built!
  *
EOF
echo -n $PKG_RELEASE > "$REPO_LOCAL_PATH/$PKG_DISTRO/RELEASE"