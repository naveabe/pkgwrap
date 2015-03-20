angular.module('ipkg.history', [])
.controller('historyController', [ 
    '$rootScope', '$scope', '$location', '$routeParams', '$http', 'BuildJobs',
    function($rootScope, $scope, $location, $routeParams, $http, BuildJobs) {

        $scope.historyHtml = "/app/history/history.html";

        $scope.buildHistory = [];

        /*
            Determine job status from container state structure.
        */
        function setJobStatus(jobHistory) {
            for(var i=0; i < jobHistory.length; i++) {

                for(var j in jobHistory[i].Containers) {
                    var state = jobHistory[i].Containers[j].State;
                    if(state.Running) {
                        state.status = "running";
                    } else if(state.ExitCode && state.ExitCode != 0) {
                        state.status = "failed";
                    } else {
                        state.status = "succeeded";
                    }
                }
            }
            return jobHistory;
        }

        $scope.fetchJobHistory = function() {
        
            //console.log($routeParams);
            if (Object.keys($routeParams).length <= 0) {
                
                BuildJobs.list(function(rslt) {
                    $scope.buildHistory = setJobStatus(rslt); 
                },
                function(err) { console.log(err); });
            } else {

                if($routeParams.version) {

                    BuildJobs.listForRepoUserProjectVersions({
                        "repository": $routeParams.repository,
                        "username"  : $routeParams.username,
                        "project"   : $routeParams.project,
                        "version"   : $routeParams.version
                    },
                    function(rslt) { $scope.buildHistory = setJobStatus(rslt); },
                    function(err)  { console.log(err); });
                
                } else if ($routeParams.project) {

                    BuildJobs.listForRepoUserProject({
                        "repository": $routeParams.repository,
                        "username"  : $routeParams.username,
                        "project"   : $routeParams.project
                    },
                    function(rslt) { $scope.buildHistory = setJobStatus(rslt); },
                    function(err)  { console.log(err); });
                
                } else if ($routeParams.username) {
                
                    BuildJobs.listForRepoUserProject({
                        "repository": $routeParams.repository,
                        "username"  : $routeParams.username
                    },
                    function(rslt) { $scope.buildHistory = setJobStatus(rslt); },
                    function(err)  { console.log(err); });

                }
            }
        }

        function init() {
            
            $scope.fetchJobHistory();
            
            $rootScope.$on('build:history:changed', function(evt, data) {
                // Account for registration
                setTimeout($scope.fetchJobHistory, 2500);
            });
        }

        init();
    }
])
.factory('BuildJobs', ['$resource', 'Configuration', function($resource, Configuration) {
    return $resource('/api/builder/:repository/:username/:project/:version', {}, {
        list: {
            method: 'GET',
            isArray: true
        },
        listForRepoUser: {
            params: {
                "repository": "@repository",
                "username": "@username"
            },
            method: 'GET',
            isArray: true
        },
        listForRepoUserProject: {
            params: {
                "repository": "@repository",
                "username": "@username",
                "project" : "@project"
            },
            method: 'GET',
            isArray: true
        },
        listForRepoUserProjectVersions: {
            params: {
                "repository": "@repository",
                "username": "@username",
                "project" : "@project",
                "version" : "@version"
            },
            method: 'GET',
            isArray: true
        }
    });
}]);