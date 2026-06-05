import Foundation
import PanchangamShared

@main
struct PanchangamSharedChecks {
    static func main() async throws {
        try checkTithiSummaryDecoding()
        try checkTithiSummaryValidation()
        try checkAPISettingsValidation()
        try checkAPISettingsParsing()
        try checkAPISettingsDevelopmentHostWarning()
        try checkAPIClientURL()
        try await checkAPIClientFetch()
        try await checkAPIClientCurrentFetchOmitsMoment()
        try await checkAPIClientInvalidSummary()
        try await checkAPIClientHTTPError()
        try checkAPIErrorMessages()
        try checkTithiCache()
        try checkSettingsSyncPayload()
        try checkSettingsStore()
        try checkStorageStatus()
        try checkTimelinePolicy()
        try checkFormatter()
    }

    private static func checkTithiSummaryDecoding() throws {
        let response = try PanchangamJSON.decoder.decode(TithiSummaryResponse.self, from: Fixtures.summaryJSON)

        try expect(response.date == "2026-06-02", "date decoded")
        try expect(response.tithi.name == "Ekadashi", "tithi name decoded")
        try expect(response.tithi.traditionalName == "Ekadashi", "traditional tithi name decoded")
        try expect(response.tithi.number == 11, "tithi number decoded")
        try expect(response.tithi.paksha == "Shukla", "paksha decoded")
        try expect(response.tithi.pakshaDay == 11, "paksha day decoded")
        try expect(response.panchaAnga.nakshatra == "Anuradha", "nakshatra decoded")
        try expect(response.day.sunriseTime == "05:42", "sunrise decoded")
        try expect(response.day.sunsetTime == "19:01", "sunset decoded")
        try expect(response.day.abhijitMuhurta.name == "Abhijit", "abhijit name decoded")
        try expect(response.day.abhijitMuhurta.startTime == "11:54", "abhijit start decoded")
        try expect(response.day.abhijitMuhurta.endTime == "12:48", "abhijit end decoded")
        try expect(response.day.abhijitMuhurta.auspicious, "abhijit auspicious decoded")
        try expect(response.calculation.timezone == "America/Los_Angeles", "calculation timezone decoded")
        try expect(response.generatedAt == Fixtures.generatedAt, "generated_at decoded")
        try expect(response.nextRefreshAt == Fixtures.nextRefreshAt, "next_refresh_at decoded")
        try response.validate()
    }

    private static func checkTithiSummaryValidation() throws {
        var lateRefreshSummary = Fixtures.summary
        lateRefreshSummary.nextRefreshAt = Fixtures.nextRefreshAt.addingTimeInterval(60)
        try expectSummaryValidationMessage(
            { try lateRefreshSummary.validate() },
            "Next refresh time must be after generated time and no later than tithi end time",
            "refresh after tithi end"
        )

        var earlyGeneratedSummary = Fixtures.summary
        earlyGeneratedSummary.generatedAt = earlyGeneratedSummary.tithi.startTime.addingTimeInterval(-60)
        try expectSummaryValidationMessage(
            { try earlyGeneratedSummary.validate() },
            "Generated time must be within the tithi start and end time",
            "generated before tithi start"
        )
    }

    private static func checkAPISettingsValidation() throws {
        try Fixtures.settings.validate()

        var invalidBaseURL = Fixtures.settings
        invalidBaseURL.apiBaseURL = URL(string: "file:///tmp/panchangam")!
        try expectValidationMessage(
            invalidBaseURL,
            "API base URL must start with http or https",
            "invalid base URL scheme"
        )

        var invalidLatitude = Fixtures.settings
        invalidLatitude.latitude = 91
        try expectValidationMessage(
            invalidLatitude,
            "Latitude must be between -90 and 90",
            "invalid latitude"
        )

        var invalidLongitude = Fixtures.settings
        invalidLongitude.longitude = -181
        try expectValidationMessage(
            invalidLongitude,
            "Longitude must be between -180 and 180",
            "invalid longitude"
        )

        var invalidTimezone = Fixtures.settings
        invalidTimezone.timezone = "Mars/Olympus"
        try expectValidationMessage(
            invalidTimezone,
            "Timezone is not recognized",
            "invalid timezone"
        )

        var invalidCalendarSystem = Fixtures.settings
        invalidCalendarSystem.calendarSystem = "Lunar"
        try expectValidationMessage(
            invalidCalendarSystem,
            "Calendar system must be Purnimanta or Amanta",
            "invalid calendar system"
        )
    }

