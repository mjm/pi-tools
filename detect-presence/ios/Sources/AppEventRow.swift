import SwiftUI
import RelaySwiftUI
import detect_presence_ios_relay_generated

private let eventFragment = graphql("""
fragment AppEventRow_event on AppEvent {
    id
    timestamp

    ...on BeaconEvent {
        action
    }
    ...on TripBeganEvent {
        trip {
            id
        }
    }
    ...on TripEndedEvent {
        queuedTrips {
            id
        }
    }
    ...on TripDiscardedEvent {
        trip {
            id
            leftAt
            returnedAt
        }
    }
    ...on RecordedTripsEvent {
        recordedTrips {
            id
        }
    }
    ...on RecordFailedEvent {
        message
    }
}
""")

struct AppEventRow: View {
    @Fragment<AppEventRow_event> var event

    var body: some View {
        if let event = event {
            VStack(alignment: .leading, spacing: 8) {
                if let beaconEvent = event.asBeaconEvent {
                    switch beaconEvent.action {
                    case .entered:
                        Text("Entered beacon region")
                            .font(.body)
                    case .exited:
                        Text("Exited beacon region")
                            .font(.body)
                    }
                }

                if let trip = event.asTripBeganEvent?.trip {
                    Text("Started trip ").font(.body) +
                        (Text(verbatim: trip.id)
                            .font(.system(.body, design: .monospaced)))
                }

                if let queuedTrips = event.asTripEndedEvent?.queuedTrips {
                    (Text("Ended trip with ") +
                        (Text("\(queuedTrips.count) trips")
                            .bold()) +
                        Text(" to record"))
                        .font(.body)
                }

                if let trip = event.asTripDiscardedEvent?.trip,
                   let leftAt = trip.leftAt.asDate,
                   let returnedAt = trip.returnedAt?.asDate {
                    (Text("Discarded ") +
                        (Text("\(Int(leftAt.distance(to: returnedAt))) second").bold()) +
                        Text(" trip"))
                        .font(.body)
                }

                if let trips = event.asRecordedTripsEvent?.recordedTrips {
                    (Text("Recorded ") +
                        (Text("\(trips.count) trips").bold()))
                        .font(.body)
                }

                if let message = event.asRecordFailedEvent?.message {
                    Text("Failed to record trips:").font(.body)
                    Text(verbatim: message).font(.system(.body, design: .monospaced))
                }

                (Text(event.timestamp.asDate!, style: .relative) + Text(" ago"))
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
        }
    }
}
