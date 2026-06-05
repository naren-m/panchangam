import Foundation

struct PanchangamAPIClient {
    var baseURL: URL
    var session: URLSession = .shared

    static func liveFromBundle() -> PanchangamAPIClient {
        let value = Bundle.main.object(forInfoDictionaryKey: "PanchangamAPIBaseURL") as? String
        let url = URL(string: value ?? "") ?? URL(string: "http://127.0.0.1:8080")!
        return PanchangamAPIClient(baseURL: url)
    }

    func fetchToday(date: Date = Date()) async throws -> PanchangamDay {
        let formatter = DateFormatter()
        formatter.calendar = Calendar(identifier: .gregorian)
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone.current
        formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssXXXXX"

        var components = URLComponents(url: baseURL.appendingPathComponent("/api/v1/tithi/current"), resolvingAgainstBaseURL: false)
        components?.queryItems = [
            URLQueryItem(name: "at", value: formatter.string(from: date)),
            URLQueryItem(name: "lat", value: "37.3382"),
            URLQueryItem(name: "lng", value: "-121.8863"),
            URLQueryItem(name: "tz", value: TimeZone.current.identifier),
            URLQueryItem(name: "region", value: "California"),
            URLQueryItem(name: "method", value: "Drik"),
            URLQueryItem(name: "locale", value: "en"),
            URLQueryItem(name: "calendar_system", value: "Purnimanta")
        ]

        guard let url = components?.url else {
            throw PanchangamAPIError.invalidURL
        }

        let (data, response) = try await session.data(from: url)
        guard let httpResponse = response as? HTTPURLResponse else {
            throw PanchangamAPIError.invalidResponse
        }
        guard (200...299).contains(httpResponse.statusCode) else {
            throw PanchangamAPIError.badStatus(httpResponse.statusCode)
        }

        return try JSONDecoder().decode(PanchangamDay.self, from: data)
    }
}

enum PanchangamAPIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case badStatus(Int)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "The Panchangam API URL is invalid."
        case .invalidResponse:
            return "The Panchangam API response was invalid."
        case .badStatus(let statusCode):
            return "The Panchangam API returned HTTP \(statusCode)."
        }
    }
}
