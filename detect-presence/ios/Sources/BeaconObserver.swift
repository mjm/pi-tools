import CoreLocation

private let beaconIdentifier = "home-beacons"

class BeaconObserver: NSObject, CLLocationManagerDelegate, ObservableObject {
    let locationManager: CLLocationManager
    let tripRecorder: TripRecorder

    @Published private(set) var status: CLRegionState = .unknown
    @Published private(set) var statusChangedTime: Date?

    init(locationManager: CLLocationManager = CLLocationManager()) {
        self.locationManager = locationManager
        self.tripRecorder = TripRecorder()
        super.init()

        self.locationManager.delegate = self
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
        case .authorizedAlways:
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
        tripRecorder.endTrip()
    }

    func locationManager(_ manager: CLLocationManager, didExitRegion region: CLRegion) {
        NSLog("Exited region: \(region)")
        tripRecorder.beginTrip()
    }

    func locationManager(_ manager: CLLocationManager, monitoringDidFailFor region: CLRegion?, withError error: Error) {
        NSLog("Monitoring failed: \(error.localizedDescription)")
    }

    func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
        NSLog("Location manager failed: \(error.localizedDescription)")
    }
}
