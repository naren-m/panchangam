import XCTest
@testable import PanchangamWatch

final class PanchangamDayTests: XCTestCase {
    func testSampleUsesSanJoseReferenceCurrentValues() {
        XCTAssertEqual(PanchangamDay.sample.date, "2026-06-02")
        XCTAssertEqual(PanchangamDay.sample.tithi, "Sasthi")
        XCTAssertEqual(PanchangamDay.sample.paksha, "Shukla")
        XCTAssertEqual(PanchangamDay.sample.nakshatra, "Pushya")
        XCTAssertEqual(PanchangamDay.sample.yoga, "Saubhagya")
        XCTAssertEqual(PanchangamDay.sample.karana, "Taitila")
        XCTAssertEqual(PanchangamDay.sample.tithiEndText, "Tithi ends at 14:32")
        XCTAssertEqual(PanchangamDay.sample.sunrise, "05:42")
        XCTAssertEqual(PanchangamDay.sample.sunset, "19:01")
        XCTAssertEqual(PanchangamDay.sample.raasiText, "Makara Raasi")
        XCTAssertEqual(PanchangamDay.sample.dateHeading, "TUE 2 JUN")
        XCTAssertEqual(PanchangamDay.sample.tithiFaceTitle, "Ṣaṣṭhī")
    }

    func testFetchTodayUsesExactCurrentInstantForCurrentEndpoint() async throws {
        let oldTimeZone = NSTimeZone.default
        NSTimeZone.default = TimeZone(identifier: "America/Los_Angeles")!
        defer { NSTimeZone.default = oldTimeZone }

        CapturingURLProtocol.capturedRequest = nil
        CapturingURLProtocol.stubData = """
        {
          "date": "2026-06-04",
          "tithi": {
            "name": "Panchami",
            "paksha": "Krishna",
            "paksha_day": 5,
            "end_time": "2026-06-05T19:03:01Z"
          },
          "pancha_anga": {
            "nakshatra": "Shravana (22)",
            "yoga": "Brahma (25)",
            "karana": "Kaulava (4)"
          },
          "day": {
            "sunrise_time": "05:47:41",
            "sunset_time": "20:23:56"
          },
          "calculation": {
            "timezone": "America/Los_Angeles"
          }
        }
        """.data(using: .utf8)!

        let configuration = URLSessionConfiguration.ephemeral
        configuration.protocolClasses = [CapturingURLProtocol.self]
        let client = PanchangamAPIClient(
            baseURL: URL(string: "http://example.test")!,
            session: URLSession(configuration: configuration)
        )
        let date = ISO8601DateFormatter().date(from: "2026-06-04T22:28:04Z")!

        _ = try await client.fetchToday(date: date)

        let requestURL = try XCTUnwrap(CapturingURLProtocol.capturedRequest?.url)
        XCTAssertEqual(requestURL.path, "/api/v1/tithi/current")
        XCTAssertEqual(queryValue("at", in: requestURL), "2026-06-04T15:28:04-07:00")
        XCTAssertEqual(queryValue("lat", in: requestURL), "37.3382")
        XCTAssertEqual(queryValue("lng", in: requestURL), "-121.8863")
        XCTAssertEqual(queryValue("tz", in: requestURL), "America/Los_Angeles")
        XCTAssertEqual(queryValue("calendar_system", in: requestURL), "Purnimanta")
    }

    func testDecodesGatewayResponseIntoWatchSummary() throws {
        let json = """
        {
          "date": "2026-06-04",
          "tithi": "Sasthi (6)",
          "nakshatra": "Pushya (8)",
          "yoga": "Saubhagya",
          "karana": "Taitila",
          "sunrise_time": "05:42:11",
          "sunset_time": "19:01:07"
        }
        """.data(using: .utf8)!

        let day = try JSONDecoder().decode(PanchangamDay.self, from: json)

        XCTAssertEqual(day.tithi, "Sasthi")
        XCTAssertEqual(day.nakshatra, "Pushya")
        XCTAssertEqual(day.yoga, "Saubhagya")
        XCTAssertEqual(day.karana, "Taitila")
        XCTAssertEqual(day.sunrise, "05:42")
        XCTAssertEqual(day.sunset, "19:01")
        XCTAssertEqual(day.masaText, "Jyeshtha Masa")
        XCTAssertEqual(day.rituText, "Grishma Ritu")
        XCTAssertEqual(day.raasiText, "Mithuna Raasi")
        XCTAssertEqual(day.samvatsaraText, "Parabhava Samvatsara")
    }

    func testDecodesLongGatewayNamesIntoCompactWatchLabels() throws {
        let json = """
        {
          "date": "2026-06-04",
          "tithi": "Chathurthi - Krishna Paksha Day 4",
          "nakshatra": "Punarvasu - Pada 4",
          "yoga": "Saubhagya",
          "karana": "Taitila",
          "sunrise_time": "05:47:22",
          "sunset_time": "20:23:12"
        }
        """.data(using: .utf8)!

        let day = try JSONDecoder().decode(PanchangamDay.self, from: json)

        XCTAssertEqual(day.tithi, "Chathurthi")
        XCTAssertEqual(day.nakshatra, "Punarvasu")
        XCTAssertEqual(day.sunrise, "05:47")
        XCTAssertEqual(day.sunset, "20:23")
    }

