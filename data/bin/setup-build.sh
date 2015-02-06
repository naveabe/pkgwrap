#! /bin/bash
#
# Environment Variables:
#
#  REPO = github.com
#  BUILD_USER = 
#  BUILD_ENV = 
#  BUILD_CMD = 
#  BUILD_DEPS = 

PROJECT=$1
TAG=$2
if [ "$TAG" == "" ]; then
    TAG=master
    #VERSION="N/A"
    #else 
    #VERSION=`echo $TAG | sed -e "s/^[^0-9]*//g" -e "s/[^0-9]*$//g"`
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


# Setup build user
( id $BUILD_USER > /dev/null 2>&1 ) || useradd $BUILD_USER

#REPO_LOCAL_PATH="/opt/pkgbuilder/repo"
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
            echo "* No build environment selected using defaults!"        
            ;;
    esac
else
    echo "Running $BUILD_TYPE build..."
    su - $BUILD_USER -c "id > /dev/null";
    su - $BUILD_USER -c "cp -a $PROJECT_PATH $BUILD_HOME_DIR/rpmbuild/SOURCES/";
fi


echo "";
echo "   Project path: $PROJECT_PATH"
echo "";
