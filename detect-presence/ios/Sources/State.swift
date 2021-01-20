import Foundation
import detect_presence_ios_relay_generated

struct State: Codable {
    var currentTrip: Trip?
    var queuedTrips: [Trip] = []
}

struct Trip: Codable {
    var id: UUID = UUID()
    var leftAt: Date = Date()
    var returnedAt: Date?

    var asInput: NewTripInput {
        NewTripInput(
            id: id.uuidString,
            leftAt: dateFormatter.string(from: leftAt),
            returnedAt: dateFormatter.string(from: returnedAt!)
        )
    }
}

private let dateFormatter = ISO8601DateFormatter()
