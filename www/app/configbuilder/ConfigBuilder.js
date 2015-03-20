"use strict";

angular.module('ipkg.configbuilder', [])
.controller('ConfigBuilderController', ['$scope', '$routeParams', 
    function($scope, $routeParams) {
        

        var init = function() {
            $scope.repository = $routeParams.repository;
            $scope.username = $routeParams.username;
            $scope.project = $routeParams.project;
        }

        init();
    }
])
.directive('configBuilder', [function() {
    return {
        restrict: 'E',
        scope: {},
        templateUrl: 'app/configbuilder/config-builder.html',
        link: function(scope, elem, attrs, ctrl) {
            
            var distroTemplate = {
                name: 'centos',
                release: '6',
                build_deps: [],
                deps: [],
                build_cmd: [],
                pre_install:[],
                post_install: [],
                pre_uninstall: [],
                post_uninstall: []
            };

            var pkgTemplate = {
                version: '',
                build_env: '',
                url: '',
                tagbranch:''
            };

            var ubuntuReleases = ['12.04','14.04'],
                centosReleases = ['6'];

            scope.support = {
                distros: [{
                    name:'centos',
                    label:'CentOS',
                    releases: centosReleases
                },{
                    name:'ubuntu',
                    label:'Ubuntu',
                    releases: ubuntuReleases
                }]
            };

            scope.distroReleases = {
                'centos': centosReleases,
                'ubuntu': ubuntuReleases
            }

            scope.pkgwrap = {
                'distributions': [ angular.copy(distroTemplate) ],
                'package': angular.copy(pkgTemplate)
            };

            scope.hasPackageData = function(obj) {
                for( var k in obj ) {
                    if ( obj[k] && obj[k] != null && obj[k] !== '' ) return true;
                }
                return false;
            }

            scope.addDistribution = function() {
                scope.pkgwrap.distributions.push(angular.copy(distroTemplate));
            }

            
            scope.onDistroChange = function(distro) {
                distro.release = scope.distroReleases[distro.name][0]
            }

        }
    }
}])
.directive('yamlList', [function() {
    return {
        restrict: 'A',
        require: '?ngModel',
        link: function(scope, elem, attrs, ctrl) {
            if(!ctrl) return;

            var lines2array = function(viewVal) {
                var arr = viewVal.split('\n');
                if (arr.length < 1) return ctrl.$modelValue;
                return arr
            }

            var array2lines = function(modelVal) {
                if (!modelVal) return "";

                return modelVal.join('\n');
            }

            var autoResize = function(evt) {
                    
                if ( evt.keyCode != 13 && evt.keyCode != 8) return;
                
                elem.height('auto');
                elem.height(evt.target.scrollHeight);
            }

            var init = function() {
                ctrl.$formatters.push(array2lines);

                ctrl.$parsers.unshift(lines2array);
                
                elem.on('keyup', autoResize);
                
            }

            init();
        }
    }
}]);