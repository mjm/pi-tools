import Foundation

struct AppEvent: Identifiable, CustomStringConvertible {
    var id = UUID()
    var timestamp = Date()
    var data: Data

    init(beaconEvent: BeaconObserver.Event) {
        self.data = .beaconEvent(beaconEvent)
    }

    init(tripsEvent: TripsController.Event) {
        self.data = .tripsEvent(tripsEvent)
    }

    init(recorderEvent: TripRecorder.Event) {
        self.data = .recorderEvent(recorderEvent)
    }

    var description: String {
        data.description
    }

    enum Data: CustomStringConvertible {
        case beaconEvent(BeaconObserver.Event)
        case tripsEvent(TripsController.Event)
        case recorderEvent(TripRecorder.Event)

        var description: String {
            switch self {
            case .beaconEvent(let evt):
                switch evt {
                case .entered:
                    return "Entered beacon region"
                case .exited:
                    return "Exited beacon region"
                }
            case .tripsEvent(let evt):
                switch evt {
                case .tripBegan(let trip):
                    return "Started trip \(trip.id)"
                case .tripEnded(let queuedTrips):
                    return "Ended trip with \(queuedTrips.count) trips to record"
                case .tripDiscarded(let trip):
                    return "Discarded \(trip.leftAt.distance(to: trip.returnedAt!)) second trip"
                }
            case .recorderEvent(let evt):
                switch evt {
                case .recorded(let trips):
                    return "Recorded \(trips.count) trips"
                case .recordFailed(let err):
                    return "Failed to record trips: \(err)"
                }
            }
        }
    }
}
