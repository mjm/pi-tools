import SwiftUI
import Combine
import Relay

class AppModel: ObservableObject {
    let tripsController: TripsController
    let tripRecorder: TripRecorder
    @Published var allEvents: [AppEvent] = []
    @Published var currentTrip: Trip?
    @Published var queuedTripCount: Int = 0
    @Published var environment: Relay.Environment?

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
        self.objectWillChange.send()
        self.environment = Relay.Environment(
            network: Network(isDevServer: useDev),
            store: Store()
        )
        self.tripRecorder.environment = self.environment
    }

    func recordQueuedTrips() {
        tripRecorder.recordTrips(tripsController.queuedTrips)
    }

    func clearQueuedTrips() {
        tripsController.clearQueue()
    }
}
