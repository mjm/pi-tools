import SwiftUI
import Combine
import RelaySwiftUI

private let beaconObserver = BeaconObserver()
private let tripsController = TripsController(events: beaconObserver.eventsPublisher())
private let tripRecorder = TripRecorder(events: tripsController.eventsPublisher())

@main
struct PresenceApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var delegate

    @AppStorage("recordToDevServer") private var recordToDevServer: Bool = false
    @StateObject var model = AppModel(
        beaconObserver: beaconObserver,
        tripsController: tripsController,
        tripRecorder: tripRecorder
    )

    var body: some Scene {
        WindowGroup {
            ContentView(model: model)
                .relayEnvironment(model.environment)
                .onChange(of: recordToDevServer) { value in
                    model.setRecordToDevServer(value)
                }
                .onAppear {
                    model.setRecordToDevServer(recordToDevServer)
                }
        }
    }
}

class AppDelegate: NSObject, UIApplicationDelegate {
    var cancellables = Set<AnyCancellable>()

    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey : Any]? = nil) -> Bool {
        tripRecorder.eventsPublisher().sink { event in
            if case .recorded(let trips) = event {
                tripsController.removeFromQueue(trips)
            }
        }.store(in: &cancellables)
        return true
    }
}
