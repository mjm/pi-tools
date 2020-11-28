import SwiftUI
import Combine

class AppModel: ObservableObject {
    enum Event: Identifiable, CustomStringConvertible {
        case beaconEvent(UUID, BeaconObserver.Event)
        case tripsEvent(UUID, TripsController.Event)
        case recorderEvent(UUID, TripRecorder.Event)

        var id: UUID {
            switch self {
            case .beaconEvent(let id, _):
                return id
            case .tripsEvent(let id, _):
                return id
            case .recorderEvent(let id, _):
                return id
            }
        }

        var description: String {
            switch self {
            case .beaconEvent(_, let evt):
                switch evt {
                case .entered:
                    return "Entered beacon region"
                case .exited:
                    return "Exited beacon region"
                }
            case .tripsEvent(_, let evt):
                switch evt {
                case .tripBegan(let trip):
                    return "Started trip \(trip.id)"
                case .tripEnded(let queuedTrips):
                    return "Ended trip with \(queuedTrips.count) trips to record"
                case .tripDiscarded(let trip):
                    return "Discarded \(trip.leftAt.distance(to: trip.returnedAt!)) second trip"
                }
            case .recorderEvent(_, let evt):
                switch evt {
                case .recorded(let trips):
                    return "Recorded \(trips.count) trips"
                case .recordFailed(let err):
                    return "Failed to record trips: \(err)"
                }
            }
        }
    }

    let tripsController: TripsController
    let tripRecorder: TripRecorder
    @Published var allEvents: [Event] = []
    @Published var currentTrip: Trip?
    @Published var queuedTripCount: Int = 0

    init(
        beaconObserver: BeaconObserver,
        tripsController: TripsController,
        tripRecorder: TripRecorder
    ) {
        self.tripsController = tripsController
        self.tripRecorder = tripRecorder

        let wrappedBeaconEvents = beaconObserver.eventsPublisher().map { Event.beaconEvent(UUID(), $0) }
        let wrappedTripsEvents = tripsController.eventsPublisher().map { Event.tripsEvent(UUID(), $0) }
        let wrappedRecorderEvents = tripRecorder.eventsPublisher().map { Event.recorderEvent(UUID(), $0) }

        wrappedBeaconEvents.merge(with: wrappedTripsEvents, wrappedRecorderEvents).scan([]) { events, nextEvent in
            [nextEvent] + events
        }.assign(to: &$allEvents)

        tripsController.$currentTrip.assign(to: &$currentTrip)
        tripsController.$queuedTrips.map(\.count).assign(to: &$queuedTripCount)
    }

    func beginTrip() {
        tripsController.beginTrip()
    }

    func endTrip() {
        tripsController.endTrip()
    }

    func setRecordToDevServer(_ useDev: Bool) {
        tripRecorder.setUpClient(useDevServer: useDev)
    }

    func recordQueuedTrips() {
        tripRecorder.recordTrips(tripsController.queuedTrips)
    }

    func clearQueuedTrips() {
        tripsController.clearQueue()
    }
}
