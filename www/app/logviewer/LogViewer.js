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
            
            // Cache content elem
            var contentElem = elem.find('pre');
            // Cache collapsable
            var contentTrigger = contentElem.parent().parent();
            
            scope.logcontent = "";

            scope.toggleLog = function() {
                contentTrigger.collapse('toggle');
                if (scope.logcontent === "") {
                    //load content
                    LogLoader.getLog(scope.distro.id)
                    .success(function(data) {
                        scope.logcontent = data;
                    }).error(function(err) {
                        console.log(err);
                    });
                }
            }
            
            function onLogcontentChange(newVal, oldVal) {
                if(!newVal) return;
                contentElem.scrollTop(contentElem[0].scrollHeight-contentElem.height());
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