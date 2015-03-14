'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:SettingsCtrl
 * @description
 * # SettingsCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('SettingsCtrl', function ($scope, $http) {
    $scope.$parent.curTab = 'settings';
    $scope.aliases = []
    $scope.newalias = {
      key: '',
      val: ''
    };
    $scope.errors = [];

    var urlBase = "/settings"


    $scope.getAliases = function() {
      var url = urlBase + "/alias";
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
          $scope.aliases = mydata;
        }).
        error(function(data, status, headers, config) {
          console.error(data);
          // called asynchronously if an error occurs
          // or server returns response with an error status.
        });
    }


    $scope.addAlias = function() {
      if(!$scope.newalias.key || !$scope.newalias.val) {
        $scope.errors.push({
          type:'danger',
          message: 'Both fields required to add an alias'
        });

        return
      }
      var url = urlBase + "/alias";
      console.log($scope.newalias)
      $http({
          method: 'POST',
          url: url,
          data: $.param($scope.newalias),
          headers: {'Content-Type': 'application/x-www-form-urlencoded'}
      }).success(function(data, status, headers, config) {
          $scope.aliases.push({ Key:$scope.newalias.key, Val: $scope.newalias.val });
          $scope.newalias = {
            key: '',
            val: ''
          };
   
        }).
        error(function(data, status, headers, config) {
          console.error(data);
          // called asynchronously if an error occurs
          // or server returns response with an error status.
        });

    }

    $scope.deleteAlias = function(index, key) {
      var url = urlBase + "/alias/" + key;
      console.log($scope.newalias)
      $http({
          method: 'DELETE',
          url: url,
          data: {}
      }).success(function(data, status, headers, config) {
          $scope.aliases.splice(index, 1) 
        }).
        error(function(data, status, headers, config) {
          console.error(data);
          // called asynchronously if an error occurs
          // or server returns response with an error status.
        });

    }

    $scope.getAliases();
  });
