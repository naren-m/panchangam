import PanchangamShared
import Foundation
import SwiftUI
import WidgetKit

private let taraSpace0 = Color(red: 0.02, green: 0.03, blue: 0.06)
private let taraStar100 = Color(red: 0.95, green: 0.96, blue: 0.99)
private let taraStar300 = Color(red: 0.77, green: 0.80, blue: 0.93)
private let taraMoon = Color(red: 0.68, green: 0.74, blue: 0.91)
private let taraGold = Color(red: 0.91, green: 0.76, blue: 0.41)

struct TithiComplicationEntryView: View {
    @Environment(\.widgetFamily) private var family

    var entry: TithiEntry
    private let formatter = TithiFormatter()

    var body: some View {
        Group {
            switch family {
            case .accessoryInline:
                inlineView
            case .accessoryCircular:
                circularView
            case .accessoryRectangular:
                rectangularView
            case .accessoryCorner:
                cornerView
            default:
                inlineView
            }
        }
        .foregroundStyle(taraStar100)
        .containerBackground(taraSpace0, for: .widget)
        .accessibilityElement(children: .ignore)
        .accessibilityLabel(accessibilityLabel)
        .widgetURL(URL(string: "panchangam://tithi/current"))
    }

    private var inlineView: some View {
        Text(primaryText)
            .foregroundStyle(taraMoon)
            .widgetAccentable()
    }

    private var circularView: some View {
        ZStack {
            mandalaComplicationRing(activeIndex: activeNakshatraIndex, progress: tithiProgress, diameter: 48, tickLength: 5)

            VStack(spacing: 2) {
                Text(circularText)
                    .font(.headline)
                    .minimumScaleFactor(0.6)
                Text(secondaryText)
                    .font(.caption2)
                    .foregroundStyle(taraMoon)
                    .minimumScaleFactor(0.6)
            }
        }
        .multilineTextAlignment(.center)
    }

    private var rectangularView: some View {
        VStack(alignment: .leading, spacing: 2) {
            ForEach(rectangularLines, id: \.self) { line in
                Text(line)
                    .font(line == rectangularLines.first ? .headline : .caption2)
                    .foregroundStyle(line == rectangularLines.first ? taraStar100 : taraStar300)
                    .lineLimit(1)
                    .minimumScaleFactor(0.7)
            }
        }
    }

    private var cornerView: some View {
        ZStack {
            mandalaComplicationRing(activeIndex: activeNakshatraIndex, progress: tithiProgress, diameter: 34, tickLength: 4)

            Text(circularText)
                .font(.caption)
                .foregroundStyle(taraGold)
                .widgetAccentable()
        }
    }

    private var activeNakshatraIndex: Int? {
        guard let summary = entry.state.summary else {
            return nil
        }
        return formatter.nakshatraIndex(for: summary)
    }

    private var tithiProgress: Double {
        guard let summary = entry.state.summary else {
            return 0
        }
        return formatter.tithiProgress(for: summary, now: entry.date)
    }

    private func mandalaComplicationRing(activeIndex: Int?, progress: Double, diameter: CGFloat, tickLength: CGFloat) -> some View {
        ZStack {
            Circle()
                .stroke(taraMoon.opacity(0.18), lineWidth: 1)
                .frame(width: diameter, height: diameter)

            Circle()
                .trim(from: 0, to: progress)
                .stroke(taraGold, style: StrokeStyle(lineWidth: 2, lineCap: .round))
                .rotationEffect(.degrees(-90))
                .frame(width: diameter, height: diameter)

            ForEach(0..<27, id: \.self) { index in
                let isMajor = index % 3 == 0
                let isActive = activeIndex == index

                Capsule()
                    .fill(isActive ? taraStar100 : taraMoon.opacity(isMajor ? 0.34 : 0.18))
                    .frame(width: isActive ? 2 : 1, height: isMajor ? tickLength : max(2, tickLength - 2))
                    .offset(y: -(diameter / 2))
                    .rotationEffect(.degrees(Double(index) * (360.0 / 27.0)))
            }
        }
        .frame(width: diameter + 10, height: diameter + 10)
    }

    private var primaryText: String {
        guard let summary = entry.state.summary else {
            return stateText
        }
        return formatter.complicationInlineText(for: summary, isStale: entry.state.isStale)
    }

    private var circularText: String {
        guard let summary = entry.state.summary else {
            return stateText
        }
        return formatter.circularText(for: summary)
    }

    private var secondaryText: String {
        guard let summary = entry.state.summary else {
            return ""
        }
        return formatter.complicationSecondaryText(for: summary, isStale: entry.state.isStale)
    }

    private var rectangularLines: [String] {
        guard let summary = entry.state.summary else {
            return [stateText]
        }
        return formatter.rectangularLines(for: summary, now: entry.date, isStale: entry.state.isStale)
    }

    private var accessibilityLabel: String {
        guard let summary = entry.state.summary else {
            return stateText
        }

        return formatter.complicationAccessibilityLabel(
            for: summary,
            now: entry.date,
            isStale: entry.state.isStale
        )
    }

    private var stateText: String {
        switch entry.state {
        case .loading:
            return "Loading"
        case .error(let message):
            return message
        case .valid, .stale:
            return primaryText
        }
    }
}

#Preview("Inline", as: .accessoryInline) {
    PanchangamComplicationWidget()
} timeline: {
    TithiEntry(date: Date(), state: .valid(.preview))
    TithiEntry(date: Date(), state: .stale(.preview))
    TithiEntry(date: Date(), state: .error("Unavailable"))
}

#Preview("Circular", as: .accessoryCircular) {
    PanchangamComplicationWidget()
} timeline: {
    TithiEntry(date: Date(), state: .loading)
    TithiEntry(date: Date(), state: .valid(.preview))
}

#Preview("Rectangular", as: .accessoryRectangular) {
    PanchangamComplicationWidget()
} timeline: {
    TithiEntry(date: Date(), state: .valid(.preview))
    TithiEntry(date: Date(), state: .error("Unavailable"))
}

#Preview("Corner", as: .accessoryCorner) {
    PanchangamComplicationWidget()
} timeline: {
    TithiEntry(date: Date(), state: .valid(.preview))
}

extension TithiSummaryResponse {
    static let preview = TithiSummaryResponse(
        date: "2026-06-02",
        tithi: TithiDetails(
            name: "Sashthi",
            traditionalName: "Sashthi",
            number: 6,
            paksha: "Shukla",
            pakshaDay: 6,
            type: "Nanda",
            startTime: Date(timeIntervalSince1970: 1_780_362_000),
            endTime: Date(timeIntervalSince1970: 1_780_452_000)
        ),
        panchaAnga: PanchaAngaSummary(nakshatra: "Pushya", yoga: "Saubhagya", karana: "Taitila", vara: "Tuesday"),
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
        generatedAt: Date(timeIntervalSince1970: 1_780_401_600),
        nextRefreshAt: Date(timeIntervalSince1970: 1_780_452_000)
    )
}
