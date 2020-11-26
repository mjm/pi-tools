import Foundation
import detect_presence_proto_trips_trips_proto

struct State: Codable {
    var currentTrip: Trip?
    var queuedTrips: [Trip] = []
}

struct Trip: Codable {
    var id: UUID = UUID()
    var leftAt: Date = Date()
    var returnedAt: Date?

    var asProto: detect_presence_proto_trips_trips_proto.Trip {
        .with { trip in
            trip.id = id.uuidString
            trip.leftAt = dateFormatter.string(from: leftAt)
            if let returnedAt = returnedAt {
                trip.returnedAt = dateFormatter.string(from: returnedAt)
            }
        }
    }
}

private let dateFormatter = ISO8601DateFormatter()
