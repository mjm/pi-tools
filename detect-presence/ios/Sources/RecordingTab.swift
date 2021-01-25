import SwiftUI
import RelaySwiftUI

struct RecordingTab: View {
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
                        Label("Simulate begin trip", systemImage: "play.fill")
                    }

                    Button {
                        model.endTrip()
                    } label: {
                        Label("Simulate end trip", systemImage: "stop.fill")
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

                    Button {
                        model.clearQueuedTrips()
                    } label: {
                        Label("Clear queued trips", systemImage: "trash.fill")
                    }
                }
            }

            Section(header: Text("All Events")) {
                ForEach(model.allEvents) { event in
                    VStack(alignment: .leading, spacing: 8) {
                        Text(event.description)
                            .font(.body)
                        (Text(event.timestamp, style: .relative) + Text(" ago"))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
            }
        }
        .navigationTitle("Presence")
    }
}
