import SwiftUI

private let tripRecorder = TripRecorder()

@main
struct PresenceApp: App {
    @StateObject var beaconObserver = BeaconObserver(tripRecorder: tripRecorder)

    var body: some Scene {
        WindowGroup {
            ContentView(tripRecorder: tripRecorder)
                .environmentObject(beaconObserver)
        }
    }
}
