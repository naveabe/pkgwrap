"use strict";

angular.module('ipkg.builder', [])
.factory('PkgBuilder', ['$http', function($http) {

    var pkgbuilder = {}

    pkgbuilder.requestPackageBuild = function(repoUrl, username, project, version, optTag) {

        var buildReq = { 
            'Package': { 
                'url': repoUrl+'/'+username+'/'+project,
                'tagbranch': (optTag !== null && optTag !== "") ? optTag : "master"
            }
        }

        return $http({
            method: 'POST',
            url: '/api/builder/'+username+'/'+project+'/'+version,
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(buildReq)
        });
    }

    return pkgbuilder;
}])
.directive('packageBuild', ['$rootScope', '$window', 'PkgBuilder', function($rootScope, $window, PkgBuilder) {
    return {
        restrict: 'E',
        replace: true, // replace the element with the directive
        scope: {
            project: '@',
            version: '@',
            tagbranch: '@',
            username: '=',
            repository: '=',
        },
        templateUrl: 'app/builder/build-request.html',
        link: function(scope, elem, attrs, ctrl) {

            scope.$watch('repository', function(newVal, oldVal) { init(); });

            var submitBuildRequest = function(tagbranch) {
                console.log('Submitting build request...');
                PkgBuilder.requestPackageBuild(
                    scope.repository,
                    scope.username,
                    scope.project,
                    scope.version, 
                    tagbranch)
                .success(function(data, status, headers, config) {
                    // UI feedback
                    console.log(data);
                    //setTimeout(function() { $window.location.reload(); }, 3000);
                    $rootScope.$broadcast('build:history:changed', data);
                })
                .error(function(data, status, headers, config) {
                    // UI feedback
                    console.error(data);
                });

                elem.modal('hide');
            };

            var init = function() {
              scope.project   = scope.project || 'Unknown';
              scope.version   = scope.version || '0.0.1';
              scope.tagbranch = scope.tagbranch || 'master';

              scope.submitBuildRequest = submitBuildRequest; 
              
              $rootScope.$on('build:request:init', function(evt, data) {
                scope.project = data.project; });
            };

        }
    }
}])
.directive('requestDialogButton', ['$rootScope', '$window', 'PkgBuilder', function($rootScope, $window, PkgBuilder) {
    return {
        restrict: 'E',
        replace: true,
        require: '?ngModel',
        templateUrl: '/app/builder/build-btn.html',
        link: function(scope, elem, attrs, ctrl) {
            if(!ctrl) return;

            scope.$watch(function() { return ctrl.$modelValue }, function(newVal, oldVal) { init(); });

            var showRequestDialog = function() {

                $('#build-request-modal').modal('show');
                $rootScope.$broadcast('build:request:init', {'project': ctrl.$modelValue});
            }

            var init = function() {
                scope.showRequestDialog = showRequestDialog;
            }
        }
    }
}]);