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
    @Query<TripsTabQuery> var query

    var body: some View {
        Group {
            switch query.get() {
            case .loading:
                Text("Loading…")
            case .failure(let error):
                Text("Error: \(error.localizedDescription)")
            case .success(let data):
                if let trips = data?.viewer?.trips {
                    List(trips) { trip in
                        TripRow(trip: trip.asFragment())
                    }
                } else {
                    Text("No trips found")
                }
            }
        }
        .navigationTitle("Your Trips")
    }
}