    private static func checkAPISettingsParsing() throws {
        let settings = try APISettings.parse(
            apiBaseURLText: " https://api.example.test ",
            latitudeText: " 37.3382 ",
            longitudeText: " -121.8863 ",
            timezoneText: " America/Los_Angeles ",
            regionText: " California ",
            methodText: " Drik ",
            localeText: " en ",
            calendarSystemText: " "
        )

        try expect(settings == APISettings(
            apiBaseURL: URL(string: "https://api.example.test")!,
            latitude: 37.3382,
            longitude: -121.8863,
            timezone: "America/Los_Angeles",
            region: "California",
            method: "Drik",
            locale: "en",
            calendarSystem: nil
        ), "parsed trimmed settings")

        let amantaSettings = try APISettings.parse(
            apiBaseURLText: "https://api.example.test",
            latitudeText: "37.3382",
            longitudeText: "-121.8863",
            timezoneText: "America/Los_Angeles",
            regionText: "Karnataka",
            methodText: "Drik",
            localeText: "en",
            calendarSystemText: " amanta "
        )
        try expect(amantaSettings.calendarSystem == "Amanta", "calendar system is canonicalized")

        try expectParsingMessage(
            latitudeText: "north",
            expectedMessage: "Latitude must be a number",
            checkName: "latitude parse error"
        )

        try expectParsingMessage(
            longitudeText: "west",
            expectedMessage: "Longitude must be a number",
            checkName: "longitude parse error"
        )

        try expectParsingMessage(
            apiBaseURLText: "http://[",
            expectedMessage: "API base URL is not valid",
            checkName: "base URL parse error"
        )
    }

    private static func checkAPISettingsDevelopmentHostWarning() throws {
        var localhostSettings = Fixtures.settings
        localhostSettings.apiBaseURL = URL(string: "http://localhost:8080")!
        try expect(
            localhostSettings.developmentHostWarning == "Use a Mac LAN address for physical iPhone or Apple Watch testing.",
            "localhost warning"
        )

        var loopbackSettings = Fixtures.settings
        loopbackSettings.apiBaseURL = URL(string: "http://127.0.0.1:8080")!
        try expect(loopbackSettings.developmentHostWarning != nil, "loopback warning")
        try expect(Fixtures.settings.developmentHostWarning == nil, "remote host has no warning")
    }

    private static func checkAPIClientURL() throws {
        let client = PanchangamAPIClient(baseURL: URL(string: "https://api.example.test")!, session: StubHTTPSession())
        let url = try client.makeTithiSummaryURL(settings: Fixtures.settings, at: Fixtures.generatedAt)
        let components = try require(URLComponents(url: url, resolvingAgainstBaseURL: false), "URL components")
        let queryItems = try require(components.queryItems, "query items")
        let query = Dictionary(uniqueKeysWithValues: queryItems.map { ($0.name, $0.value ?? "") })

        try expect(components.scheme == "https", "URL scheme")
        try expect(components.host == "api.example.test", "URL host")
        try expect(components.path == "/api/v1/tithi/current", "URL path")
        try expect(query["at"] == "2026-06-02T12:00:00Z", "at query")
        try expect(query["lat"] == "37.3382", "lat query")
        try expect(query["lng"] == "-121.8863", "lng query")
        try expect(query["tz"] == "America/Los_Angeles", "timezone query")
        try expect(query["region"] == "California", "region query")
        try expect(query["method"] == "Drik", "method query")
        try expect(query["locale"] == "en", "locale query")
        try expect(query["calendar_system"] == "Purnimanta", "calendar query")
    }

