angular.module('zfh').controller('ctrl.login', ['$scope', '$http', '$timeout', function($scope, $http, $timeout){
    var vm = $scope.vm = {};

    $.extend(vm, {
        cclogin : {
            name : "",
            pass : ""
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

        checkError : function(){
            var error = false;
            $.each(vm.cclogin, function(k, v){
                if($.trim(v).length == 0){
                    error = vm.config["null"];
                    return false;
                }
            });
            return error;
        },

        login : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }

            vm.doAjax('/login/login', function(data){
                setTimeout(function(){
                    location.href = data.id == -1 ? '/month' : '/day';
                }, 0);
            });
            return false;
        },

        doAjax : function(url, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(vm.cclogin),
                responseType : "json"
            }).success(function(data){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    cb && cb(data);
                    vm.ajaxdo("false");
                }
            });
        },

        config : {
            "null" : "每项都需要填写"
        }
    });
}]);