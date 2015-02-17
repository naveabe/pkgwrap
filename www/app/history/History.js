angular.module('ipkg.history', [])
.controller('historyController', [ 
    '$scope', '$location', '$routeParams', '$http', 'Authenticator', 'PkgWrapRepo', 'PkgWrapJobs',
    function($scope, $location, $routeParams, $http, Authenticator, PkgWrapRepo, PkgWrapJobs) {

        $scope.historyHtml = "/app/history/history.html";

        $scope.username = $routeParams.username;
        $scope.project = $routeParams.project;
        $scope.buildHistory = [];

        function getLogUrl(id, follow) { 
           
            if(follow) {
                return '/api/jobs/' +
                    $scope.username+'/'+$scope.project+'/'+$scope.version+'/'+id+'/log?follow=1';
            } else {
                return '/api/jobs/' +
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

        function init() {
            PkgWrapJobs.listJobsForProject({
                username: $scope.username,
                project: $scope.project
            },
            function(rslt) {
                $scope.buildHistory = rslt;
            },
            function(err) { console.log(err); });
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

        init();
    }
]);