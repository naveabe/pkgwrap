#! /bin/bash

. /opt/pkgwrap/bin/setup-build.sh

#
# build/$name contains vfs tarball
#

# Install deps as build happens outside of .spec file.
install_deps() {
    yum -y update
    
    if [ "$BUILD_DEPS" != "" ]; then
        for pkg in $BUILD_DEPS; do 
            yum -y install "$pkg" || exit 1;
        done
    fi
}

install_deps;

# Copy spec file from repo
su - $BUILD_USER -c "cp $REPO_LOCAL_PATH/$PKG_DISTRO/$PROJECT.spec ~/rpmbuild/SPECS/" || exit 2

if [ "$BUILD_TYPE" == "source" ]; then
    # Build source.
    if [ "$BUILD_CMD" != "" ]; then
        su - $BUILD_USER -c "cd $PROJECT_PATH && $BUILD_CMD" || exit 3
        # Copy package data to rpm SOURCES destination
        su - $BUILD_USER -c "cp -a $PROJECT_PATH/build/$PROJECT ~/rpmbuild/SOURCES/" || exit 4
        # Write file list to spec when being built.
        su - $BUILD_USER -c "cd ~/rpmbuild/SOURCES/$PROJECT && ( find . -type f | sed s/^\.//g >> ~/rpmbuild/SPECS/$PROJECT.spec ) && cd -" || exit 5
        # Copy updated spec back to repo after file list update.
        cp $BUILD_HOME_DIR/rpmbuild/SPECS/$PROJECT.spec $REPO_LOCAL_PATH/$PKG_DISTRO/
    else
        echo " ** No build command specified! **"
    fi
    
    #else
    #su - $BUILD_USER -c "cp -a $PROJECT_PATH ~/rpmbuild/SOURCES/" || exit 3
fi

# Build spec
# QA_RPATHS=$[ 0x0001|0x0010 ] : Ignore check-rpath warning 
su - $BUILD_USER -c "QA_RPATHS=$[ 0x0001|0x0010 ] rpmbuild -ba ~/rpmbuild/SPECS/$PROJECT.spec" || exit 6

# Copy RPM back to repo
#[ -d "$REPO_LOCAL_PATH/$PKG_DISTRO" ] || mkdir -p "$REPO_LOCAL_PATH/$PKG_DISTRO"
find $BUILD_HOME_DIR/rpmbuild/RPMS/ -name "$PROJECT*.rpm" -exec cp -v '{}' $REPO_LOCAL_PATH/$PKG_DISTRO/ \; || exit 7
cat <<EOF
  
  *
  * RPM successfully built!
  *
  * Attempting to install RPM...
  *
EOF
echo -n $PKG_RELEASE > "$REPO_LOCAL_PATH/$PKG_DISTRO/RELEASE"

# Install produced rpm
for pkg in `ls $REPO_LOCAL_PATH/$PKG_DISTRO/ | egrep "^$PROJECT-$PKG_VERSION-$PKG_RELEASE.*\.rpm"`; do
    yum -y install $REPO_LOCAL_PATH/$PKG_DISTRO/$pkg;
done

echo "";
echo "  ** DONE **"
echo "";