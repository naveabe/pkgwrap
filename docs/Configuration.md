Configuration
=============

#### Configuration Files
There are 2 configuration files - for the backend and frontend:

- etc/pkgwrap/pkgwrap.conf.json
- www/conf/conf.json

#### Building Images
Before you can begin using the system, the docker containers need to be built.  This can be done running the following command:

    ./scripts/rebuild-images.sh

This will build all of the provided Dockerfiles.  Once this has successfully completed you can begin using the system.