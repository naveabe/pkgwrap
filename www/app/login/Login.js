angular.module('ipkg.login', [])
.factory('Authenticator', ['$window', '$http', '$location', function($window, $http, $location) {
    
    function _sessionIsAuthenticated() {
        if($window.sessionStorage['credentials']) {

            var creds = JSON.parse($window.sessionStorage['credentials']);
            if(creds.username && creds.username !== "" && creds.password && creds.password !== "") {
                // do custom checking here
                return true
            }
        }
        return false;
    }
    
    function _login(creds) {
        // do actual auth here //
        //if(creds.username === "guest" && creds.password === "guest") {
        if(creds.password === "guest") {
            $window.sessionStorage['credentials'] = JSON.stringify(creds);
            return true;
        }
        return false;
    }
    
    function _logout() {
        var sStor = $window.sessionStorage;
        if(sStor['credentials']) {
            delete sStor['credentials'];
        }
        $location.url("/login");
    }
    
    var Authenticator = {
        login                 : _login,
        logout                : _logout,
        sessionIsAuthenticated: _sessionIsAuthenticated,
        checkAuthOrRedirect   : function(redirectTo) {
            
            if(!_sessionIsAuthenticated()){
            //    if(redirectTo) $location.url("/login?redirect="+redirectTo);   
            //    else $location.url("/login");
            }
        }
    };

    return (Authenticator);
}])
.controller('loginController', [
    '$scope', '$window', '$routeParams', '$location', 'Authenticator', 'Configuration',
    function($scope, $window, $routeParams, $location, Authenticator, Configuration) {
        
        var defaultPage = "/";
        
        $scope.credentials = { username: "", password: "guest" };
        $scope.supportedRepos = Configuration.repos;
        $scope.selectedRepo = Configuration.repos[0];
        //$scope.selectedRepo

        $scope.attemptLogin = function() {
            if(Authenticator.login($scope.credentials)) {

                if($routeParams.redirect) {
                    $location.url($routeParams.redirect); 
                } else {
                    var newUrl = defaultPage+$scope.selectedRepo.repo+'/'+$scope.credentials.username;
                    //$location.url(defaultPage+$scope.credentials.username);
                    $location.url(newUrl);
                }
            } else {

                $("#login-window-header").html("<span>Auth failed!</span>");
            }
        }

        function _initialize() {
            if($window.sessionStorage['credentials']) {

                var creds = JSON.parse($window.sessionStorage['credentials']);
                if(creds.username && creds.username !== "" && creds.password && creds.password !== "") {

                    $scope.credentials = creds;
                    $scope.attemptLogin();
                }
            }
        }

        _initialize();
    }
]);