#
# Name: buildsys-ubuntu:14.04
#

FROM ubuntu:14.04

MAINTAINER pkgwrap

RUN apt-get -y update
RUN apt-get -y install cmake pkg-config build-essential dh-make debhelper devscripts

RUN mkdir -p /opt/pkgwrap/bin

RUN apt-get -y install git-core

RUN echo "ALL            ALL = (ALL) NOPASSWD: ALL" >> /etc/sudoers

# Go
RUN apt-get -y install golang mercurial
RUN echo "export GOPATH=~/gopath:\$GOPATH" >> /etc/skel/.profile
RUN echo "[ -e ~/gopath ] || mkdir -p ~/gopath/{src,bin,pkg}" >> /etc/skel/.profile
RUN echo "export GOPATH=~/gopath:\$GOPATH" >> /etc/skel/.bashrc
RUN echo "[ -e ~/gopath ] || mkdir -p ~/gopath/{src,bin,pkg}" >> /etc/skel/.bashrc

# Python
RUN apt-get -y install python-setuptools python-dev libpython-dev python-setuptools-git python-pip
