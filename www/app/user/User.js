"use strict;"

angular.module('ipkg.user', [])
.controller('userController', [ 
    '$scope', '$routeParams', 'Authenticator', 'PkgWrapRepo', 'GithubRepo', 'SupportedVCs',
    function($scope, $routeParams, Authenticator, PkgWrapRepo, GithubRepo, SupportedVCs) {

        $scope.repository = $routeParams.repository;
        $scope.username = $routeParams.username;

        $scope.repositoryDetails = SupportedVCs.getDetails($scope.repository);

        Authenticator.checkAuthOrRedirect("/"+$routeParams.username,
                                            $scope.repositoryDetails);

        $scope.userProjects = [];


        function setActiveProjects(projList) {
            for( var p=0; p < projList.length; p++ ) {

                for( var g=0; g < $scope.userProjects.length; g++ ) {

                    if(projList[p] === $scope.userProjects[g].name) {
                        $scope.userProjects[g].pkgwrapd = true;
                        break;
                    }
                }
            }
        }

        function loadGithubUserProjects() {
            GithubRepo.userRepos({
                "username": $scope.username
            }, 
            function(rslt) { 
                
                for( var g=0; g < rslt.length; g++ ) {
                    rslt.pkgwrapd = false;
                }
                $scope.userProjects = rslt;
            
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

        function loadPkgwrapProjects() {
            PkgWrapRepo.listUserProjects(
                {
                    "repo": $scope.repository,
                    "username": $scope.username
                }, function(rslt) { 
                    for(var i=0; i< rslt.length; i++) {
                        rslt[i] = {name: rslt[i]};
                    }
                    $scope.userProjects = rslt;
                }, function(err) { 
                    console.log(err); 
                }
            );
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
                    loadPkgwrapProjects();
                    break;
            }
        }

        init();
    }
]);