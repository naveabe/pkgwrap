
var appControllers = angular.module('appControllers', []);

appControllers.controller('rootController', [ '$window', '$location', '$scope', 'Authenticator',
	function($window, $location, $scope, Authenticator) {
		
		//Authenticator.checkAuthOrRedirect("/");
	}
]);

appControllers.controller('defaultController', [ 
	'$window', '$location', '$routeParams', '$scope', '$rootScope', 'Authenticator',
	function($window, $location, $routeParams, $scope, $rootScope, Authenticator) {

		$scope.authedUser = "";
		$scope.isGuest = true;

		$scope.repository = "";

		$scope.logout = function() {
			Authenticator.logout();
	    }

	    function init() {

	    	//console.log($routeParams);
			
			$rootScope.$on('user:auth:success', function(evt, data) {
				
				$scope.authedUser = data.username;
				$scope.isGuest = false;

				$scope.repository = data.repo;
			});

			$rootScope.$on('user:unauth', function(evt, data) {
				$scope.authedUser = "";
				$scope.isGuest = true;
				$scope.repository = "";
			});
	    }

	    init();
	}
]);
