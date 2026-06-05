import Foundation
import PanchangamShared
import SwiftUI

private let taraSpace0 = Color(red: 0.02, green: 0.03, blue: 0.06)
private let taraSpace1 = Color(red: 0.04, green: 0.06, blue: 0.13)
private let taraSpace2 = Color(red: 0.07, green: 0.10, blue: 0.20)
private let taraStar100 = Color(red: 0.95, green: 0.96, blue: 0.99)
private let taraStar300 = Color(red: 0.77, green: 0.80, blue: 0.93)
private let taraStar500 = Color(red: 0.55, green: 0.58, blue: 0.74)
private let taraMoon = Color(red: 0.68, green: 0.74, blue: 0.91)
private let taraGold = Color(red: 0.91, green: 0.76, blue: 0.41)
private let taraGood = Color(red: 0.45, green: 0.82, blue: 0.66)
private let watchStars: [(x: CGFloat, y: CGFloat, size: CGFloat, opacity: Double)] = [
    (0.16, 0.12, 1.6, 0.40),
    (0.32, 0.22, 1.1, 0.26),
    (0.58, 0.11, 1.4, 0.32),
    (0.78, 0.24, 1.2, 0.24),
    (0.18, 0.48, 1.3, 0.34),
    (0.82, 0.48, 1.5, 0.36),
    (0.28, 0.74, 1.1, 0.24),
    (0.66, 0.78, 1.3, 0.30)
]

private let taraWash = LinearGradient(
    colors: [Color(red: 0.10, green: 0.14, blue: 0.31), taraSpace1, taraSpace0],
    startPoint: .top,
    endPoint: .bottom
)

struct WatchContentView: View {
    @ObservedObject var state: WatchAppState
    @ObservedObject var receiver: WatchSettingsReceiver

    @State private var selectedSummary: SelectedTithiSummary?
    @State private var opensDetailsWhenSummaryArrives = false

    private let formatter = TithiFormatter()

    var body: some View {
        ZStack {
            taraWash
                .ignoresSafeArea()
            watchStarField
                .ignoresSafeArea()

            if let summary = state.summary {
                mandalaFace(summary)
            } else {
                emptyMandala
            }
        }
        .sheet(item: $selectedSummary) { selected in
            WatchTithiDetailView(summary: selected.summary, statusText: selected.statusText)
        }
        .onChange(of: state.summary) { _, summary in
            guard opensDetailsWhenSummaryArrives, let summary else {
                return
            }

            opensDetailsWhenSummaryArrives = false
            openDetails(for: summary)
        }
        .onOpenURL { url in
            openTithiDetails(from: url)
        }
        ._statusBarHidden(true)
        .persistentSystemOverlays(.hidden)
    }

    private func mandalaFace(_ summary: TithiSummaryResponse) -> some View {
        let activeNakshatraIndex = formatter.nakshatraIndex(for: summary)

        return ZStack {
            Group {
                Circle()
                    .stroke(taraMoon.opacity(0.14), lineWidth: 1)
                    .frame(width: 154, height: 154)

                ForEach(0..<27, id: \.self) { index in
                    mandalaTick(index: index, activeIndex: activeNakshatraIndex)
                }
            }
            .offset(y: -18)

            VStack(spacing: 0) {
                header(summary)
                    .offset(y: -5)
                Spacer(minLength: 8)
                center(summary)
                    .offset(y: -10)
                Spacer(minLength: 10)
                bottomDetails(summary)
                abhijitPill(summary)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 5)
        }
        .contentShape(Rectangle())
        .onTapGesture {
            openDetails(for: summary)
        }
        .accessibilityElement(children: .combine)
        .accessibilityLabel("Open Mandala tithi details")
        .accessibilityHint("Shows current tithi details")
        .accessibilityAddTraits(.isButton)
        .accessibilityAction {
            openDetails(for: summary)
        }
    }

    private var watchStarField: some View {
        GeometryReader { proxy in
            ForEach(Array(watchStars.enumerated()), id: \.offset) { _, star in
                Circle()
                    .fill(taraMoon.opacity(star.opacity))
                    .frame(width: star.size, height: star.size)
                    .position(x: proxy.size.width * star.x, y: proxy.size.height * star.y)
            }
        }
        .allowsHitTesting(false)
    }

