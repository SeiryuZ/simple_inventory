'use strict';

angular.module('simpleInventoryApp')
  .factory('Auth', function ($http, $rootScope, $cookieStore) {


    var currentUser = $cookieStore.get('user') || {username: ''};

    // Public API here
    return {

      isLoggedIn: function(user) {
        if (user === undefined){
          user = $rootScope.user;
          if (user === undefined){
            return false;
          }
        }

        return true;
      },

      login: function(user, success, error) {
        console.log( 'LOGGING IN', user);
        $http.post('/api/login', user).success(function(user){
          $rootScope.user = user;
          success(user);
        }).error(error);
      },

      user: currentUser
    };
  });
