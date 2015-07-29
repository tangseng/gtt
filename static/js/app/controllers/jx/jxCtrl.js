angular.module('zfh')
    .controller('ctrl.jx', ['$scope', '$http', '$timeout', 'myMonth', 'myMjx', 'myTime', function($scope, $http, $timeout, myMonth, myMjx, myTime){
        var vm = $scope.vm = {};

        $.extend(vm, {
            days : [],
            inited : false,
            ccJxday : [],
            qian : "",
            hou : "",
            mjxs : [],

            init : function(){
                if(Y){
                    vm.days = myMonth(Y, M, D);
                    vm.inited = true;
                    vm.qian = vm.days[0].date;
                    vm.hou = vm.days[vm.days.length - 1].date;

                    myMjx({
                        qian : myTime.time(vm.qian + ' ' + GlobalConfig.start),
                        hou : myTime.time(vm.hou + ' ' + GlobalConfig.end)
                    }, function(data){
                        $.each(data, function(_, dd){
                            vm.mjxs.push($.extend({}, dd, {'date' : myTime.timeToStr(dd.date * 1000)}));
                        });
                    });
                }
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
                        d: d
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
            }
        }
    }])
    .factory('myJxDay', ['$http', '$cacheFactory', function($http, $cacheFactory){
        var cache = $cacheFactory('jxday');
        return function(params, cb){
            var cc = cache.get(params.date);
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
                    cache.put(params.date, data);
                }
            });
        }
    }])
    .factory('myMjx', ['$http', function($http){
        return function(params, cb){
            $http({
                method : "GET",
                url : '/jx/getMjx',
                params : params,
                responseType : "json"
            }).success(function(data){
                if(data && !data.error){
                    cb && cb(data, false);
                }
            });
        }
    }]);