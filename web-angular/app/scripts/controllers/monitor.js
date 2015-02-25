'use strict';

/**
 * @ngdoc function
 * @name webAngularApp.controller:MonitorCtrl
 * @description
 * # MonitorCtrl
 * Controller of the webAngularApp
 */
angular.module('webAngularApp')
  .controller('MonitorCtrl', function ($scope, $websocket) {
    var ws = null;
    $scope.$parent.curTab = 'monitor';
    $scope.messages = [];
    $scope.watch = '';
    $scope.limit = 300;
    $scope.connected = false;
    $scope.err = '';

    $scope.connect = function() {
      if($scope.watch == '') {
        $scope.watch = 'invite'
      }
      ws = $websocket('ws://12.0.0.2:8000/ws/' + $scope.watch);

      ws.onOpen(function(event){
        $scope.connected = true;
      });

      ws.onError(function(event){
        $scope.connected = false;
        $scope.err = event;
      });

      ws.onMessage(function(event) {
        console.log(event)
        $scope.messages.push(JSON.parse(event.data))
        if($scope.messages.length > 300) {
          $scope.message.shift()
        }
      })
      ws.onClose(function() {
        $scope.disconnect()
      })
    }

    $scope.disconnect = function() {
      $scope.connected = false;
      if (ws != null) {
        ws.close();
      }
    }
    


  });
