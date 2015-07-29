angular.module('zfh').controller('ctrl.day.app', ['$scope', '$http', '$timeout', 'myPlan', 'myTime', function($scope, $http, $timeout, myPlan, myTime){
    var vm = $scope.vm = {};
    var Apps = [];
    $.extend(vm, {
        apps : [],

        cappid : 0,
        cappname : '',
        open : false,
        ajaxok : 0,

        init : function(reset){
            if(!reset){
                $scope.$watch("$parent.vm.appid", function(newValue){
                    if(newValue == 0 && vm.cappid != 0){
                        vm.cappid = 0;
                        vm.cappname = '';
                    }
                });
            }
        },

        ajaxing : "false",
        ajaxdo : function(status){
            vm.ajaxing = status;
            if(status != "true" || status != "false") {
                $timeout(function(){
                    vm.ajaxing = "false";
                }, 0);
            }
        },

        ajax : function(){
            if(vm.ajaxok == 2){
                vm.open = true;
                return;
            }
            if(vm.ajaxok == 1){
                return;
            }
            vm.ajaxdo("true");
            vm.ajaxok = 1;
            $http({
                method : "GET",
                url : '/app/getApp',
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    data.sort(function(a, b){
                        return b.id - a.id;
                    });

                    Apps = data;
                    vm.apps = Apps.concat();
                    vm.open = true;
                    vm.ajaxok = 2;
                    vm.init(true);
                    vm.ajaxdo("false");
                }
            });
        },

        filter : function(key){
            key = $.trim(key);
            vm.apps = [];
            if(key.length) {
                $.each(Apps, function (_, app) {
                    if ((app.name || '').indexOf(key) != -1) {
                        vm.apps.push(app);
                    }
                });
            } else {
                vm.apps = Apps;
            }
        },

        select : function(id, name){
            if(vm.cappid == id){
                id = 0;
                name = '';
            }
            $scope.$parent.vm.appid = vm.cappid = id;
            vm.cappname = name;
            if(id != 0){
                vm.close();
            }
        },

        close : function(){
            vm.open = false;
        }
    });
    vm.init();
}]);