    func testDecodesCurrentTithiResponseIntoUsefulWatchLabels() throws {
        let json = """
        {
          "date": "2026-06-04",
          "tithi": {
            "name": "Chathurthi",
            "paksha": "Krishna",
            "paksha_day": 4,
            "end_time": "2026-06-05T13:18:00Z"
          },
          "pancha_anga": {
            "nakshatra": "Punarvasu - Pada 4",
            "yoga": "Saubhagya",
            "karana": "Taitila"
          },
          "day": {
            "sunrise_time": "05:47:22",
            "sunset_time": "20:23:12"
          },
          "calculation": {
            "timezone": "America/Los_Angeles"
          },
          "masa": "Jyeshtha",
          "ritu": "Grishma",
          "raasi": "Makara",
          "samvatsara": "Parabhava"
        }
        """.data(using: .utf8)!

        let day = try JSONDecoder().decode(PanchangamDay.self, from: json)

        XCTAssertEqual(day.tithi, "Chathurthi")
        XCTAssertEqual(day.tithiEndText, "Tithi ends at 06:18")
        XCTAssertEqual(day.nakshatra, "Punarvasu")
        XCTAssertEqual(day.masaText, "Jyeshtha Masa")
        XCTAssertEqual(day.rituText, "Grishma Ritu")
        XCTAssertEqual(day.raasiText, "Makara Raasi")
        XCTAssertEqual(day.samvatsaraText, "Parabhava Samvatsara")
    }

    func testBuildsReferenceFaceLabelsFromCurrentResponse() throws {
        let json = """
        {
          "date": "2026-06-02",
          "tithi": {
            "name": "Sasthi",
            "paksha": "Shukla",
            "paksha_day": 6,
            "end_time": "2026-06-02T21:32:00Z"
          },
          "pancha_anga": {
            "nakshatra": "Pushya (8)",
            "yoga": "Saubhagya",
            "karana": "Taitila"
          },
          "day": {
            "sunrise_time": "05:42:11",
            "sunset_time": "19:01:07"
          },
          "calculation": {
            "timezone": "America/Los_Angeles"
          },
          "masa": "Jyeshtha",
          "ritu": "Grishma",
          "raasi": "Makara",
          "samvatsara": "Parabhava"
        }
        """.data(using: .utf8)!

        let day = try JSONDecoder().decode(PanchangamDay.self, from: json)

        XCTAssertEqual(day.dateHeading, "TUE 2 JUN")
        XCTAssertEqual(day.weekdayText, "Maṅgalavāra")
        XCTAssertEqual(day.pakshaText, "Shukla Pakṣa")
        XCTAssertEqual(day.tithiFaceTitle, "Ṣaṣṭhī")
        XCTAssertEqual(day.tithiLocalName, "षष्ठी")
        XCTAssertEqual(day.nakshatraLocalName, "पुष्य")
        XCTAssertEqual(day.tithiEndFaceText, "Tithi ends 14:32")
        XCTAssertEqual(day.tithiFaceDetailText, "Tithi ends 14:32 • Jyeshtha • Grishma")
        XCTAssertEqual(day.sunlightText, "05:42 - 19:01")
        XCTAssertEqual(day.abhijitWindowText, "11:58 - 12:46")
        XCTAssertEqual(day.rahuStartText, "Rahu 14:01")
    }

    func testSelectsNextGoodPeriodFromKnownDailyWindows() {
        let beforeSunrise = PanchangamDay.sample.nextGoodPeriod(atClockMinutes: (4 * 60) + 30)
        XCTAssertEqual(beforeSunrise.name, "Brahma Muhūrta")
        XCTAssertEqual(beforeSunrise.windowText, "04:06 - 05:42")

        let afterSunrise = PanchangamDay.sample.nextGoodPeriod(atClockMinutes: (6 * 60) + 30)
        XCTAssertEqual(afterSunrise.name, "Pratah Kaal")
        XCTAssertEqual(afterSunrise.windowText, "05:42 - 08:42")

        let lateMorning = PanchangamDay.sample.nextGoodPeriod(atClockMinutes: (9 * 60) + 30)
        XCTAssertEqual(lateMorning.name, "Abhijit Muhūrta")
        XCTAssertEqual(lateMorning.windowText, "11:58 - 12:46")

        let evening = PanchangamDay.sample.nextGoodPeriod(atClockMinutes: (17 * 60) + 30)
        XCTAssertEqual(evening.name, "Sandhya Kaal")
        XCTAssertEqual(evening.windowText, "18:01 - 19:01")

        let afterAllPeriods = PanchangamDay.sample.nextGoodPeriod(atClockMinutes: 23 * 60)
        XCTAssertEqual(afterAllPeriods.name, "Brahma Muhūrta")
        XCTAssertEqual(afterAllPeriods.windowText, "04:06 - 05:42")
    }
}

private final class CapturingURLProtocol: URLProtocol {
    static var capturedRequest: URLRequest?
    static var stubData = Data()

    override class func canInit(with request: URLRequest) -> Bool {
        true
    }

    override class func canonicalRequest(for request: URLRequest) -> URLRequest {
        request
    }

    override func startLoading() {
        Self.capturedRequest = request
        let response = HTTPURLResponse(
            url: request.url!,
            statusCode: 200,
            httpVersion: nil,
            headerFields: ["Content-Type": "application/json"]
        )!
        client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: Self.stubData)
        client?.urlProtocolDidFinishLoading(self)
    }

    override func stopLoading() {}
}

private func queryValue(_ name: String, in url: URL) -> String? {
    URLComponents(url: url, resolvingAgainstBaseURL: false)?
        .queryItems?
        .first { $0.name == name }?
        .value
}
