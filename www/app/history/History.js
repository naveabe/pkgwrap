angular.module('ipkg.history', [])
.controller('historyController', [ 
    '$rootScope', '$scope', '$location', '$routeParams', '$http', 'BuildJobs',
    function($rootScope, $scope, $location, $routeParams, $http, BuildJobs) {

        $scope.historyHtml = "/app/history/history.html";

        $scope.buildHistory = [];
        /*
        function getLogUrl(id, follow) { 
           
            if(follow) {
                return '/api/jobs/' + $scope.repository + '/' +
                    $scope.username+'/'+$scope.project+'/'+$scope.version+'/'+id+'/log?follow=1';
            } else {
                return '/api/jobs/' + $scope.repository + '/' +
                    $scope.username+'/'+$scope.project+'/'+$scope.version+'/'+id+'/log';
            }
        }

        function transferComplete(evt) {
            console.log(evt);
        }
        function transferFailed(evt) {
            console.log(evt);
        }
        function transferCanceled(evt) {
            console.log(evt);
        }

        function tailLog(bJob) {
            var oReq = new XMLHttpRequest();

            oReq.addEventListener("progress", function(evt) {
                //console.log(evt.target.responseText);
                $scope.$apply(function() {bJob.logContent = evt.target.responseText;});
            }, false);
            oReq.addEventListener("load", transferComplete, false);
            oReq.addEventListener("error", transferFailed, false);
            oReq.addEventListener("abort", transferCanceled, false);
            oReq.open('GET', getLogUrl(bJob.id, true));
            oReq.send();
        }
        */

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
        /*
        $scope.refreshHistory = function() {
            fetchJobHistory();
        }

        $scope.loadBuildLog = function(bJob) {
            if(bJob.logContent && bJob.logContent != "" ) {
                return false;
            } else {
                bJob.logContent = "";
                
                if(bJob.status === 'started') {
                    tailLog(bJob);
                } else {
                    // Set loading
                    $http.get(getLogUrl(bJob.id))
                        .success(function(rslt) {
                            //$scope.$apply(function() {
                            bJob.logContent = rslt;
                            // Unset loading
                        })
                        .error(function(err) {
                            console.log(err);
                            // Unset loading
                        });
                }
            }
        }
        */
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