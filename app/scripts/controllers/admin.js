'use strict';

angular.module('simpleInventoryApp')
  .controller('AdminCtrl', function ($scope, Auth, $http) {

    $scope.user = {username: '', password: ''};

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
        $http.post('/api/logout', {'test': 'qw'}).success(function(){
          console.log('succes logout');
        }).error(function(){
            console.log('error logout');
          });
      };

    $scope.login = function() {
        console.log('trying to login');
        Auth.login($scope.user, function(){
            console.log('success login');
          }, function(){
            console.log('error login');
          });
      };
  });
