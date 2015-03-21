Usage
=====

The system can perform packaging of prebuilt data as well as compile then perform packaging.  The pre-built method can be particularly useful for arbitrary data and script/s.

### Pre-built packaging
To package pre-built data an archive must be provided containing a directory structure with the project name as the root containing data as it would be layed out on the filesystem.

Example directory structure:
    
    pkgwrap/
        /usr/local/bin/ppkgwrapd
        /etc/pkgwrap/pkgwrap.conf

To submit a build request you can submit a request as follows:

    curl -v -XPOST http://localhost:6565/api/builder/myproject/0.0.1 \
        -F package=@myproject.tgz 
        -F conf=@.pkgwrap.json

**package**:

File to upload.

This is a compressed tarball with files layed out as they would be on the filesystem.

**conf**:

Build configuration file in json format used to perform builds.


### Source packaging
To perform source packaging a configruation file **.pkgwrap.yml** must be included in the root of your project.  A sample can be found at the root of this project.  The project repository must be a git repository.

#### .pkgwrap.yml

This is the configruation file used to build the package. The available properties a shown below:

##### Distributions
This section contains distribution specific details to build package.

**name (required)**:

Name of the distribution to build package for

Example:

    centos

**release (required)**:

The release specific to the distribution.

Example:
    
    6

**deps**:

Dependencies needed to install the package.

**build_deps**:

Dependencies required to build the package.

**build_cmd (required)**:

This is the command/s to build the package. This is required for all source builds.  The result of this command should produce a directory ./build/project_name_here.  The build process looks for this directory which should contain files as they would be layed out on the filesystem.

**pre_install**:

Commands to run before installing package on the target system.

**post_install**:

Commands to run after installing package on the target system.

**pre_uninstall**:

Commands to run before uninstalling package on the target system.

**post_uninstall**:

Commands to run after uninstalling package on the target system.

##### Package
This section contains build information pertaining to the package.

**version**:

Version of the package being built.  This information can come from a a tag assuming it contains a version string i.e. x.x.x.  If this cannot be extrapolated from the tag it should be provided in the config.

Example:
    
    0.0.2


**url (required)**: 

URL to git repository

Example:

    https://github.com/naveable/pkgwrap


**packager**: 

This is the repo user.  It is also the user the package will be built with.

Example:

    naveabe  

**build_env**:

Environment to be used for build.

Example:

    go

**tagbranch**:

The tag or branch to use for the build.  If named appropriately, the version will automatically be extracted from the name.

Example:

    v0.0.1
    v0.0.3-dev


### Package Builds

Aside from builds being triggered from events or repository webhooks, manual builds can also be issued.


**Manual Builds**

Builds can be manually triggered by issuing the following command:

    curl -XPOST http://localhost:6565/api/builder/github.com/naveabe/pkgwrap/0.0.1 -d '{
        "Package": {
            "url": "https://github.com/naveabe/pkgwrap"
        }
    }'

Options for the body of the request are same ones mentioned above.

The endpoint format should be as follows:

    /api/builder/:repository/:username/:project/:version

Any configurations specified during a manual build will override the settings in the .pkgwrap.yml configuration file.

