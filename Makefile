
SHELL = /bin/bash

APP_HOME = /opt/pkgwrap

INSTALL_DIR = build/pkgwrap

.clean:
	rm -rf ./build
	go clean -i ./...

.deb_deps:
	apt-get install cmake pkg-config build-essential libgit2-0 libgit2-dev golang docker.io

# Build git2go w/ libgit2 next branch
.git2go:
	go get -d -u github.com/libgit2/git2go
	cd "../../libgit2/git2go" && git checkout next && git submodule update --init && make install

.deps: .git2go
	go get -d -v ./...

.build:
	go install -v ./...
	
	[ -d "$(INSTALL_DIR)" ] || mkdir -p $(INSTALL_DIR)
	
	[ -d "$(INSTALL_DIR)/usr/local/bin" ] || mkdir -p $(INSTALL_DIR)/usr/local/bin/
	cp $$GOPATH/bin/pkgwrap $(INSTALL_DIR)/usr/local/bin/
	
	[ -d "$(INSTALL_DIR)$(APP_HOME)/data" ] || mkdir -p "$(INSTALL_DIR)$(APP_HOME)/data/repository"
	cp -a data/bin "$(INSTALL_DIR)/$(APP_HOME)/data/"
	cp -a data/templates "$(INSTALL_DIR)$(APP_HOME)/data/"
	cp -a data/imagefiles "$(INSTALL_DIR)$(APP_HOME)/data/"
	cp -a scripts "$(INSTALL_DIR)$(APP_HOME)/"
	cp -a etc "$(INSTALL_DIR)/"
	cp -a www $(INSTALL_DIR)$(APP_HOME)

all: .clean .deps .build

