import SwiftUI

struct ContentView: View {
    @ObservedObject var model: AppModel
    @AppStorage("recordToDevServer") private var recordToDevServer: Bool = false

    var body: some View {
        List {
            Section {
                Toggle("Record to development server", isOn: $recordToDevServer)

                if recordToDevServer {
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
                }

                if let trip = model.currentTrip {
                    Text("Current trip started ") + Text(trip.leftAt, style: .relative) + Text(" ago")
                } else {
                    Text("Not currently on a trip")
                }

                if model.queuedTripCount > 0 {
                    Button {
                        model.recordQueuedTrips()
                    } label: {
                        Label("Record \(model.queuedTripCount) queued trips", systemImage: "icloud.and.arrow.up.fill")
                    }
                }
            }

            Section(header: Text("All Events")) {
                ForEach(model.allEvents) { event in
                    Text(event.description)
                }
            }
        }
        .navigationTitle("Presence")
        .onAppear {
            NSLog("setting record to dev server to \(recordToDevServer)")
            model.setRecordToDevServer(recordToDevServer)
        }
        .onChange(of: recordToDevServer) { useDev in
            NSLog("setting record to dev server to \(useDev)")
            model.setRecordToDevServer(useDev)
        }
    }
}
