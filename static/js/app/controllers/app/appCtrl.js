angular.module('zfh')
    .controller('ctrl.apps', ['$scope', '$filter', function($scope, $filter){
        var vm = $scope.vm = {};
        var timer = null;
        $.extend(vm, {
            apps : {},
            cyear : '',

            init : function(){
                vm.guila(Apps);
                vm.cyear = new Date().getFullYear();
            },

            search : function(value){
                if(timer){
                    clearTimeout(timer);
                    timer = null;
                }
                timer = setTimeout(function() {
                    if (value.length > 0) {
                        var apps = [];
                        $.each(Apps, function (_, app) {
                            if ((app.name || '').indexOf(value) != -1) {
                                apps.push(app);
                            }
                        });
                        vm.guila(apps);
                    } else {
                        vm.init();
                    }
                    $scope.$digest();
                }, 100);
            },

            cy : function(year){
                vm.cyear = year;
            },

            guila : function(apps){
                vm.apps = {};
                var dateFilter = $filter('date');
                $.each(apps, function(_, app){
                    var time = app.create_time * 1000;
                    var year = dateFilter(time, 'yyyy');
                    vm.apps[year] = vm.apps[year] || {};
                    var month = dateFilter(time, 'M');
                    vm.apps[year][month] = vm.apps[year][month] || {};
                    vm.apps[year][month][app.id] = app;
                });

                var tmp = {};
                $.each(vm.apps, function(year, appss){
                    var tt = [];
                    var count = 0;
                    $.each(appss, function(month, apps){
                        $.each(apps, function(){
                            count++;
                        });
                        tt.push({
                            month : month,
                            apps : apps
                        });
                    });
                    tt.reverse();
                    tmp[year] = {
                        apps : tt,
                        count : count
                    };
                });
                vm.apps = tmp;
            }
        });
        vm.init();
    }])
    .controller('ctrl.app', ['$scope', '$http', '$timeout', 'myTime', function($scope, $http, $timeout, myTime){
        var vm = $scope.vm = {};

        $.extend(vm, {
            logs : [],
            cclog : {},
            offset : 0,
            pageNum : 10,
            canP : false,
            canN : false,
            open : true,
            appid : 0,

            cyear : '',

            searchInit : false,

            init : function(reset){
                vm.cclog = {
                    desc : ""
                };

                if(!reset){
//                    $scope.$watch("$parent.vm.cyear", function(newValue){
//                        console.log(newValue);
//                    });
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

            doOpen : function(){
                vm.open = !vm.open;
                if(vm.open && !vm.searchInit) {
                    vm.search();
                    vm.searchInit = true;
                }
            },

            checkError : function(){
                var error = false;
                $.each(['desc'], function(_, sx){
                    if(!vm.cclog[sx].length){
                        error = vm.config["null"];
                        return false;
                    }
                });
                vm.cclog.appid = vm.appid;
                return error;
            },

            doAdd : function(){
                var error = vm.checkError();
                if(error){
                    vm.ajaxdo(error);
                    return false;
                }
                vm.doAjax('/app/add', function(data){
                    vm.logs.unshift($.extend({}, vm.cclog, {id : data.id, time : +new Date() / 1000, pname : User.name}));
                });
                return false;
            },

            doAjax : function(url, cb, method, postdata){
                vm.ajaxdo("true");
                var dd = {
                    method : method || "POST",
                    url : url,

                    responseType : "json"
                };
                if(dd.method == 'GET'){
                    dd["params"] = postdata;
                } else {
                    dd.data = $.param(vm.cclog);
                }
                $http(dd).success(function(data, status){
                    if(data.error){
                        vm.ajaxdo(data.error);
                    } else {
                        cb && cb(data);
                        vm.init(true);
                        vm.ajaxdo("false");
                    }
                });
            },

            foucs : function(){
                if(vm.foucsInit){
                    return;
                }
                vm.foucsInit = true;
                vm.search();
            },

            search : function(){
                vm.doAjax('/app/search', function(data){
                    vm.logs = [];
                    $.each(data.logs, function(_, log){
                        vm.logs.push($.extend({}, log, {pname : data.pps[log.uid]}));
                    });
                    if(vm.logs.length == vm.pageNum){
                        vm.canN = true;
                    } else {
                        vm.canN = false;
                    }
                    if(vm.offset > 0){
                        vm.canP = true;
                    } else {
                        vm.canP = false;
                    }
                }, "GET", {
                    appid : vm.appid,
                    offset : vm.offset,
                    page : vm.pageNum
                });
            },

            prev : function(){
                vm.offset -= vm.pageNum;
                if(vm.offset < 0){
                    vm.offset = 0;
                }
                this.search();
            },

            next : function(){
                vm.offset += vm.pageNum;
                this.search();
            },

            config : {
                "null" : "日志必须填写点内容"
            }
        });

        vm.init();
    }]);