angular.module('zfh').controller('ctrl.gg', ['$scope', '$http', 'myTime', function($scope, $http, myTime){
    var vm = $scope.vm = {};

    $.extend(vm, {
        ggs : [],
        offset : 0,

        init : function(){
            this.ajax(0);
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
        }
    });

    vm.init();
}]);