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
    $scope.$parent.curTab = 'search';
    $scope.callId = $routeParams.callId;
    $scope.messages = [];
    var settings = {
      headerSpacing: 50,
      headerWidth: 200,
      headerHeight: 50,
      rowHeight: 60,
      arrowSize: 5,
      textHeight: 13,
      bottomPadding: 20,
    }
    var height = 280;
    var width = 1080;
    var movie;


    var generateImage = function(width, height) {
      movie = bonsai.run(document.getElementById('sipimage'), {
        code: function() {
          stage.on('message', function(data) {
            var columnHeaders = {}
            var messages = data.messages;
            var sortIndex = 0;
            var headerSpacing = data.settings.headerSpacing;
            var headerWidth = data.settings.headerWidth;
            var headerHeight = data.settings.headerHeight;
            var rowHeight = data.settings.rowHeight;
            var arrowSize = data.settings.arrowSize;
            var textHeight = data.settings.textHeight;
            var bottomPadding = data.settings.bottomPadding;

            messages.forEach(function(message, index) {
              if(!columnHeaders.hasOwnProperty(message.SourceIp)) {
                columnHeaders[message.SourceIp] = sortIndex;
                sortIndex++;
              }
              if(!columnHeaders.hasOwnProperty(message.DestinationIp)) {
                columnHeaders[message.DestinationIp] = sortIndex;
                sortIndex++;
              }
            });
            var columnHeadersArray = []
            for (var header in columnHeaders) {
              columnHeadersArray[columnHeaders[header]] = header;
            }
            columnHeadersArray.forEach(function(header, index){
              // Draw header boxes (servers)
              // var headerSquare = new Rect(index * (headerSpacing + headerWidth), 0, headerWidth, headerHeight, 5);
              // headerSquare.stroke('#000', 2);
              // headerSquare.addTo(stage)
              

              var headerText = new Text(header).attr({
                'x': (index * (headerSpacing + headerWidth)) + (headerWidth/2),
                'y': 20,
                textAlign: 'center'
              });
              headerText.addTo(stage);

              // Draw vertical lines
              var lineStartX = (index * (headerSpacing + headerWidth)) + (headerWidth/2);

              var verticalLine = new Path()
                .moveTo(lineStartX, headerHeight)
                .lineTo(lineStartX, headerHeight + (messages.length * rowHeight) + bottomPadding)
                .closePath()
                .stroke('gray', 1)
                .addTo(stage);

            });

            var headerLine = new Path()
                .moveTo(0, headerHeight)
                .lineTo(columnHeadersArray.length * (headerWidth + headerSpacing), headerHeight)
                .stroke('gray',2)
                .addTo(stage)

            messages.forEach(function(message, index) {
              var lineStartX = (columnHeaders[message.SourceIp] * (headerSpacing + headerWidth)) + (headerWidth/2)
              var lineEndX = (columnHeaders[message.DestinationIp] * (headerSpacing + headerWidth)) + (headerWidth/2)
              var lineStartY = index * rowHeight + headerHeight + rowHeight;
              var leftOrRight = 1

              if(lineStartX > lineEndX) {
                leftOrRight = -1;
              }

              var messageLine = new Path()
                .moveTo(lineStartX, lineStartY)
                .lineTo(lineEndX, lineStartY)
                .closePath()
                .stroke('black', 1)
                .fill('black')
                .addTo(stage);

              var arrow = new Path()
                .moveTo(lineEndX, lineStartY)
                .lineTo(lineEndX - arrowSize*leftOrRight, lineStartY + arrowSize)
                .lineTo(lineEndX - arrowSize*leftOrRight, lineStartY - arrowSize)
                .lineTo(lineEndX, lineStartY)
                .closePath()
                .stroke('blue', 1)
                .fill('blue')
                .addTo(stage)

              var sourcePortText = new Text(message.SourcePort).attr({
                'x': lineStartX - 5*leftOrRight,
                'y': lineStartY - textHeight/2,
                textAlign: (leftOrRight === 1)?'right':'left',
                textFillColor: 'gray',
                fontSize: textHeight,
              }).addTo(stage);

              var destinationPortText = new Text(message.DestinationPort).attr({
                'x': lineEndX + 5*leftOrRight,
                'y': lineStartY - textHeight/2,
                textAlign: (leftOrRight === -1)?'right':'left',
                textFillColor: 'gray',
                fontSize: textHeight,
              }).addTo(stage);

              var methodTextColor = 'blue';
              methodInt = parseInt(message.Method);
                if(methodInt >= 100 && methodInt <= 299) {
                  methodTextColor = 'green';
                }
                else if (methodInt >= 400) {
                  methodTextColor = 'red';
                }
              

              var methodText = new Text(message.Method + " " + message.ReplyReason).attr({
                'x': ((leftOrRight === 1)?lineStartX:lineEndX) + 30,
                'y': lineStartY - textHeight,
                textAlign: 'left',
                textFillColor: methodTextColor,
                fontSize: textHeight,
              }).addTo(stage);

              var timeText = new Text(message.Datetime).attr({
                'x': lineStartX + ((headerWidth+headerSpacing)/2)*leftOrRight,
                'y': lineStartY + 4,
                textAlign: 'center',
                textFillColor: "#555",
                fontSize: textHeight,
              }).addTo(stage);

            });
          });
          stage.sendMessage('ready', {});
        },
        width: width,
        height: height
      });
      return movie;
    }
    
    $http({method: 'GET', url: '/call/' + $scope.callId}).
      success(function(data, status, headers, config) {
        console.log(data);
        $scope.messages = data;
        height = data.length * settings.rowHeight + settings.headerHeight + settings.bottomPadding;
        generateImage(1080,height)
        movie.on('load', function() {
        // receive event from the runner context
        movie.on('message:ready', function() {
          movie.sendMessage({
            messages: $scope.messages,
            settings: settings
          });
        });
      });
      }).
      error(function(data, status, headers, config) {
        console.error(data);
        // called asynchronously if an error occurs
        // or server returns response with an error status.
      });

    
  });
