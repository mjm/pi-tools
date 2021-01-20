import Relay

// language=GraphQL
private let mutation = graphql("""
mutation RecordTripsMutation($input: RecordTripsInput!) {
    recordTrips(input: $input) {
        recordedTrips {
            id
            leftAt
            returnedAt
        }
        failures {
            tripID
            message
        }
    }
}
""")
