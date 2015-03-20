'use strict';
/*
    This is a work in progress as the OPTIONS call server side is the blocker.
*/
angular.module('repositories')
.factory('Gitlab', ['$http', 'Authenticator', 
    function($http, Authenticator) {
    
        var gitlab = {},
            _projList = [],     // cached
            prefix = '/api/v3';

        var _callHttp = function(url, method, data, headers) {
            var req = {
                url: url,
                method: method,
                headers: { 'PRIVATE-TOKEN': Authenticator.getCreds().private_token }
            }
            if (data) req.data = data;
            if (headers) req.headers = $.extend({}, req.headers, headers);

            return $http(req);
        }

        gitlab.listProjects = function() {
            var d = $.Deferred();

            if(Authenticator.getCreds() == null ) {
                d.reject('Not authenticated!');
            } else if( _projList.length < 1 ) {
            
                _callHttp(prefix+'/projects', 'GET', null, null)
                .success(function(data) {
                    _projList = data;
                    d.resolve(_projList);
                }).error(function(err) { d.reject(err); });
            } else {
                d.resolve(_projList);
            }

            return d;
        }

        return gitlab;
    }
]);