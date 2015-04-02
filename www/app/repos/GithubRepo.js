'use strict';

angular.module('repositories', [])
.factory('Github', ['$http', 'AccessToken', 'Authenticator',
    /* Github OAuth API */
    function($http, AccessToken, Authenticator) {
    
        var baseUrl = "https://api.github.com",
            github = {},
            _userInfo = null,
            _userOrgs = null,
            _projList = [];

        var oauthToken = AccessToken.set();
        /*
            Generic http call wrapper

            Args:
                method : GET, POST, PUT
                data   : JSON stringify'able object
        */
        var _callHttp = function(url, method, data) {
            var hdrs = { 'Authorization': 'token ' + AccessToken.get().access_token };
            if ( data ) {
                hdrs['Content-Type'] = 'application/json';
                return $http({
                    method: method,
                    url: url,
                    headers: hdrs,
                    data: angular.toJson(data)
                });
            }

            return $http({
                method: method,
                url: url,
                headers: hdrs
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

        github.createFile = function(project, path, b64content, commitMsg, branch) {
            // Format: '/repos/:owner/:repo/contents/:path'
            var dfd = $.Deferred(),
                creds = Authenticator.getCreds(),
                url = baseUrl + '/repos/',
                payload;

            if ( creds == null ) {
                dfd.reject({"error": "No credentials!"});
                return;
            }

            url += creds.username + '/' + project + '/contents/' + path;
            
            payload = {
                message: commitMsg,
                content: b64content,
                committer: {
                    name: "ipkg UI",
                    email: "builder@ipkg.io"
                }
            }

            _callHttp(url, 'PUT', payload)
            .success(function(data) {
                dfd.resolve(data);
            }).error(function(err) {
                dfd.reject({"error": err});
            });

            return dfd;
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

        github.createWebhook = function(project) {
            // push is default
            var hookData = {
                name: 'ipkg.io',
                active: true,
                events: ["create", "release"],
                config: {
                    insecure_ssl: true,
                    url: '',
                    content_type: 'json',
                }
            };

            return _callHttp(baseUrl+'/repos/'+Authenticator.getCreds().username+'/'+project+'/hooks' , 'POST', hookData)
        }

        /* In progress */
        github.disableWebhook = function(project) {
            var hookData = {
                active: false
            };
        }

        return github;
    }
])
.factory('GithubPublic', ['$resource', function($resource) {
    /* github public (unauthenticated) API */ 
    return $resource('https://api.github.com/users/:username/:qtype', {}, {
        userRepos: {
            params: {"username": "@username", "qtype": "repos"},
            method: 'GET',
            isArray: true
        }
    });
}]);