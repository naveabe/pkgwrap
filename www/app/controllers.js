
var appControllers = angular.module('appControllers', []);

appControllers.controller('rootController', [ 
	'$window', '$location', '$scope', 'Authenticator',
	function($window, $location, $scope, Authenticator) {
		
		//Authenticator.checkAuthOrRedirect("/");
	}
]);

appControllers.controller('defaultController', [ 
	'$window', '$location', '$routeParams', '$scope', '$rootScope', 'Authenticator',
	function($window, $location, $routeParams, $scope, $rootScope, Authenticator) {

		$scope.username = "";
		$scope.userIcon = ""
		$scope.isGuest = true;

		$scope.repository = "";

		$scope.userData = {};

		$scope.logout = Authenticator.logout;

	    var _resetUser = function() {
	    	$scope.username = "";
	    	$scope.userIcon = "";
			$scope.isGuest = true;
			$scope.repository = "";
			$scope.userData = {};
	    }

		var setUser = function(creds) {
			$scope.username = creds.username; 
	    	$scope.isGuest = false;
	    	$scope.userIcon = creds.icon;
	    	$scope.repository = creds.repo;
		}

	    function init() {
	    	var creds = Authenticator.getCreds();
	    	if ( creds ) {
	    		setUser(creds);
	    	}
	    	
			$rootScope.$on('user:auth:success', function(evt, data) {		
				setUser(data);
				if(data.repo) $scope.repository = data.repo;
				$scope.userData = data;
			});

			$rootScope.$on('user:unauth', function(evt, data) {
				_resetUser();
				$location.url("/");
			});
	    }

	    init();
	}
]);
