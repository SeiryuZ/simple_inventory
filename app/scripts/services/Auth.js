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
        $http.post('/login', user).success(function(user){
          $rootScope.user = user;
          success(user);
        }).error(error);
      },

      logout: function(user, success, error) {
        console.log( 'LOGGING OUT');
        $http.get('/logout').success(success).error(error);
      },


      register: function(user, success, error) {
        console.log( 'REGISTERING', user);
        $http.post('/register', user).success(success).error(error);
      },

      user: currentUser
    };
  });
