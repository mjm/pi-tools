import SwiftUI

struct ContentView: View {
    @ObservedObject var model: AppModel

    var body: some View {
        List {
            Section {
                Button("Simulate Begin Trip") {
                    model.beginTrip()
                }

                Button("Simulate End Trip") {
                    model.endTrip()
                }

                if let trip = model.currentTrip {
                    Text("Current trip started ") + Text(trip.leftAt, style: .relative) + Text(" ago")
                } else {
                    Text("Not currently on a trip")
                }
            }

            Section(header: Text("All Events")) {
                ForEach(model.allEvents) { event in
                    Text(event.description)
                }
            }
        }
        .navigationTitle("Presence")
    }
}