    private static func checkAPIClientFetch() async throws {
        let session = StubHTTPSession(data: Fixtures.summaryJSON, statusCode: 200)
        let client = PanchangamAPIClient(baseURL: URL(string: "https://api.example.test")!, session: session)

        let response = try await client.currentTithi(settings: Fixtures.settings, at: Fixtures.generatedAt)

        try expect(response.tithi.name == "Ekadashi", "API response decoded")
        try expect(response.nextRefreshAt == Fixtures.nextRefreshAt, "API next refresh decoded")
        try expect(session.lastRequest?.url?.path == "/api/v1/tithi/current", "API request path")
    }

    private static func checkAPIClientCurrentFetchOmitsMoment() async throws {
        let session = StubHTTPSession(data: Fixtures.summaryJSON, statusCode: 200)
        let client = PanchangamAPIClient(baseURL: URL(string: "https://api.example.test")!, session: session)

        _ = try await client.currentTithi(settings: Fixtures.settings)

        let requestURL = try require(session.lastRequest?.url, "current API request URL")
        let components = try require(URLComponents(url: requestURL, resolvingAgainstBaseURL: false), "current API URL components")
        let queryItems = try require(components.queryItems, "current API query items")
        let query = Dictionary(uniqueKeysWithValues: queryItems.map { ($0.name, $0.value ?? "") })

        try expect(query["at"] == nil, "current API omits at by default")
        try expect(query["lat"] == "37.3382", "current API still sends latitude")
    }

    private static func checkAPIClientInvalidSummary() async throws {
        var invalidSummary = Fixtures.summary
        invalidSummary.tithi.name = " "
        let invalidData = try PanchangamJSON.encoder.encode(invalidSummary)
        let session = StubHTTPSession(data: invalidData, statusCode: 200)
        let client = PanchangamAPIClient(baseURL: URL(string: "https://api.example.test")!, session: session)

        do {
            _ = try await client.currentTithi(settings: Fixtures.settings, at: Fixtures.generatedAt)
            throw CheckFailure("invalid API summary")
        } catch let error as PanchangamAPIError {
            try expect(error.userMessage == "Tithi name is required", "invalid API summary")
        }
    }

