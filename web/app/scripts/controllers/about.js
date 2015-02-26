'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:AboutCtrl
 * @description
 * # AboutCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('AboutCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
