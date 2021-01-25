import SwiftUI

struct ContentView: View {
    @ObservedObject var model: AppModel

    var body: some View {
        TabView {
            NavigationView {
                TripsTab()
            }
            .tabItem {
                Image(systemName: "figure.walk")
                Text("Trips")
            }
            NavigationView {
                RecordingTab(model: model)
            }
            .tabItem {
                Image(systemName: "antenna.radiowaves.left.and.right")
                Text("Recording")
            }
        }
    }
}
