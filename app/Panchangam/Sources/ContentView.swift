import Combine
import CoreLocation
import PanchangamShared
import SwiftUI

private let taraSpace0 = Color(red: 0.02, green: 0.03, blue: 0.06)
private let taraSpace1 = Color(red: 0.04, green: 0.06, blue: 0.13)
private let taraSpace2 = Color(red: 0.07, green: 0.10, blue: 0.20)
private let taraSpace3 = Color(red: 0.10, green: 0.14, blue: 0.28)
private let taraStar100 = Color(red: 0.95, green: 0.96, blue: 0.99)
private let taraStar300 = Color(red: 0.77, green: 0.80, blue: 0.93)
private let taraStar500 = Color(red: 0.55, green: 0.58, blue: 0.74)
private let taraMoon = Color(red: 0.68, green: 0.74, blue: 0.91)
private let taraGold = Color(red: 0.91, green: 0.76, blue: 0.41)
private let taraGood = Color(red: 0.45, green: 0.82, blue: 0.66)
private let taraBad = Color(red: 0.88, green: 0.44, blue: 0.54)
private let phoneStars: [(x: CGFloat, y: CGFloat, size: CGFloat, opacity: Double)] = [
    (0.12, 0.08, 2, 0.42),
    (0.28, 0.18, 1.4, 0.28),
    (0.58, 0.10, 1.8, 0.34),
    (0.84, 0.20, 1.3, 0.26),
    (0.16, 0.44, 1.6, 0.36),
    (0.78, 0.42, 2, 0.38),
    (0.34, 0.68, 1.5, 0.24),
    (0.66, 0.76, 1.6, 0.32),
    (0.88, 0.64, 1.2, 0.24)
]

struct ContentView: View {
    @ObservedObject var state: PanchangamAppState
    @ObservedObject var locationService: LocationService
    @ObservedObject var watchSync: WatchSettingsSync

    @State private var didRunStartupTask = false

    private let formatter = TithiFormatter()