    private static func checkAPIClientHTTPError() async throws {
        let errorJSON = Data(#"{"error":{"code":"SERVICE_UNAVAILABLE","message":"Service is temporarily unavailable"}}"#.utf8)
        let session = StubHTTPSession(data: errorJSON, statusCode: 503)
        let client = PanchangamAPIClient(baseURL: URL(string: "https://api.example.test")!, session: session)

        do {
            _ = try await client.currentTithi(settings: Fixtures.settings, at: Fixtures.generatedAt)
            throw CheckFailure("Expected PanchangamAPIError")
        } catch let error as PanchangamAPIError {
            try expect(error.userMessage == "HTTP 503: Service is temporarily unavailable", "HTTP API error message")
        }
    }

    private static func checkAPIErrorMessages() throws {
        try expect(PanchangamAPIError.invalidURL.userMessage == "Invalid API base URL", "invalid URL error message")
        try expect(
            PanchangamAPIError.missingHTTPResponse.userMessage == "API did not return an HTTP response",
            "missing HTTP response error message"
        )
        try expect(
            PanchangamAPIError.httpStatus(code: 404, message: nil).userMessage == "HTTP 404 from Panchangam API",
            "HTTP status fallback error message"
        )
    }

    private static func checkTithiCache() throws {
        let fileURL = FileManager.default.temporaryDirectory
            .appendingPathComponent(UUID().uuidString)
            .appendingPathExtension("json")
        let cache = TithiCache(fileURL: fileURL)
        defer { try? FileManager.default.removeItem(at: fileURL) }

        try cache.save(Fixtures.summary, fetchedAt: Fixtures.generatedAt)
        let cached = try require(try cache.load(), "cached summary")

        try expect(cached.summary.tithi.name == "Ekadashi", "cached tithi")
        try expect(cached.fetchedAt == Fixtures.generatedAt, "cached fetch date")

        try cache.clear()
        try expect(try cache.load() == nil, "cache cleared")

        var invalidSummary = Fixtures.summary
        invalidSummary.tithi.name = " "
        try expectSummaryValidationMessage(
            { try cache.save(invalidSummary, fetchedAt: Fixtures.generatedAt) },
            "Tithi name is required",
            "invalid cache save"
        )

        let invalidCached = CachedTithiSummary(summary: invalidSummary, fetchedAt: Fixtures.generatedAt)
        let invalidCacheData = try PanchangamJSON.encoder.encode(invalidCached)
        try FileManager.default.createDirectory(
            at: fileURL.deletingLastPathComponent(),
            withIntermediateDirectories: true
        )
        try invalidCacheData.write(to: fileURL, options: [.atomic])
        try expectSummaryValidationMessage(
            { _ = try cache.load() },
            "Tithi name is required",
            "invalid cache load"
        )
    }

    private static func checkSettingsSyncPayload() throws {
        let payload = SettingsSyncPayload(settings: Fixtures.settings, latestSummary: Fixtures.summary)
        let dictionary = try payload.dictionary()
        let decoded = try SettingsSyncPayload(dictionary: dictionary)

        try expect(decoded.settings.apiBaseURL == Fixtures.settings.apiBaseURL, "settings base URL")
        try expect(decoded.settings.latitude == Fixtures.settings.latitude, "settings latitude")
        try expect(decoded.settings.longitude == Fixtures.settings.longitude, "settings longitude")
        try expect(decoded.settings.timezone == Fixtures.settings.timezone, "settings timezone")
        try expect(decoded.settings.calendarSystem == Fixtures.settings.calendarSystem, "settings calendar")
        try expect(decoded.latestSummary?.tithi.name == "Ekadashi", "synced summary tithi")
        try expect(decoded.latestSummary?.nextRefreshAt == Fixtures.nextRefreshAt, "synced summary refresh")

        var invalidSettings = Fixtures.settings
        invalidSettings.latitude = 91
        try expectSyncPayloadEncodingValidation(
            SettingsSyncPayload(settings: invalidSettings),
            "Latitude must be between -90 and 90",
            "invalid synced latitude encoding"
        )

        let invalidDictionary: [String: Any] = [
            "settings": [
                "apiBaseURL": Fixtures.settings.apiBaseURL.absoluteString,
                "latitude": 91.0,
                "longitude": Fixtures.settings.longitude,
                "timezone": Fixtures.settings.timezone,
                "region": Fixtures.settings.region,
                "method": Fixtures.settings.method,
                "locale": Fixtures.settings.locale,
                "calendarSystem": Fixtures.settings.calendarSystem as Any
            ]
        ]
        try expectSyncPayloadValidation(
            invalidDictionary,
            "Latitude must be between -90 and 90",
            "invalid synced latitude"
        )

        var invalidSummary = Fixtures.summary
        invalidSummary.tithi.name = " "
        try expectSummaryValidationMessage(
            { _ = try SettingsSyncPayload(settings: Fixtures.settings, latestSummary: invalidSummary).dictionary() },
            "Tithi name is required",
            "invalid synced summary encoding"
        )

        var invalidSummaryDictionary = try payload.dictionary()
        var latestSummary = try require(
            invalidSummaryDictionary["latestSummary"] as? [String: Any],
            "latest summary dictionary"
        )
        var tithi = try require(latestSummary["tithi"] as? [String: Any], "latest summary tithi dictionary")
        tithi["name"] = " "
        latestSummary["tithi"] = tithi
        invalidSummaryDictionary["latestSummary"] = latestSummary
        try expectSummaryValidationMessage(
            { _ = try SettingsSyncPayload(dictionary: invalidSummaryDictionary) },
            "Tithi name is required",
            "invalid synced summary decoding"
        )
    }

    private static func checkSettingsStore() throws {
        let suiteName = UUID().uuidString
        let defaults = try require(UserDefaults(suiteName: suiteName), "temporary defaults")
        defer { defaults.removePersistentDomain(forName: suiteName) }
        let store = SettingsStore(defaults: defaults)

        try expect(!store.hasSavedSettings(), "new settings store starts empty")

        try store.save(Fixtures.settings)
        let loaded = try store.load()

        try expect(loaded == Fixtures.settings, "stored settings round trip")
        try expect(store.hasSavedSettings(), "saved settings are detectable")

        var invalidSettings = Fixtures.settings
        invalidSettings.latitude = 91
        try expectSettingsStoreSaveValidation(
            store,
            invalidSettings,
            "Latitude must be between -90 and 90",
            "invalid settings save"
        )

        let invalidData = try JSONEncoder().encode(invalidSettings)
        defaults.set(invalidData, forKey: store.key)
        try expectSettingsStoreLoadValidation(
            store,
            "Latitude must be between -90 and 90",
            "invalid settings load"
        )
        try expect(!store.hasSavedSettings(), "invalid saved settings are not usable")
    }

    private static func checkStorageStatus() throws {
        try expect(
            !PanchangamStorage.appGroupContainerAvailable(appGroupIdentifier: nil),
            "nil app group is unavailable"
        )
        try expect(
            PanchangamStorage.appGroupStatusText(appGroupIdentifier: nil) == "App Group unavailable",
            "missing app group status"
        )
    }

    private static func checkTimelinePolicy() throws {
        let policy = TithiTimelinePolicy()
        try expect(policy.refreshDate(for: Fixtures.summary, now: Fixtures.generatedAt) == Fixtures.nextRefreshAt, "future refresh date")
        try expect(policy.entryDates(for: Fixtures.summary, now: Fixtures.generatedAt).prefix(3) == [
            Fixtures.generatedAt,
            Fixtures.generatedAt.addingTimeInterval(30 * 60),
            Fixtures.generatedAt.addingTimeInterval(60 * 60)
        ], "current timeline entry cadence")

        var staleSummary = Fixtures.summary
        staleSummary.nextRefreshAt = Fixtures.generatedAt.addingTimeInterval(-60)
        try expect(
            policy.refreshDate(for: staleSummary, now: Fixtures.generatedAt) == Fixtures.generatedAt.addingTimeInterval(15 * 60),
            "stale refresh date"
        )
        try expect(policy.entryDates(for: staleSummary, now: Fixtures.generatedAt) == [
            Fixtures.generatedAt
        ], "stale timeline uses one entry")
    }

    private static func checkFormatter() throws {
        let formatter = TithiFormatter()
        var alternateNameSummary = Fixtures.summary
        alternateNameSummary.tithi.traditionalName = "Ekadashi Traditional"

        try expect(formatter.inlineText(for: Fixtures.summary) == "Ekadashi, Shukla 11", "inline text")
        try expect(
            formatter.inlineText(for: alternateNameSummary) == "Ekadashi (Ekadashi Traditional), Shukla 11",
            "inline text includes distinct traditional name"
        )
        try expect(
            formatter.complicationInlineText(for: Fixtures.summary, isStale: false) == "Ekadashi, Shukla 11",
            "current complication inline text"
        )
        try expect(
            formatter.complicationInlineText(for: Fixtures.summary, isStale: true) == "Cached: Ekadashi, Shukla 11",
            "cached complication inline text"
        )
        try expect(
            formatter.complicationSecondaryText(for: Fixtures.summary, isStale: false) == "Anuradha",
            "current complication secondary text"
        )
        try expect(
            formatter.complicationSecondaryText(for: Fixtures.summary, isStale: true) == "Cached",
            "cached complication secondary text"
        )
        try expect(formatter.circularText(for: Fixtures.summary) == "Eka 11", "circular text")
        try expect(formatter.nakshatraIndex(for: Fixtures.summary) == 16, "nakshatra mandala index")
        try expect(formatter.tithiProgress(for: Fixtures.summary, now: Fixtures.summary.tithi.startTime.addingTimeInterval(-60)) == 0, "tithi progress before start")
        try expect(formatter.tithiProgress(for: Fixtures.summary, now: Fixtures.generatedAt) == 0.44, "tithi progress at generated time")
        try expect(formatter.tithiProgress(for: Fixtures.summary, now: Fixtures.nextRefreshAt.addingTimeInterval(60)) == 1, "tithi progress after end")
        try expect(formatter.progressText(for: Fixtures.summary, now: Fixtures.generatedAt) == "44% elapsed", "progress text")
        try expect(formatter.remainingText(for: Fixtures.summary, now: Fixtures.generatedAt) == "14h 0m left", "remaining text")
        try expect(
            formatter.remainingText(
                for: Fixtures.summary,
                now: Fixtures.nextRefreshAt.addingTimeInterval(5 * 60)
            ) == "ended 5m ago",
            "ended tithi text"
        )
        try expect(formatter.rectangularLines(for: Fixtures.summary, now: Fixtures.generatedAt) == [
            "Ekadashi, Shukla 11",
            "Nakshatra: Anuradha",
            "Yoga/Karana: Siddha / Gara",
            "Abhijit: 11:54-12:48",
            "Remaining: 14h 0m left"
        ], "rectangular lines")
        try expect(formatter.rectangularLines(for: Fixtures.summary, now: Fixtures.generatedAt, isStale: true) == [
            "Ekadashi, Shukla 11",
            "Nakshatra: Anuradha",
            "Yoga/Karana: Siddha / Gara",
            "Abhijit: 11:54-12:48",
            "Cached: 14h 0m left"
        ], "stale rectangular lines")

        let details = Dictionary(uniqueKeysWithValues: formatter.detailRows(for: Fixtures.summary, now: Fixtures.generatedAt))
        try expect(details["Tithi Number"] == "11", "tithi number detail")
        try expect(details["Traditional Name"] == "Ekadashi", "traditional tithi name detail")
        try expect(details["Paksha"] == "Shukla", "paksha detail")
        try expect(details["Paksha Day"] == "11", "paksha day detail")
        try expect(details["Tithi Type"] == "Nanda", "tithi type detail")
        try expect(details["Starts"] == "2026-06-01 18:00 PDT", "localized start time")
        try expect(details["Ends"] == "2026-06-02 19:00 PDT", "localized end time")
        try expect(details["Progress"] == "44% elapsed", "progress detail")
        try expect(details["Remaining"] == "14h 0m left", "remaining detail")
        try expect(details["Generated"] == "2026-06-02 05:00 PDT", "generated detail")
        try expect(details["Next Refresh"] == "2026-06-02 19:00 PDT", "next refresh detail")
        try expect(details["Sunrise"] == "05:42", "sunrise detail")
        try expect(details["Sunset"] == "19:01", "sunset detail")
        try expect(details["Abhijit"] == "11:54-12:48", "abhijit detail")
        try expect(details["Region"] == "California", "region detail")
        try expect(details["Timezone"] == "America/Los_Angeles", "timezone detail")
        try expect(details["Calendar System"] == "Purnimanta", "calendar detail")
        try expect(details["Method"] == "Drik", "method detail")
        try expect(details["Locale"] == "en", "locale detail")
        try expect(formatter.abhijitText(for: Fixtures.summary) == "11:54-12:48", "abhijit compact text")
        try expect(
            formatter.complicationAccessibilityLabel(
                for: Fixtures.summary,
                now: Fixtures.generatedAt,
                isStale: false
            ) == "Current tithi Ekadashi, Shukla 11, nakshatra Anuradha, yoga Siddha, karana Gara, abhijit 11:54-12:48, remaining 14h 0m left",
            "current complication accessibility label"
        )
        try expect(
            formatter.complicationAccessibilityLabel(
                for: Fixtures.summary,
                now: Fixtures.generatedAt,
                isStale: true
            ) == "Cached tithi Ekadashi, Shukla 11, nakshatra Anuradha, yoga Siddha, karana Gara, abhijit 11:54-12:48, remaining 14h 0m left",
            "cached complication accessibility label"
        )
        try expect(
            formatter.complicationAccessibilityLabel(
                for: Fixtures.summary,
                now: Fixtures.nextRefreshAt.addingTimeInterval(5 * 60),
                isStale: true
            ) == "Cached tithi Ekadashi, Shukla 11, nakshatra Anuradha, yoga Siddha, karana Gara, abhijit 11:54-12:48, ended 5m ago",
            "expired complication accessibility label"
        )
    }

    private static func expect(_ condition: Bool, _ message: String) throws {
        if !condition {
            throw CheckFailure(message)
        }
    }

    private static func require<T>(_ value: T?, _ message: String) throws -> T {
        guard let value else {
            throw CheckFailure(message)
        }
        return value
    }

    private static func expectValidationMessage(_ settings: APISettings, _ expectedMessage: String, _ checkName: String) throws {
        do {
            try settings.validate()
            throw CheckFailure(checkName)
        } catch let error as APISettingsValidationError {
            try expect(error.userMessage == expectedMessage, checkName)
        }
    }

    private static func expectParsingMessage(
        apiBaseURLText: String = "https://api.example.test",
        latitudeText: String = "37.3382",
        longitudeText: String = "-121.8863",
        expectedMessage: String,
        checkName: String
    ) throws {
        do {
            _ = try APISettings.parse(
                apiBaseURLText: apiBaseURLText,
                latitudeText: latitudeText,
                longitudeText: longitudeText,
                timezoneText: "America/Los_Angeles",
                regionText: "California",
                methodText: "Drik",
                localeText: "en",
                calendarSystemText: "Purnimanta"
            )
            throw CheckFailure(checkName)
        } catch let error as APISettingsValidationError {
            try expect(error.userMessage == expectedMessage, checkName)
        }
    }

    private static func expectSyncPayloadValidation(
        _ dictionary: [String: Any],
        _ expectedMessage: String,
        _ checkName: String
    ) throws {
        do {
            _ = try SettingsSyncPayload(dictionary: dictionary)
            throw CheckFailure(checkName)
        } catch let error as APISettingsValidationError {
            try expect(error.userMessage == expectedMessage, checkName)
        }
    }

    private static func expectSyncPayloadEncodingValidation(
        _ payload: SettingsSyncPayload,
        _ expectedMessage: String,
        _ checkName: String
    ) throws {
        do {
            _ = try payload.dictionary()
            throw CheckFailure(checkName)
        } catch let error as APISettingsValidationError {
            try expect(error.userMessage == expectedMessage, checkName)
        }
    }

    private static func expectSettingsStoreSaveValidation(
        _ store: SettingsStore,
        _ settings: APISettings,
        _ expectedMessage: String,
        _ checkName: String
    ) throws {
        do {
            try store.save(settings)
            throw CheckFailure(checkName)
        } catch let error as APISettingsValidationError {
            try expect(error.userMessage == expectedMessage, checkName)
        }
    }

    private static func expectSettingsStoreLoadValidation(
        _ store: SettingsStore,
        _ expectedMessage: String,
        _ checkName: String
    ) throws {
        do {
            _ = try store.load()
            throw CheckFailure(checkName)
        } catch let error as APISettingsValidationError {
            try expect(error.userMessage == expectedMessage, checkName)
        }
    }

    private static func expectSummaryValidationMessage(
        _ operation: () throws -> Void,
        _ expectedMessage: String,
        _ checkName: String
    ) throws {
        do {
            try operation()
            throw CheckFailure(checkName)
        } catch let error as LocalizedError {
            try expect(error.errorDescription == expectedMessage, checkName)
        }
    }
}

struct CheckFailure: Error, CustomStringConvertible {
    let description: String

