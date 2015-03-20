'use strict';

angular.module('ipkg.login', [])
.factory('Authenticator', [
    '$window', '$http', '$location', '$rootScope', 'AccessToken',
    function($window, $http, $location, $rootScope, AccessToken) {

        function _sessionIsAuthenticated(evtArgs) {
            var creds = _getCreds();

            if(creds && creds.username && creds.username !== "" &&
                            creds.password && creds.password !== "") {
                
                var eArgs = $.extend({}, {username:creds.username}, evtArgs);
                // do custom checking here
                $rootScope.$emit('user:auth:success', eArgs);
                return true
            }
            return false;
        }
        
        function _login(creds, args) {
            var eArgs = $.extend({}, creds, args);

            _setCreds(eArgs);

            $rootScope.$emit('user:auth:success', eArgs);
                
            return true;
        }
        
        /*
            Wrapper to multiple auth types
        */
        function _logout() {

            var sStor = $window.sessionStorage;
            if(sStor['credentials']) {
                delete sStor['credentials'];
            }
        
            AccessToken.destroy();
            
            $rootScope.$emit('user:unauth', {});
        }

        var _getCreds = function() {
            
            if($window.sessionStorage["credentials"]) {
                return JSON.parse($window.sessionStorage["credentials"]);
            }
            return null;
        }

        var _setCreds = function(args) {
            $window.sessionStorage['credentials'] = angular.toJson(args);
        }
        
        var Authenticator = {
            login                 : _login,
            logout                : _logout,
            getCreds              : _getCreds,
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
.factory('GitlabAuth', ['$http', 'Authenticator', 
    function($http, Authenticator) {
        /*
            Since OPTIONS call does not return ACL this 
            does not currently work.
        */
        var GitlabAuth = {};
        var _userInfo, 
            prefix = '/api/v3';
        
        /* Contains token and user info */
        GitlabAuth.getPrivateToken = function(url, creds) {
            
            var dfd = $.Deferred();

            var auth = Authenticator.getCreds();
            if( auth && auth.private_token ) {
                dfd.resolve(auth);
            } else if ( _userInfo == null ) {
                
                $http({
                    url: url + prefix + '/session',
                    method: 'POST',
                    headers: { 'Content-Type': 'applicaiton/x-www-form-urlencoded' },
                    data: 'login='+creds.username+'&password='+creds.password,
                }).success(function(data, resp, hdrs, cfg) {
                
                    _userInfo = data;            
                    Authenticator.login(_userInfo);
                    dfd.resolve(_userInfo);
                
                }).error(function(err) { dfd.reject(err); });

            } else {
                dfd.resolve(_userInfo);
            }

            return dfd;
        }

        return GitlabAuth;
    }
])
.controller('loginController', [
    '$scope', 'Configuration',
    function($scope, Configuration) {

        function init() {

            $scope.supportedRepos = Configuration.repos;
        }

        init();
    }
])
.controller('repoLoginController', [
    '$scope', '$routeParams', '$location', 'Authenticator', 'SupportedVCs', 'Github', 'AccessToken', 'GitlabAuth',
    function($scope, $routeParams, $location, Authenticator, SupportedVCs, Github, AccessToken, GitlabAuth) {

        var oauthToken;

        var basicAuthLogin = function() {
            if ( $scope.credentials.username != '' ) {

                Authenticator.login($scope.credentials);
                /*
                // Enable this once OPTIONS call has been sorted //
                GitlabAuth.getPrivateToken( $scope.repoDetails.url, $scope.credentials )
                .then(function(data, resp, hdrs, cfg) {
                    
                    $location.url('/'+$scope.repository+'/'+$scope.credentials.username);
                
                }, function(err) { console.log(err); });
                */
                $location.url('/'+$scope.repository+'/'+$scope.credentials.username);
            }
        }

        var getUserInfo = function() {
            
            Github.getUserInfo().then(function(data) {
                data.repo = $scope.repository; 
                data.icon = data.avatar_url;
                Authenticator.login({ username: data.login }, data);    
                
                $location.url('/' + $scope.repoDetails.repo + '/' + data.login);
            }, function(err) {
                $("#login-window-header").html(err);
            });
        }

        var init = function() {
            $scope.repository  = $routeParams.repository;
            $scope.repoDetails = SupportedVCs.getDetails($scope.repository);
            // Basic auth
            $scope.credentials = { username: "", password: "" };
            $scope.basicAuthLogin = basicAuthLogin;

            oauthToken = AccessToken.set();
            //console.log(oauthToken);
            if ( oauthToken != null ) {    
                 getUserInfo();
            }
        }


        init();
    }
]);




