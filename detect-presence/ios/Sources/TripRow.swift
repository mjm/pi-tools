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
            VStack(alignment: .leading, spacing: 8) {
                if let returnedAt = trip.returnedAt?.asDate {
                    Text(DateInterval(start: leftAt, end: returnedAt), formatter: dateIntervalFormatter)
                        .font(.headline)
                } else {
                    (Text(leftAt, style: .relative) + Text(" ago"))
                        .font(.headline)
                }

                if !trip.tags.isEmpty {
                    FlowLayout(mode: .scrollable, items: trip.tags) { tag in
                        TagPill(text: tag)
                    }
                    .padding(-4)
                }
            }
            .padding(.vertical, 4)
        }
    }
}

private let dateIntervalFormatter: DateIntervalFormatter = {
    let formatter = DateIntervalFormatter()
    formatter.dateStyle = .medium
    formatter.timeStyle = .short
    return formatter
}()
