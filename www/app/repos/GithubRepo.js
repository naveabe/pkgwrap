angular.module('repositories', [])
.factory('Github', ['$http', 'AccessToken', 
    function($http, AccessToken) {
    
        var baseUrl = "https://api.github.com",
            github = {},
            _userInfo = null;
            _userOrgs = null;
            _projList = [];

        var oauthToken = AccessToken.set();

        var _callHttp = function(url, method) {
            return $http({
                method: method,
                url: url,
                headers: {
                    'Authorization': 'token ' + AccessToken.get().access_token
                }
            });
        }

        github.getUserInfo = function() {

            var _dfd = $.Deferred();
            
            if ( oauthToken == null ) {

                _dfd.reject("Not authorized!");

            } else if ( _userInfo == null ) {
                console.log('live');
                _callHttp(baseUrl+'/user', 'GET')
                .success(function(data, resp, hdrs, cfg) {
                    _userInfo = data;
                    _dfd.resolve(data);
                }).error(function(err) { _dfd.reject(err); });
            } else {
                console.log('cache');
                _dfd.resolve(_userInfo);
            }

            return _dfd;
        }
        
        github.getUserOrgs = function() {
            var _dfd = $.Deferred();

            if ( oauthToken == null ) {
                
                _dfd.reject("Not authorized!");

            } else if ( _userOrgs == null ) {
                
                console.log('live');
                _callHttp(baseUrl+'/user/orgs', 'GET')
                .success(function(data, resp, hdrs, cfg) {
                    
                    _userOrgs = data;
                    for(var i=0; i < _userOrgs.length; i++) {
                        _userOrgs[i].icon = _userOrgs[i].avatar_url;
                    }
                    _dfd.resolve(data);
                
                }).error(function(err) { _dfd.reject(err); });
            
            } else {
            
                console.log('cache');
                _dfd.resolve(_userOrgs);
            }
            return _dfd;
        }

        github.listProjects = function() {
            var _dfd = $.Deferred();

            if ( oauthToken == null ) {
                
                _dfd.reject("Not authorized!");

            } else if ( _projList.length < 1 ) {

                _callHttp(baseUrl + '/user/repos', 'GET')
                .success(function(data) {
                    _projList = data;
                    _dfd.resolve(_projList);
                }).error(function(err) {
                    _dfd.reject(err);
                });
            } else {
                _dfd.resolve(_projList);
            }
            return _dfd;
        }

        return github;
    }
]);