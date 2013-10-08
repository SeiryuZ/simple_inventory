'use strict';

angular.module('simpleInventoryApp', ['restangular', 'ngCookies'])
  .config(function ($routeProvider, $provide,  $httpProvider, RestangularProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/admin', {
        templateUrl: 'views/admin.html',
        controller: 'AdminCtrl'
      })
      .when('/admin/products', {
        templateUrl: 'views/product.html',
        controller: 'ProductCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });



    $provide.factory('myHttpInterceptor', function($q, $location) {
      return function(promise) {
        return promise.then(function(response) {
          // do nothing on success
          return response;
        }, function(response) {

          // redirect to login page on forbidden / unauthorized error
          if (response.status === 401 || response.status === 403) {
            $location.path('/admin');
          }
          return $q.reject(response);

        });
      };
    });

    $httpProvider.responseInterceptors.push('myHttpInterceptor');


    RestangularProvider.setRestangularFields({
      id: 'ID'
    });
  });



