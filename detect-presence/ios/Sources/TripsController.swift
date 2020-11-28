import Foundation
import Combine

class TripsController {
    enum Event {
        case tripBegan(Trip)
        case tripEnded([Trip])
        case tripDiscarded(Trip)
    }

    @Published var currentTrip: Trip?
    @Published var queuedTrips: [Trip] = []

    @Published private var state = State()
    private let eventsSubject = PassthroughSubject<Event, Never>()
    private var cancellables = Set<AnyCancellable>()

    init<P: Publisher>(events: P) where P.Output == BeaconObserver.Event, P.Failure == Never {
        restoreState()

        events.sink { [weak self] event in
            guard let self = self else { return }

            switch event {
            case .entered:
                self.endTrip()
            case .exited:
                self.beginTrip()
            }
        }.store(in: &cancellables)

        $state.map(\.currentTrip).assign(to: &$currentTrip)
        $state.map(\.queuedTrips).assign(to: &$queuedTrips)
    }

    func eventsPublisher() -> AnyPublisher<Event, Never> {
        eventsSubject.eraseToAnyPublisher()
    }

    private func restoreState() {
        do {
            let stateFileURL = try savedStateURL()
            let savedStateData = try Data(contentsOf: stateFileURL)
            state = try PropertyListDecoder().decode(State.self, from: savedStateData)
            NSLog("Restored state: \(state)")
        } catch {
            NSLog("Could not restore state: \(error.localizedDescription)")
        }
    }

    private func saveState() {
        do {
            let stateFileURL = try savedStateURL()
            let stateData = try PropertyListEncoder().encode(state)
            try stateData.write(to: stateFileURL, options: .atomicWrite)
        } catch {
            NSLog("Could not save state: \(error.localizedDescription)")
        }
    }

    private func savedStateURL() throws -> URL {
        let url = try FileManager.default.url(for: .applicationSupportDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
        return url.appendingPathComponent("SavedState.plist")
    }

    func beginTrip() {
        guard state.currentTrip == nil else {
            NSLog("Already have a current trip, so not beginning a new one.")
            return
        }

        state.currentTrip = Trip()
        saveState()

        eventsSubject.send(.tripBegan(state.currentTrip!))
    }

    func endTrip() {
        guard var trip = state.currentTrip else {
            NSLog("No trip in progress, nothing to end.")
            return
        }

        trip.returnedAt = Date()
        state.currentTrip = nil

        let duration = trip.leftAt.distance(to: trip.returnedAt!)
        guard duration > 60 else {
            NSLog("Discarding trip that only lasted \(duration) seconds")
            saveState()
            eventsSubject.send(.tripDiscarded(trip))
            return
        }

        state.queuedTrips.append(trip)
        saveState()

        eventsSubject.send(.tripEnded(state.queuedTrips))
    }

    func removeFromQueue(_ trips: [Trip]) {
        state.queuedTrips = state.queuedTrips.filter { trip in
            !trips.contains(where: { $0.id == trip.id })
        }
        saveState()
    }

    func clearQueue() {
        state.queuedTrips.removeAll()
        saveState()
    }
}
