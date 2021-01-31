import SwiftUI
import Combine
import Relay

class AppModel: ObservableObject {
    let tripsController: TripsController
    let tripRecorder: TripRecorder
    @Published var queuedTripCount: Int = 0
    @Published var environment: Relay.Environment!

    private var cancellables = Set<AnyCancellable>()

    init(
        beaconObserver: BeaconObserver,
        tripsController: TripsController,
        tripRecorder: TripRecorder
    ) {
        self.tripsController = tripsController
        self.tripRecorder = tripRecorder

        self.setRecordToDevServer(false)

        beaconObserver.eventsPublisher().sink { [weak self] event in
            self?.record(event)
        }.store(in: &cancellables)
        tripsController.eventsPublisher().sink { [weak self] event in
            self?.record(event)
        }.store(in: &cancellables)
        tripRecorder.eventsPublisher().sink { [weak self] event in
            self?.record(event)
        }.store(in: &cancellables)

        tripsController.$currentTrip.sink { [weak self] trip in
            self?.setCurrentTrip(trip)
        }.store(in: &cancellables)

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
        self.environment.commitUpdate { store in
            store.root.setLinkedRecords("queuedTrips", records: [])
            store.root.setLinkedRecords("appEvents", records: [])
        }

        self.tripRecorder.environment = self.environment
    }

    func recordQueuedTrips() {
        tripRecorder.recordTrips(tripsController.queuedTrips)
    }

    func clearQueuedTrips() {
        tripsController.clearQueue()
    }

    private func record(_ beaconEvent: BeaconObserver.Event) {
        environment.commitUpdate { store in
            let event = store.createEvent(typeName: "BeaconEvent")
            switch beaconEvent {
            case .entered:
                event["action"] = "ENTERED"
            case .exited:
                event["action"] = "EXITED"
            }
            store.prependEvent(event)
        }
    }

    private func record(_ tripsEvent: TripsController.Event) {
        environment.commitUpdate { store in
            let event: RecordProxy

            switch tripsEvent {
            case .tripBegan(let trip):
                event = store.createEvent(typeName: "TripBeganEvent")
                event.setLinkedRecord("trip", record: store.upsert(trip))
            case .tripEnded(let queuedTrips):
                event = store.createEvent(typeName: "TripEndedEvent")
                event.setLinkedRecords("queuedTrips", records: queuedTrips.map {
                    store.upsert($0)
                })
            case .tripDiscarded(let trip):
                event = store.createEvent(typeName: "TripDiscardedEvent")
                event.setLinkedRecord("trip", record: store.upsert(trip))
            }

            store.prependEvent(event)
        }
    }

    private func record(_ recordEvent: TripRecorder.Event) {
        environment.commitUpdate { store in
            let event: RecordProxy

            switch recordEvent {
            case .recorded(let trips):
                event = store.createEvent(typeName: "RecordedTripsEvent")
                event.setLinkedRecords("recordedTrips", records: trips.map { store.upsert($0) })
            case .recordFailed(let message):
                event = store.createEvent(typeName: "RecordFailedEvent")
                event["message"] = message
            }

            store.prependEvent(event)
        }
    }

    private func setCurrentTrip(_ trip: Trip?) {
        environment.commitUpdate { store in
            if let trip = trip {
                store.root.setLinkedRecord("currentTrip", record: store.upsert(trip))
            } else {
                store.root["currentTrip"] = NSNull()
            }
        }
    }
}

private extension RecordSourceProxy {
    mutating func createEvent(typeName: String) -> RecordProxy {
        let eventID = UUID().uuidString
        let event = create(dataID: DataID(eventID), typeName: typeName)
        event["id"] = eventID
        event["timestamp"] = Date().asString
        return event
    }

    mutating func upsert(_ trip: Trip) -> RecordProxy {
        let tripID = trip.id.uuidString
        let tripRecord = self[DataID(tripID)] ?? create(dataID: DataID(tripID), typeName: "Trip")
        tripRecord["id"] = trip.id.uuidString
        tripRecord["leftAt"] = trip.leftAt.asString
        if let returnedAt = trip.returnedAt?.asString {
            tripRecord["returnedAt"] = returnedAt
        }
        return tripRecord
    }

    mutating func prependEvent(_ event: RecordProxy) {
        var events = root.getLinkedRecords("appEvents") ?? []
        events.insert(event, at: 0)
        root.setLinkedRecords("appEvents", records: events)
    }
}
