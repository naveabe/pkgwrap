
var appDirectives = angular.module('appDirectives', [])

appDirectives.directive('appStage', ['$window', function($window) {
	return {
		restrict: 'A',
		require: '?ngModel',
		link: function(scope, elem, attrs, ctrl) {
			if(!ctrl) return;

			var stage = { jDom: $(elem) };
			stage.fillScreen = function() {

				if(stage.jDom) {
					
					stage.jDom.css("width", $window.innerWidth - stage.jDom.scrollWidth);
					stage.jDom.css("height", $window.innerHeight);
				}
			};

			function init() {
			
				stage.fillScreen();
				$window.addEventListener("resize", function(event) { $stage.fillScreen() });
			}

			init();
		}
	}
}]);