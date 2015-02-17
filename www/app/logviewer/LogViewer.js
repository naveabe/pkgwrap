angular.module('ipkg.logviewer', [])
.directive('logTailer', ['Configuration', function(Configuration) {
    return {
        restrict: 'A',
        require: '?ngModel',
        templateUrl: 'app/logviewer/log-viewer.html',
        link: function(scope, elem, attrs, ctrl) {
            if(!ctrl) return;

            var logContent;

/*
            var oReq = new XMLHttpRequest();

            function getLogUrl(id) { 
                console.log(Configuration.pkgwrap.url+'/api/jobs/' +
                    scope.username+'/'+scope.project+'/'+scope.version+'/'+id+'/log');
                return Configuration.pkgwrap.url+'/api/jobs/' +
                    scope.username+'/'+scope.project+'/'+scope.version+'/'+id+'/log';
            }

            function updateProgress(evt) {
                scope.$apply(function() {
                    ctrl.$modelValue.logContent += evt.target.responseText;
                });
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
*/
            function init() {
                scope.$watch(
                    function() { return ctrl.$modelValue.id },
                    function(newVal, oldVal) {
                        logContent = $('[data-log-id='+newVal+'] .log-content')[0];
                        //$('[data-log-id='+ctrl.$modelValue.id+']').collapse({toggle:false});
                        //console.log(logContent);
                    }, true);

                scope.$watch(
                    function() { return ctrl.$modelValue.logContent },
                    function(newVal, oldVal) {
                        //logContent. scroll down
                    }, true);
            }

            init();

        }
    }
}]);