import Foundation

public struct APISettings: Codable, Equatable, Sendable {
    public var apiBaseURL: URL
    public var latitude: Double
    public var longitude: Double
    public var timezone: String
    public var region: String
    public var method: String
    public var locale: String
    public var calendarSystem: String?

    public init(
        apiBaseURL: URL = URL(string: "http://localhost:8080")!,
        latitude: Double = 37.3382,
        longitude: Double = -121.8863,
        timezone: String = TimeZone.current.identifier,
        region: String = "global",
        method: String = "Drik",
        locale: String = "en",
        calendarSystem: String? = nil
    ) {
        self.apiBaseURL = apiBaseURL
        self.latitude = latitude
        self.longitude = longitude
        self.timezone = timezone
        self.region = region
        self.method = method
        self.locale = locale
        self.calendarSystem = calendarSystem
    }

    public static func parse(
        apiBaseURLText: String,
        latitudeText: String,
        longitudeText: String,
        timezoneText: String,
        regionText: String,
        methodText: String,
        localeText: String,
        calendarSystemText: String
    ) throws -> APISettings {
        let trimmedBaseURL = apiBaseURLText.trimmingCharacters(in: .whitespacesAndNewlines)
        let trimmedLatitude = latitudeText.trimmingCharacters(in: .whitespacesAndNewlines)
        let trimmedLongitude = longitudeText.trimmingCharacters(in: .whitespacesAndNewlines)
        let trimmedCalendarSystem = calendarSystemText.trimmingCharacters(in: .whitespacesAndNewlines)
        let calendarSystem: String?
        switch trimmedCalendarSystem.lowercased() {
        case "":
            calendarSystem = nil
        case "purnimanta":
            calendarSystem = "Purnimanta"
        case "amanta":
            calendarSystem = "Amanta"
        default:
            calendarSystem = trimmedCalendarSystem
        }

        guard let apiBaseURL = URL(string: trimmedBaseURL) else {
            throw APISettingsValidationError.invalidBaseURLFormat
        }

        guard let latitude = Double(trimmedLatitude) else {
            throw APISettingsValidationError.invalidLatitudeText
        }

        guard let longitude = Double(trimmedLongitude) else {
            throw APISettingsValidationError.invalidLongitudeText
        }

        let settings = APISettings(
            apiBaseURL: apiBaseURL,
            latitude: latitude,
            longitude: longitude,
            timezone: timezoneText.trimmingCharacters(in: .whitespacesAndNewlines),
            region: regionText.trimmingCharacters(in: .whitespacesAndNewlines),
            method: methodText.trimmingCharacters(in: .whitespacesAndNewlines),
            locale: localeText.trimmingCharacters(in: .whitespacesAndNewlines),
            calendarSystem: calendarSystem
        )
        try settings.validate()
        return settings
    }

    public func validate() throws {
        guard let scheme = apiBaseURL.scheme?.lowercased(), scheme == "http" || scheme == "https" else {
            throw APISettingsValidationError.invalidBaseURLScheme
        }

        guard apiBaseURL.host != nil else {
            throw APISettingsValidationError.invalidBaseURL
        }

        guard (-90...90).contains(latitude) else {
            throw APISettingsValidationError.invalidLatitude
        }

        guard (-180...180).contains(longitude) else {
            throw APISettingsValidationError.invalidLongitude
        }

        guard TimeZone(identifier: timezone) != nil else {
            throw APISettingsValidationError.invalidTimezone
        }

        guard !region.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw APISettingsValidationError.emptyRegion
        }

        guard !method.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw APISettingsValidationError.emptyMethod
        }

        guard !locale.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw APISettingsValidationError.emptyLocale
        }

        if let calendarSystem, calendarSystem != "Purnimanta" && calendarSystem != "Amanta" {
            throw APISettingsValidationError.invalidCalendarSystem
        }
    }

    public var developmentHostWarning: String? {
        guard let host = apiBaseURL.host?.lowercased(),
              host == "localhost" || host == "127.0.0.1" || host == "::1" else {
            return nil
        }

        return "Use a Mac LAN address for physical iPhone or Apple Watch testing."
    }
}

public enum APISettingsValidationError: Error, Equatable, Sendable {
    case invalidBaseURL
    case invalidBaseURLFormat
    case invalidBaseURLScheme
    case invalidLatitudeText
    case invalidLatitude
    case invalidLongitudeText
    case invalidLongitude
    case invalidTimezone
    case emptyRegion
    case emptyMethod
    case emptyLocale
    case invalidCalendarSystem

    public var userMessage: String {
        switch self {
        case .invalidBaseURL:
            return "API base URL must include a host"
        case .invalidBaseURLFormat:
            return "API base URL is not valid"
        case .invalidBaseURLScheme:
            return "API base URL must start with http or https"
        case .invalidLatitudeText:
            return "Latitude must be a number"
        case .invalidLatitude:
            return "Latitude must be between -90 and 90"
        case .invalidLongitudeText:
            return "Longitude must be a number"
        case .invalidLongitude:
            return "Longitude must be between -180 and 180"
        case .invalidTimezone:
            return "Timezone is not recognized"
        case .emptyRegion:
            return "Region is required"
        case .emptyMethod:
            return "Calculation method is required"
        case .emptyLocale:
            return "Locale is required"
        case .invalidCalendarSystem:
            return "Calendar system must be Purnimanta or Amanta"
        }
    }
}