    init(_ description: String) {
        self.description = description
    }
}

enum Fixtures {
    static let generatedAt = Date(timeIntervalSince1970: 1_780_401_600)
    static let nextRefreshAt = Date(timeIntervalSince1970: 1_780_452_000)

    static let settings = APISettings(
        apiBaseURL: URL(string: "https://api.example.test")!,
        latitude: 37.3382,
        longitude: -121.8863,
        timezone: "America/Los_Angeles",
        region: "California",
        method: "Drik",
        locale: "en",
        calendarSystem: "Purnimanta"
    )

    static var summary: TithiSummaryResponse {
        TithiSummaryResponse(
            date: "2026-06-02",
            tithi: TithiDetails(
                name: "Ekadashi",
                traditionalName: "Ekadashi",
                number: 11,
                paksha: "Shukla",
                pakshaDay: 11,
                type: "Nanda",
                startTime: Date(timeIntervalSince1970: 1_780_362_000),
                endTime: nextRefreshAt
            ),
            panchaAnga: PanchaAngaSummary(
                nakshatra: "Anuradha",
                yoga: "Siddha",
                karana: "Gara",
                vara: "Tuesday"
            ),
            day: DaySummary(
                sunriseTime: "05:42",
                sunsetTime: "19:01",
                abhijitMuhurta: TimeWindow(
                    name: "Abhijit",
                    startTime: "11:54",
                    endTime: "12:48",
                    auspicious: true
                )
            ),
            calculation: TithiCalculationSummary(
                timezone: "America/Los_Angeles",
                region: "California",
                calendarSystem: "Purnimanta",
                method: "Drik",
                locale: "en"
            ),
            generatedAt: generatedAt,
            nextRefreshAt: nextRefreshAt
        )
    }

