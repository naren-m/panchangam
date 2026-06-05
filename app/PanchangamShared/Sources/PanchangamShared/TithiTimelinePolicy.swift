import Foundation

public struct TithiTimelinePolicy: Sendable {
    public var staleBackoff: TimeInterval
    public var emptyBackoff: TimeInterval
    public var entryCadence: TimeInterval

    public init(
        staleBackoff: TimeInterval = 15 * 60,
        emptyBackoff: TimeInterval = 30 * 60,
        entryCadence: TimeInterval = 30 * 60
    ) {
        self.staleBackoff = staleBackoff
        self.emptyBackoff = emptyBackoff
        self.entryCadence = entryCadence
    }

    public func refreshDate(for summary: TithiSummaryResponse?, now: Date = Date()) -> Date {
        guard let summary else {
            return now.addingTimeInterval(emptyBackoff)
        }

        if summary.nextRefreshAt > now {
            return summary.nextRefreshAt
        }

        return now.addingTimeInterval(staleBackoff)
    }

    public func entryDates(for summary: TithiSummaryResponse?, now: Date = Date()) -> [Date] {
        guard let summary, summary.nextRefreshAt > now else {
            return [now]
        }

        var dates = [now]
        var nextDate = now.addingTimeInterval(entryCadence)
        while nextDate < summary.nextRefreshAt {
            dates.append(nextDate)
            nextDate = nextDate.addingTimeInterval(entryCadence)
        }
        return dates
    }
}
