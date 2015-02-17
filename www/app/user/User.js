"use strict;"

angular.module('ipkg.user', [])
.controller('userController', [ 
    '$scope', '$routeParams', 'Authenticator', 'PkgWrapRepo', 'GithubRepo',
    function($scope, $routeParams, Authenticator, PkgWrapRepo, GithubRepo) {
        
        Authenticator.checkAuthOrRedirect("/"+$routeParams.username);
        
        $scope.pageHeaderHtml = "/partials/page-header.html";

        $scope.repository = $routeParams.repository;
        $scope.username = $routeParams.username;

        $scope.userRepos = [];

        function setActiveProjects(ghlist) {
            var out = [];
            
            for( var g=0; g < ghlist.length; g++ ) {
                
                var found = false;
                for( var p=0; p < $scope.userRepos.length; p++ ) {
                

                    if(ghlist[g].name === $scope.userRepos[p].name) {
                        ghlist[g].pkgwrapd = true;
                        found = true;
                        break;
                    }
                    //out.push(ghlist[g]);
                }

            }
            return out;
        }

        function loadGithubUserProjects() {
            GithubRepo.userRepos({
                "username": $scope.username
            }, 
            function(rslt) { 
                
                for( var g=0; g < rslt.length; g++ ) {
                    rslt.pkgwrapd = false;
                }
                $scope.userRepos = rslt;
            
                PkgWrapRepo.listUserProjects({
                    "repo": $scope.repository,
                    "username": $scope.username
                },
                function(rslt) { 
                    setActiveProjects(rslt);
                }, 
                function(err) { 
                    console.log(err); 
                });

            }, 
            function(err) { console.log(err); });
        }

        $scope.projectActivationChanged = function(usrRepo) {
            if(usrRepo.pkgwrapd === true) {
                console.log('Activate');
            } else {
                console.log('De-activate');
            }
        }

        function init() {
            switch($scope.repository) {
                case "github.com":
                    loadGithubUserProjects()
                    break;
                default:
                    break;
            }
        }

        init();
    }
]);