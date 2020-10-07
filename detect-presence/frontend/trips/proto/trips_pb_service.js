// package: 
// file: trips.proto

var trips_pb = require("./trips_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var TripsService = (function () {
  function TripsService() {}
  TripsService.serviceName = "TripsService";
  return TripsService;
}());

TripsService.ListTrips = {
  methodName: "ListTrips",
  service: TripsService,
  requestStream: false,
  responseStream: false,
  requestType: trips_pb.ListTripsRequest,
  responseType: trips_pb.ListTripsResponse
};

exports.TripsService = TripsService;

function TripsServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

TripsServiceClient.prototype.listTrips = function listTrips(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TripsService.ListTrips, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.TripsServiceClient = TripsServiceClient;

