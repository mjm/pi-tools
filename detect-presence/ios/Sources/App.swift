import SwiftUI
import Combine

private let beaconObserver = BeaconObserver()
private let tripsController = TripsController(events: beaconObserver.eventsPublisher())
private let tripRecorder = TripRecorder(events: tripsController.eventsPublisher())

@main
struct PresenceApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var delegate

    var body: some Scene {
        WindowGroup {
            ContentView(tripsController: tripsController)
                .environmentObject(beaconObserver)
        }
    }
}

class AppDelegate: NSObject, UIApplicationDelegate {
    var cancellables = Set<AnyCancellable>()

    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey : Any]? = nil) -> Bool {
        tripRecorder.eventsPublisher().sink { event in
            switch event {
            case .recorded(let trips):
                tripsController.removeFromQueue(trips)
            }
        }.store(in: &cancellables)
        return true
    }
}
