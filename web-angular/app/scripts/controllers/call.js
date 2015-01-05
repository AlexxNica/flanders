'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:CallCtrl
 * @description
 * # CallCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('CallCtrl', function ($scope, $routeParams, $http) {

    $scope.callId = $routeParams.callId;
    $scope.messages = [];


    $http({method: 'GET', url: 'http://12.0.0.2:8000/call/' + $scope.callId}).
      success(function(data, status, headers, config) {
        console.log(data);
        $scope.messages = data;
      }).
      error(function(data, status, headers, config) {
        console.error(data);
        // called asynchronously if an error occurs
        // or server returns response with an error status.
      });

    bonsai.run(document.getElementById('movie'), {
      code: function() {
        new Rect(10, 10, 100, 100)
          .addTo(stage)
          .attr('fillColor', 'green');
      },
      width: 500,
      height: 400
    });


  });
