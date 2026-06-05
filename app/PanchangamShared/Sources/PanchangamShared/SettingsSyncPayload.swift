import Foundation

public struct SettingsSyncPayload: Codable, Equatable, Sendable {
    public var settings: APISettings
    public var latestSummary: TithiSummaryResponse?

    public init(settings: APISettings, latestSummary: TithiSummaryResponse? = nil) {
        self.settings = settings
        self.latestSummary = latestSummary
    }

    public init(dictionary: [String: Any]) throws {
        let data = try JSONSerialization.data(withJSONObject: dictionary)
        let payload = try JSONDecoder().decode(SettingsSyncPayload.self, from: data)
        try payload.settings.validate()
        try payload.latestSummary?.validate()
        self = payload
    }

    public func dictionary() throws -> [String: Any] {
        try settings.validate()
        try latestSummary?.validate()
        let data = try JSONEncoder().encode(self)
        let object = try JSONSerialization.jsonObject(with: data)
        guard let dictionary = object as? [String: Any] else {
            return [:]
        }
        return dictionary
    }
}
