'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:SearchCtrl
 * @description
 * # SearchCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('SearchCtrl', function ($scope, $http) {
    $scope.$parent.curTab = 'search';
    $scope.filter = {
      startDate: '',
      endDate: ''
    };

    $http({method: 'GET', url: 'http://12.0.0.2:8000/search?limit=50&startDate=' + $scope.filter.startDate + '&endDate=' + $scope.filter.endDate}).
      success(function(data, status, headers, config) {
        console.log(data);
        $scope.messages = data;
      }).
      error(function(data, status, headers, config) {
        console.error(data);
        // called asynchronously if an error occurs
        // or server returns response with an error status.
      });



  });
