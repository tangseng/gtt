angular.module('zfh').controller('ctrl.month', ['$scope', '$http', '$timeout', 'myMonth', 'myTime', function($scope, $http, $timeout, myMonth, myTime){
    var vm = $scope.vm = {};

    $.extend(vm, {
        persons : {},
        tmpPersons : {},
        group : {},
        days : null,
        plans : [],

        init : function(){
            if(Group){
                vm.group = Group;
                if(Persons){
                    $.each(vm.group, function(g, _){
                        if(!vm.persons[g]){
                            vm.persons[g] = {};
                        }
                    });
                    $.each(Persons, function(_, person){
                        vm.tmpPersons[person.id] = person;
                        vm.persons[person.group][person.id] = person;
                    });
                }
            }
            if(Y){
                var days = myMonth(Y, M, D);
                vm.days = days;
                vm.count = days.length;
                vm.qian = days[0].date;
                vm.hou = days[days.length - 1].date;

                if(Persons) {
                    vm.http("/month/plan", {start: myTime.time(vm.qian), end: myTime.time(vm.hou)}, function (data) {
                        var plans = {};
                        data && !data.error && $.each(data, function(_, dd){
                            if(!plans[dd.uid]){
                                plans[dd.uid] = {};
                            }
                            plans[dd.uid][dd.id] = dd;
                        });
                        var gplans = {};
                        $.each(vm.persons, function(g, gp){
                            if(!gplans[g]){
                                gplans[g] = {};
                            }
                            !$.isEmptyObject(gp) && $.each(gp, function(_, p){
                                var _plans = plans[p.id];
                                _plans && $.each(_plans, function(index, plan){
                                    _plans[index]['color'] = p.color;
                                });
                                gplans[g][p.id] = _plans;
                            });
                        });
                        $.each(gplans, function(_, gp){
                            gp && $.each(gp, function(_, p){
                                p && $.each(p, function(_, _p){
                                    vm.plans.push(_p);
                                });
                            });
                        });
                    });

                }
            }
        },

        http : function(url, params, cb){
            $http({
                method : "GET",
                url : url,
                params : params,
                responseType : "json"
            }).success(function(data, status){
                if(data.error){

                } else {
                    cb && cb(data);
                }
            });
        }
    });

    vm.init();
}]);