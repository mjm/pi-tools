import SwiftUI

struct ContentView: View {
    enum Tab: Hashable {
        case trips
        case recording
    }

    @ObservedObject var model: AppModel
    @State private var selectedItem: Tab = .trips
    @State private var fetchKey = UUID()

    @Environment(\.scenePhase) var scenePhase

    var body: some View {
        TabView(selection: $selectedItem) {
            NavigationView {
                TripsTab(fetchKey: fetchKey)
            }
            .tabItem {
                Image(systemName: "figure.walk")
                Text("Trips")
            }
            .tag(Tab.trips)

            NavigationView {
                RecordingTab(model: model)
            }
            .tabItem {
                Image(systemName: "antenna.radiowaves.left.and.right")
                Text("Recording")
            }
            .tag(Tab.recording)
        }
        .onChange(of: selectedItem) { newValue in
            if newValue == .trips {
                fetchKey = UUID()
            }
        }
        .onChange(of: scenePhase) { newValue in
            if newValue == .active {
                fetchKey = UUID()
            }
        }
    }
}
