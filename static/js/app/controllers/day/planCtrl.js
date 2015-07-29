angular.module('zfh').controller('ctrl.plan', ['$scope', '$http', '$timeout', function($scope, $http, $timeout){
    var vm = $scope.vm = {};

    $.extend(vm, {
        gstart : GlobalConfig.start,
        gend : GlobalConfig.end,

        plans : $scope.$parent.plans,
        completePlans : $scope.$parent.completePlans,
        completePlansCount : 0,
        ccplan : {},
        postplan : {},


        init : function(reset){
            vm.ccplan = {
                id : 0,
                content : "",
                startDate : "",
                startTime : 0,
                endDate : "",
                endTime : 0
            };
            if(!reset){
                $scope.$watch(
                    'vm.completePlans',
                    function (newValue, oldValue) {
                        vm.completePlansCount = !$.isEmptyObject(newValue) ? 1 : 0;
                    },
                    true
                );
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

        stime : function(date){
            return +new Date(date) / 1000;
        },

        checkError : function(){
            var error = false;
            var content = $.trim(vm.ccplan.content);
            var startTime = vm.stime(vm.ccplan.startDate);
            var endTime = vm.stime(vm.ccplan.endDate);
            if(!content.length){
                error = vm.config["content"];
            }else if(endTime <= startTime) {
                error = vm.config["end"];
            }else{
                vm.postplan = {
                    id : vm.ccplan.id,
                    content : content,
                    startTime : startTime,
                    endTime : endTime
                }
            }
            return error;
        },

        doAdd : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/plan/add', function(data){
                vm.plans[data.id] = $.extend({}, vm.postplan, {
                    id : data.id,
                    status : 0,
                    startDate : vm.ccplan.startDate,
                    endDate : vm.ccplan.endDate
                });
            });
            return false;
        },

        update : function(id){
            id = parseInt(id) || 0;
            $.each(vm.plans, function(_, plan){
                if(plan.id == id){
                    vm.ccplan = $.extend({}, plan);
                }
            });
            return false;
        },

        doUpdate : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/plan/update', function(data){
                $.extend(vm.plans[vm.postplan.id], vm.postplan, {
                    startDate : vm.ccplan.startDate,
                    endDate : vm.ccplan.endDate
                });
            });
            return false;
        },

        doAjax : function(url, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(vm.postplan),
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
                url : '/plan/delete',
                data : $.param({id : id}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    delete vm.plans[id]
                }
            });
            return false;
        },

        config : {
            "content" : "请填写工作任务",
            "end" : "结束日期不能早于开始日期"
        }
    });

    vm.init();
}]);

angular.module('zfh').factory('myPlan', ['$http', function($http){
    function http(url, cb){
        $http({
            method : "GET",
            url : url,
            responseType : "json"
        }).success(function(data, status){
            if(data && !data.error){
                cb && cb(data);
            }
        });
    }
    return {
        get : function(cb){
            http('/plan/', cb);
        },

        getComplete : function(cb){
            http('/plan/complete', cb);
        }
    }
}]);