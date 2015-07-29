angular.module('zfh').controller('ctrl.base', ['$scope', '$http', '$timeout', 'myPlan', 'myTime', function($scope, $http, $timeout, myPlan, myTime){
    $scope.MyPlan = function(sx, method){
        myPlan[method](function(data){
            $.each(data, function(_, dd){
                $scope[sx][dd.id] = $.extend({}, dd, {
                    startDate : myTime.timeToStr(dd.startTime * 1000, true),
                    endDate : myTime.timeToStr(dd.endTime * 1000, true),
                    realDate : dd.realTime ? myTime.timeToStr(dd.realTime * 1000, true) : ""
                });
            });
        });
    };
    $.each({plans : 'get', completePlans : 'getComplete'}, function(sx, method){
        $scope[sx] = {};
        $scope.MyPlan(sx, method);
    });
}]);

angular.module('zfh').controller('ctrl.day', ['$scope', '$http', '$timeout', 'myPlan', 'myTime', function($scope, $http, $timeout, myPlan, myTime){
    var vm = $scope.vm = {};

    $.extend(vm, {
        y : Today.y,
        m : Today.m,
        d : Today.d,
        works : [],
        ccwork : {},
        postwork : {},
        plans : $scope.$parent.plans,

        appid : 0,

        init : function(reset){
            vm.ccwork = {
                id : 0,
                content : "",
                startTimeH : "09",
                startTimeI : "00",
                endTimeH : "12",
                endTimeI : "00",
                status : 0,
                planId : 0
            };

            vm.appid = 0;

            if(!reset){
                Day && Day.works && $.each(Day.works, function(_, work){
                    vm.works.push($.extend({}, work, {
                        start : myTime.hm(work.startTime),
                        end : myTime.hm(work.endTime),
                        duration : myTime.duration(work.startTime, work.endTime)
                    }));
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

        checkError : function(){
            var error = false;
            var content = $.trim(vm.ccwork.content);
            if(!content.length){
                error = vm.config["content"];
            }else if(vm.ccwork.endTimeH < vm.ccwork.startTimeH || (vm.ccwork.endTimeH == vm.ccwork.startTimeH && vm.ccwork.endTimeI <= vm.ccwork.startTimeI)) {
                error = vm.config["end"];
            }else if(!vm.ccwork.status){
                error = vm.config["status"];
            }else{
                vm.postwork = {
                    content : content,
                    startTime : myTime.offsetTime(vm.ccwork.startTimeH, vm.ccwork.startTimeI),
                    endTime : myTime.offsetTime(vm.ccwork.endTimeH, vm.ccwork.endTimeI),
                    status : vm.ccwork.status,
                    planId : parseInt(vm.ccwork.planId),
                    appId : vm.appid
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
            vm.doAjax('/day/add', function(data){
                vm.works.push($.extend({}, vm.postwork, {
                    id : data.id,
                    start : vm.ccwork.startTimeH + ':' + vm.ccwork.startTimeI,
                    end : vm.ccwork.endTimeH + ':' + vm.ccwork.endTimeI,
                    duration : myTime.duration(vm.postwork.startTime, vm.postwork.endTime)
                }));
                var planId = vm.postwork.planId;
                if(planId > 0){
                    if (vm.postwork.status == 100) {
                        delete vm.plans[planId];
                        $scope.$parent.MyPlan('completePlans', 'getComplete');
                    } else {
                        vm.plans[planId].status = vm.postwork.status;
                    }
                }
            });
            return false;
        },

        doAjax : function(url, cb){
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : url,
                data : $.param(vm.postwork),
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
                url : '/day/delete',
                data : $.param({id : id}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    $.each(vm.works, function(index, work){
                        if(work && work.id == id){
                            vm.works.splice(index, 1);
                            return false;
                        }
                    });
                }
            });
            return false;
        },

        config : {
            "content" : "请填写工作内容",
            "end" : "结束时间不能少于开始时间",
            "status" : "请填写完成度"
        }
    });

    vm.init();
}]);