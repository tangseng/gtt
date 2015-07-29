angular.module('zfh').controller('ctrl.day.admin', ['$scope', '$http', '$timeout', 'myTime', function($scope, $http, $timeout, myTime){
    var vm = $scope.vm = {};

    $.extend(vm, {
        ccwork : {},
        sid : "",
        end : GlobalConfig.end,

        init : function(reset){
            vm.sid = "";
            vm.ccwork = {
                id : 0,
                content : "",
                startTimeH : "09",
                startTimeI : "00",
                endTimeH : "12",
                endTimeI : "00",
                date : "",
                status : 0
            };
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
                    id : vm.ccwork.id,
                    content : content,
                    startTime : myTime.offsetTime(vm.ccwork.startTimeH, vm.ccwork.startTimeI),
                    endTime : myTime.offsetTime(vm.ccwork.endTimeH, vm.ccwork.endTimeI),
                    status : vm.ccwork.status,
                    date : myTime.time(vm.ccwork.date)
                }
            }
            return error;
        },

        search : function(){
            if(!(vm.sid > 0)){
                vm.ajaxdo(vm.config["sid"]);
                return false;
            }
            vm.ajaxdo("true");
            $http({
                method : "GET",
                url : '/day/search',
                params : {id : vm.sid},
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    vm.ccwork = data;
                    vm.ccwork.date = myTime.timeToStr(vm.ccwork.time * 1000);
                    var startHM = myTime.h_m(vm.ccwork.startTime);
                    vm.ccwork.startTimeH = startHM.h;
                    vm.ccwork.startTimeI = startHM.m;
                    var endHM = myTime.h_m(vm.ccwork.endTime);
                    vm.ccwork.endTimeH = endHM.h;
                    vm.ccwork.endTimeI = endHM.m;
                    vm.ajaxdo("false");
                }
            });
        },

        doUpdate : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.ajaxdo("true");
            $http({
                method : "POST",
                url : '/day/update',
                data : $.param(vm.postwork),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    vm.ajaxdo(data.error);
                } else {
                    vm.ajaxdo('调整成功');
                    vm.init(true);
                    $timeout(function(){
                        vm.ajaxdo("false");
                    }, 1000);
                }
            });
        },

        doDelete : function(){
            if(!confirm('确定要删除么？')){
                return false;
            }
            $http({
                method : "POST",
                url : '/day/delete',
                data : $.param({id : vm.ccwork.id}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    vm.init(true);
                }
            });
            return false;
        },

        config : {
            "content" : "请填写工作内容",
            "end" : "结束时间不能少于开始时间",
            "status" : "请填写完成度",
            "sid" : "输入id后再查找"
        }
    });

    vm.init();
}]);