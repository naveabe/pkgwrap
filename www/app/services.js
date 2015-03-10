angular.module('appServices', [])
.factory('PkgWrapRepo', ['$resource', 'Configuration', function($resource, Configuration) {
    return $resource('/api/repo/:repo/:username/:project/:version/:distro', {}, {
        listUserProjects: {
            params: {
                "repo": "@repo",
                "username": "@username"
            },
            method: 'GET',
            isArray: true
        },
        listProjectVersions: {
            params: {
                "repo": "@repo",
                "username": "@username",
                "project" : "@project"
            },
            method: 'GET',
            isArray: true
        },
        listPkgDistros: {
            params: {
                "repo": "@repo",
                "username": "@username",
                "project" : "@project",
                "version" : "@version"
            },
            method: 'GET',
            isArray: true
        },
        listDistroPkgs: {
            params: {
                "repo"    : "@repo",
                "username": "@username",
                "project" : "@project",
                "version" : "@version",
                "distro"  : "@distro"
            },
            method: 'GET',
            isArray: true    
        }
    });
}])
.factory('GithubRepo', ['$resource', function($resource) {
    return $resource('https://api.github.com/users/:username/:qtype', {}, {
        userRepos: {
            params: {"username": "@username", "qtype": "repos"},
            method: 'GET',
            isArray: true
        }
    });
}])
.factory('SupportedVCs', ['Configuration', function(Configuration) {
    
    var index = {};

    var supportedVCs = {};

    supportedVCs.getDetails = function(repo) {

        if (Object.keys(index).length < 1) {
            for( var i=0; i < Configuration.repos.length; i++ ) {
                index[Configuration.repos[i].repo] = Configuration.repos[i];
            }
        }
        return index[repo];
    }

    return supportedVCs;
}]);
/*
.factory('GitlabRepo', ['Configuration', function(Configuration) {
    
    var GitlabRepo = {};

    GitlabRepo.list = function() {
        var repoObj = null;
        for(var r=0; r < Configuration.repos.length; r++) {
            if ( Configuration.repos[r].type == "gitlab" ) {
                repoObj = Configuration.repos[r];
                break;
            }
        }
        if ( repoObj == null) {
            // error
            //return
        }

        return $http({
            method: 'GET',
            url: repoObj.url+'/api/v3/projects',
            headers: {"PRIVATE-TOKEN": ""}
        });
    }

    return GitlabRepo;
}])
*/