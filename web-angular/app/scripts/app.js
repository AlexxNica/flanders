'use strict';

/**
 * @ngdoc overview
 * @name webAngularApp
 * @description
 * # webAngularApp
 *
 * Main module of the application.
 */
angular
  .module('webAngularApp', [
    'ngAnimate',
    'ngCookies',
    'ngMessages',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/search.html',
        controller: 'SearchCtrl'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
      })
      .when('/monitor', {
        templateUrl: 'views/monitor.html',
        controller: 'MonitorCtrl'
      })
      .when('/call/:callId', {
        templateUrl: 'views/call.html',
        controller: 'CallCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });
