angular.module('zfh').controller('ctrl.tb', ['$scope', '$http', '$timeout', 'myTime', function($scope, $http, $timeout, myTime){
    var vm = $scope.vm = {};

    $.extend(vm, {
        tbs : [],
        cctb : {},
        shtb : {},
        offset : 0,
        pageNum : 30,
        canP : false,
        canN : false,
        open : true,

        init : function(reset){
            vm.cctb = {
                id : 0,
                app : "",
                title : "",
                desc : "",
                why : "",
                fix : "",
                custom : "",
                dateString : "",
                date : 0
            };

            if(typeof PageNum != undefined){
                vm.pageNum = PageNum;
            }

            if(!reset){
                vm.shtb = {
                    app : "",
                    custom : "",
                    key : ""
                };

                this.search();
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
        },

        checkError : function(){
            var error = false;
            $.each(['title', 'fix'], function(_, sx){
                if(!vm.cctb[sx].length){
                    error = vm.config["null"];
                    return false;
                }
            });
            vm.cctb.date = myTime.time(vm.cctb.dateString);
            return error;
        },

        doAdd : function(){
            var error = vm.checkError();
            if(error){
                vm.ajaxdo(error);
                return false;
            }
            vm.doAjax('/tb/add', function(data){
                vm.tbs.unshift($.extend({}, vm.cctb, {id : data.id, pname : User.name}));
            });
            return false;
        },

        update : function(id){
            id = parseInt(id) || 0;
            $.each(vm.tbs, function(_, tb){
                if(tb.id == id){
                    vm.cctb = $.extend({}, tb);
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
            vm.doAjax('/tb/update', function(data){
                $.each(vm.tbs, function(index, tb){
                    if(tb.id == vm.cctb.id){
                        vm.tbs.splice(index, 1, $.extend({}, vm.cctb));
                    }
                });
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
                dd.data = $.param(vm.cctb);
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

        search : function(){
            vm.doAjax('/tb/search', function(data){
                vm.tbs = [];
                $.each(data.tbs, function(_, tb){
                    vm.tbs.push($.extend({}, tb, {'dateString' : tb.date > 0 ? myTime.timeToStr(tb.date * 1000) : "", pname : data.pps[tb.uid]}));
                });
                if(vm.tbs.length == vm.pageNum){
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
                app : vm.shtb.app,
                custom : vm.shtb.custom,
                key : vm.shtb.key,
                offset : vm.offset
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

        doDelete : function(id){
            if(!confirm('确定要删除么？')){
                return false;
            }
            $http({
                method : "POST",
                url : '/tb/delete',
                data : $.param({id : id}),
                responseType : "json"
            }).success(function(data, status){
                if(data.error){
                    alert(data.error);
                } else {
                    $.each(vm.tbs, function(index, tb){
                        if(tb && tb.id == id){
                            vm.tbs.splice(index, 1);
                            return false;
                        }
                    });
                }
            });
            return false;
        },

        config : {
            "null" : "标题与解决方案是必须填写的"
        }
    });

    vm.init();
}]);