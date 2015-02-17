angular.module('ipkg.user', [])
.controller('userController', [ 
    '$scope', '$routeParams', 'Authenticator', 'GithubRepo', 'PkgWrapRepo',
    function($scope, $routeParams, Authenticator, GithubRepo, PkgWrapRepo) {
        
        Authenticator.checkAuthOrRedirect("/"+$routeParams.username);
        
        $scope.pageHeaderHtml = "/partials/page-header.html";

        $scope.repository = $routeParams.repository;
        $scope.username = $routeParams.username;

        $scope.userRepos = [];
        $scope.userOrgs = [];


        function setActiveProjects(projList) {
            for( var p=0; p < projList.length; p++ ) {

                for( var g=0; g < $scope.userRepos.length; g++ ) {

                    if(projList[p] === $scope.userRepos[g].name) {
                        $scope.userRepos[g].pkgwrapd = true;
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
            if($scope.repository == 'github.com') {
                loadGithubUserProjects();
            }
        }
        init();
    }
]);