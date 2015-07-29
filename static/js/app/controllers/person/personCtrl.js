angular.module('zfh').controller('ctrl.person', ['$scope', '$http', '$timeout', function($scope, $http, $timeout){
    var vm = $scope.vm = {};

    $.extend(vm, {
        persons : [],
        group : Group,
        ccperson : {},

        init : function(reset){
            vm.ccperson = {
                id : 0,
                name : "",
                loginName : "",
                loginPass : "",
                color : "",
                group : ""
            };

            !reset && $.each(Persons, function(_, person){
                vm.persons.push(person);
            });
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

        checkError : function(update){
            var error = false;
            $.each(vm.ccperson, function(k, v){
                if($.trim(v).length == 0){
                    error = vm.config["null"];
                    return false;
                }
            });
            !update && $.each(vm.persons, function(_, person){
                if(vm.ccperson.name == person.name){
                    error = vm.config["cf"];
                    return false;
                }
                if(vm.ccperson.loginName == person.loginName){
                    error = vm.config["login"];
                    return false;
                }
            });
            return error;
        },

        doAdd : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/person/add', function(data){
                vm.persons.push($.extend({}, vm.ccperson, {id : data.id}));
            });
            return false;
        },

        update : function(id){
            id = parseInt(id) || 0;
            $.each(vm.persons, function(_, person){
                if(person.id == id){
                    vm.ccperson = $.extend({}, person);
                }
            });
            return false;
        },

        doUpdate : function(){
            var error = vm.checkError(true);
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/person/update', function(data){
                $.each(vm.persons, function(index, person){
                    if(person.id == vm.ccperson.id){
                        vm.persons.splice(index, 1, $.extend({}, vm.ccperson));
                    }
                });
            });
            return false;
        },

        doAjax : function(url, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(vm.ccperson),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    cb && cb(data);
                    vm.init(true);
                    vm.ajaxdo("false");
                }
            });
        },

        doDelete : function(id){
            if(!confirm('确定要删除么？')){
                return false;
            }
            $http({
                method : "POST",
                url : '/person/delete',
                data : $.param({id : id}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    $.each(vm.persons, function(index, person){
                        if(person && person.id == id){
                            vm.persons.splice(index, 1);
                            return false;
                        }
                    });
                }
            });
            return false;
        },

        config : {
            "null" : "每项都需要填写",
            "cf" : "已经有重复的用户存在",
            "login" : "登录名不能相同"
        }
    });

    vm.init();
}]);