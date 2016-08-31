'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:SearchCtrl
 * @description
 * # SearchCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('SearchCtrl', function ($scope, $http, $location) {
    $scope.$parent.curTab = 'search';
    var searchParams = $location.search();
    $scope.filter = {
      startDate: searchParams.startDate || '',
      endDate: searchParams.endDate || '',
      touser: searchParams.touser || '',
      fromuser: searchParams.fromuser || '',
      sourceip: searchParams.sourceip || '',
      destip: searchParams.destip || '',
      callId: searchParams.callId || ''
    };

    var urlBase = "/search?limit=100";

    $scope.search = function() {
      var url = urlBase;
      for(var key in $scope.filter) {
        if($scope.filter[key] != '') {
          url += "&" + key + "=" + $scope.filter[key];
        }
      }
      $http({ method: 'GET', url: url }).
        success(function(data, status, headers, config) {
          var mydata;
          if(!data || data == 'null') {
            mydata = new Array();
          }
          else {
            mydata = data;
          }
          console.log(mydata);
          $scope.messages = mydata;
        }).
        error(function(data, status, headers, config) {
          console.error(data);
          // called asynchronously if an error occurs
          // or server returns response with an error status.
        });
    };

    $scope.search();

    $scope.insertPlusOne = function() {
      //Adding the string for +1 into text field by button click.
      $(function () {
        $('#button').on('click', function () {
          var text = $('#text');
          text.val(text.val() + '%2b1');
        });
      });
    };

    $scope.insertPlusOne();
  });
