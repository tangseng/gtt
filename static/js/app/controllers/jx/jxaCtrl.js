angular.module('zfh')
    .controller('ctrl.jxa', ['$scope', 'myXX', 'myTime', function($scope, myXX, myTime){
        var vm = $scope.vm = {};

        $.extend(vm, {
            all : {},
            date : '',

            init : function(){
                $scope.$watch('vm.date', function(newValue){
                    var time = new Date(newValue);
                    var y = time.getFullYear();
                    var m = time.getMonth() + 1;
                    var d = time.getDate();
                    vm.ajax(y, m, d);
                });
                if(Y){
                    vm.date = Y + '-' + M + '-' + D;
                }
            },

            ajax : function(y, m, d){
                myXX({
                    y : y,
                    m : m,
                    d : d
                }, function(data){
                    var tmp = [];
                    var self = null;
                    $.each(data, function(_, day){
                        var jx = [];
                        day.works && $.each(day.works, function(_, w){
                            jx.push({
                                content : w.content,
                                start : myTime.hm(w.startTime),
                                end : myTime.hm(w.endTime),
                                status : w.status,
                                time : myTime.timeToStr(w.time * 1000, true)
                            });
                        });
                        var item = {
                            person : Persons[day.uid],
                            jx : jx
                        };
                        if(day.uid == User.id){
                            self = item;
                            return;
                        }
                        tmp.push(item);
                    });
                    self && tmp.push(self);
                    vm.all = tmp;
                });
            }
        });

        vm.init();
    }])
    .factory('myXX', ['$http', function($http){
        return function(params, cb){
            $http({
                method : "GET",
                url : '/jx/xxs',
                params : params,
                responseType : "json"
            }).success(function(data){
                if(data && !data.error){
                    cb && cb(data, false);
                }
            });
        }
    }]);