    static let summaryJSON = Data("""
    {
      "date": "2026-06-02",
      "tithi": {
        "name": "Ekadashi",
        "traditional_name": "Ekadashi",
        "number": 11,
        "paksha": "Shukla",
        "paksha_day": 11,
        "type": "Nanda",
        "start_time": "2026-06-02T01:00:00Z",
        "end_time": "2026-06-03T02:00:00Z"
      },
      "pancha_anga": {
        "nakshatra": "Anuradha",
        "yoga": "Siddha",
        "karana": "Gara",
        "vara": "Tuesday"
      },
      "day": {
        "sunrise_time": "05:42",
        "sunset_time": "19:01",
        "abhijit_muhurta": {
          "name": "Abhijit",
          "start_time": "11:54",
          "end_time": "12:48",
          "auspicious": true
        }
      },
      "calculation": {
        "timezone": "America/Los_Angeles",
        "region": "California",
        "calendar_system": "Purnimanta",
        "method": "Drik",
        "locale": "en"
      },
      "generated_at": "2026-06-02T12:00:00Z",
      "next_refresh_at": "2026-06-03T02:00:00Z"
    }
    """.utf8)
}

final class StubHTTPSession: PanchangamHTTPSession {
    var lastRequest: URLRequest?
    var data: Data
    var statusCode: Int

    init(data: Data = Data(), statusCode: Int = 200) {
        self.data = data
        self.statusCode = statusCode
    }

    func data(for request: URLRequest) async throws -> (Data, URLResponse) {
        lastRequest = request
        let response = HTTPURLResponse(
            url: request.url!,
            statusCode: statusCode,
            httpVersion: nil,
            headerFields: nil
        )!
        return (data, response)
    }
}
