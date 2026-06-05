const sample = {
  tithiName: "Sashthi",
  traditionalName: "Sashthi",
  tithiNumber: 6,
  paksha: "Shukla",
  pakshaDay: 6,
  tithiType: "Nanda",
  nakshatra: "Pushya",
  yoga: "Saubhagya",
  karana: "Taitila",
  vara: "Tuesday",
  sunrise: "05:42",
  sunset: "19:01",
  abhijit: "11:54-12:48",
  startTime: "2026-06-02 01:40 UTC",
  endTime: "2026-06-03 02:00 UTC",
  progressText: "44% elapsed",
  remainingText: "14h 0m left",
  region: "California",
  timezone: "America/Los_Angeles",
  calendarSystem: "Purnimanta",
  method: "Drik",
  locale: "en",
  generatedAt: "2026-06-02 12:00 UTC",
  nextRefreshAt: "2026-06-03 02:00 UTC",
};

const stateCopy = {
  valid: {
    statusText: "Loaded at 2026-06-02 12:00 UTC",
    watchAppStatusText: "Loaded",
    watchSyncText: "Watch settings synced at 12:00",
    receiverStatusText: "Settings received",
    complicationText: null,
  },
  stale: {
    statusText: "Showing cached result",
    watchAppStatusText: "Showing cached result",
    watchSyncText: "Watch sync ready",
    receiverStatusText: "Watch sync ready",
    complicationText: null,
  },
  loading: {
    statusText: "Loading",
    watchAppStatusText: "Loading",
    watchSyncText: "Watch sync inactive",
    receiverStatusText: "Waiting for iPhone settings",
    complicationText: "Loading",
  },
  error: {
    statusText: "Backend unavailable",
    watchAppStatusText: "Backend unavailable",
    watchSyncText: "Watch sync failed",
    receiverStatusText: "Settings sync failed",
    complicationText: "Unavailable",
  },
};

const nakshatraNames = [
  "Ashwini",
  "Bharani",
  "Krittika",
  "Rohini",
  "Mrigashira",
  "Ardra",
  "Punarvasu",
  "Pushya",
  "Ashlesha",
  "Magha",
  "Purva Phalguni",
  "Uttara Phalguni",
  "Hasta",
  "Chitra",
  "Swati",
  "Vishakha",
  "Anuradha",
  "Jyeshtha",
  "Mula",
  "Purva Ashadha",
  "Uttara Ashadha",
  "Shravana",
  "Dhanishta",
  "Shatabhisha",
  "Purva Bhadrapada",
  "Uttara Bhadrapada",
  "Revati",
];

function normalizedNakshatraName(name) {
  return name.toLowerCase().replace(/[^a-z0-9]/g, "");
}

function nakshatraIndex() {
  const currentName = normalizedNakshatraName(sample.nakshatra);
  return nakshatraNames.findIndex((name) => normalizedNakshatraName(name) === currentName);
}

function normalizedTithiName(name) {
  return name.toLowerCase().replace(/[^a-z0-9]/g, "");
}

function tithiDisplayName() {
  if (!sample.traditionalName || normalizedTithiName(sample.traditionalName) === normalizedTithiName(sample.tithiName)) {
    return sample.tithiName;
  }

  return `${sample.tithiName} (${sample.traditionalName})`;
}

function inlineText() {
  return `${tithiDisplayName()}, ${sample.paksha} ${sample.pakshaDay}`;
}

function detailRows() {
  return [
    ["Tithi", inlineText()],
    ["Tithi Number", String(sample.tithiNumber)],
    ["Traditional Name", sample.traditionalName],
    ["Paksha", sample.paksha],
    ["Paksha Day", String(sample.pakshaDay)],
    ["Tithi Type", sample.tithiType],
    ["Nakshatra", sample.nakshatra],
    ["Yoga", sample.yoga],
    ["Karana", sample.karana],
    ["Vara", sample.vara],
    ["Sunrise", sample.sunrise],
    ["Sunset", sample.sunset],
    ["Abhijit", sample.abhijit],
    ["Starts", sample.startTime],
    ["Ends", sample.endTime],
    ["Progress", sample.progressText],
    ["Remaining", sample.remainingText],
    ["Generated", sample.generatedAt],
    ["Next Refresh", sample.nextRefreshAt],
    ["Region", sample.region],
    ["Timezone", sample.timezone],
    ["Calendar System", sample.calendarSystem],
    ["Method", sample.method],
    ["Locale", sample.locale],
  ];
}

