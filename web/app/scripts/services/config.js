'use strict';

/**
 * @ngdoc service
 * @name webAngularApp.Config
 * @description
 * # Config
 * Service in the webAngularApp.
 */
angular.module('webAngularApp')
  .service('Config', function () {
    var config = {
      urlBase: 'http://12.0.0.2:8000' 

    }
    return config;
  });
