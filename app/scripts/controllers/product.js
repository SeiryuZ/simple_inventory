'use strict';
/* global _ */

angular.module('simpleInventoryApp')
  .controller('ProductCtrl', function ($scope, Restangular, $timeout) {

    var baseProducts = Restangular.all('api/products');


    // var validate_product = function () {
    //     $scope.product
    // }

    var hideAlert = function() {
        $timeout(function(){
            $scope.showAlert = false;
          }, 2000);
      };


    $scope.products = [];
    $scope.newProduct = {};

    $scope.showAddForm = false;

    $scope.showAlert = false;
    $scope.alertType = '';
    $scope.alert = '';

    $scope.listProduct = function() {
        console.log('list product');
        baseProducts.getList().then(function(products){
          $scope.products = products;
        });
      };

    $scope.createNewProduct  = function() {
      console.log('creating new product');
      baseProducts.post($scope.newProduct).then(function(product) {
        $scope.showAddForm = false;

        $scope.products.push(product);
        $scope.newProduct = {};

        $scope.alert = 'Product has been created';
        $scope.alertType = 'success';
        $scope.showAlert = true;
        hideAlert();

      }, function() {

        $scope.alert = 'There\'s something wrong';
        $scope.alertType = 'alert';
        $scope.showAlert = true;
        hideAlert();

      });
    };

    $scope.showAddProductForm = function() {
        $scope.showAddForm = !$scope.showAddForm;
      };

    $scope.deleteProduct = function(product) {
      var confirm = window.confirm('Are you sure? Product will be deleted');

      if (confirm) {
        var targetProduct = _.find($scope.products, function(searchedProduct){
          return searchedProduct.ID === product.ID;
        });

        targetProduct.id = targetProduct.ID;
        targetProduct.remove().then(function(){
          $scope.products = _.without($scope.products, targetProduct);
        });

        $scope.alert = 'Product has been deleted';
        $scope.alertType = 'success';
        $scope.showAlert = true;
        hideAlert();

      }
    };

    $scope.listProduct();
  });
