#
# Name: buildsys-ubuntu
#

FROM ubuntu

MAINTAINER pkgwrap

RUN DEBIAN_FRONTEND=noninteractive apt-get -y update
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install cmake pkg-config build-essential dh-make debhelper devscripts

RUN mkdir -p /opt/pkgwrap/bin

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install git-core sudo

RUN echo "ALL            ALL = (ALL) NOPASSWD: ALL" >> /etc/sudoers

# Go
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install golang mercurial
RUN echo "export GOPATH=~/gopath:\$GOPATH" >> /etc/skel/.profile
RUN echo "[ -e ~/gopath ] || mkdir -p ~/gopath/{src,bin,pkg}" >> /etc/skel/.profile

# Python
RUN apt-get -y install python-setuptools python-dev libpython-dev python-setuptools-git python-pip
