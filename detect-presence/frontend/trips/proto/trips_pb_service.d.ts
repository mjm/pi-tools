// package: 
// file: trips.proto

import * as trips_pb from "./trips_pb";
import {grpc} from "@improbable-eng/grpc-web";

type TripsServiceListTrips = {
  readonly methodName: string;
  readonly service: typeof TripsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof trips_pb.ListTripsRequest;
  readonly responseType: typeof trips_pb.ListTripsResponse;
};

export class TripsService {
  static readonly serviceName: string;
  static readonly ListTrips: TripsServiceListTrips;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class TripsServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listTrips(
    requestMessage: trips_pb.ListTripsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: trips_pb.ListTripsResponse|null) => void
  ): UnaryResponse;
  listTrips(
    requestMessage: trips_pb.ListTripsRequest,
    callback: (error: ServiceError|null, responseMessage: trips_pb.ListTripsResponse|null) => void
  ): UnaryResponse;
}

