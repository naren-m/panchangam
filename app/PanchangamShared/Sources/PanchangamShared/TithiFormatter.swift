import Foundation

public struct TithiFormatter {
    private static let nakshatraNames = [
        "Ashwini",
        "Bharani",
        "Krittika",
        "Rohini",
        "Mrigashira",
        "Ardra",
        "Punarvasu",
        "Pushya",
        "Ashlesha",
        "Magha",
        "Purva Phalguni",
        "Uttara Phalguni",
        "Hasta",
        "Chitra",
        "Swati",
        "Vishakha",
        "Anuradha",
        "Jyeshtha",
        "Mula",
        "Purva Ashadha",
        "Uttara Ashadha",
        "Shravana",
        "Dhanishta",
        "Shatabhisha",
        "Purva Bhadrapada",
        "Uttara Bhadrapada",
        "Revati"
    ]

    public init() {}

    public func inlineText(for summary: TithiSummaryResponse) -> String {
        "\(tithiDisplayName(for: summary)), \(summary.tithi.paksha) \(summary.tithi.pakshaDay)"
    }

    public func complicationInlineText(for summary: TithiSummaryResponse, isStale: Bool) -> String {
        let text = inlineText(for: summary)
        return isStale ? "Cached: \(text)" : text
    }

    public func complicationSecondaryText(for summary: TithiSummaryResponse, isStale: Bool) -> String {
        isStale ? "Cached" : summary.panchaAnga.nakshatra
    }

    public func circularText(for summary: TithiSummaryResponse) -> String {
        "\(summary.tithi.name.prefix(3)) \(summary.tithi.pakshaDay)"
    }

    public func tithiProgress(for summary: TithiSummaryResponse, now: Date) -> Double {
        let duration = summary.tithi.endTime.timeIntervalSince(summary.tithi.startTime)
        guard duration > 0 else {
            return 1
        }

        let elapsed = now.timeIntervalSince(summary.tithi.startTime)
        let progress = min(1, max(0, elapsed / duration))
        return (progress * 100).rounded() / 100
    }

    public func progressText(for summary: TithiSummaryResponse, now: Date) -> String {
        let percent = Int((tithiProgress(for: summary, now: now) * 100).rounded())
        return "\(percent)% elapsed"
    }

    public func nakshatraIndex(for summary: TithiSummaryResponse) -> Int? {
        let currentName = normalizedNakshatraName(summary.panchaAnga.nakshatra)
        return Self.nakshatraNames.firstIndex { name in
            normalizedNakshatraName(name) == currentName
        }
    }

    public func rectangularLines(for summary: TithiSummaryResponse) -> [String] {
        rectangularLines(for: summary, now: summary.generatedAt)
    }

    public func rectangularLines(for summary: TithiSummaryResponse, now: Date) -> [String] {
        rectangularLines(for: summary, now: now, isStale: false)
    }

    public func rectangularLines(for summary: TithiSummaryResponse, now: Date, isStale: Bool) -> [String] {
        let timeLabel = isStale ? "Cached" : "Remaining"
        return [
            inlineText(for: summary),
            "Nakshatra: \(summary.panchaAnga.nakshatra)",
            "Yoga/Karana: \(summary.panchaAnga.yoga) / \(summary.panchaAnga.karana)",
            "\(summary.day.abhijitMuhurta.name): \(abhijitText(for: summary))",
            "\(timeLabel): \(remainingText(for: summary, now: now))"
        ]
    }

    public func remainingText(for summary: TithiSummaryResponse, now: Date) -> String {
        let secondsLeft = summary.tithi.endTime.timeIntervalSince(now)
        if secondsLeft == 0 {
            return "ending now"
        }

        let prefix = secondsLeft > 0 ? "" : "ended "
        let suffix = secondsLeft > 0 ? " left" : " ago"
        let totalMinutes = max(1, Int(ceil(abs(secondsLeft) / 60)))
        let durationText = durationText(minutes: totalMinutes)

        return "\(prefix)\(durationText)\(suffix)"
    }

