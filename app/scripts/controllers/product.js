'use strict';

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

    $scope.listProduct = function () {
        console.log('list product');
        $scope.products = baseProducts.getList();
      };

    $scope.createNewProduct  = function () {
      console.log('creating new product');
      baseProducts.post($scope.newProduct).then(function() {
        $scope.showAddForm = false;

        $scope.products.push($scope.newProduct);
        $scope.newProduct = {};

        $scope.alert = 'Product has been created';
        $scope.alertType = 'success';
        $scope.showAlert = true;
        hideAlert();

      }, function () {

        $scope.alert = 'There\'s something wrong';
        $scope.alertType = 'alert';
        $scope.showAlert = true;
        hideAlert();

      });
    };

    $scope.showAddProductForm = function () {
        $scope.showAddForm = !$scope.showAddForm;
      };

    $scope.listProduct();
  });
