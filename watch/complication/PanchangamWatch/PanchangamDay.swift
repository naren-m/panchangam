import Foundation

struct GoodTimePeriod: Equatable {
    let name: String
    let windowText: String
}

struct PanchangamDay: Decodable, Equatable {
    let date: String
    let tithi: String
    let tithiEndText: String
    let paksha: String
    let nakshatra: String
    let yoga: String
    let karana: String
    let sunrise: String
    let sunset: String
    let masaText: String
    let rituText: String
    let raasiText: String
    let samvatsaraText: String

    enum CodingKeys: String, CodingKey {
        case date
        case tithi
        case paksha
        case nakshatra
        case yoga
        case karana
        case sunriseTime = "sunrise_time"
        case sunsetTime = "sunset_time"
        case panchaAnga = "pancha_anga"
        case day
        case calculation
        case masa
        case maasa
        case lunarMonth = "lunar_month"
        case ritu
        case ruthu
        case rashi
        case raasi
        case samvatsara
        case samvathsara
    }

    init(
        date: String,
        tithi: String,
        tithiEndText: String,
        paksha: String = "",
        nakshatra: String,
        yoga: String,
        karana: String,
        sunrise: String,
        sunset: String,
        masaText: String,
        rituText: String,
        raasiText: String,
        samvatsaraText: String
    ) {
        self.date = date
        self.tithi = tithi
        self.tithiEndText = tithiEndText
        self.paksha = paksha
        self.nakshatra = nakshatra
        self.yoga = yoga
        self.karana = karana
        self.sunrise = sunrise
        self.sunset = sunset
        self.masaText = masaText
        self.rituText = rituText
        self.raasiText = raasiText
        self.samvatsaraText = samvatsaraText
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        let panchaAnga = try? container.decode(CurrentPanchaAnga.self, forKey: .panchaAnga)
        let currentDay = try? container.decode(CurrentDay.self, forKey: .day)
        let calculation = try? container.decode(CurrentCalculation.self, forKey: .calculation)
        let currentTithi = try? container.decode(CurrentTithi.self, forKey: .tithi)

        date = try container.decode(String.self, forKey: .date)
        let fallbackContext = PanchangamDay.calendarContextFallback(for: date)

        if let currentTithi {
            tithi = PanchangamDay.cleanName(currentTithi.name)
            paksha = PanchangamDay.cleanName(currentTithi.paksha)
            tithiEndText = PanchangamDay.timeLabel(
                prefix: "Tithi ends at",
                value: currentTithi.endTime,
                timezoneIdentifier: calculation?.timezone
            )
        } else {
            let rawTithi = try container.decode(String.self, forKey: .tithi)
            tithi = PanchangamDay.cleanName(rawTithi)
            paksha = PanchangamDay.cleanName(
                (try? container.decodeIfPresent(String.self, forKey: .paksha)) ??
                    PanchangamDay.pakshaFromName(rawTithi)
            )
            tithiEndText = "Tithi end unavailable"
        }

        let rawNakshatra: String
        let rawYoga: String
        let rawKarana: String
        let rawSunrise: String
        let rawSunset: String

        if let panchaAnga {
            rawNakshatra = panchaAnga.nakshatra
            rawYoga = panchaAnga.yoga
            rawKarana = panchaAnga.karana
        } else {
            rawNakshatra = try container.decode(String.self, forKey: .nakshatra)
            rawYoga = try container.decode(String.self, forKey: .yoga)
            rawKarana = try container.decode(String.self, forKey: .karana)
        }

        if let currentDay {
            rawSunrise = currentDay.sunriseTime
            rawSunset = currentDay.sunsetTime
        } else {
            rawSunrise = try container.decode(String.self, forKey: .sunriseTime)
            rawSunset = try container.decode(String.self, forKey: .sunsetTime)
        }

        nakshatra = PanchangamDay.cleanName(rawNakshatra)
        yoga = PanchangamDay.cleanName(rawYoga)
        karana = PanchangamDay.cleanName(rawKarana)
        sunrise = PanchangamDay.shortTime(rawSunrise, timezoneIdentifier: calculation?.timezone)
        sunset = PanchangamDay.shortTime(rawSunset, timezoneIdentifier: calculation?.timezone)
        masaText = PanchangamDay.suffixedLabel(
            value: PanchangamDay.decodeFirstString(from: container, keys: [.masa, .maasa, .lunarMonth]) ?? fallbackContext.masa,
            suffix: "Masa",
            fallback: "Masa unavailable"
        )
        rituText = PanchangamDay.suffixedLabel(
            value: PanchangamDay.decodeFirstString(from: container, keys: [.ritu, .ruthu]) ?? fallbackContext.ritu,
            suffix: "Ritu",
            fallback: "Ritu unavailable"
        )
        raasiText = PanchangamDay.suffixedLabel(
            value: PanchangamDay.decodeFirstString(from: container, keys: [.raasi, .rashi]) ?? fallbackContext.raasi,
            suffix: "Raasi",
            fallback: "Raasi unavailable"
        )
        samvatsaraText = PanchangamDay.suffixedLabel(
            value: PanchangamDay.decodeFirstString(from: container, keys: [.samvatsara, .samvathsara]) ?? fallbackContext.samvatsara,
            suffix: "Samvatsara",
            fallback: "Samvatsara unavailable"
        )
    }

