#! /bin/bash
#
# https://wiki.debian.org/IntroDebianPackaging
#
. /opt/pkgwrap/bin/setup-build.sh

apt-get update
install_deps "apt-get";


su - $BUILD_USER -c "[ -d ~/debuild ] || mkdir -p ~/debuild" || exit 2;

su - $BUILD_USER -c "cp -a $REPO_LOCAL_PATH/$PKG_DISTRO/debian/debian-binary ~/debuild/" || exit 2
su - $BUILD_USER -c "cp -a $REPO_LOCAL_PATH/$PKG_DISTRO/debian/control.tar.gz ~/debuild/" || exit 2;

if [ "$BUILD_TYPE" == "source" ]; then
    # Build source.
    if [ "$BUILD_CMD" != "" ]; then
        # User build command
        su - $BUILD_USER -c "cd $PROJECT_PATH && $BUILD_CMD" || exit 3
        # TODO: fire - user-build-success event    
        
        copy_startup "$PROJECT_PATH/build/$PROJECT"

        # Create data tarball
        su - $BUILD_USER -c "cd ~/debuild && tar zcvf data.tar.gz -C $PROJECT_PATH/build/$PROJECT ." || exit 4; 
    else
        echo " ** WARNING: No build command specified! **"
    fi
else
    # Binary (pre-compiled)
    copy_startup "$PROJECT_PATH"

    su - $BUILD_USER -c "cd ~/debuild && tar czvf data.tar.gz -C $PROJECT_PATH ." || exit 3;
fi

# Make .deb (i.e. ar -r ...)
su - $BUILD_USER -c "cd ~/debuild && ar -r ${PROJECT}_${PKG_VERSION}-${PKG_RELEASE}_amd64.deb debian-binary control.tar.gz data.tar.gz" || exit 5;

# Copy .deb back to repo
add_pkg_to_repo "$BUILD_HOME_DIR/debuild/"
# Install build pkg (i.e. test)
install_built_pkg "dpkg -i"