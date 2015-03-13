'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:SettingsCtrl
 * @description
 * # SettingsCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('SettingsCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
