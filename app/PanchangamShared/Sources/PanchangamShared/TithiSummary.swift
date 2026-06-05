import Foundation

public enum TithiSummaryValidationError: Error, Equatable, LocalizedError, Sendable {
    case missingField(String)
    case invalidTithiNumber
    case invalidPakshaDay
    case invalidTithiWindow
    case invalidGeneratedWindow
    case invalidRefreshWindow
    case invalidTimezone

    public var errorDescription: String? {
        userMessage
    }

    public var userMessage: String {
        switch self {
        case .missingField(let field):
            return "\(field) is required"
        case .invalidTithiNumber:
            return "Tithi number must be between 1 and 30"
        case .invalidPakshaDay:
            return "Paksha day must be between 1 and 15"
        case .invalidTithiWindow:
            return "Tithi end time must be after start time"
        case .invalidGeneratedWindow:
            return "Generated time must be within the tithi start and end time"
        case .invalidRefreshWindow:
            return "Next refresh time must be after generated time and no later than tithi end time"
        case .invalidTimezone:
            return "Timezone is not recognized"
        }
    }
}

public struct TithiSummaryResponse: Codable, Equatable, Sendable {
    public var date: String
    public var tithi: TithiDetails
    public var panchaAnga: PanchaAngaSummary
    public var day: DaySummary
    public var calculation: TithiCalculationSummary
    public var generatedAt: Date
    public var nextRefreshAt: Date

    public init(
        date: String,
        tithi: TithiDetails,
        panchaAnga: PanchaAngaSummary,
        day: DaySummary,
        calculation: TithiCalculationSummary,
        generatedAt: Date,
        nextRefreshAt: Date
    ) {
        self.date = date
        self.tithi = tithi
        self.panchaAnga = panchaAnga
        self.day = day
        self.calculation = calculation
        self.generatedAt = generatedAt
        self.nextRefreshAt = nextRefreshAt
    }

    private enum CodingKeys: String, CodingKey {
        case date
        case tithi
        case panchaAnga = "pancha_anga"
        case day
        case calculation
        case generatedAt = "generated_at"
        case nextRefreshAt = "next_refresh_at"
    }

    public func validate() throws {
        try requireText(date, "Date")
        try tithi.validate()
        try panchaAnga.validate()
        try day.validate()
        try calculation.validate()

        guard generatedAt >= tithi.startTime, generatedAt < tithi.endTime else {
            throw TithiSummaryValidationError.invalidGeneratedWindow
        }

        guard nextRefreshAt > generatedAt, nextRefreshAt <= tithi.endTime else {
            throw TithiSummaryValidationError.invalidRefreshWindow
        }
    }

    private func requireText(_ value: String, _ field: String) throws {
        guard !value.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw TithiSummaryValidationError.missingField(field)
        }
    }
}

public struct TithiDetails: Codable, Equatable, Sendable {
    public var name: String
    public var traditionalName: String
    public var number: Int
    public var paksha: String
    public var pakshaDay: Int
    public var type: String
    public var startTime: Date
    public var endTime: Date

    public init(
        name: String,
        traditionalName: String,
        number: Int,
        paksha: String,
        pakshaDay: Int,
        type: String,
        startTime: Date,
        endTime: Date
    ) {
        self.name = name
        self.traditionalName = traditionalName
        self.number = number
        self.paksha = paksha
        self.pakshaDay = pakshaDay
        self.type = type
        self.startTime = startTime
        self.endTime = endTime
    }

    private enum CodingKeys: String, CodingKey {
        case name
        case traditionalName = "traditional_name"
        case number
        case paksha
        case pakshaDay = "paksha_day"
        case type
        case startTime = "start_time"
        case endTime = "end_time"
    }

    public func validate() throws {
        try requireText(name, "Tithi name")
        try requireText(traditionalName, "Traditional tithi name")
        try requireText(paksha, "Paksha")
        try requireText(type, "Tithi type")

        guard (1...30).contains(number) else {
            throw TithiSummaryValidationError.invalidTithiNumber
        }

        guard (1...15).contains(pakshaDay) else {
            throw TithiSummaryValidationError.invalidPakshaDay
        }

        guard endTime > startTime else {
            throw TithiSummaryValidationError.invalidTithiWindow
        }
    }

    private func requireText(_ value: String, _ field: String) throws {
        guard !value.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw TithiSummaryValidationError.missingField(field)
        }
    }
}

public struct PanchaAngaSummary: Codable, Equatable, Sendable {
    public var nakshatra: String
    public var yoga: String
    public var karana: String
    public var vara: String

