"use strict";

angular.module('ipkg.configbuilder', [])
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
                build_cmd: "",
                pre_install:[],
                post_install: [],
                pre_uninstall: [],
                post_uninstall: []
            };

            scope.support = {
                distros: [
                    {name:'ubuntu',label:'Ubuntu'},
                    {name:'centos',label:'CentOS'},
                ],
            };

            scope.centosReleases = [ '6' ];
            scope.ubuntuReleases = [ '12.04', '14.04' ];

            scope.distributions = [ angular.copy(distroTemplate) ];

            scope.addDistribution = function() {
                scope.distributions.push(angular.copy(distroTemplate));
            }
        }
    }
}]);