    var dateHeading: String {
        guard let date = Self.parsedDate(date) else {
            return date
        }
        let formatter = DateFormatter()
        formatter.calendar = Self.gregorianCalendar
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        formatter.dateFormat = "E d MMM"
        return formatter.string(from: date).uppercased()
    }

    var weekdayText: String {
        guard let index = Self.weekdayIndex(for: date) else {
            return "Vāra unavailable"
        }
        return Self.weekdayNames[index - 1]
    }

    var pakshaText: String {
        guard !paksha.isEmpty else {
            return "Pakṣa"
        }
        return "\(paksha) Pakṣa"
    }

    var tithiFaceTitle: String {
        Self.romanTithiName(for: tithi)
    }

    var tithiLocalName: String {
        Self.devanagariTithiName(for: tithi)
    }

    var nakshatraLocalName: String {
        Self.devanagariNakshatraName(for: nakshatra)
    }

    var tithiEndFaceText: String {
        tithiEndText.replacingOccurrences(of: "Tithi ends at", with: "Tithi ends")
    }

    var tithiFaceDetailText: String {
        "\(tithiEndFaceText) • \(Self.trimSuffix(masaText, suffix: "Masa")) • \(Self.trimSuffix(rituText, suffix: "Ritu"))"
    }

    var sunlightText: String {
        "\(sunrise) - \(sunset)"
    }

    var abhijitWindowText: String {
        guard let period = goodPeriodWindows.first(where: { $0.name == "Abhijit Muhūrta" }) else {
            return "--:-- - --:--"
        }
        return period.windowText
    }

    var rahuStartText: String {
        guard
            let weekdayIndex = Self.weekdayIndex(for: date),
            let sunriseMinutes = Self.clockMinutes(sunrise),
            let sunsetMinutes = Self.clockMinutes(sunset),
            let rahuPart = Self.rahuPartByWeekday[weekdayIndex]
        else {
            return "Rahu --:--"
        }
        let partDuration = Double(sunsetMinutes - sunriseMinutes) / 8
        let start = Double(sunriseMinutes) + (Double(rahuPart - 1) * partDuration)
        return "Rahu \(Self.clockText(minutes: start))"
    }

    func nextGoodPeriod(at date: Date) -> GoodTimePeriod {
        nextGoodPeriod(atClockMinutes: Self.clockMinutes(from: date))
    }

    func nextGoodPeriod(atClockMinutes clockMinutes: Int) -> GoodTimePeriod {
        let periods = goodPeriodWindows
        guard let firstPeriod = periods.first else {
            return GoodTimePeriod(name: "Good Time", windowText: "--:-- - --:--")
        }
        let normalizedClock = Double(max(0, min(clockMinutes, 1439)))
        let period = periods.first { $0.endMinutes > normalizedClock } ?? firstPeriod
        return GoodTimePeriod(name: period.name, windowText: period.windowText)
    }

    private var goodPeriodWindows: [GoodPeriodWindow] {
        guard let sunriseMinutes = Self.clockMinutes(sunrise), let sunsetMinutes = Self.clockMinutes(sunset) else {
            return []
        }

        let dayLength = sunsetMinutes - sunriseMinutes
        guard dayLength > 0 else {
            return []
        }

        let solarNoon = Double(sunriseMinutes) + Double(dayLength) / 2
        let periods = [
            GoodPeriodWindow(name: "Brahma Muhūrta", startMinutes: Double(sunriseMinutes - 96), endMinutes: Double(sunriseMinutes)),
            GoodPeriodWindow(name: "Pratah Kaal", startMinutes: Double(sunriseMinutes), endMinutes: Double(sunriseMinutes + 180)),
            GoodPeriodWindow(name: "Abhijit Muhūrta", startMinutes: solarNoon - 24, endMinutes: solarNoon + 24),
            GoodPeriodWindow(name: "Sandhya Kaal", startMinutes: Double(sunsetMinutes - 60), endMinutes: Double(sunsetMinutes))
        ]

        return periods.sorted { $0.startMinutes < $1.startMinutes }
    }