    public init(nakshatra: String, yoga: String, karana: String, vara: String) {
        self.nakshatra = nakshatra
        self.yoga = yoga
        self.karana = karana
        self.vara = vara
    }

    public func validate() throws {
        try requireText(nakshatra, "Nakshatra")
        try requireText(yoga, "Yoga")
        try requireText(karana, "Karana")
        try requireText(vara, "Vara")
    }

    private func requireText(_ value: String, _ field: String) throws {
        guard !value.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw TithiSummaryValidationError.missingField(field)
        }
    }
}

public struct DaySummary: Codable, Equatable, Sendable {
    public var sunriseTime: String
    public var sunsetTime: String
    public var abhijitMuhurta: TimeWindow

    public init(sunriseTime: String, sunsetTime: String, abhijitMuhurta: TimeWindow) {
        self.sunriseTime = sunriseTime
        self.sunsetTime = sunsetTime
        self.abhijitMuhurta = abhijitMuhurta
    }

    private enum CodingKeys: String, CodingKey {
        case sunriseTime = "sunrise_time"
        case sunsetTime = "sunset_time"
        case abhijitMuhurta = "abhijit_muhurta"
    }

    public func validate() throws {
        try requireText(sunriseTime, "Sunrise")
        try requireText(sunsetTime, "Sunset")
        try abhijitMuhurta.validate()
    }

    private func requireText(_ value: String, _ field: String) throws {
        guard !value.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw TithiSummaryValidationError.missingField(field)
        }
    }
}

public struct TimeWindow: Codable, Equatable, Sendable {
    public var name: String
    public var startTime: String
    public var endTime: String
    public var auspicious: Bool

    public init(name: String, startTime: String, endTime: String, auspicious: Bool) {
        self.name = name
        self.startTime = startTime
        self.endTime = endTime
        self.auspicious = auspicious
    }

    private enum CodingKeys: String, CodingKey {
        case name
        case startTime = "start_time"
        case endTime = "end_time"
        case auspicious
    }

    public func validate() throws {
        try requireText(name, "Time window name")
        try requireText(startTime, "Time window start")
        try requireText(endTime, "Time window end")
    }

    private func requireText(_ value: String, _ field: String) throws {
        guard !value.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw TithiSummaryValidationError.missingField(field)
        }
    }
}

public struct TithiCalculationSummary: Codable, Equatable, Sendable {
    public var timezone: String
    public var region: String
    public var calendarSystem: String
    public var method: String
    public var locale: String

    public init(timezone: String, region: String, calendarSystem: String, method: String, locale: String) {
        self.timezone = timezone
        self.region = region
        self.calendarSystem = calendarSystem
        self.method = method
        self.locale = locale
    }

    private enum CodingKeys: String, CodingKey {
        case timezone
        case region
        case calendarSystem = "calendar_system"
        case method
        case locale
    }

    public func validate() throws {
        try requireText(timezone, "Timezone")
        try requireText(region, "Region")
        try requireText(calendarSystem, "Calendar system")
        try requireText(method, "Method")
        try requireText(locale, "Locale")

        guard TimeZone(identifier: timezone) != nil else {
            throw TithiSummaryValidationError.invalidTimezone
        }
    }

    private func requireText(_ value: String, _ field: String) throws {
        guard !value.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            throw TithiSummaryValidationError.missingField(field)
        }
    }
}

#if DEBUG
public extension TithiSummaryResponse {
    static var watchSimulatorPreview: TithiSummaryResponse {
        let generatedAt = Date()
        let startTime = generatedAt.addingTimeInterval(-6 * 60 * 60)
        let endTime = generatedAt.addingTimeInterval(6 * 60 * 60)

        return TithiSummaryResponse(
            date: "2026-06-03",
            tithi: TithiDetails(
                name: "Sashthi",
                traditionalName: "Sashthi",
                number: 6,
                paksha: "Shukla",
                pakshaDay: 6,
                type: "Nanda",
                startTime: startTime,
                endTime: endTime
            ),
            panchaAnga: PanchaAngaSummary(nakshatra: "Pushya", yoga: "Saubhagya", karana: "Taitila", vara: "Wednesday"),
            day: DaySummary(
                sunriseTime: "05:42",
                sunsetTime: "19:01",
                abhijitMuhurta: TimeWindow(name: "Abhijit", startTime: "11:54", endTime: "12:48", auspicious: true)
            ),
            calculation: TithiCalculationSummary(
                timezone: "America/Los_Angeles",
                region: "California",
                calendarSystem: "Purnimanta",
                method: "Drik",
                locale: "en"
            ),
            generatedAt: generatedAt,
            nextRefreshAt: generatedAt.addingTimeInterval(60 * 60)
        )
    }
}
#endif
