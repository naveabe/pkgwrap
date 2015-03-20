angular.module('ipkg.project', [])
.controller('projectController', [ 
    '$scope', '$location', '$routeParams', 'Authenticator', 'PkgWrapRepo', 'SupportedVCs',
    function($scope, $location, $routeParams, Authenticator, PkgWrapRepo, SupportedVCs) {
        
        $scope.repository = $routeParams.repository;
        $scope.repositoryDetails = SupportedVCs.getDetails($scope.repository);
        
        Authenticator.checkAuthOrRedirect("/"+$routeParams.username+"/"+$routeParams.project,
            $scope.repositoryDetails);

        $scope.username = $routeParams.username;
        $scope.project = $routeParams.project;
        $scope.version = $routeParams.version;

        $scope.projectVersions = [];

        $scope.selectedVersion = {};

        $scope.loadVersionDistros = function(version) {
            $scope.selectedVersion = version;

            PkgWrapRepo.listPkgDistros({
                "repo": $scope.repository,
                "username": $scope.username,
                "project": $scope.project,
                "version": version
            },
            function(distros) {
                //console.log(distros);
                for(var i=0; i < $scope.projectVersions.length; i++) {
                    if($scope.projectVersions[i].version == version) {
                        
                        $scope.selectedVersion = $scope.projectVersions[i];
                        $scope.selectedVersion.distros = [];
                        for(var d=0; d< distros.length; d++) {
                            var bd = {name: distros[d]};
                            $scope.selectedVersion.distros.push(bd);
                            $scope.loadDistroPkgs(bd);
                        }
                        
                        break;
                    }
                }
            },
            function(err) { console.log(err) });
        }

        function pkgsSort(a, b) {
            if (a.mtime < b.mtime) return -1;
            if (a.mtime > b.mtime) return 1;
            return 0;
        }

        $scope.loadDistroPkgs = function(distro) {
            //$scope.selectedDistro = distro;

            PkgWrapRepo.listDistroPkgs({
                "repo": $scope.repository,
                "username": $scope.username,
                "project": $scope.project,
                "version": $scope.selectedVersion.version,
                "distro": distro.name
            },
            function(pkgs) {
                for(var i=0; i < $scope.selectedVersion.distros.length; i++) {
                    if($scope.selectedVersion.distros[i].name == distro.name) {
                        $scope.selectedVersion.distros[i].packages = pkgs.sort(pkgsSort);
                        break;
                    }
                }
            },
            function(err) { console.log(err) });   
        }

        function init() {

            PkgWrapRepo.listProjectVersions({
                "repo": $scope.repository,
                "username": $scope.username,
                "project": $scope.project
                },
                function(vlist) {
                    for( var i=0; i < vlist.length; i++ ) {
                        $scope.projectVersions.unshift({version: vlist[i]});
                    }
                    
                    if(!$routeParams.version || $routeParams.version == '') {
                        // Re-route to latest version
                        $location.url("/"+$scope.repository+"/"+$scope.username+"/"+
                            $scope.project+"/"+$scope.projectVersions[0].version);
                    } else {
                        // Set selected based on routeParams
                        for(var j=0 ; j < $scope.projectVersions.length; j++ ) {
                            if($scope.projectVersions[j].version === $scope.version) {
                                $scope.selectedVersion = $scope.projectVersions[j]; 
                                break;
                            }
                        }
                        $scope.loadVersionDistros($routeParams.version);
                    }
                }, function(err) { 
                    console.log(err); 
                });
        }

        init();
    }
]);