    static let sample = PanchangamDay(
        date: "2026-06-02",
        tithi: "Sasthi",
        tithiEndText: "Tithi ends at 14:32",
        paksha: "Shukla",
        nakshatra: "Pushya",
        yoga: "Saubhagya",
        karana: "Taitila",
        sunrise: "05:42",
        sunset: "19:01",
        masaText: "Jyeshtha Masa",
        rituText: "Grishma Ritu",
        raasiText: "Makara Raasi",
        samvatsaraText: "Parabhava Samvatsara"
    )

    private static func cleanName(_ value: String) -> String {
        let withoutCount = value.split(separator: "(", maxSplits: 1).first.map(String.init) ?? value
        return withoutCount.split(separator: "-", maxSplits: 1).first.map {
            String($0).trimmingCharacters(in: .whitespacesAndNewlines)
        } ?? withoutCount.trimmingCharacters(in: .whitespacesAndNewlines)
    }

    private static func pakshaFromName(_ value: String) -> String {
        if value.localizedCaseInsensitiveContains("Shukla") {
            return "Shukla"
        }
        if value.localizedCaseInsensitiveContains("Krishna") {
            return "Krishna"
        }
        return ""
    }

    private static func shortTime(_ value: String, timezoneIdentifier: String?) -> String {
        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        if let date = ISO8601DateFormatter().date(from: trimmed) {
            let formatter = DateFormatter()
            formatter.calendar = Calendar(identifier: .gregorian)
            formatter.locale = Locale(identifier: "en_US_POSIX")
            formatter.dateFormat = "HH:mm"
            if let timezoneIdentifier, let timezone = TimeZone(identifier: timezoneIdentifier) {
                formatter.timeZone = timezone
            }
            return formatter.string(from: date)
        }
        guard trimmed.count >= 5 else {
            return trimmed
        }
        return String(trimmed.prefix(5))
    }

    private static func timeLabel(prefix: String, value: String, timezoneIdentifier: String?) -> String {
        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else {
            return "Tithi end unavailable"
        }
        return "\(prefix) \(shortTime(trimmed, timezoneIdentifier: timezoneIdentifier))"
    }

    private static func suffixedLabel(value: String?, suffix: String, fallback: String) -> String {
        guard let value, !value.isEmpty else {
            return fallback
        }
        if value.localizedCaseInsensitiveContains(suffix) {
            return value
        }
        return "\(value) \(suffix)"
    }

    private static func trimSuffix(_ value: String, suffix: String) -> String {
        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        let marker = " \(suffix)"
        guard trimmed.lowercased().hasSuffix(marker.lowercased()) else {
            return trimmed
        }
        return String(trimmed.dropLast(marker.count))
    }

    private static func parsedDate(_ value: String) -> Date? {
        let formatter = DateFormatter()
        formatter.calendar = gregorianCalendar
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        formatter.dateFormat = "yyyy-MM-dd"
        return formatter.date(from: String(value.prefix(10)))
    }

    private static func weekdayIndex(for value: String) -> Int? {
        guard let date = parsedDate(value) else {
            return nil
        }
        return gregorianCalendar.component(.weekday, from: date)
    }

    private static func clockMinutes(_ value: String) -> Int? {
        let parts = value.prefix(5).split(separator: ":")
        guard
            parts.count == 2,
            let hours = Int(parts[0]),
            let minutes = Int(parts[1]),
            (0...23).contains(hours),
            (0...59).contains(minutes)
        else {
            return nil
        }
        return (hours * 60) + minutes
    }

    private static func clockText(minutes: Double) -> String {
        let rounded = Int(minutes.rounded())
        let total = ((rounded % 1440) + 1440) % 1440
        return String(format: "%02d:%02d", total / 60, total % 60)
    }

    static func clockTextForWindow(_ minutes: Double) -> String {
        clockText(minutes: minutes)
    }

    private static func clockMinutes(from date: Date) -> Int {
        let components = Calendar.current.dateComponents([.hour, .minute], from: date)
        return ((components.hour ?? 0) * 60) + (components.minute ?? 0)
    }

