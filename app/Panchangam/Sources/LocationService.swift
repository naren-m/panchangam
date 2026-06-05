import Combine
@preconcurrency import CoreLocation
import Foundation

@MainActor
final class LocationService: NSObject, ObservableObject, CLLocationManagerDelegate {
    @Published var lastLocation: CLLocation?
    @Published var authorizationText: String?

    private let manager = CLLocationManager()

    override init() {
        super.init()
        manager.delegate = self
        manager.desiredAccuracy = kCLLocationAccuracyKilometer
        authorizationText = Self.text(for: manager.authorizationStatus)
    }

    func requestLocation() {
        switch manager.authorizationStatus {
        case .notDetermined:
            manager.requestWhenInUseAuthorization()
        case .authorizedAlways, .authorizedWhenInUse:
            manager.requestLocation()
        case .denied, .restricted:
            authorizationText = "Location permission denied"
        @unknown default:
            authorizationText = "Location permission unavailable"
        }
    }

    nonisolated func locationManagerDidChangeAuthorization(_ manager: CLLocationManager) {
        let status = manager.authorizationStatus
        Task { @MainActor [weak self] in
            self?.authorizationText = Self.text(for: status)
            self?.requestLocationIfAllowed(status: status)
        }
    }

    nonisolated func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        guard let location = locations.last else {
            return
        }
        let latitude = location.coordinate.latitude
        let longitude = location.coordinate.longitude
        let altitude = location.altitude
        let horizontalAccuracy = location.horizontalAccuracy
        let verticalAccuracy = location.verticalAccuracy
        let timestamp = location.timestamp

        Task { @MainActor [weak self] in
            self?.lastLocation = CLLocation(
                coordinate: CLLocationCoordinate2D(latitude: latitude, longitude: longitude),
                altitude: altitude,
                horizontalAccuracy: horizontalAccuracy,
                verticalAccuracy: verticalAccuracy,
                timestamp: timestamp
            )
            self?.authorizationText = "Location updated"
        }
    }

    nonisolated func locationManager(_ manager: CLLocationManager, didFailWithError _: Error) {
        Task { @MainActor [weak self] in
            self?.authorizationText = "Location unavailable"
        }
    }

    private func requestLocationIfAllowed(status: CLAuthorizationStatus) {
        guard status == .authorizedAlways || status == .authorizedWhenInUse else {
            return
        }
        manager.requestLocation()
    }

    private static func text(for status: CLAuthorizationStatus) -> String? {
        switch status {
        case .notDetermined:
            return nil
        case .authorizedAlways, .authorizedWhenInUse:
            return "Location allowed"
        case .denied, .restricted:
            return "Location permission denied"
        @unknown default:
            return "Location permission unavailable"
        }
    }
}
