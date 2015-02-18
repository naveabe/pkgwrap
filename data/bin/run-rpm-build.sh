#! /bin/bash
#
# build/$name contains vfs tarball
#
. /opt/pkgwrap/bin/setup-build.sh

yum -y update
install_deps "yum"


# Copy spec file from repo
su - $BUILD_USER -c "cp $REPO_LOCAL_PATH/$PKG_DISTRO/$PROJECT.spec ~/rpmbuild/SPECS/" || exit 2

if [ "$BUILD_TYPE" == "source" ]; then
    # Build source.
    if [ "$BUILD_CMD" != "" ]; then
        su - $BUILD_USER -c "cd $PROJECT_PATH && $BUILD_CMD" || exit 3
        # Copy package data to rpm SOURCES destination
        su - $BUILD_USER -c "cp -a $PROJECT_PATH/build/$PROJECT ~/rpmbuild/SOURCES/" || exit 4
        
        copy_startup "$BUILD_HOME_DIR/rpmbuild/SOURCES/$PROJECT"
        
        # Write file list to spec when being built.
        su - $BUILD_USER -c "cd ~/rpmbuild/SOURCES/$PROJECT && ( find . -type f | sed s/^\.//g >> ~/rpmbuild/SPECS/$PROJECT.spec ) && cd -" || exit 5
        # Copy updated spec back to repo after file list update.
        cp $BUILD_HOME_DIR/rpmbuild/SPECS/$PROJECT.spec $REPO_LOCAL_PATH/$PKG_DISTRO/
    else
        echo " ** No build command specified! **"
    fi
else
  # Binary
  copy_startup "$BUILD_HOME_DIR/rpmbuild/SOURCES/$PROJECT"
  
  su - $BUILD_USER -c "cp -a $PROJECT_PATH $BUILD_HOME_DIR/rpmbuild/SOURCES/";
fi

# Build spec
# QA_RPATHS=$[ 0x0001|0x0010 ] : Ignore check-rpath warning 
su - $BUILD_USER -c "QA_RPATHS=$[ 0x0001|0x0010 ] rpmbuild -ba ~/rpmbuild/SPECS/$PROJECT.spec" || exit 6

# Copy RPM back to repo
add_pkg_to_repo "$BUILD_HOME_DIR/rpmbuild/RPMS/"
# Install produced rpm
install_built_pkg "yum -y install"