    private static func romanTithiName(for value: String) -> String {
        romanTithiNames[normalizedKey(value)] ?? value
    }

    private static func devanagariTithiName(for value: String) -> String {
        devanagariTithiNames[normalizedKey(value)] ?? ""
    }

    private static func devanagariNakshatraName(for value: String) -> String {
        devanagariNakshatraNames[normalizedKey(value)] ?? ""
    }

    private static func normalizedKey(_ value: String) -> String {
        let folded = value.folding(options: [.diacriticInsensitive, .caseInsensitive], locale: Locale(identifier: "en_US_POSIX"))
        let scalars = folded.unicodeScalars.filter { CharacterSet.letters.contains($0) }
        return String(String.UnicodeScalarView(scalars)).lowercased()
    }

    private static func decodeFirstString(
        from container: KeyedDecodingContainer<CodingKeys>,
        keys: [CodingKeys]
    ) -> String? {
        for key in keys {
            if let value = try? container.decodeIfPresent(String.self, forKey: key) {
                let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
                if !trimmed.isEmpty {
                    return trimmed
                }
            }
        }
        return nil
    }

    private static func calendarContextFallback(for date: String) -> CalendarContext {
        let parts = date.prefix(10).split(separator: "-")
        guard
            parts.count >= 2,
            let year = Int(parts[0]),
            let month = Int(parts[1]),
            (1...12).contains(month)
        else {
            return CalendarContext(masa: nil, ritu: nil, raasi: nil, samvatsara: nil)
        }

        let masaByMonth = [
            "Pausha", "Magha", "Phalguna", "Chaitra",
            "Vaishakha", "Jyeshtha", "Ashadha", "Shravana",
            "Bhadrapada", "Ashwin", "Kartika", "Margashirsha"
        ]
        let rituByMonth = [
            "Shishira", "Shishira", "Vasanta", "Vasanta",
            "Grishma", "Grishma", "Varsha", "Varsha",
            "Sharad", "Sharad", "Hemanta", "Hemanta"
        ]
        let raasiByMonth = [
            "Makara", "Kumbha", "Meena", "Mesha",
            "Vrishabha", "Mithuna", "Karka", "Simha",
            "Kanya", "Tula", "Vrischika", "Dhanu"
        ]

        return CalendarContext(
            masa: masaByMonth[month - 1],
            ritu: rituByMonth[month - 1],
            raasi: raasiByMonth[month - 1],
            samvatsara: samvatsaraName(year: year, month: month)
        )
    }

    private static func samvatsaraName(year: Int, month: Int) -> String {
        let names = [
            "Prabhava", "Vibhava", "Shukla", "Pramodoota", "Prajothpatti",
            "Angirasa", "Srimukha", "Bhava", "Yuva", "Dhata",
            "Ishvara", "Bahudhanya", "Pramadi", "Vikrama", "Vrisha",
            "Chitrabhanu", "Svabhanu", "Tarana", "Parthiva", "Vyaya",
            "Sarvajit", "Sarvadhari", "Virodhi", "Vikruti", "Khara",
            "Nandana", "Vijaya", "Jaya", "Manmatha", "Durmukhi",
            "Hevilambi", "Vilambi", "Vikari", "Sharvari", "Plava",
            "Shubhakrith", "Shobhakrith", "Krodhi", "Vishvavasu", "Parabhava",
            "Plavanga", "Kilaka", "Saumya", "Sadharana", "Virodhikrit",
            "Paridhavi", "Pramadeecha", "Ananda", "Rakshasa", "Nala",
            "Pingala", "Kalayukthi", "Siddharthi", "Raudra", "Durmathi",
            "Dundubhi", "Rudhirodgari", "Raktakshi", "Krodhana", "Akshaya"
        ]
        let cycleYear = month < 4 ? year - 1 : year
        return names[(cycleYear + 53) % names.count]
    }

    private static let gregorianCalendar: Calendar = {
        var calendar = Calendar(identifier: .gregorian)
        calendar.timeZone = TimeZone(secondsFromGMT: 0)!
        return calendar
    }()

    private static let weekdayNames = [
        "Ravivāra",
        "Somavāra",
        "Maṅgalavāra",
        "Budhavāra",
        "Guruvāra",
        "Śukravāra",
        "Śanivāra"
    ]

    private static let rahuPartByWeekday = [
        1: 7,
        2: 1,
        3: 6,
        4: 4,
        5: 3,
        6: 2,
        7: 5
    ]

