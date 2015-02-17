angular.module('ipkg.logviewer', [])
.directive('logTailer', ['Configuration', function(Configuration) {
    return {
        restrict: 'A',
        require: '?ngModel',
        templateUrl: 'app/logviewer/log-viewer.html',
        link: function(scope, elem, attrs, ctrl) {
            if(!ctrl) return;

            var jElem = $(elem[0]);
            var contentElem;

            function init() {
                scope.$watch(
                    function() { return ctrl.$modelValue.id },
                    function(newVal, oldVal) {
                         contentElem = $(jElem.find('[data-log-content='+newVal+']')[0]);

                    }, true);

                scope.$watch(
                    function() { return ctrl.$modelValue.logContent },
                    function(newVal, oldVal) {
                        if(!newVal) return;

                        contentElem.scrollTop(contentElem[0].scrollHeight-contentElem.height());
                    }, true);
            }

            init();

        }
    }
}]);