    private func header(_ summary: TithiSummaryResponse?) -> some View {
        HStack(alignment: .top) {
            cornerBlock("Sunrise", shortTime(summary?.day.sunriseTime ?? "--:--"), tint: taraGold)
            Spacer(minLength: 10)
            cornerBlock("Sunset", shortTime(summary?.day.sunsetTime ?? "--:--"), tint: taraMoon)
        }
    }

    private func shortTime(_ value: String) -> String {
        let pieces = value.split(separator: ":")
        guard pieces.count >= 2 else {
            return value
        }

        return "\(pieces[0]):\(pieces[1])"
    }

    private func abhijitText(_ summary: TithiSummaryResponse) -> String {
        "\(shortTime(summary.day.abhijitMuhurta.startTime))-\(shortTime(summary.day.abhijitMuhurta.endTime))"
    }

    private func center(_ summary: TithiSummaryResponse) -> some View {
        VStack(spacing: 5) {
            moonPhase

            TimelineView(.periodic(from: Date(), by: 60)) { context in
                Text(context.date, style: .time)
                    .font(.system(size: 40, weight: .regular, design: .rounded))
                    .foregroundStyle(taraStar100)
                    .lineLimit(1)
                    .minimumScaleFactor(0.62)
            }

            Text(formatter.inlineText(for: summary))
                .font(.caption.weight(.bold))
                .foregroundStyle(taraStar100)
                .lineLimit(1)
                .minimumScaleFactor(0.55)
                .frame(maxWidth: 132)

            Text("Nakshatra - \(summary.panchaAnga.nakshatra)")
                .font(.system(size: 9, weight: .bold))
                .tracking(1.1)
                .textCase(.uppercase)
                .foregroundStyle(taraStar500)
                .lineLimit(1)
                .minimumScaleFactor(0.55)
                .frame(maxWidth: 154)

            TimelineView(.periodic(from: Date(), by: 60)) { context in
                Text("\(summary.tithi.paksha) Paksha - \(formatter.remainingText(for: summary, now: context.date))")
                    .font(.system(size: 8, weight: .semibold))
                    .foregroundStyle(taraStar500)
                    .lineLimit(1)
                    .minimumScaleFactor(0.6)
            }
        }
        .frame(maxWidth: .infinity)
    }

    private func bottomDetails(_ summary: TithiSummaryResponse) -> some View {
        HStack(alignment: .bottom) {
            cornerBlock("Yoga", summary.panchaAnga.yoga, tint: taraMoon)
            Spacer(minLength: 10)
            cornerBlock("Karana", summary.panchaAnga.karana, tint: taraMoon)
                .multilineTextAlignment(.trailing)
        }
    }

    private func abhijitPill(_ summary: TithiSummaryResponse) -> some View {
        HStack(spacing: 5) {
            Circle()
                .fill(taraGood)
                .frame(width: 6, height: 6)
            Text(summary.day.abhijitMuhurta.name)
                .font(.system(size: 9, weight: .bold))
            Text(abhijitText(summary))
                .font(.system(size: 9, weight: .semibold))
        }
        .padding(.horizontal, 9)
        .padding(.vertical, 5)
        .background(taraGood.opacity(0.22), in: Capsule())
        .foregroundStyle(taraStar100)
        .accessibilityLabel("\(summary.day.abhijitMuhurta.name) \(abhijitText(summary))")
        .padding(.top, 7)
    }

    private var emptyMandala: some View {
        VStack(spacing: 10) {
            header(nil)
            Spacer()
            moonPhase
            Text("No tithi loaded")
                .font(.headline.weight(.semibold))
                .foregroundStyle(taraStar100)
                .multilineTextAlignment(.center)
            Text(state.status.text)
                .font(.caption2)
                .foregroundStyle(taraStar300)
            Text(receiver.statusText)
                .font(.caption2)
                .foregroundStyle(taraMoon)
            Button {
                Task { await state.refresh() }
            } label: {
                Label("Refresh", systemImage: "arrow.clockwise")
            }
            .buttonStyle(.borderedProminent)
            .tint(taraMoon)
            .foregroundStyle(taraSpace1)
            .disabled(state.status == .loading)
            Spacer()
        }
        .padding(12)
    }

