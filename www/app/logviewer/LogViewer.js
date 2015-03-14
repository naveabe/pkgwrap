angular.module('ipkg.logviewer', [])
.directive('logTailer', [ '$rootScope', 'LogLoader', function($rootScope, LogLoader) {
    return {
        restrict: 'EA',
        scope: {
            distro: "=",
            state: "="
        },
        templateUrl: 'app/logviewer/log-viewer.html',
        link: function(scope, elem, attrs, ctrl) {
            scope.logcontent = "";
            // Cache content elem
            var contentElem = elem.find('pre');
            // Cache collapsable
            var contentTrigger = contentElem.parent().parent();
            // Scroll to the bottom on new log data
            var onLogcontentChange = function(newVal, oldVal) {
                if(!newVal) return;
                contentElem.scrollTop(contentElem[0].scrollHeight-contentElem.height());
            }
            /*
                Trigger histroy reload only when tailing the log.
                This only gets called when state == running.
            */
            var onLogTailComplete = function(evt) {
                $rootScope.$broadcast('build:history:changed', {});
            }

            var tailLog = function() {
                var oReq = new XMLHttpRequest();

                oReq.addEventListener("progress", function(evt) {
                    scope.$apply(function() { 
                        scope.logcontent = evt.target.responseText; 
                    });
                }, false);

                oReq.addEventListener("load", onLogTailComplete, false);
                oReq.addEventListener("error", function(e) { console.log('log error'); }, false);
                oReq.addEventListener("abort", function(e) { console.log('log cancelled'); }, false);
                oReq.open('GET', '/api/logs/' + scope.distro.id + '?follow');
                oReq.send();
            }


            

            scope.toggleLog = function() {
                contentTrigger.collapse('toggle');
                if (scope.logcontent === "") {
                    if ( scope.state.status == 'running' ) {
                        
                        // TODO; make initial regular call for current data
                        // tail log
                        console.log("tailing log...");
                        tailLog();
                    } else {
                        // load complete log 
                        LogLoader.getLog(scope.distro.id)
                            .success(function(data) { scope.logcontent = data;})
                            .error(function(err) { console.log(err); });
                    }
                }
            }

            function init() {
                scope.$watch(function() { return scope.logcontent },
                    onLogcontentChange, true);
            }

            init();
        }
    }
}])
.factory('LogLoader', ['$http', function($http) {
    return {
        getLog: function(containerId) {
            return $http({
                url: '/api/logs/' + containerId,
                method: 'GET',
                headers: { 'Content-Type': 'text/plain' },
                transformResponse: function(value) {
                    // Since this is plain text
                    return value;
                }
            });
        }
    }
}]);