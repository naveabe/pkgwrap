angular.module('ipkg.login', [])
.factory('Authenticator', [
    '$window', '$http', '$location', '$rootScope',
    function($window, $http, $location, $rootScope) {
    
        function _sessionIsAuthenticated(evtArgs) {
            if($window.sessionStorage['credentials']) {

                var creds = JSON.parse($window.sessionStorage['credentials']);
                if(creds.username && creds.username !== "" && creds.password && creds.password !== "") {
                    var eArgs = $.extend({}, {username:creds.username}, evtArgs);
                    // do custom checking here
                    $rootScope.$emit('user:auth:success', eArgs);
                    return true
                }
            }
            return false;
        }
        
        function _login(creds, args) {
            // do actual auth here //
            //if(creds.username === "guest" && creds.password === "guest") {
            if(creds.password === "guest") {
                $window.sessionStorage['credentials'] = JSON.stringify(creds);
                //console.log(args);
                var eArgs = $.extend({}, {username:creds.username}, args);
                $rootScope.$emit('user:auth:success', eArgs);
                
                return true;
            }
            return false;
        }
        
        function _logout() {
            var sStor = $window.sessionStorage;
            if(sStor['credentials']) {
                delete sStor['credentials'];
            }
            $rootScope.$emit('user:unauth', {});
            $location.url("/login");
        }
        
        var Authenticator = {
            login                 : _login,
            logout                : _logout,
            sessionIsAuthenticated: _sessionIsAuthenticated,
            checkAuthOrRedirect   : function(redirectTo, evtArgs) {
                
                if( !_sessionIsAuthenticated(evtArgs) ){
                //    if(redirectTo) $location.url("/login?redirect="+redirectTo);   
                //    else $location.url("/login");
                }
            }
        };

        return (Authenticator);
    }
])
.controller('loginController', [
    '$scope', '$window', '$routeParams', '$location', 'Authenticator', 'Configuration',
    function($scope, $window, $routeParams, $location, Authenticator, Configuration) {
        
        var defaultPage = "/";
        
        $scope.credentials = { username: "", password: "guest" };
        $scope.supportedRepos = Configuration.repos;
        $scope.selectedRepo = Configuration.repos[0];
        //$scope.selectedRepo

        $scope.attemptLogin = function() {
            if(Authenticator.login($scope.credentials, $scope.selectedRepo)) {

                if($routeParams.redirect) {
                    $location.url($routeParams.redirect); 
                } else {
                    $location.url(defaultPage + $scope.selectedRepo.repo + 
                        '/' + $scope.credentials.username);
                }
            } else {

                $("#login-window-header").html("<span>Auth failed!</span>");
            }
        }

        $scope.onRepoSelectionChange = function() {
            $location.url('/' + $scope.selectedRepo.repo + '/login')
        }

        function _initialize() {
            if(!$routeParams.repository) {
                $location.url($scope.selectedRepo.repo+'/login');
            }

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