    private var moonPhase: some View {
        ZStack {
            Circle()
                .fill(taraMoon)
            Circle()
                .fill(taraSpace0)
                .offset(x: -5)
        }
        .frame(width: 21, height: 21)
        .clipShape(Circle())
        .overlay(Circle().stroke(taraMoon.opacity(0.28), lineWidth: 1))
    }

    private func cornerBlock(_ label: String, _ value: String, tint: Color) -> some View {
        VStack(alignment: label == "Sunset" || label == "Karana" ? .trailing : .leading, spacing: 2) {
            Text(label)
                .font(.system(size: 8, weight: .bold))
                .tracking(0.8)
                .textCase(.uppercase)
                .foregroundStyle(taraStar500)
            Text(value)
                .font(.system(size: 12, weight: .semibold))
                .foregroundStyle(label == "Sunrise" ? tint : taraStar100)
                .lineLimit(1)
                .minimumScaleFactor(0.65)
        }
    }

    private func mandalaTick(index: Int, activeIndex: Int?) -> some View {
        let isMajor = index % 3 == 0
        let isActive = activeIndex == index

        return Capsule()
            .fill(isActive ? taraStar100 : taraMoon.opacity(isMajor ? 0.36 : 0.22))
            .frame(width: isActive ? 2 : 1, height: isMajor ? 11 : 7)
            .offset(y: -76)
            .rotationEffect(.degrees(Double(index) * (360.0 / 27.0)))
    }

    private func openTithiDetails(from url: URL) {
        guard url.scheme == "panchangam",
              url.host == "tithi",
              url.path == "/current" else {
            return
        }

        if let summary = state.summary {
            openDetails(for: summary)
            return
        }

        opensDetailsWhenSummaryArrives = true
        state.loadCachedSummary()
        if let summary = state.summary {
            opensDetailsWhenSummaryArrives = false
            openDetails(for: summary)
            return
        }

        Task { await state.refresh() }
    }

    private func openDetails(for summary: TithiSummaryResponse) {
        selectedSummary = SelectedTithiSummary(summary: summary, statusText: state.status.text)
    }
}

private struct SelectedTithiSummary: Identifiable {
    let summary: TithiSummaryResponse
    let statusText: String

    var id: String {
        "\(summary.date)-\(summary.tithi.number)-\(summary.generatedAt.timeIntervalSince1970)"
    }
}

private struct WatchTithiDetailView: View {
    @Environment(\.dismiss) private var dismiss

    let summary: TithiSummaryResponse
    let statusText: String
    private let formatter = TithiFormatter()

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 10) {
                Text("Current Tithi")
                    .font(.headline.weight(.bold))
                    .foregroundStyle(taraStar100)

                Text(statusText)
                    .font(.caption2.weight(.semibold))
                    .foregroundStyle(taraMoon)

                TimelineView(.periodic(from: Date(), by: 60)) { context in
                    ForEach(formatter.detailRows(for: summary, now: context.date), id: \.0) { label, value in
                        VStack(alignment: .leading, spacing: 2) {
                            Text(label)
                                .font(.system(size: 9, weight: .bold))
                                .textCase(.uppercase)
                                .foregroundStyle(taraStar500)
                            Text(value)
                                .font(.system(size: 12, weight: .semibold))
                                .foregroundStyle(taraStar100)
                                .lineLimit(2)
                                .minimumScaleFactor(0.7)
                        }
                    }
                }

                Button("Done") {
                    dismiss()
                }
                .buttonStyle(.bordered)
                .tint(taraMoon)
            }
            .padding(12)
        }
        .background(taraWash)
    }
}

#Preview("Loaded") {
    WatchContentView(state: WatchAppState.previewLoaded, receiver: WatchSettingsReceiver())
}

#Preview("Error") {
    WatchContentView(state: WatchAppState(), receiver: WatchSettingsReceiver())
}

#Preview("Waiting for iPhone") {
    WatchContentView(state: WatchAppState.previewWaitingForSettings, receiver: WatchSettingsReceiver())
}

private extension WatchAppState {
    static var previewLoaded: WatchAppState {
        let state = WatchAppState()
        state.summary = .preview
        state.status = .loaded
        return state
    }

    static var previewWaitingForSettings: WatchAppState {
        let state = WatchAppState()
        state.status = .waitingForSettings
        return state
    }
}

private extension TithiSummaryResponse {
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
