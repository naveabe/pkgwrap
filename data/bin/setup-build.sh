#! /bin/bash

#
# Common for both .deb and .rpm builds
#

PROJECT=$1
TAG=$2
if [ "$TAG" == "" ]; then
    TAG=master
fi

if [[ ( "$BUILD_USER" == "" ) || ( "$REPO" == "" ) ]]; then
    echo "BUILD_USER and/or REPO not specified!";
    exit 1;
fi

BUILD_HOME_DIR="/home/$BUILD_USER";

echo "";
echo " Environment:";
echo "";
echo "   Distro     : $PKG_DISTRO";
echo "   Package    : $PKG_TYPE";
echo "";
echo "   Project    : $PROJECT";
echo "   Tag        : $TAG";
echo "   Version    : $PKG_VERSION";
echo "   Release    : $PKG_RELEASE";
echo "";
echo "   Build Type : $BUILD_TYPE";
echo "";
echo "   Env        : $BUILD_ENV";
echo "   User       : $BUILD_USER";
echo "   Repo       : $REPO";
echo "";
echo "   Build Cmd  : $BUILD_CMD";
echo "   Build Deps : $BUILD_DEPS";
echo "   Deps       : $PKG_DEPS";
echo "";
echo "   Build Home : $BUILD_HOME_DIR";
echo "";


# Setup build user (-m needed by ubuntu to create home dir)
( id $BUILD_USER > /dev/null 2>&1 ) || useradd -m $BUILD_USER

REPO_LOCAL_PATH="/opt/pkgwrap/repo"

# Initial clone puts the project at the root of user homedir
PROJECT_PATH="$REPO_LOCAL_PATH/$PROJECT"

if [ "$BUILD_TYPE" == "source" ]; then
    # The first su - call initialized the environment.
    case "$BUILD_ENV" in
        go)
            su - $BUILD_USER -c "[ -e $BUILD_HOME_DIR/gopath/src/$REPO/$BUILD_USER ] || mkdir -p $BUILD_HOME_DIR/gopath/src/$REPO/$BUILD_USER"
            su - $BUILD_USER -c "cp -a $PROJECT_PATH $BUILD_HOME_DIR/gopath/src/$REPO/$BUILD_USER/" || exit 1
            PROJECT_PATH="$BUILD_HOME_DIR/gopath/src/$REPO/$BUILD_USER/$PROJECT"
            ;;
        *)
            su - $BUILD_USER -c "cp -a $PROJECT_PATH $BUILD_HOME_DIR/"
            PROJECT_PATH="$BUILD_HOME_DIR/$PROJECT"
            echo "    ** No build environment selected using defaults! **"
            ;;
    esac
fi

echo "";
echo "   Project path: $PROJECT_PATH"
echo "";
echo "Running $BUILD_TYPE build..."

# rhel does not immediately setup the user (first login)
su - $BUILD_USER -c "id > /dev/null";


##### Helper functions #####
install_deps() {
    pkg_mgr="$1"
    if [ "$BUILD_DEPS" != "" ]; then
        for pkg in $BUILD_DEPS; do 
            $pkg_mgr -y install "$pkg" || exit 1;
        done
    fi
}


add_pkg_to_repo() {
    base_dir="$1"
    find $base_dir -name "$PROJECT*.$PKG_TYPE" -exec cp -v '{}' $REPO_LOCAL_PATH/$PKG_DISTRO/ \; || {
        echo "** Failed to add package to repo! **"
        exit 7;
    }
cat <<EOF

  *
  * $PKG_TYPE successfully built!
  *

EOF

    echo -n $PKG_RELEASE > "$REPO_LOCAL_PATH/$PKG_DISTRO/RELEASE"
    echo "Release Updated!"
    # TODO: fire - added-to-repo event
}


# Install built package on the build system.
# i.e. test
install_built_pkg() {
    pkg_mgr="$1"

    for pkg in `ls $REPO_LOCAL_PATH/$PKG_DISTRO/ | egrep "^$PROJECT.*$PKG_VERSION-$PKG_RELEASE.*\.$PKG_TYPE"`; do
        echo "-> $pkg_mgr $REPO_LOCAL_PATH/$PKG_DISTRO/$pkg"
        $pkg_mgr $REPO_LOCAL_PATH/$PKG_DISTRO/$pkg;
    done

    echo "";
    echo "  ** DONE **"
    echo "";
    # TODO: fire - installed-built-pkg event
}

copy_startup() {
    DST="${1}/etc/init.d"

    if [ -e "$REPO_LOCAL_PATH/${PROJECT}.service" ]; then
        [ -d "$DST" ] || mkdir -p "${DST}"
        ( cp $REPO_LOCAL_PATH/${PROJECT}.service ${DST}/${PROJECT} && chmod 755 ${DST}/${PROJECT} ) || {
            echo "Failed to copy startup script: ${PROJECT}.service"
        }
    fi
}
