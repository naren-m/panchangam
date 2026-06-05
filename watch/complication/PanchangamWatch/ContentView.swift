import SwiftUI

struct ContentView: View {
    @State private var day = PanchangamDay.sample
    @State private var isLoading = false
    @State private var statusText = "Demo"

    private let client = PanchangamAPIClient.liveFromBundle()

    var body: some View {
        GeometryReader { geometry in
            let width = geometry.size.width
            let height = geometry.size.height
            let horizontalPadding = max(width * 0.085, 16)
            let contentWidth = width - (horizontalPadding * 2)

            ZStack {
                FaceBackground()
                StarField()

                VStack(alignment: .leading, spacing: 4) {
                    Text(day.dateHeading)
                        .font(.system(size: 10, weight: .bold, design: .rounded))
                        .tracking(1.4)
                        .foregroundStyle(.white.opacity(0.72))

                    Text(day.weekdayText)
                        .font(.system(size: 12, weight: .semibold, design: .rounded))
                        .foregroundStyle(.white.opacity(0.58))
                }
                .frame(width: contentWidth, alignment: .leading)
                .position(x: width / 2, y: height * 0.105)

                VStack(alignment: .trailing, spacing: 7) {
                    MoonMark(size: 28)

                    Text(day.pakshaText)
                        .font(.system(size: 11, weight: .semibold, design: .rounded))
                        .foregroundStyle(.white.opacity(0.58))
                        .lineLimit(1)
                        .minimumScaleFactor(0.72)
                }
                .frame(width: contentWidth, alignment: .trailing)
                .position(x: width / 2, y: height * 0.118)

                TimelineView(.periodic(from: Date(), by: 60)) { context in
                    Text(Self.timeFormatter.string(from: context.date))
                        .font(.system(size: min(max(width * 0.215, 42), 50), weight: .bold, design: .rounded))
                        .monospacedDigit()
                        .foregroundStyle(.white.opacity(0.96))
                        .lineLimit(1)
                        .minimumScaleFactor(0.76)
                }
                .frame(width: contentWidth, alignment: .leading)
                .position(x: width / 2, y: height * 0.285)

                HStack(alignment: .firstTextBaseline, spacing: 7) {
                    Text(day.tithiFaceTitle)
                        .font(.system(size: min(max(width * 0.078, 19), 24), weight: .bold, design: .rounded))
                        .foregroundStyle(.white.opacity(0.96))
                        .lineLimit(1)
                        .minimumScaleFactor(0.62)

                    if !day.tithiLocalName.isEmpty {
                        Text(day.tithiLocalName)
                            .font(.system(size: 12, weight: .bold, design: .rounded))
                            .foregroundStyle(.white.opacity(0.70))
                            .lineLimit(1)
                    }
                }
                .frame(width: contentWidth, alignment: .leading)
                .position(x: width / 2, y: height * 0.420)

                Text(day.tithiFaceDetailText)
                    .font(.system(size: 10, weight: .medium, design: .rounded))
                    .foregroundStyle(.white.opacity(0.57))
                    .lineLimit(1)
                    .minimumScaleFactor(0.58)
                    .frame(width: contentWidth, alignment: .leading)
                    .position(x: width / 2, y: height * 0.475)

                PanchangamMetricRow(
                    title: "YOGA",
                    value: day.yoga,
                    localValue: "",
                    trailing: ""
                )
                .frame(width: contentWidth)
                .position(x: width / 2, y: height * 0.585)

                TimelineView(.periodic(from: Date(), by: 60)) { context in
                    GoodPeriodCard(period: day.nextGoodPeriod(at: context.date))
                }
                .frame(width: contentWidth)
                .position(x: width / 2, y: height * 0.755)

                FooterTimes(sunlight: day.sunlightText, rahu: day.rahuStartText)
                    .frame(width: contentWidth)
                    .position(x: width / 2, y: height * 0.890)
            }
        }
        .ignoresSafeArea()
        .persistentSystemOverlays(.hidden)
        ._statusBarHidden()
        .task {
            await refresh()
        }
    }

    private func refresh() async {
        isLoading = true
        defer { isLoading = false }

        do {
            day = try await client.fetchToday()
            statusText = "Live"
        } catch {
            day = .sample
            statusText = "Offline"
        }
    }

    private static let timeFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "H:mm"
        return formatter
    }()
}

private struct FaceBackground: View {
    var body: some View {
        LinearGradient(
            colors: [
                Color(red: 0.015, green: 0.018, blue: 0.040),
                Color(red: 0.045, green: 0.060, blue: 0.150),
                Color(red: 0.070, green: 0.090, blue: 0.205)
            ],
            startPoint: .bottom,
            endPoint: .top
        )
        .ignoresSafeArea()
    }
}

private struct PanchangamMetricRow: View {
    let title: String
    let value: String
    let localValue: String
    let trailing: String

