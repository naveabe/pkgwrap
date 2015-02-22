
var app = angular.module('app', [
	'ngRoute',
    'ngResource',
	'appDirectives',
	'appControllers',
	'appServices',
    'ipkg.login',
	'ipkg.user',
    'ipkg.project',
    'ipkg.history',
    'ipkg.logviewer',
    'ipkg.builder'
]);

(function() {
	// Bootstrap the app with the config fetched via http //
	var configConstant = "Configuration";
	var configUrl      = "/conf/config.json";

    function fetchAndInjectConfig() {
        var initInjector = angular.injector(["ng"]);
        var $http = initInjector.get("$http");

        return $http.get(configUrl).then(function(response) {
            app.constant(configConstant, response.data);
        }, function(errorResponse) {
            // Handle error case
            console.log(errorResponse);
        });
    }

    function bootstrapApplication() {
        angular.element(document).ready(function() {
            angular.bootstrap(document, ["app"]);
        });
    }

    fetchAndInjectConfig().then(bootstrapApplication);
    
}());

app.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.when('/login', {
			templateUrl: 'app/login/login.html',
			controller: 'loginController'
		}).when('/:repository/:username/:project/:version/:distro', {
            templateUrl: 'app/project/project.html',
            controller: 'projectController'
        }).when('/:repository/:username/:project/:version', {
            templateUrl: 'app/project/project.html',
            controller: 'projectController'
        }).when('/:repository/:username/:project', {
			templateUrl: 'app/project/project.html',
			controller: 'projectController'
		}).when('/:repository/:username', {
            templateUrl: 'app/user/user.html',
            controller: 'userController'
        }).otherwise({
			redirectTo: '/login'
		});
	}
]);

app.filter('objectLength', function() {
	return function(obj) {
    	return Object.keys(obj).length;
	};
})
.filter('dotsToUnderscores', function() {
    return function(str) {
        return str.replace(/\./, '_');
    };
})
.filter('valueUnit', function() {
    return function(fileSize) {
        var kb  = fileSize/1024;
        if(kb < 1024) {
            return kb.toFixed(2).toString() +" KB";
        }

        var mb = kb/1024;
        if(mb < 1024) {
            return mb.toFixed(2).toString() +" MB";
        }

        return (mb/1024).toFixed(2).toString()+" GB";
    }
})
.filter('distroLabel', function() {
    /* Ignores everything after the 2nd dash */
    return function(dlabel) {
        var arr = dlabel.split("-")
        return arr[0].replace(/\w\S*/g, 
            function(txt) {
                return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
            }
        ) + " " + arr[1];

    }
});
