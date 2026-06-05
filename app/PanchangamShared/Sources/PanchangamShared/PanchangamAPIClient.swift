import Foundation

public protocol PanchangamHTTPSession {
    func data(for request: URLRequest) async throws -> (Data, URLResponse)
}

extension URLSession: PanchangamHTTPSession {}

public enum PanchangamAPIError: Error {
    case invalidURL
    case missingHTTPResponse
    case httpStatus(code: Int, message: String?)
    case invalidSummary(String)
    case decoding(Error)

    public var userMessage: String {
        switch self {
        case .invalidURL:
            return "Invalid API base URL"
        case .missingHTTPResponse:
            return "API did not return an HTTP response"
        case .httpStatus(let code, let message):
            if let message, !message.isEmpty {
                return "HTTP \(code): \(message)"
            }
            return "HTTP \(code) from Panchangam API"
        case .invalidSummary(let message):
            return message
        case .decoding:
            return "Could not read tithi response"
        }
    }
}

public struct PanchangamAPIClient {
    public var baseURL: URL

    private let session: any PanchangamHTTPSession

    public init(baseURL: URL = URL(string: "http://localhost:8080")!, session: any PanchangamHTTPSession = URLSession.shared) {
        self.baseURL = baseURL
        self.session = session
    }

    public func currentTithi(settings: APISettings, at: Date? = nil) async throws -> TithiSummaryResponse {
        let url = try makeTithiSummaryURL(settings: settings, at: at)
        var request = URLRequest(url: url)
        request.timeoutInterval = 15

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else {
            throw PanchangamAPIError.missingHTTPResponse
        }

        guard 200..<300 ~= httpResponse.statusCode else {
            throw PanchangamAPIError.httpStatus(
                code: httpResponse.statusCode,
                message: decodeErrorMessage(from: data)
            )
        }

        do {
            let summary = try PanchangamJSON.decoder.decode(TithiSummaryResponse.self, from: data)
            try summary.validate()
            return summary
        } catch let error as TithiSummaryValidationError {
            throw PanchangamAPIError.invalidSummary(error.userMessage)
        } catch {
            throw PanchangamAPIError.decoding(error)
        }
    }

    public func makeTithiSummaryURL(settings: APISettings, at: Date? = nil) throws -> URL {
        try settings.validate()

        var components = URLComponents(url: settings.apiBaseURL, resolvingAgainstBaseURL: false)
        if components == nil {
            components = URLComponents(url: baseURL, resolvingAgainstBaseURL: false)
        }

        guard var components else {
            throw PanchangamAPIError.invalidURL
        }

        let basePath = components.path.trimmingCharacters(in: CharacterSet(charactersIn: "/"))
        if basePath.isEmpty {
            components.path = "/api/v1/tithi/current"
        } else {
            components.path = "/\(basePath)/api/v1/tithi/current"
        }

        var queryItems = [
            URLQueryItem(name: "lat", value: decimalString(settings.latitude)),
            URLQueryItem(name: "lng", value: decimalString(settings.longitude)),
            URLQueryItem(name: "tz", value: settings.timezone),
            URLQueryItem(name: "region", value: settings.region),
            URLQueryItem(name: "method", value: settings.method),
            URLQueryItem(name: "locale", value: settings.locale)
        ]

        if let at {
            queryItems.insert(URLQueryItem(name: "at", value: PanchangamJSON.format(at)), at: 0)
        }

        if let calendarSystem = settings.calendarSystem, !calendarSystem.isEmpty {
            queryItems.append(URLQueryItem(name: "calendar_system", value: calendarSystem))
        }

        components.queryItems = queryItems

        guard let url = components.url else {
            throw PanchangamAPIError.invalidURL
        }
        return url
    }

    private func decimalString(_ value: Double) -> String {
        let formatter = NumberFormatter()
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.numberStyle = .decimal
        formatter.maximumFractionDigits = 8
        formatter.minimumFractionDigits = 0
        formatter.usesGroupingSeparator = false
        return formatter.string(from: NSNumber(value: value)) ?? String(value)
    }

    private func decodeErrorMessage(from data: Data) -> String? {
        let decoder = JSONDecoder()
        return try? decoder.decode(APIErrorEnvelope.self, from: data).error.message
    }
}

private struct APIErrorEnvelope: Decodable {
    let error: ErrorBody

    struct ErrorBody: Decodable {
        let message: String
    }
}
