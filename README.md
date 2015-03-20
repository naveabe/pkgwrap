pkgwrap
=======
A simplified docker based package builder for binary (pre-built) as well as source data.  

pkgwrap is system specific package builder.  It is meant to be integrated with other continous integration platforms.  Packages are built in docker containers based on the requested distributions and versions.

### Features:

- Build .rpm and .deb packages from pre-built binarys, arbitrarty data and/or scripts as well as building packages from source.  
- Packages built in their own distro specific container.
- Aims to integrates with other CI systems.


### Supported Distributions:

- CentOS
    - 6
    - 7

- Ubuntu
    - 12.04
    - 14.04

Many many more distro's to come...

Suggestions, ideas, pull requests etc. are all always welcomed.