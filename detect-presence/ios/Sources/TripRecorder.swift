import Foundation
import detect_presence_proto_trips_trips_proto
import detect_presence_proto_trips_trips_swift_proto_grpc_client

private let dateFormatter = ISO8601DateFormatter()

class TripRecorder {
    private let client: TripsServiceService

    private var currentTrip: Trip?

    init() {
//        client = TripsServiceServiceClient(address: "100.117.39.47:2121", secure: false)
        client = TripsServiceServiceClient(address: "detect-presence-grpc.homelab", certificates: homelabCA)
    }

    func beginTrip() {
        guard currentTrip == nil else {
            NSLog("Already have a current trip, so not beginning a new one.")
            return
        }

        currentTrip = .with {
            $0.id = UUID().uuidString
            $0.leftAt = dateFormatter.string(from: Date())
        }
        NSLog("starting trip \(currentTrip!)")
    }

    func endTrip() {
        guard var trip = currentTrip else {
            NSLog("No trip in progress, nothing to end.")
            return
        }

        trip.returnedAt = dateFormatter.string(from: Date())
        NSLog("ending trip \(trip)")

        do {
            try client.recordTrips(.with { $0.trips = [trip] }) { _, result in
                if result.success {
                    NSLog("successfully recorded trip")
                } else {
                    NSLog("failed to record trip: \(result)")
                }
            }
        } catch {
            NSLog("error trying to record trips: \(error)")
        }
    }
}

// probably a better way to do this but it works!
private let homelabCA = """
-----BEGIN CERTIFICATE-----
MIIDVjCCAj6gAwIBAgIBATANBgkqhkiG9w0BAQsFADBQMRowGAYDVQQDDBFNYXR0
J3MgSG9tZWxhYiBDQTELMAkGA1UEBhMCVVMxJTAjBgkqhkiG9w0BCQEWFm1hdHRA
bWF0dG1vcmlhcml0eS5jb20wHhcNMjAxMTA4MTY0OTA4WhcNMjExMTA4MTY0OTA4
WjBQMRowGAYDVQQDDBFNYXR0J3MgSG9tZWxhYiBDQTELMAkGA1UEBhMCVVMxJTAj
BgkqhkiG9w0BCQEWFm1hdHRAbWF0dG1vcmlhcml0eS5jb20wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDRl2HGjFNsNVVs/l3tRWPd4C1UonuBFgPYNnkv
ZR7mxovM0lCgVVv9RSiqEK6E5gUO3US383G09VO04fITwGbpLRaYDMqfRfFkayWb
Hqtxu/oNeNlw8upL/dByvICSlLcSlcpWWqLkrYVTCU2LripHMv2Kqsbcyps8seM/
F9Ie5Mj3jBn01FJ9veq+a9kkW1CzcpnzBXWgqmczXlJ69iLTJq8qQ7O64RVMel74
kBaxvUT6exkgjFcRjvXrecRxNEOEiKIhIruDxtJIqgx5sZDVdOoPBYUKdpIW2bhF
9ax720UV1phcAul8gBMevus0S0y5o38yLjcSs4UlDVFtERy1AgMBAAGjOzA5MA8G
A1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgKkMBYGA1UdJQEB/wQMMAoGCCsG
AQUFBwMBMA0GCSqGSIb3DQEBCwUAA4IBAQDOw07kOn9BI0OMwvHiua5H0vLPM58s
7jCfUZRWmXHwStq948QslNWVzbIbn7xj7EDUIIciu3ubAWMx/hY/6diq640JzTN1
50U6ntoSYQcTmbAGBrM6VhXwNJSLmAB9LdVwh6a3tXVqgKZgqdYFHmbMxg/NfGln
WjzK7/ytsObOLhq9xGwhJdmKlGwUuQydbRtl8hqJ63toJ1xingH+/nhaiJRA+3Uu
J8N0sUUTv25QS7WKKSJWRcgX3GmEP8IJtRICcz/buPCGoreLD3Wcl171iDYIEYHq
9+mDsO5/lPY1EGk09S0hqgq65mNhZIYq9jpPlQ2ZkzP+2ZXgNKBi1cNQ
-----END CERTIFICATE-----
"""
