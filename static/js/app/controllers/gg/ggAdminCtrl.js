angular.module('zfh').controller('ctrl.gg.admin', ['$scope', '$http', '$timeout', 'myTime', function($scope, $http, $timeout, myTime){
    var vm = $scope.vm = {};

    $.extend(vm, {
        ggs : [],
        ccgg : {},
        offset : 0,
        pageNum : PageNum,

        init : function(reset){
            vm.ccgg = {
                id : 0,
                title : "",
                desc : "",
                class : "default",
                dateString : "",
                date : 0
            };

            if(!reset){
               this.ajax(0);
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

        checkError : function(update){
            var error = false;
            $.each(vm.ccgg, function(k, v){
                if($.trim(v).length == 0){
                    error = vm.config["null"];
                    return false;
                }
            });
            if(!error){
                vm.ccgg.date = myTime.time(vm.ccgg.dateString);
            }
            return error;
        },

        ajax : function(more){
            if(more){
                vm.offset += vm.pageNum;
            }
            $http({
                method : "GET",
                url : '/gg/ajax',
                params : {offset : vm.offset},
                responseType : "json"
            }).success(function(data, status){
                $.each(data, function(_, gg){
                    vm.ggs.push($.extend({}, gg, {dateString : myTime.timeToStr(gg.date * 1000)}));
                });
            });
        },

        doAdd : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/gg/add', function(data){
                vm.ggs.push($.extend({}, vm.ccgg, {id : data.id}));
            });
            return false;
        },

        update : function(id){
            id = parseInt(id) || 0;
            $.each(vm.ggs, function(_, gg){
                if(gg.id == id){
                    vm.ccgg = $.extend({}, gg);
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
            vm.doAjax('/gg/update', function(data){
                $.each(vm.ggs, function(index, gg){
                    if(gg.id == vm.ccgg.id){
                        vm.ggs.splice(index, 1, $.extend({}, vm.ccgg));
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
                data : $.param(vm.ccgg),
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
                url : '/gg/delete',
                data : $.param({id : id}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    $.each(vm.ggs, function(index, gg){
                        if(gg && gg.id == id){
                            vm.ggs.splice(index, 1);
                            return false;
                        }
                    });
                }
            });
            return false;
        },

        config : {
            "null" : "每项都需要填写"
        }
    });

    vm.init();
}]);