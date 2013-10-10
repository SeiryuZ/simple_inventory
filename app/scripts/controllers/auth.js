'use strict';

angular.module('simpleInventoryApp')
  .controller('AuthCtrl', function ($scope, Auth, $http, $location, $rootScope) {
    $scope.user = {username: '', password: ''};

    if ($rootScope.loggedIn === true) {
      $location.path('/admin/products');
    }

    $scope.test = function () {
        console.log('testing restricted resource');
        $http.post('/api/products', {'test': 'qw'}).success(function(){
          console.log('succes restricted');
        }).error(function(){
            console.log('error restricted');
          });
      };

    $scope.logout = function () {
        console.log('logout');
        Auth.logout(function(){
          console.log('succes logout');
          $rootScope.loggedIn = false;
          $location.path('/admin');
        },function(){
            console.log('error logout');
          });
      };

    $scope.login = function() {
        console.log('trying to login');
        Auth.login($scope.user, function(){
            console.log('success login');
            $rootScope.loggedIn = true;
            $location.path('/admin/products');

          }, function(){
            console.log('error login');
          });
      };

    $scope.register = function() {
        console.log('trying to register');
        Auth.register($scope.user, function(){
            console.log('success register');
          }, function(error){
            console.log('error register' + error);
          });
      };
  });
