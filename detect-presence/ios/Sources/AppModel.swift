import SwiftUI
import Combine

class AppModel: ObservableObject {
    let tripsController: TripsController
    let tripRecorder: TripRecorder
    @Published var allEvents: [AppEvent] = []
    @Published var currentTrip: Trip?
    @Published var queuedTripCount: Int = 0

    init(
        beaconObserver: BeaconObserver,
        tripsController: TripsController,
        tripRecorder: TripRecorder
    ) {
        self.tripsController = tripsController
        self.tripRecorder = tripRecorder

        let wrappedBeaconEvents = beaconObserver.eventsPublisher().map(AppEvent.init)
        let wrappedTripsEvents = tripsController.eventsPublisher().map(AppEvent.init)
        let wrappedRecorderEvents = tripRecorder.eventsPublisher().map(AppEvent.init)

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
