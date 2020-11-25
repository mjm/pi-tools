import CoreLocation
import detect_presence_proto_trips_trips_proto
import detect_presence_proto_trips_trips_swift_proto_grpc_client

private let beaconIdentifier = "home-beacons"

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

class BeaconObserver: NSObject, CLLocationManagerDelegate, ObservableObject {
    var locationManager: CLLocationManager
    var tripsClient: TripsServiceService

    @Published private(set) var status: CLRegionState = .unknown
    @Published private(set) var statusChangedTime: Date?

    init(locationManager: CLLocationManager = CLLocationManager()) {
        self.locationManager = locationManager
        self.tripsClient = TripsServiceServiceClient(address: "detect-presence-grpc.homelab", certificates: homelabCA)
        super.init()

        self.locationManager.delegate = self

        try! self.tripsClient.listTrips(ListTripsRequest()) { response, result in
            NSLog("Response: \(response)")
            NSLog("Result: \(result)")
        }
    }

    func startObserving() {
        guard CLLocationManager.isMonitoringAvailable(for: CLBeaconRegion.self) else {
            NSLog("Monitoring is not available. :(")
            return
        }

        let region = CLBeaconRegion(uuid: UUID(uuidString: "7298c12b-f658-445f-b1f2-5d6d582f0fb0")!,
                                    identifier: beaconIdentifier)
        locationManager.startMonitoring(for: region)
        NSLog("Started observing region")
    }

    func locationManagerDidChangeAuthorization(_ manager: CLLocationManager) {
        NSLog("Got authorization status: \(manager.authorizationStatus.rawValue)")
        switch manager.authorizationStatus {
        case .notDetermined:
            manager.requestWhenInUseAuthorization()
        case .authorizedWhenInUse:
            manager.requestAlwaysAuthorization()
            startObserving()
        default:
            break
        }
    }

    func locationManager(_ manager: CLLocationManager, didDetermineState state: CLRegionState, for region: CLRegion) {
        switch state {
        case .inside:
            NSLog("Device appears to be at home")
        case .outside:
            NSLog("Device appears to be away from home")
        case .unknown:
            NSLog("Device location is unknown")
        }

        status = state
        statusChangedTime = Date()
    }

    func locationManager(_ manager: CLLocationManager, didEnterRegion region: CLRegion) {
        NSLog("Entered region: \(region)")
    }

    func locationManager(_ manager: CLLocationManager, didExitRegion region: CLRegion) {
        NSLog("Exited region: \(region)")
    }

    func locationManager(_ manager: CLLocationManager, monitoringDidFailFor region: CLRegion?, withError error: Error) {
        NSLog("Monitoring failed: \(error.localizedDescription)")
    }

    func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
        NSLog("Location manager failed: \(error.localizedDescription)")
    }
}
