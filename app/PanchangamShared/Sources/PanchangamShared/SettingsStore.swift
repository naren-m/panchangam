import Foundation

public struct SettingsStore {
    public var defaults: UserDefaults
    public var key: String

    public init(defaults: UserDefaults = .standard, key: String = "panchangam.api.settings") {
        self.defaults = defaults
        self.key = key
    }

    public func load() throws -> APISettings {
        guard let data = defaults.data(forKey: key) else {
            return APISettings()
        }
        let settings = try JSONDecoder().decode(APISettings.self, from: data)
        try settings.validate()
        return settings
    }

    public func hasSavedSettings() -> Bool {
        guard defaults.data(forKey: key) != nil else {
            return false
        }
        return (try? load()) != nil
    }

    public func save(_ settings: APISettings) throws {
        try settings.validate()
        let data = try JSONEncoder().encode(settings)
        defaults.set(data, forKey: key)
    }
}