    var body: some View {
        HStack(alignment: .bottom, spacing: 10) {
            VStack(alignment: .leading, spacing: 5) {
                Text(title)
                    .font(.system(size: 9, weight: .bold, design: .rounded))
                    .tracking(1.5)
                    .foregroundStyle(.white.opacity(0.60))

                HStack(alignment: .firstTextBaseline, spacing: 5) {
                    Text(value)
                        .font(.system(size: 15, weight: .bold, design: .rounded))
                        .foregroundStyle(.white.opacity(0.96))
                        .lineLimit(1)
                        .minimumScaleFactor(0.66)

                    if !localValue.isEmpty {
                        Text(localValue)
                            .font(.system(size: 11, weight: .bold, design: .rounded))
                            .foregroundStyle(.white.opacity(0.66))
                            .lineLimit(1)
                    }
                }
            }

            Spacer(minLength: 8)

            if !trailing.isEmpty {
                Text(trailing)
                    .font(.system(size: 14, weight: .semibold, design: .rounded))
                    .monospacedDigit()
                    .foregroundStyle(.white.opacity(0.64))
                    .lineLimit(1)
            }
        }
    }
}

private struct GoodPeriodCard: View {
    let period: GoodTimePeriod

    var body: some View {
        HStack(alignment: .center, spacing: 9) {
            Circle()
                .fill(Color(red: 0.45, green: 0.90, blue: 0.74))
                .frame(width: 8, height: 8)

            VStack(alignment: .leading, spacing: 3) {
                Text(period.windowText)
                    .font(.system(size: 9, weight: .bold, design: .rounded))
                    .monospacedDigit()
                    .foregroundStyle(.white.opacity(0.78))
                    .lineLimit(1)
                    .minimumScaleFactor(0.62)

                Text(period.name)
                    .font(.system(size: 12, weight: .bold, design: .rounded))
                    .foregroundStyle(.white.opacity(0.94))
                    .lineLimit(1)
                    .minimumScaleFactor(0.62)
                    .frame(maxWidth: .infinity, alignment: .leading)
            }
        }
        .padding(.horizontal, 12)
        .padding(.vertical, 10)
        .background(Color(red: 0.10, green: 0.27, blue: 0.25).opacity(0.86))
        .clipShape(RoundedRectangle(cornerRadius: 14, style: .continuous))
    }
}

private struct FooterTimes: View {
    let sunlight: String
    let rahu: String

    var body: some View {
        HStack(alignment: .center) {
            HStack(spacing: 4) {
                Image(systemName: "sunrise.fill")
                    .font(.system(size: 10, weight: .semibold))
                    .foregroundStyle(Color(red: 0.98, green: 0.76, blue: 0.30))

                Text(sunlight)
                    .font(.system(size: 11, weight: .semibold, design: .rounded))
                    .monospacedDigit()
                    .foregroundStyle(.white.opacity(0.66))
                    .lineLimit(1)
            }

            Spacer(minLength: 8)

            Text(rahu)
                .font(.system(size: 11, weight: .bold, design: .rounded))
                .monospacedDigit()
                .foregroundStyle(Color(red: 1.0, green: 0.46, blue: 0.62))
                .lineLimit(1)
        }
    }
}

private struct MoonMark: View {
    let size: CGFloat

    var body: some View {
        ZStack(alignment: .trailing) {
            Circle()
                .fill(Color(red: 0.86, green: 0.89, blue: 1.0))
                .frame(width: size, height: size)

            Circle()
                .fill(Color(red: 0.030, green: 0.035, blue: 0.080))
                .frame(width: size * 0.82, height: size)
                .offset(x: -size * 0.24)
        }
        .shadow(color: Color(red: 0.64, green: 0.72, blue: 1.0).opacity(0.45), radius: 8)
    }
}

private struct StarField: View {
    private let points = [
        CGPoint(x: 0.15, y: 0.17),
        CGPoint(x: 0.30, y: 0.25),
        CGPoint(x: 0.50, y: 0.21),
        CGPoint(x: 0.78, y: 0.31),
        CGPoint(x: 0.92, y: 0.43),
        CGPoint(x: 0.17, y: 0.58),
        CGPoint(x: 0.44, y: 0.52),
        CGPoint(x: 0.63, y: 0.61),
        CGPoint(x: 0.83, y: 0.70),
        CGPoint(x: 0.25, y: 0.78),
        CGPoint(x: 0.56, y: 0.84),
        CGPoint(x: 0.72, y: 0.91)
    ]

    var body: some View {
        GeometryReader { geometry in
            ForEach(points.indices, id: \.self) { index in
                Circle()
                    .fill(.white.opacity(index.isMultiple(of: 3) ? 0.34 : 0.18))
                    .frame(width: index.isMultiple(of: 4) ? 2 : 1.3, height: index.isMultiple(of: 4) ? 2 : 1.3)
                    .position(
                        x: points[index].x * geometry.size.width,
                        y: points[index].y * geometry.size.height
                    )
            }
        }
        .ignoresSafeArea()
    }
}

#Preview {
    ContentView()
}
