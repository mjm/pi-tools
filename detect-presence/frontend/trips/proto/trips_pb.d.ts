// package: 
// file: trips.proto

import * as jspb from "google-protobuf";

export class ListTripsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTripsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTripsRequest): ListTripsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListTripsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTripsRequest;
  static deserializeBinaryFromReader(message: ListTripsRequest, reader: jspb.BinaryReader): ListTripsRequest;
}

export namespace ListTripsRequest {
  export type AsObject = {
  }
}

export class ListTripsResponse extends jspb.Message {
  clearTripsList(): void;
  getTripsList(): Array<ListTripsResponse.Trip>;
  setTripsList(value: Array<ListTripsResponse.Trip>): void;
  addTrips(value?: ListTripsResponse.Trip, index?: number): ListTripsResponse.Trip;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTripsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTripsResponse): ListTripsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListTripsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTripsResponse;
  static deserializeBinaryFromReader(message: ListTripsResponse, reader: jspb.BinaryReader): ListTripsResponse;
}

export namespace ListTripsResponse {
  export type AsObject = {
    tripsList: Array<ListTripsResponse.Trip.AsObject>,
  }

  export class Trip extends jspb.Message {
    getId(): string;
    setId(value: string): void;

    getLeftAt(): string;
    setLeftAt(value: string): void;

    getReturnedAt(): string;
    setReturnedAt(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Trip.AsObject;
    static toObject(includeInstance: boolean, msg: Trip): Trip.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Trip, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Trip;
    static deserializeBinaryFromReader(message: Trip, reader: jspb.BinaryReader): Trip;
  }

  export namespace Trip {
    export type AsObject = {
      id: string,
      leftAt: string,
      returnedAt: string,
    }
  }
}

