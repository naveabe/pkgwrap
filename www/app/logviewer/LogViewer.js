angular.module('ipkg.logviewer', [])
.directive('logTailer', ['Configuration', function(Configuration) {
    return {
        restrict: 'A',
        require: '?ngModel',
        templateUrl: 'app/logviewer/log-viewer.html',
        link: function(scope, elem, attrs, ctrl) {
            if(!ctrl) return;
            /*
            var logContent;

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
            */
        }
    }
}]);