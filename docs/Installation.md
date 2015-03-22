Installation
============

#### System Requirements
Although pkgwrap will run on any system, the following minimum requirements are recommended:

- 2 CPU/Cores (64 bit)
- 2 GB Memory

#### Stack Requirements
The following software technoligies are required:

- Elasticsearch >= 1.4
- A working go environmnt (to build this project)

#### Building Images
Before you can begin using the system, the docker containers need to be built.  This can be done running the following command:

    ./scripts/rebuild-images.sh

This will build all of the provided Dockerfiles.  Once this has successfully completed you can begin using the system.