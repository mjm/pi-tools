import SwiftUI

@main
struct PresenceApp: App {
    @StateObject var beaconObserver = BeaconObserver()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(beaconObserver)
        }
    }
}