    private static let romanTithiNames = [
        "pratipada": "Pratipadā",
        "dwitiya": "Dvitīyā",
        "dvitiya": "Dvitīyā",
        "tritiya": "Tṛtīyā",
        "trutiya": "Tṛtīyā",
        "chaturthi": "Chaturthī",
        "chathurthi": "Chaturthī",
        "panchami": "Pañchamī",
        "sasthi": "Ṣaṣṭhī",
        "shashti": "Ṣaṣṭhī",
        "sashti": "Ṣaṣṭhī",
        "saptami": "Saptamī",
        "ashtami": "Aṣṭamī",
        "navami": "Navamī",
        "dashami": "Daśamī",
        "ekadashi": "Ekādaśī",
        "dwadashi": "Dvādaśī",
        "trayodashi": "Trayodaśī",
        "chaturdashi": "Chaturdaśī",
        "purnima": "Pūrṇimā",
        "amavasya": "Amāvasyā"
    ]

    private static let devanagariTithiNames = [
        "pratipada": "प्रतिपदा",
        "dwitiya": "द्वितीया",
        "dvitiya": "द्वितीया",
        "tritiya": "तृतीया",
        "trutiya": "तृतीया",
        "chaturthi": "चतुर्थी",
        "chathurthi": "चतुर्थी",
        "panchami": "पञ्चमी",
        "sasthi": "षष्ठी",
        "shashti": "षष्ठी",
        "sashti": "षष्ठी",
        "saptami": "सप्तमी",
        "ashtami": "अष्टमी",
        "navami": "नवमी",
        "dashami": "दशमी",
        "ekadashi": "एकादशी",
        "dwadashi": "द्वादशी",
        "trayodashi": "त्रयोदशी",
        "chaturdashi": "चतुर्दशी",
        "purnima": "पूर्णिमा",
        "amavasya": "अमावस्या"
    ]

    private static let devanagariNakshatraNames = [
        "ashwini": "अश्विनी",
        "bharani": "भरणी",
        "krittika": "कृत्तिका",
        "rohini": "रोहिणी",
        "mrigashira": "मृगशीर्षा",
        "ardra": "आर्द्रा",
        "punarvasu": "पुनर्वसु",
        "pushya": "पुष्य",
        "ashlesha": "आश्लेषा",
        "magha": "मघा",
        "purvaphalguni": "पूर्व फाल्गुनी",
        "uttaraphalguni": "उत्तर फाल्गुनी",
        "hasta": "हस्त",
        "chitra": "चित्रा",
        "swati": "स्वाती",
        "vishakha": "विशाखा",
        "anuradha": "अनुराधा",
        "jyeshtha": "ज्येष्ठा",
        "mula": "मूल",
        "purvashada": "पूर्वाषाढा",
        "uttarashada": "उत्तराषाढा",
        "shravana": "श्रवण",
        "dhanishta": "धनिष्ठा",
        "shatabhisha": "शतभिषा",
        "purvabhadra": "पूर्व भाद्रपदा",
        "uttarabhadra": "उत्तर भाद्रपदा",
        "revati": "रेवती"
    ]
}

private struct GoodPeriodWindow {
    let name: String
    let startMinutes: Double
    let endMinutes: Double

    var windowText: String {
        "\(PanchangamDay.clockTextForWindow(startMinutes)) - \(PanchangamDay.clockTextForWindow(endMinutes))"
    }
}

private struct CalendarContext {
    let masa: String?
    let ritu: String?
    let raasi: String?
    let samvatsara: String?
}

private struct CurrentTithi: Decodable {
    let name: String
    let paksha: String
    let endTime: String

    enum CodingKeys: String, CodingKey {
        case name
        case paksha
        case traditionalName = "traditional_name"
        case endTime = "end_time"
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        name = (try? container.decodeIfPresent(String.self, forKey: .name)) ??
            (try? container.decodeIfPresent(String.self, forKey: .traditionalName)) ??
            "Tithi"
        paksha = (try? container.decodeIfPresent(String.self, forKey: .paksha)) ?? ""
        endTime = (try? container.decodeIfPresent(String.self, forKey: .endTime)) ?? ""
    }
}

private struct CurrentPanchaAnga: Decodable {
    let nakshatra: String
    let yoga: String
    let karana: String
}

private struct CurrentDay: Decodable {
    let sunriseTime: String
    let sunsetTime: String

    enum CodingKeys: String, CodingKey {
        case sunriseTime = "sunrise_time"
        case sunsetTime = "sunset_time"
    }
}

private struct CurrentCalculation: Decodable {
    let timezone: String
}
