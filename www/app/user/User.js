"use strict;"

angular.module('ipkg.user', [])
.controller('userController', [ 
    '$scope', '$routeParams', 'Authenticator', 'PkgWrapRepo', 'GithubPublic', 'SupportedVCs', 'Github', 'Gitlab',
    function($scope, $routeParams, Authenticator, PkgWrapRepo, GithubPublic, SupportedVCs, Github, Gitlab) {

        $scope.repository = $routeParams.repository;
        $scope.username = $routeParams.username;

        $scope.repositoryDetails = SupportedVCs.getDetails($scope.repository);

        Authenticator.checkAuthOrRedirect("/"+$routeParams.username,
                                            $scope.repositoryDetails);

        $scope.userProjects = [];
        $scope.orgMembership = [];

        var setActiveProjects = function(projList) {
            for( var p=0; p < projList.length; p++ ) {

                for( var g=0; g < $scope.userProjects.length; g++ ) {

                    if(projList[p] === $scope.userProjects[g].name) {
                        $scope.userProjects[g].pkgwrapd = true;
                        break;
                    }
                }
            }
        }

        var loadGithubUserProjects = function() {
            GithubPublic.userRepos({
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

        var loadPkgwrapProjects = function() {
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
                    console.log('No packages in repo:', err); 
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

        var loadUserOrgs = function() {
            Github.getUserOrgs().then(function(data) {
                $scope.orgMembership = data;
            }, function(err) {
                console.log(err);
            });
        }

        var loadGitlabProjects = function() {
            Gitlab.listProjects().success(function(data) {
                $scope.userProjects = data;
                console.log(data);
            }).error(function(err) {console.log(err); });
        }

        function init() {
            //console.log();
            switch($scope.repositoryDetails.type) {
            case "github":
                loadGithubUserProjects();
                loadUserOrgs();
                break;
            //case "gitlab":
            //    loadGitlabProjects();
            //    break;
            default:
                loadPkgwrapProjects();
                break;
            }
        }

        init();
    }
]);