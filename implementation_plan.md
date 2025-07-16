# Panchangam Project: Implementation Plan

This file outlines the major tasks and sub-tasks for the Panchangam project. Each sub-task is suitable to be tracked as an individual GitHub issue.

---

## 1. Core Astronomical Calculations

> **All astronomical and calendar calculations must use latitude, longitude, timezone, region, and calculation method from the request to ensure location- and tradition-accurate results.**

### 1.1 Tithi Calculation
- Implement Moon and Sun longitude calculation (use request's location and method fields)
- Compute tithi using (Moon_longitude - Sun_longitude) ÷ 12°
- Handle edge cases: multiple tithis in a day, tithis spanning days
- Categorize tithis (Nanda, Bhadra, Jaya, Rikta, Purna)

### 1.2 Vara (Weekday) Calculation
- Calculate sunrise times for any location and date
- Assign weekdays based on sunrise
- Integrate planetary associations and hora system
- Handle time zones and daylight saving

### 1.3 Nakshatra Calculation
- Divide zodiac into 27 segments (13°20' each)
- Track Moon’s position using ephemeris
- Implement pada (quarter) calculations
- Store nakshatra attributes (deity, planetary lord, symbolism)

### 1.4 Yoga Calculation
- Calculate as (Sun_longitude + Moon_longitude) ÷ 13°20'
- Flag inauspicious and auspicious yogas
- Provide visual/tabular representation of yoga transitions

### 1.5 Karana Calculation
- Divide each tithi into two karanas (6° each)
- Implement the 11 karana cycle (7 movable, 4 fixed)
- Special handling for Vishti (Bhadra) karana

---

## 1a. Panchangam Event Generation
- Generate events such as Rahu Kalam, Yamagandam, and festivals using request parameters
- Map each event to the structured `event_type` field in the proto (e.g., RAHU_KALAM, YAMAGANDAM, FESTIVAL, etc.)

---

## 2. Regional Variations & Calendar Systems
- Support Amanta (South Indian) and Purnimanta (North Indian) month endings
- Allow user selection for Drik Ganita or Vakya calculation methods
- Implement logic for Tamil Nadu, Kerala, Bengal, Gujarat, Maharashtra
- Support Naazhikai time units and local language preferences

---

## 3. Astronomical Data Handling
- Integrate with a modern ephemeris (e.g., JPL, Swiss Ephemeris)
- Implement interpolation methods for planetary positions
- Handle retrograde motion and planetary stations

---

## 4. User Interface and Experience
- Design progressive disclosure UI (summary first, drill-down for details)
- Visualize planetary positions, lunar phases, muhurta qualities
- Allow user selection for location, calculation method, language
- Present both tabular and graphical data

---

## 5. Festival and Muhurta Calculations
- Implement logic for major Hindu festivals (regional/method differences)
- Calculate and display daily muhurta windows (auspicious times)

---

## 6. Testing and Validation
- Cross-verify calculations with established Panchangam sources
- Implement unit and integration tests for all core modules
- Increase test coverage for the observability package, including error and edge cases

---

## 7. Documentation and Extensibility
- Document algorithms, regional logic, and data sources
- Provide clear APIs for extension

---

## 8. Observability and Logging
- Instrument all service endpoints and core calculation steps with context-aware logging and OpenTelemetry tracing
- Ensure errors and important events are logged and span events are created for traceability

---

### [GitHub Issue Suggestions]

Each of the above sub-tasks can be created as an individual GitHub issue. For example:

1. Implement Moon and Sun longitude calculation
2. Compute tithi using (Moon_longitude - Sun_longitude) ÷ 12°
3. Handle tithi edge cases (multiple tithis/day, spanning days)
4. Categorize tithis into Nanda, Bhadra, Jaya, Rikta, Purna
5. Calculate sunrise times for any location and date
6. Assign weekdays based on sunrise
7. Integrate planetary associations and hora system
8. Handle time zones and daylight saving in weekday calculation
9. Divide zodiac into 27 segments and track Moon’s position
10. Implement pada calculations for nakshatra
11. Store nakshatra attributes (deity, planetary lord, symbolism)
12. Calculate yogas and flag auspicious/inauspicious periods
13. Visualize yoga transitions
14. Divide each tithi into two karanas and implement cycle
15. Special handling for Vishti (Bhadra) karana
16. Support Amanta and Purnimanta month endings
17. Allow user selection for Drik Ganita or Vakya methods
18. Implement regional logic for Tamil Nadu, Kerala, Bengal, etc.
19. Support Naazhikai time units and local language preferences
20. Integrate with modern ephemeris
21. Implement interpolation for planetary positions
22. Handle retrograde motion and planetary stations
23. Design progressive disclosure UI
24. Visualize planetary/lunar/muhurta data
25. Allow user selection for location, method, language
26. Present data in tabular and graphical formats
27. Implement festival calculation logic
28. Calculate/display daily muhurta windows
29. Cross-verify with established Panchangam sources
30. Implement unit/integration tests
31. Document algorithms, regional logic, data sources
32. Provide extensible APIs

Each issue should include acceptance criteria and reference relevant sections of this plan.


This file outlines the major tasks and sub-tasks for the Panchangam project. Each sub-task is suitable to be tracked as an individual GitHub issue.

---

## 1. Core Astronomical Calculations

### 1.1 Tithi Calculation

- Implement Moon and Sun longitude calculation
- Compute tithi using (Moon_longitude - Sun_longitude) ÷ 12°
- Handle edge cases: multiple tithis in a day, tithis spanning days
- Categorize tithis (Nanda, Bhadra, Jaya, Rikta, Purna)

### 1.2 Vara (Weekday) Calculation

- Calculate sunrise times for any location and date
- Assign weekdays based on sunrise
- Integrate planetary associations and hora system
- Handle time zones and daylight saving

### 1.3 Nakshatra Calculation

- Divide zodiac into 27 segments (13°20' each)
- Track Moon’s position using ephemeris
- Implement pada (quarter) calculations
- Store nakshatra attributes (deity, planetary lord, symbolism)

### 1.4 Yoga Calculation

- Calculate as (Sun_longitude + Moon_longitude) ÷ 13°20'
- Flag inauspicious and auspicious yogas
- Provide visual/tabular representation of yoga transitions

### 1.5 Karana Calculation

- Divide each tithi into two karanas (6° each)
- Implement the 11 karana cycle (7 movable, 4 fixed)
- Special handling for Vishti (Bhadra) karana

---

## 2. Regional Variations & Calendar Systems

- Support Amanta (South Indian) and Purnimanta (North Indian) month endings
- Allow user selection for Drik Ganita or Vakya calculation methods
- Implement logic for Tamil Nadu, Kerala, Bengal, Gujarat, Maharashtra
- Support Naazhikai time units and local language preferences

---

## 3. Astronomical Data Handling

- Integrate with a modern ephemeris (e.g., JPL, Swiss Ephemeris)
- Implement interpolation methods for planetary positions
- Handle retrograde motion and planetary stations

---

## 4. User Interface and Experience

- Design progressive disclosure UI (summary first, drill-down for details)
- Visualize planetary positions, lunar phases, muhurta qualities
- Allow user selection for location, calculation method, language
- Present both tabular and graphical data

---

## 5. Festival and Muhurta Calculations

- Implement logic for major Hindu festivals (regional/method differences)
- Calculate and display daily muhurta windows (auspicious times)

---

## 6. Testing and Validation

- Cross-verify calculations with established Panchangam sources
- Implement unit and integration tests for all core modules

---

## 7. Documentation and Extensibility

- Document algorithms, regional logic, and data sources
- Provide clear APIs for extension

---

### [GitHub Issue Suggestions]

Each of the above sub-tasks can be created as an individual GitHub issue. For example:

1. Implement Moon and Sun longitude calculation
2. Compute tithi using (Moon_longitude - Sun_longitude) ÷ 12°
3. Handle tithi edge cases (multiple tithis/day, spanning days)
4. Categorize tithis into Nanda, Bhadra, Jaya, Rikta, Purna
5. Calculate sunrise times for any location and date
6. Assign weekdays based on sunrise
7. Integrate planetary associations and hora system
8. Handle time zones and daylight saving in weekday calculation
9. Divide zodiac into 27 segments and track Moon’s position
10. Implement pada calculations for nakshatra
11. Store nakshatra attributes (deity, planetary lord, symbolism)
12. Calculate yogas and flag auspicious/inauspicious periods
13. Visualize yoga transitions
14. Divide each tithi into two karanas and implement cycle
15. Special handling for Vishti (Bhadra) karana
16. Support Amanta and Purnimanta month endings
17. Allow user selection for Drik Ganita or Vakya methods
18. Implement regional logic for Tamil Nadu, Kerala, Bengal, etc.
19. Support Naazhikai time units and local language preferences
20. Integrate with modern ephemeris
21. Implement interpolation for planetary positions
22. Handle retrograde motion and planetary stations
23. Design progressive disclosure UI
24. Visualize planetary/lunar/muhurta data
25. Allow user selection for location, method, language
26. Present data in tabular and graphical formats
27. Implement festival calculation logic
28. Calculate/display daily muhurta windows
29. Cross-verify with established Panchangam sources
30. Implement unit/integration tests
31. Document algorithms, regional logic, data sources
32. Provide extensible APIs

Each issue should include acceptance criteria and reference relevant sections of this plan.