    var body: some View {
        NavigationStack {
            ZStack {
                taraBackground

                ScrollView {
                    VStack(alignment: .leading, spacing: 16) {
                        heroSection
                        currentSection
                        locationSection
                        backendSection
                        calculationSection
                    }
                    .padding(18)
                }
            }
            .navigationTitle("Mandala")
            .toolbarColorScheme(.dark, for: .navigationBar)
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button {
                        Task { await state.refresh(sync: watchSync) }
                    } label: {
                        Label("Refresh", systemImage: "arrow.clockwise")
                    }
                    .tint(taraMoon)
                    .disabled(state.status == .loading)
                }
            }
            .task {
                guard !didRunStartupTask else {
                    return
                }

                didRunStartupTask = true
                state.loadCachedSummary()
                state.syncCachedToWatch(sync: watchSync)
                if state.hasSavedSettings {
                    await state.refresh(sync: watchSync)
                }
            }
            .onReceive(locationService.$lastLocation.compactMap { $0 }) { location in
                applyAndRefresh(location: location)
            }
        }
    }

    private var taraBackground: some View {
        ZStack {
            taraSpace0.ignoresSafeArea()
            LinearGradient(
                colors: [Color(red: 0.10, green: 0.14, blue: 0.31), taraSpace1, taraSpace0],
                startPoint: .top,
                endPoint: .bottom
            )
            .ignoresSafeArea()
            RadialGradient(
                colors: [taraMoon.opacity(0.24), .clear],
                center: .top,
                startRadius: 8,
                endRadius: 420
            )
            .ignoresSafeArea()
            phoneStarField
                .ignoresSafeArea()
        }
    }

    private var phoneStarField: some View {
        GeometryReader { proxy in
            ForEach(Array(phoneStars.enumerated()), id: \.offset) { _, star in
                Circle()
                    .fill(taraMoon.opacity(star.opacity))
                    .frame(width: star.size, height: star.size)
                    .position(x: proxy.size.width * star.x, y: proxy.size.height * star.y)
            }
        }
        .allowsHitTesting(false)
    }

    private func applyAndRefresh(location: CLLocation) {
        state.apply(location: location)
        Task { await state.refresh(sync: watchSync) }
    }

    private var heroSection: some View {
        VStack(alignment: .leading, spacing: 14) {
            HStack(alignment: .center) {
                sectionLabel("Mandala")
                Spacer()
                HStack(spacing: 6) {
                    Image(systemName: "moonphase.first.quarter")
                    Text(state.status.text)
                }
                .font(.footnote.weight(.semibold))
                .foregroundStyle(statusColor)
                .lineLimit(1)
                .minimumScaleFactor(0.75)
            }

            ZStack {
                Circle()
                    .stroke(taraMoon.opacity(0.14), lineWidth: 1)
                    .frame(width: 236, height: 236)

                ForEach(0..<27, id: \.self) { index in
                    phoneMandalaTick(index: index, activeIndex: state.summary.flatMap { formatter.nakshatraIndex(for: $0) })
                }

                VStack(spacing: 7) {
                    ZStack {
                        Circle()
                            .fill(taraMoon)
                        Circle()
                            .fill(taraSpace0)
                            .offset(x: -6)
                    }
                    .frame(width: 24, height: 24)
                    .clipShape(Circle())
                    .overlay(Circle().stroke(taraMoon.opacity(0.28), lineWidth: 1))

                    TimelineView(.periodic(from: Date(), by: 60)) { context in
                        Text(context.date, style: .time)
                            .font(.system(size: 52, weight: .regular, design: .rounded))
                            .foregroundStyle(taraStar100)
                            .lineLimit(1)
                            .minimumScaleFactor(0.65)
                    }

                    if let summary = state.summary {
                        Text(formatter.inlineText(for: summary))
                            .font(.headline.weight(.bold))
                            .foregroundStyle(taraStar100)
                            .lineLimit(1)
                            .minimumScaleFactor(0.7)

                        Text("Nakshatra - \(summary.panchaAnga.nakshatra)")
                            .font(.caption2.weight(.bold))
                            .tracking(1.1)
                            .textCase(.uppercase)
                            .foregroundStyle(taraStar500)
                            .lineLimit(1)
                            .minimumScaleFactor(0.65)
                    } else {
                        Text("No tithi loaded")
                            .font(.headline.weight(.semibold))
                            .foregroundStyle(taraStar100)
                    }
                }
            }

            if let summary = state.summary {
                HStack(alignment: .bottom) {
                    mandalaCorner("Yoga", summary.panchaAnga.yoga)
                    Spacer()
                    mandalaCorner("Karana", summary.panchaAnga.karana)
                        .multilineTextAlignment(.trailing)
                }
            }

            HStack(spacing: 6) {
                Circle()
                    .fill(taraGood)
                    .frame(width: 8, height: 8)
                Text(state.summary?.day.abhijitMuhurta.name ?? "Watch Sync")
                    .font(.caption.weight(.bold))
                Text(state.summary.map { formatter.abhijitText(for: $0) } ?? "Waiting")
                    .font(.caption.weight(.semibold))
                Spacer()
                Text(watchSync.lastSyncText ?? "Watch Sync")
                    .font(.caption2.weight(.semibold))
                    .foregroundStyle(taraMoon)
                    .lineLimit(1)
                    .minimumScaleFactor(0.7)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 8)
            .background(taraGood.opacity(0.18), in: Capsule())
            .foregroundStyle(taraStar100)
        }
        .padding(18)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(taraSpace2.opacity(0.82), in: RoundedRectangle(cornerRadius: 24))
        .overlay(
            RoundedRectangle(cornerRadius: 24)
                .stroke(taraMoon.opacity(0.18), lineWidth: 1)
        )
    }

    private var currentSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            sectionLabel("Current Tithi")

            if let summary = state.summary {
                TimelineView(.periodic(from: Date(), by: 60)) { context in
                    VStack(spacing: 0) {
                        ForEach(formatter.detailRows(for: summary, now: context.date), id: \.0) { label, value in
                            TaraRow(label: label, value: value)
                        }
                    }
                }
            } else {
                Text("No tithi loaded")
                    .foregroundStyle(taraStar500)
            }
        }
        .taraPanel()
    }

    private var locationSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            sectionLabel("Location")

            Toggle("Manual Location", isOn: $state.usesManualLocation)
                .tint(taraMoon)
                .foregroundStyle(taraStar300)

            HStack(spacing: 8) {
                styledTextField("Latitude", text: $state.latitudeText)
                    .keyboardType(.numbersAndPunctuation)
                styledTextField("Longitude", text: $state.longitudeText)
                    .keyboardType(.numbersAndPunctuation)
            }
            .disabled(!state.usesManualLocation)

            styledTextField("Timezone", text: $state.timezoneText)
                .textInputAutocapitalization(.never)
                .autocorrectionDisabled()

            Button {
                locationService.requestLocation()
            } label: {
                Label("Use Current Location", systemImage: "location")
                    .frame(maxWidth: .infinity)
            }
            .buttonStyle(.borderedProminent)
            .tint(taraMoon)
            .foregroundStyle(taraSpace1)

            if let authorizationText = locationService.authorizationText {
                Text(authorizationText)
                    .font(.footnote)
                    .foregroundStyle(taraStar500)
            }
        }
        .taraPanel()
    }

    private var backendSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            sectionLabel("Backend")

            styledTextField("API Base URL", text: $state.apiBaseURLText)
                .textInputAutocapitalization(.never)
                .autocorrectionDisabled()
                .keyboardType(.URL)
            if let developmentHostWarning = state.developmentHostWarning {
                TaraRow(label: "Backend Host", value: developmentHostWarning, valueColor: taraGold)
            }
            if let storageStatusText = state.storageStatusText {
                TaraRow(label: "Storage", value: storageStatusText, valueColor: taraGold)
            }
            if let syncText = watchSync.lastSyncText {
                TaraRow(label: "Watch Sync", value: syncText, valueColor: taraMoon)
            }
        }
        .taraPanel()
    }

    private var calculationSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            sectionLabel("Calculation")

            styledTextField("Region", text: $state.regionText)
            styledTextField("Method", text: $state.methodText)
            styledTextField("Locale", text: $state.localeText)
                .textInputAutocapitalization(.never)

            VStack(alignment: .leading, spacing: 8) {
                Text("Calendar System")
                    .font(.caption.weight(.semibold))
                    .foregroundStyle(taraStar500)

                Picker("Calendar System", selection: $state.calendarSystemText) {
                    Text("Auto").tag("")
                    Text("Purnimanta").tag("Purnimanta")
                    Text("Amanta").tag("Amanta")
                }
                .pickerStyle(.segmented)
                .tint(taraMoon)
            }
        }
        .taraPanel()
    }

    private var statusColor: Color {
        switch state.status {
        case .loaded:
            return taraGood
        case .loading:
            return taraGold
        case .stale:
            return taraMoon
        case .failed:
            return taraBad
        case .idle:
            return taraStar300
        }
    }

    private func sectionLabel(_ text: String) -> some View {
        Text(text)
            .font(.caption2.weight(.bold))
            .tracking(1.5)
            .textCase(.uppercase)
            .foregroundStyle(taraStar500)
    }

    private func mandalaCorner(_ label: String, _ value: String) -> some View {
        VStack(alignment: label == "Karana" ? .trailing : .leading, spacing: 3) {
            Text(label)
                .font(.caption2.weight(.bold))
                .tracking(1.0)
                .textCase(.uppercase)
                .foregroundStyle(taraStar500)
            Text(value)
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(taraStar100)
                .lineLimit(1)
                .minimumScaleFactor(0.75)
        }
    }

    private func phoneMandalaTick(index: Int, activeIndex: Int?) -> some View {
        let isMajor = index % 3 == 0
        let isActive = activeIndex == index

        return Capsule()
            .fill(isActive ? taraStar100 : taraMoon.opacity(isMajor ? 0.34 : 0.22))
            .frame(width: isActive ? 2 : 1, height: isMajor ? 13 : 8)
            .offset(y: -118)
            .rotationEffect(.degrees(Double(index) * (360.0 / 27.0)))
    }

    private func styledTextField(_ placeholder: String, text: Binding<String>) -> some View {
        TextField(placeholder, text: text)
            .textFieldStyle(.plain)
            .foregroundStyle(taraStar100)
            .padding(11)
            .background(taraSpace3.opacity(0.74), in: RoundedRectangle(cornerRadius: 12))
            .overlay(
                RoundedRectangle(cornerRadius: 12)
                    .stroke(taraMoon.opacity(0.18), lineWidth: 1)
            )
    }
}

