import Foundation

public struct CachedTithiSummary: Codable, Equatable, Sendable {
    public var summary: TithiSummaryResponse
    public var fetchedAt: Date

    public init(summary: TithiSummaryResponse, fetchedAt: Date) {
        self.summary = summary
        self.fetchedAt = fetchedAt
    }
}

public struct TithiCache: Sendable {
    public var fileURL: URL

    public init(fileURL: URL) {
        self.fileURL = fileURL
    }

    public func save(_ summary: TithiSummaryResponse, fetchedAt: Date = Date()) throws {
        try summary.validate()
        let cached = CachedTithiSummary(summary: summary, fetchedAt: fetchedAt)
        let data = try PanchangamJSON.encoder.encode(cached)
        try FileManager.default.createDirectory(
            at: fileURL.deletingLastPathComponent(),
            withIntermediateDirectories: true
        )
        try data.write(to: fileURL, options: [.atomic])
    }

    public func load() throws -> CachedTithiSummary? {
        guard FileManager.default.fileExists(atPath: fileURL.path) else {
            return nil
        }

        let data = try Data(contentsOf: fileURL)
        let cached = try PanchangamJSON.decoder.decode(CachedTithiSummary.self, from: data)
        try cached.summary.validate()
        return cached
    }

    public func clear() throws {
        guard FileManager.default.fileExists(atPath: fileURL.path) else {
            return
        }
        try FileManager.default.removeItem(at: fileURL)
    }
}

public enum PanchangamStorage {
    public static let appGroupIdentifier = "group.app.panchangam"

    public static func appGroupContainerAvailable(
        fileManager: FileManager = .default,
        appGroupIdentifier: String? = appGroupIdentifier
    ) -> Bool {
        appGroupContainerURL(fileManager: fileManager, appGroupIdentifier: appGroupIdentifier) != nil
    }

    public static func appGroupStatusText(
        fileManager: FileManager = .default,
        appGroupIdentifier: String? = appGroupIdentifier
    ) -> String? {
        appGroupContainerAvailable(fileManager: fileManager, appGroupIdentifier: appGroupIdentifier)
            ? nil
            : "App Group unavailable"
    }

    public static func cacheFileURL(fileManager: FileManager = .default, appGroupIdentifier: String? = appGroupIdentifier) -> URL {
        if let containerURL = appGroupContainerURL(fileManager: fileManager, appGroupIdentifier: appGroupIdentifier) {
            return containerURL.appendingPathComponent("tithi-summary-cache.json")
        }

        let cachesURL = fileManager.urls(for: .cachesDirectory, in: .userDomainMask).first ?? fileManager.temporaryDirectory
        return cachesURL.appendingPathComponent("tithi-summary-cache.json")
    }

    private static func appGroupContainerURL(
        fileManager: FileManager,
        appGroupIdentifier: String?
    ) -> URL? {
        guard let appGroupIdentifier else {
            return nil
        }
        return fileManager.containerURL(forSecurityApplicationGroupIdentifier: appGroupIdentifier)
    }
}
