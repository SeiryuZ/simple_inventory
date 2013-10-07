'use strict';

angular.module('simpleInventoryApp')
  .controller('AdminCtrl', function ($scope, Restangular) {

    var baseProducts = Restangular.all('api/products');

    $scope.products = [];

    $scope.newProduct = {};


    $scope.refresh = function () {
      console.log('refreshing list');
      $scope.products = baseProducts.getList();
    };

    $scope.createNewProduct  = function () {
      console.log('creating new product');
      baseProducts.post($scope.newProduct, function(error){
        console.log('success')
        console.log(error);
      }, function(error){
        console.log('error')
                console.log(error);

      });
    };
  });
