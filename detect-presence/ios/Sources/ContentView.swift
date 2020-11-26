import SwiftUI

struct ContentView: View {
    @EnvironmentObject var beaconObserver: BeaconObserver
    let tripRecorder: TripRecorder

    var body: some View {
        VStack(spacing: 8) {
            switch beaconObserver.status {
            case .unknown:
                Text("Not sure if you're home")
            case .inside:
                Text("Looks like you're home!")
            case .outside:
                Text("Looks like you're away from home!")
            }

            if let changedTime = beaconObserver.statusChangedTime {
                Text("Transitioned ") + Text(changedTime, style: .relative) + Text(" ago")
            }

            Button("Simulate Begin Trip") {
                tripRecorder.beginTrip()
            }

            Button("Simulate End Trip") {
                tripRecorder.endTrip()
            }
        }
    }
}
