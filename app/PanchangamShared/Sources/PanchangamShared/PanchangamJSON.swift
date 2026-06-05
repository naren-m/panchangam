import Foundation

public enum PanchangamJSON {
    public static let decoder: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let value = try container.decode(String.self)
            guard let date = parseDate(value) else {
                throw DecodingError.dataCorruptedError(
                    in: container,
                    debugDescription: "Invalid RFC3339 date: \(value)"
                )
            }
            return date
        }
        return decoder
    }()

    public static let encoder: JSONEncoder = {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .custom { date, encoder in
            var container = encoder.singleValueContainer()
            try container.encode(format(date))
        }
        return encoder
    }()

    public static func format(_ date: Date) -> String {
        makeRFC3339Formatter(fractionalSeconds: false).string(from: date)
    }

    public static func parseDate(_ value: String) -> Date? {
        makeRFC3339Formatter(fractionalSeconds: false).date(from: value) ??
            makeRFC3339Formatter(fractionalSeconds: true).date(from: value)
    }

    private static func makeRFC3339Formatter(fractionalSeconds: Bool) -> ISO8601DateFormatter {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = fractionalSeconds ? [.withInternetDateTime, .withFractionalSeconds] : [.withInternetDateTime]
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        return formatter
    }
}
