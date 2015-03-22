
SHELL = /bin/bash

APP_HOME = /opt/pkgwrap

INSTALL_DIR = build/pkgwrap

.clean:
	rm -rf ./build
	go clean -i ./...

# Build git2go w/ libgit2 next branch
.git2go:
	go get -d -u github.com/libgit2/git2go
	cd "../../libgit2/git2go" && git checkout next && git submodule update --init && make install

# Go deps
.deps: .git2go
	go get -d -v ./...

.build:
	go install -v ./...
	
	[ -d "$(INSTALL_DIR)" ] || mkdir -p $(INSTALL_DIR)
	
	[ -d "$(INSTALL_DIR)/usr/local/bin" ] || mkdir -p $(INSTALL_DIR)/usr/local/bin/
	cp ../../../../bin/pkgwrap $(INSTALL_DIR)/usr/local/bin/
	
	[ -d "$(INSTALL_DIR)$(APP_HOME)/data" ] || mkdir -p "$(INSTALL_DIR)$(APP_HOME)/data/repository"
	# Copy data files
	cp -a data/bin "$(INSTALL_DIR)/$(APP_HOME)/data/"
	cp -a data/templates "$(INSTALL_DIR)$(APP_HOME)/data/"
	cp -a data/imagefiles "$(INSTALL_DIR)$(APP_HOME)/data/"
	cp -a scripts "$(INSTALL_DIR)$(APP_HOME)/"
	cp -a etc "$(INSTALL_DIR)/"
	cp -a www $(INSTALL_DIR)$(APP_HOME)

all: .clean .deps .build

