import SwiftUI
import RelaySwiftUI
import detect_presence_ios_relay_generated

private let query = graphql("""
query TripsTabQuery {
    viewer {
        trips(first: 30) @connection(key: "TripsPageQuery_trips") {
            edges {
                node {
                    id
                    ...TripRow_trip
                }
            }
        }
    }
}
""")

struct TripsTab: View {
    @Query<TripsTabQuery> private var query

    @State private var isLoginPresented = false

    var model: AppModel
    var fetchKey: UUID

    var body: some View {
        Group {
            switch query.get(fetchKey: fetchKey) {
            case .loading:
                Text("Loading…")
            case .failure(let error):
                Text("Error: \(error.localizedDescription)")
            case .success(let data):
                if let trips = data?.viewer?.trips {
                    List(trips) { trip in
                        TripRow(trip: trip.asFragment())
                    }
                    .listStyle(PlainListStyle())
                } else {
                    Text("No trips found")
                }
            }
        }
        .navigationTitle("Your Trips")
        .navigationBarItems(trailing: Group {
            Button {
                isLoginPresented = true
            } label: {
                Text("Log In")
            }
            .sheet(isPresented: $isLoginPresented) {
                LoginView(model: model)
            }
        })
    }
}
