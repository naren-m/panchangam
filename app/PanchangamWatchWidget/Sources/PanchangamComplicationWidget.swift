import SwiftUI
import WidgetKit

@main
struct PanchangamComplicationWidget: Widget {
    let kind = "PanchangamTithiComplication"

    var body: some WidgetConfiguration {
        StaticConfiguration(kind: kind, provider: TithiTimelineProvider()) { entry in
            TithiComplicationEntryView(entry: entry)
        }
        .configurationDisplayName("Panchangam Tithi")
        .description("Shows the current tithi and pancha anga details as a watch face complication.")
        .supportedFamilies([
            .accessoryInline,
            .accessoryCircular,
            .accessoryRectangular,
            .accessoryCorner
        ])
    }
}
