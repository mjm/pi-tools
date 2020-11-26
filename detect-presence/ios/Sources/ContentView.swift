import SwiftUI

struct ContentView: View {
    @ObservedObject var model: AppModel

    var body: some View {
        List {
            Section {
                Button {
                    model.beginTrip()
                } label: {
                    Label("Simulate Begin Trip", systemImage: "play.fill")
                }

                Button {
                    model.endTrip()
                } label: {
                    Label("Simulate End Trip", systemImage: "stop.fill")
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
