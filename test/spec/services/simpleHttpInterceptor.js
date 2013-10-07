'use strict';

describe('Service: simpleHttpInterceptor', function () {

  // load the service's module
  beforeEach(module('simpleInventoryApp'));

  // instantiate service
  var simpleHttpInterceptor;
  beforeEach(inject(function (_simpleHttpInterceptor_) {
    simpleHttpInterceptor = _simpleHttpInterceptor_;
  }));

  it('should do something', function () {
    expect(!!simpleHttpInterceptor).toBe(true);
  });

});