private struct TaraRow: View {
    var label: String
    var value: String
    var valueColor: Color = taraStar100

    var body: some View {
        HStack(alignment: .firstTextBaseline, spacing: 16) {
            Text(label)
                .font(.caption.weight(.semibold))
                .foregroundStyle(taraStar500)
            Spacer(minLength: 8)
            Text(value)
                .font(.subheadline.weight(.semibold))
                .foregroundStyle(valueColor)
                .multilineTextAlignment(.trailing)
        }
        .padding(.vertical, 9)
        .overlay(alignment: .bottom) {
            Rectangle()
                .fill(taraMoon.opacity(0.10))
                .frame(height: 1)
        }
    }
}

private extension View {
    func taraPanel() -> some View {
        self
            .padding(16)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(taraSpace2.opacity(0.80), in: RoundedRectangle(cornerRadius: 22))
            .overlay(
                RoundedRectangle(cornerRadius: 22)
                    .stroke(taraMoon.opacity(0.14), lineWidth: 1)
            )
    }
}

#Preview("Loaded") {
    ContentView(
        state: PanchangamAppState.previewLoaded,
        locationService: LocationService(),
        watchSync: WatchSettingsSync()
    )
}

#Preview("Empty") {
    ContentView(
        state: PanchangamAppState(),
        locationService: LocationService(),
        watchSync: WatchSettingsSync()
    )
}

private extension PanchangamAppState {
    static var previewLoaded: PanchangamAppState {
        let state = PanchangamAppState()
        state.summary = .preview
        state.status = .loaded(Date(timeIntervalSince1970: 1_780_401_600))
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
            panchaAnga: PanchaAngaSummary(
                nakshatra: "Pushya",
                yoga: "Saubhagya",
                karana: "Taitila",
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
        generatedAt: Date(timeIntervalSince1970: 1_780_401_600),
        nextRefreshAt: Date(timeIntervalSince1970: 1_780_452_000)
    )
}
