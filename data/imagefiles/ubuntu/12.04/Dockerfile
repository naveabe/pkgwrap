#
# Name: buildsys-ubuntu:12.04
#

FROM ubuntu:12.04

MAINTAINER pkgwrap

RUN apt-get -y update
RUN apt-get -y install cmake pkg-config build-essential dh-make debhelper devscripts

RUN mkdir -p /opt/pkgwrap/bin

RUN apt-get -y install git-core

RUN echo "ALL            ALL = (ALL) NOPASSWD: ALL" >> /etc/sudoers
RUN sed -i -e "s/SHELL=.*/SHELL=\/bin\/bash/" /etc/default/useradd

# Go
RUN apt-get -y install golang mercurial
RUN echo "export GOPATH=~/gopath:\$GOPATH" >> /etc/skel/.profile
RUN echo "[ -e ~/gopath ] || mkdir -p ~/gopath/{src,bin,pkg}" >> /etc/skel/.profile
RUN echo "export GOPATH=~/gopath:\$GOPATH" >> /etc/skel/.bashrc
RUN echo "[ -e ~/gopath ] || mkdir -p ~/gopath/{src,bin,pkg}" >> /etc/skel/.bashrc

# Python
RUN apt-get -y install python-setuptools python-dev python-pip
