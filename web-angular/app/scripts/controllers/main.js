'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('MainCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
