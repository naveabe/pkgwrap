#
# Name: buildsys-centos:latest
#

FROM centos:7

MAINTAINER pkgwrap

RUN rpm -qa | grep epel || yum -y install "http://mirror.pnl.gov/epel/7/x86_64/e/epel-release-7-5.noarch.rpm"
RUN rpm -qa | grep rpmforge || yum -y install "http://pkgs.repoforge.org/rpmforge-release/rpmforge-release-0.5.3-1.el7.rf.x86_64.rpm"

RUN yum clean all
RUN yum -y update

RUN yum -y groupinstall "Development Tools"

RUN yum -y install which rpmdevtools git cmake pkgconfig tar
RUN echo -e '[ -e "~/rpmbuild" ] || rpmdev-setuptree\n' >> /etc/skel/.bashrc

RUN mkdir -p /opt/pkgwrap/repo; mkdir -p /opt/pkgwrap/bin;

# Go
RUN yum -y install golang mercurial
RUN echo -e 'export GOPATH=~/gopath\n[ -e "~/gopath" ] || mkdir -p ~/gopath/{src,bin,pkg}\n' >> /etc/skel/.bashrc

# Python
RUN yum -y install python-setuptools python-devel
RUN which pip ||  easy_install pip

# Extras
RUN yum -y install which sudo
RUN echo "ALL            ALL = (ALL) NOPASSWD: ALL" >> /etc/sudoers
RUN sed -i -e "/Defaults \s*requiretty/d" /etc/sudoers