    private func durationText(minutes totalMinutes: Int) -> String {
        let hours = totalMinutes / 60
        let minutes = totalMinutes % 60

        if hours > 0 {
            return "\(hours)h \(minutes)m"
        }

        return "\(minutes)m"
    }

    public func detailRows(for summary: TithiSummaryResponse) -> [(String, String)] {
        detailRows(for: summary, now: summary.generatedAt)
    }

    public func detailRows(for summary: TithiSummaryResponse, now: Date) -> [(String, String)] {
        let timeFormatter = timeFormatter(for: summary.calculation.timezone)
        return [
            ("Tithi", inlineText(for: summary)),
            ("Tithi Number", String(summary.tithi.number)),
            ("Traditional Name", summary.tithi.traditionalName),
            ("Paksha", summary.tithi.paksha),
            ("Paksha Day", String(summary.tithi.pakshaDay)),
            ("Tithi Type", summary.tithi.type),
            ("Nakshatra", summary.panchaAnga.nakshatra),
            ("Yoga", summary.panchaAnga.yoga),
            ("Karana", summary.panchaAnga.karana),
            ("Vara", summary.panchaAnga.vara),
            ("Sunrise", summary.day.sunriseTime),
            ("Sunset", summary.day.sunsetTime),
            ("Abhijit", abhijitText(for: summary)),
            ("Starts", timeFormatter.string(from: summary.tithi.startTime)),
            ("Ends", timeFormatter.string(from: summary.tithi.endTime)),
            ("Progress", progressText(for: summary, now: now)),
            ("Remaining", remainingText(for: summary, now: now)),
            ("Generated", timeFormatter.string(from: summary.generatedAt)),
            ("Next Refresh", timeFormatter.string(from: summary.nextRefreshAt)),
            ("Region", summary.calculation.region),
            ("Timezone", summary.calculation.timezone),
            ("Calendar System", summary.calculation.calendarSystem),
            ("Method", summary.calculation.method),
            ("Locale", summary.calculation.locale)
        ]
    }

    public func abhijitText(for summary: TithiSummaryResponse) -> String {
        "\(summary.day.abhijitMuhurta.startTime)-\(summary.day.abhijitMuhurta.endTime)"
    }

    public func complicationAccessibilityLabel(
        for summary: TithiSummaryResponse,
        now: Date,
        isStale: Bool
    ) -> String {
        let statusText = isStale ? "Cached tithi" : "Current tithi"
        return [
            "\(statusText) \(inlineText(for: summary))",
            "nakshatra \(summary.panchaAnga.nakshatra)",
            "yoga \(summary.panchaAnga.yoga)",
            "karana \(summary.panchaAnga.karana)",
            "abhijit \(abhijitText(for: summary))",
            timeStatusText(for: summary, now: now)
        ].joined(separator: ", ")
    }

    private func timeStatusText(for summary: TithiSummaryResponse, now: Date) -> String {
        let text = remainingText(for: summary, now: now)
        if text == "ending now" || text.hasPrefix("ended ") {
            return text
        }
        return "remaining \(text)"
    }

    private func timeFormatter(for timezone: String) -> DateFormatter {
        let formatter = DateFormatter()
        formatter.calendar = Calendar(identifier: .gregorian)
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(identifier: timezone) ?? TimeZone(secondsFromGMT: 0)
        formatter.dateFormat = "yyyy-MM-dd HH:mm zzz"
        return formatter
    }

    private func tithiDisplayName(for summary: TithiSummaryResponse) -> String {
        let name = summary.tithi.name.trimmingCharacters(in: .whitespacesAndNewlines)
        let traditionalName = summary.tithi.traditionalName.trimmingCharacters(in: .whitespacesAndNewlines)

        guard !traditionalName.isEmpty,
              normalizedName(traditionalName) != normalizedName(name) else {
            return name
        }

        return "\(name) (\(traditionalName))"
    }

    private func normalizedNakshatraName(_ name: String) -> String {
        normalizedName(name)
    }

    private func normalizedName(_ name: String) -> String {
        name.lowercased().filter { $0.isLetter || $0.isNumber }
    }
}
