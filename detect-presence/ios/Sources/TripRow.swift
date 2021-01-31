import SwiftUI
import RelaySwiftUI
import detect_presence_ios_relay_generated

private let tripFragment = graphql("""
    fragment TripRow_trip on Trip {
        id
        leftAt
        returnedAt
        tags
    }
""")

struct TripRow: View {
    @Fragment<TripRow_trip> var trip

    var body: some View {
        if let trip = trip, let leftAt = trip.leftAt.asDate {
            VStack(alignment: .leading) {
                if let returnedAt = trip.returnedAt?.asDate {
                    Text(DateInterval(start: leftAt, end: returnedAt), formatter: dateIntervalFormatter)
                } else {
                    Text(leftAt, style: .relative) + Text(" ago")
                }

                if !trip.tags.isEmpty {
                    Text(verbatim: trip.tags.joined(separator: ", "))
                }
            }
        }
    }
}

private let dateIntervalFormatter: DateIntervalFormatter = {
    let formatter = DateIntervalFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short
    return formatter
}()
