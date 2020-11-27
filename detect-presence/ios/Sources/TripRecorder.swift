import Foundation
import Combine
import detect_presence_proto_trips_trips_proto
import detect_presence_proto_trips_trips_swift_proto_grpc_client

class TripRecorder {
    enum Event {
        case recorded([Trip])
        case recordFailed(String)
    }

    private var client: TripsServiceService!
    private let eventsSubject = PassthroughSubject<Event, Never>()
    private var cancellables = Set<AnyCancellable>()

    init<P: Publisher>(events: P) where P.Output == TripsController.Event, P.Failure == Never {
        events.sink { [weak self] event in
            guard let self = self else { return }

            switch event {
            case .tripBegan(let trip):
                NSLog("starting trip \(trip)")
            case .tripEnded(let queuedTrips):
                self.recordTrips(queuedTrips)
            }
        }.store(in: &cancellables)
    }

    func eventsPublisher() -> AnyPublisher<Event, Never> {
        eventsSubject.eraseToAnyPublisher()
    }

    func recordTrips(_ trips: [Trip]) {
        do {
            let request = RecordTripsRequest.with {
                $0.trips = trips.map(\.asProto)
            }
            try client.recordTrips(request) { _, result in
                DispatchQueue.main.sync {
                    if result.success && result.statusCode == .ok {
                        NSLog("successfully recorded trip")
                        self.eventsSubject.send(.recorded(trips))
                    } else {
                        NSLog("failed to record trip: \(result)")
                        self.eventsSubject.send(.recordFailed(result.description))
                    }
                }
            }
        } catch {
            NSLog("error trying to record trips: \(error)")
            self.eventsSubject.send(.recordFailed(error.localizedDescription))
        }
    }

    func setUpClient(useDevServer: Bool = false) {
        if useDevServer {
            NSLog("creating dev server client")
            client = TripsServiceServiceClient(address: "100.117.39.47:2121", secure: false)
        } else {
            NSLog("creating real client")
            client = TripsServiceServiceClient(address: "detect-presence-grpc.homelab", certificates: homelabCA)
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