function setText(selector, value) {
  document.querySelectorAll(selector).forEach((node) => {
    node.textContent = value;
  });
}

function renderDetails() {
  const details = document.getElementById("iphone-details");
  if (!details) {
    return;
  }

  details.innerHTML = "";
  detailRows().forEach(([label, value]) => {
    const row = document.createElement("div");
    const term = document.createElement("dt");
    const description = document.createElement("dd");
    term.textContent = label;
    description.textContent = value;
    row.append(term, description);
    details.append(row);
  });
}

function renderWatchDetails() {
  const details = document.getElementById("watch-details");
  if (!details) {
    return;
  }

  details.innerHTML = "";
  detailRows().forEach(([label, value]) => {
    const row = document.createElement("div");
    const term = document.createElement("dt");
    const description = document.createElement("dd");
    term.textContent = label;
    description.textContent = value;
    row.append(term, description);
    details.append(row);
  });
}

function openWatchDetails() {
  const sheet = document.getElementById("watch-detail-sheet");
  const closeButton = document.querySelector("[data-action='close-watch-details']");
  if (!sheet || !closeButton) {
    return;
  }

  sheet.hidden = false;
  closeButton.focus();
}

function closeWatchDetails() {
  const sheet = document.getElementById("watch-detail-sheet");
  const face = document.querySelector(".mandala-face");
  if (!sheet) {
    return;
  }

  sheet.hidden = true;
  if (face) {
    face.focus();
  }
}

function renderMandalaTicks() {
  const ring = document.getElementById("mandala-ring");
  if (!ring) {
    return;
  }

  const activeIndex = nakshatraIndex();
  ring.innerHTML = "";
  for (let index = 0; index < 27; index += 1) {
    const tick = document.createElement("span");
    tick.className = index === activeIndex ? "mandala-tick is-active" : "mandala-tick";
    tick.style.setProperty("--tick-angle", `${index * (360 / 27)}deg`);
    ring.append(tick);
  }
}

function renderState(state) {
  const copy = stateCopy[state];
  const complicationText = copy.complicationText;
  document.body.dataset.previewState = state;

  setText("[data-field='sunriseText']", sample.sunrise);
  setText("[data-field='sunsetText']", sample.sunset);
  setText("[data-field='inlineText']", complicationText || inlineText());
  setText("[data-field='mandalaTithi']", inlineText());
  setText("[data-field='mandalaNakshatra']", `Nakshatra - ${sample.nakshatra}`);
  setText("[data-field='phoneSubtitle']", `${sample.nakshatra} - ${sample.remainingText}`);
  setText("[data-field='yogaText']", sample.yoga);
  setText("[data-field='karanaText']", sample.karana);
  setText("[data-field='statusText']", copy.statusText);
  setText("[data-field='watchAppStatusText']", copy.watchAppStatusText);
  setText("[data-field='watchSyncText']", copy.watchSyncText);
  setText("[data-field='receiverStatusText']", copy.receiverStatusText);
  setText("[data-field='backendHostWarning']", "Use a Mac LAN address for physical iPhone or Apple Watch testing.");
  document.querySelectorAll("[data-field='statusText']").forEach((node) => {
    node.dataset.statusTone = state;
  });
}

const mandalaFace = document.querySelector(".mandala-face");
if (mandalaFace) {
  mandalaFace.addEventListener("click", openWatchDetails);
  mandalaFace.addEventListener("keydown", (event) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      openWatchDetails();
    }
  });
}

const watchDetailSheet = document.getElementById("watch-detail-sheet");
if (watchDetailSheet) {
  watchDetailSheet.addEventListener("click", (event) => {
    if (event.target === watchDetailSheet) {
      closeWatchDetails();
    }
  });
}

const closeWatchDetailsButton = document.querySelector("[data-action='close-watch-details']");
if (closeWatchDetailsButton) {
  closeWatchDetailsButton.addEventListener("click", closeWatchDetails);
}

document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && watchDetailSheet && !watchDetailSheet.hidden) {
    closeWatchDetails();
  }
});

renderMandalaTicks();
renderDetails();
renderWatchDetails();
renderState("valid");
