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
                centosReleases = ['6', '7'];

            var sanitizeConfig = function(data) {
                /* Remove empty fields from the pkgwrap config */
                var cfg = angular.copy(data);
                for(var i in cfg.Package) {
                    if ( cfg.Package[i] == '' ) delete cfg.Package[i];
                }
                if ( Object.keys(cfg.Package).length < 1) delete cfg.Package;

                for(var i=0; i < cfg.Distributions.length; i++) {

                    var distro = cfg.Distributions[i];
                    for(var j in  distro ) {
                        if ( j == 'name' || j == 'release' ) continue;
                        if( distro[j].length < 1 ) delete distro[j];
                    }
                }

                return cfg;
            }

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

            scope.addDistribution = function() {
                scope.pkgwrap.Distributions.push(angular.copy(distroTemplate));
            }

            
            scope.onDistroChange = function(distro) {
                distro.release = scope.distroReleases[distro.name][0]
            }

            var init = function() {
                scope.pkgwrap = {
                    'Distributions': [ angular.copy(distroTemplate) ],
                    'Package': angular.copy(pkgTemplate)
                };
                
                scope.$watch(function() {return scope.pkgwrap}, function(newVal) { 
                    scope.yamlConfig = YAML.stringify(sanitizeConfig(scope.pkgwrap)); 
                }, true);
            }

            init();
        }
    }
}])
.directive('yamlList', [function() {
    /* Converts newline delimted input to array */
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