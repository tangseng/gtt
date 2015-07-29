angular.module('zfh')
    .controller('ctrl.jx.admin', ['$scope', '$http', '$timeout', 'myMonth', 'myMjx', 'myTime', function($scope, $http, $timeout, myMonth, myMjx, myTime){
        var vm = $scope.vm = {};

        $.extend(vm, {
            persons : {},
            uid : "0",
            days : [],
            inited : false,
            ccJxday : [],
            qian : "",
            hou : "",
            mjxs : {},

            init : function(){
                vm.persons = Persons
                if(Y){
                    vm.days = myMonth(Y, M, D);
                    vm.inited = true;
                    vm.qian = vm.days[0].date;
                    vm.hou = vm.days[vm.days.length - 1].date;
                    $scope.$watch("vm.uid", function (newValue) {
                        vm.ccJxday = [];
                    });
                }
            },

            select : function(uid){
                vm.uid = uid;
            }
        });

        vm.init();
    }])
    .controller('ctrl.mjx', ['$scope', '$timeout', 'myTime', 'myMjx', function($scope, $timeout, myTime, myMjx){
        var vm = $scope.vm = {};
        $.extend(vm, {
            mjxs : [],
            ccmjx : {},

            init : function(reset) {
                vm.ccmjx = {
                    id: 0,
                    date: "",
                    content: "",
                    score: 10
                };

                if (!reset){
                    $scope.$parent.$watch("vm.uid", function (newValue) {
                        vm.mjxs = [];
                        if (newValue > 0) {
                            myMjx.get({
                                uid: $scope.$parent.vm.uid,
                                qian: myTime.time($scope.$parent.vm.qian + ' ' + GlobalConfig.start),
                                hou: myTime.time($scope.$parent.vm.hou + ' ' + GlobalConfig.end)
                            }, function (data) {
                                $.each(data, function (_, dd) {
                                    vm.mjxs.push($.extend({}, dd, {'date': myTime.timeToStr(dd.date * 1000)}));
                                });
                            });
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

            checkError : function(){
                var error = false;
                var content = $.trim(vm.ccmjx.content);
                var uid = $scope.$parent.vm.uid;
                if(uid <= 0){
                    error = vm.config["uid"];
                } else if(!content.length) {
                    error = vm.config["content"];
                } else if(vm.ccmjx.date == "") {
                    error = vm.config["date"];
                } else {
                    vm.postmjx = {
                        id : vm.ccmjx.id,
                        uid : uid,
                        content : content,
                        date : myTime.time(vm.ccmjx.date),
                        score : vm.ccmjx.score
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

                vm.ajaxdo("true");
                myMjx.add(vm.postmjx, function(data){
                    if(data.error){
                        vm.ajaxdo(data.error);
                    } else {
                        vm.mjxs.push($.extend({}, vm.postmjx, {
                            id : data.id,
                            date : vm.ccmjx.date
                        }));
                        vm.init(true);
                        vm.ajaxdo("false");
                    }
                });

                return false;
            },

            update : function(id){
                id = parseInt(id) || 0;
                $.each(vm.mjxs, function(_, mjx){
                    if(mjx.id == id){
                        vm.ccmjx = $.extend({}, mjx);
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

                vm.ajaxdo("true");
                myMjx.update(vm.postmjx, function(data){
                    if(data.error){
                        vm.ajaxdo(data.error);
                    } else {
                        $.each(vm.mjxs, function(index, mjx){
                            if(mjx.id == vm.ccmjx.id){
                                vm.mjxs.splice(index, 1, $.extend({}, vm.ccmjx));
                            }
                        });
                        vm.init(true);
                        vm.ajaxdo("false");
                    }
                });
                return false;
            },

            doDelete : function(id){
                if(!confirm('确定要删除么？')){
                    return false;
                }
                myMjx.delete(id, function(data){
                    if(data.error){
                        alert(data.error);
                    } else {
                        $.each(vm.mjxs, function(index, mjx){
                            if(mjx && mjx.id == id){
                                vm.mjxs.splice(index, 1);
                                return false;
                            }
                        });
                    }
                });
                return false;
            },

            config : {
                'uid' : '请先选择人员',
                'content' : '绩效内容不能为空',
                'date' : '日期需要选择'
            }
        });

        vm.init();
    }])
    .directive('myJx', ['myJxDay', 'myTime', function(myJxDay, myTime){
        return {
            restrict : "A",
            link : function(scope, element, atts){

                function filterJxday(data){
                    var tmp = [];
                    data && $.each(data, function(_, dd){
                        tmp.push({
                            id : dd.id,
                            content : dd.content,
                            start : myTime.hm(dd.startTime),
                            end : myTime.hm(dd.endTime),
                            duration : myTime.duration(dd.startTime, dd.endTime),
                            date : myTime.timeToStr(dd.time * 1000, true),
                            status : dd.status,
                            score : dd.score
                        });
                    });
                    return tmp;
                }

                element.find('.nav').on('show.bs.tab', 'a', function(){
                    if(!(scope.vm.uid > 0)){
                        return;
                    }
                    var $this = $(this);
                    var $tab = $this.parent();
                    var date = $tab.attr('title');
                    var y = $tab.data('y');
                    var m = $tab.data('m');
                    var d = $tab.data('d');
                    myJxDay({
                        date: date,
                        y: y,
                        m: m,
                        d: d,
                        uid : scope.vm.uid
                    }, function (data, needApply) {
                        if(needApply){
                            scope.$apply(function(){
                                scope.vm.ccJxday = filterJxday(data);
                            });
                        } else {
                            scope.vm.ccJxday = filterJxday(data);
                        }
                    });
                });

                scope.$watch(atts.inited, function(newValue){
                    if(newValue == true){
                        //这边的处理有瑕疵，会出问题，再想
                        setTimeout(function(){
                            element.find('.nav .init a').tab('show');
                        }, 0);
                    }
                });

                scope.$watch("vm.uid", function (newValue) {
                    if(newValue > 0){
                        //这边的处理有瑕疵，会出问题，再想
                        setTimeout(function(){
                            element.find('.nav .active a').trigger('show.bs.tab');
                        }, 0);
                    }
                });
            }
        }
    }])
    .directive('myJxScore', ['$http', function($http){
        return {
            restrict : "A",
            link : function(scope, element, atts){
                element.on({
                    'click':function(event){
                        var $this = $(this);
                        var score = parseInt($this.prev().val()) || 0;
                        $this.prop('disabled', true);
                        $http({
                            method : 'POST',
                            url : '/jx/score',
                            data : $.param({score : score, id : $this.data('id')}),
                            responseType : "json"
                        }).success(function(data, status){
                            if(data.error){
                                alert(data.error);
                            } else {
                                $this.prop('disabled', false);
                            }
                        });
                    }
                });
            }
        }
    }])
    .factory('myJxDay', ['$http', '$cacheFactory', function($http, $cacheFactory){
        var cache = $cacheFactory('jxday');
        return function(params, cb){
            var key = params.date + '-' + params.uid;
            var cc = cache.get(key);
            if(cc != undefined){
                cb && cb(cc, true);
                return;
            }
            $http({
                method : "GET",
                url : '/jx/getDay',
                params : params,
                responseType : "json"
            }).success(function(data, status){
                if(data && !data.error){
                    cb && cb(data, false);
                    cache.put(key, data);
                }
            });
        }
    }])
    .factory('myMjx', ['$http', function($http){
        var http = function(method, url, params, postdata, cb){
            $http({
                method: method,
                url: url,
                params: params || {},
                data : postdata ? $.param(postdata) : null,
                responseType: "json"
            }).success(function (data) {
                cb && cb(data);
            });
        };
        return {
            'get' : function(params, cb){
                http('GET', '/jx/getMjx', params, null, cb);
            },

            'add' : function(postdata, cb){
                http('POST', '/jx/addMjx', null, postdata, cb);
            },

            'update' : function(postdata, cb){
                http('POST', '/jx/updateMjx', null, postdata, cb);
            },

            'delete' : function(id, cb){
                http('POST', '/jx/deleteMjx', {id : id}, null, cb);
            }
        }
    }]);