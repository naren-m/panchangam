/**
 * Annotation Helper for Educational Content
 *
 * Generates informative annotations that explain Panchangam elements
 * to users in an educational, accessible way.
 *
 * Annotations are shown when users hover over or click on elements,
 * providing context about what each element means and how it's calculated.
 */

import {
  Annotation,
  PanchangamElements,
  TithiInfo,
  NakshatraInfo,
  YogaInfo,
  KaranaInfo,
  RashiInfo,
  EclipticBeltDimensions,
  ScreenPosition,
} from '../types/eclipticBelt';
import { longitudeToX, getZoneCenterY } from './eclipticLayout';

// ============================================================================
// Annotation Generators
// ============================================================================

/**
 * Generate annotation for current Tithi
 */
export function generateTithiAnnotation(
  tithi: TithiInfo,
  sunLongitude: number,
  moonLongitude: number,
  dimensions: EclipticBeltDimensions
): Annotation {
  const pakshaName = tithi.paksha === 'Shukla'
    ? 'Shukla Paksha (Bright Fortnight - Waxing Moon)'
    : 'Krishna Paksha (Dark Fortnight - Waning Moon)';

  const dayInPaksha = tithi.paksha === 'Shukla' ? tithi.number : tithi.number - 15;

  const detail = [
    `**${tithi.name}** - Day ${dayInPaksha} of ${pakshaName}`,
    '',
    `**What is Tithi?**`,
    `Tithi is the lunar day, measured by the angular distance between`,
    `the Moon and Sun. Each tithi spans 12° of separation.`,
    '',
    `**Current Calculation:**`,
    `Moon (${moonLongitude.toFixed(1)}°) - Sun (${sunLongitude.toFixed(1)}°) = ${tithi.angle.toFixed(1)}°`,
    `${tithi.angle.toFixed(1)}° ÷ 12° = Tithi ${tithi.number}`,
    '',
    `**Progress:** ${tithi.percentComplete.toFixed(0)}% of ${tithi.name} elapsed`,
    `**Presiding Deity:** ${tithi.deity}`
  ].join('\n');

  // Position annotation at the midpoint between Sun and Moon
  const midX = (longitudeToX(sunLongitude, dimensions) + longitudeToX(moonLongitude, dimensions)) / 2;
  const y = getZoneCenterY('annotation', dimensions);

  return {
    id: 'tithi-annotation',
    type: 'tithi',
    content: `${tithi.paksha} ${tithi.name}`,
    detail,
    position: { x: midX, y },
    highlight: {
      startX: longitudeToX(sunLongitude, dimensions),
      endX: longitudeToX(moonLongitude, dimensions)
    }
  };
}

/**
 * Generate annotation for current Nakshatra
 */
export function generateNakshatraAnnotation(
  nakshatra: NakshatraInfo,
  moonLongitude: number,
  dimensions: EclipticBeltDimensions
): Annotation {
  const detail = [
    `**${nakshatra.name}** - Nakshatra #${nakshatra.number} of 27`,
    '',
    `**What is Nakshatra?**`,
    `Nakshatras are 27 lunar mansions that divide the ecliptic into`,
    `13°20' segments. The Moon visits each nakshatra monthly.`,
    '',
    `**Current Position:**`,
    `Moon at ${moonLongitude.toFixed(2)}° is in ${nakshatra.name}`,
    `(${nakshatra.startDegree.toFixed(1)}° - ${nakshatra.endDegree.toFixed(1)}°)`,
    '',
    `**Pada (Quarter):** ${nakshatra.pada} of 4`,
    `**Presiding Deity:** ${nakshatra.deity}`,
    `**Symbol:** ${nakshatra.symbol}`,
    '',
    `*Each nakshatra spans 3°20' per pada (quarter)*`
  ].join('\n');

  const x = longitudeToX(moonLongitude, dimensions);
  const y = getZoneCenterY('annotation', dimensions);

  return {
    id: 'nakshatra-annotation',
    type: 'nakshatra',
    content: `${nakshatra.name} (Pada ${nakshatra.pada})`,
    detail,
    position: { x, y },
    highlight: {
      startX: longitudeToX(nakshatra.startDegree, dimensions),
      endX: longitudeToX(nakshatra.endDegree, dimensions)
    }
  };
}

/**
 * Generate annotation for current Yoga
 */
export function generateYogaAnnotation(
  yoga: YogaInfo,
  sunLongitude: number,
  moonLongitude: number,
  dimensions: EclipticBeltDimensions
): Annotation {
  const natureEmoji = yoga.nature === 'Auspicious' ? '\u2728' :
    yoga.nature === 'Inauspicious' ? '\u26A0\uFE0F' : '\u2696\uFE0F';

  const detail = [
    `**${yoga.name}** - Yoga #${yoga.number} of 27 ${natureEmoji}`,
    '',
    `**What is Yoga?**`,
    `Yoga is calculated from the SUM of Sun and Moon longitudes,`,
    `divided into 27 parts of 13°20' each.`,
    '',
    `**Current Calculation:**`,
    `Sun (${sunLongitude.toFixed(1)}°) + Moon (${moonLongitude.toFixed(1)}°) = ${yoga.combinedLongitude.toFixed(1)}°`,
    `${yoga.combinedLongitude.toFixed(1)}° ÷ 13.33° = Yoga ${yoga.number}`,
    '',
    `**Meaning:** "${yoga.meaning}"`,
    `**Nature:** ${yoga.nature}`,
    '',
    `*Yoga indicates the combined influence of Sun and Moon*`
  ].join('\n');

  // Position in annotation zone
  const x = dimensions.padding.left + (dimensions.width - dimensions.padding.left - dimensions.padding.right) * 0.25;
  const y = getZoneCenterY('annotation', dimensions);

  return {
    id: 'yoga-annotation',
    type: 'yoga',
    content: `${yoga.name} (${yoga.nature})`,
    detail,
    position: { x, y }
  };
}

/**
 * Generate annotation for current Karana
 */
export function generateKaranaAnnotation(
  karana: KaranaInfo,
  tithi: TithiInfo,
  dimensions: EclipticBeltDimensions
): Annotation {
  const isFirstHalf = (karana.number % 2) === 1;
  const halfLabel = isFirstHalf ? 'First Half' : 'Second Half';

  const detail = [
    `**${karana.name}** - Karana #${karana.number} of 60`,
    '',
    `**What is Karana?**`,
    `Karana is half of a Tithi. Each tithi has 2 karanas,`,
    `so there are 60 karanas in a lunar month.`,
    '',
    `**Current Position:**`,
    `${halfLabel} of ${tithi.name} (Tithi ${tithi.number})`,
    '',
    `**Type:** ${karana.type}`,
    `**Nature:** ${karana.nature}`,
    '',
    `*There are 11 karana types: 4 fixed + 7 movable*`,
    `*Vishti (Bhadra) karana is considered inauspicious*`
  ].join('\n');

  // Position in annotation zone
  const x = dimensions.padding.left + (dimensions.width - dimensions.padding.left - dimensions.padding.right) * 0.5;
  const y = getZoneCenterY('annotation', dimensions);

  return {
    id: 'karana-annotation',
    type: 'karana',
    content: `${karana.name} (${karana.type})`,
    detail,
    position: { x, y }
  };
}

/**
 * Generate annotation for a Rashi (Zodiac Sign)
 */
export function generateRashiAnnotation(
  rashi: RashiInfo,
  celestialBody: 'sun' | 'moon',
  longitude: number,
  dimensions: EclipticBeltDimensions
): Annotation {
  const bodyName = celestialBody === 'sun' ? 'Sun' : 'Moon';
  const degreeInSign = longitude - rashi.startDegree;

  const detail = [
    `**${rashi.symbol} ${rashi.name}** (${rashi.westernName})`,
    '',
    `**What is Rashi?**`,
    `Rashi means "zodiac sign" in Sanskrit. The ecliptic is divided`,
    `into 12 rashis of 30° each.`,
    '',
    `**${bodyName}'s Position:**`,
    `${longitude.toFixed(2)}° = ${degreeInSign.toFixed(1)}° into ${rashi.name}`,
    '',
    `**Element:** ${rashi.element}`,
    `**Ruling Planet:** ${rashi.ruler}`,
    `**Span:** ${rashi.startDegree}° - ${rashi.endDegree}°`,
    '',
    `*Indian astrology uses sidereal zodiac (fixed stars)*`,
    `*Western astrology uses tropical zodiac (seasons)*`
  ].join('\n');

  const x = longitudeToX(longitude, dimensions);
  const y = getZoneCenterY('annotation', dimensions);

  return {
    id: `rashi-${celestialBody}-annotation`,
    type: 'rashi',
    content: `${rashi.symbol} ${rashi.name}`,
    detail,
    position: { x, y },
    highlight: {
      startX: longitudeToX(rashi.startDegree, dimensions),
      endX: longitudeToX(rashi.endDegree, dimensions)
    }
  };
}

// ============================================================================
// Info Card Content
// ============================================================================

/**
 * Generate summary card content for all panchangam elements
 */
export function generatePanchangamSummary(panchangam: PanchangamElements): string {
  const { tithi, nakshatra, yoga, karana, rashi, sunRashi, sunPosition, moonPosition } = panchangam;

  return [
    `✨ **Panchangam Summary**`,
    '',
    `🌙 **Tithi:** ${tithi.paksha} ${tithi.name} (${tithi.percentComplete.toFixed(0)}%)`,
    `⭐ **Nakshatra:** ${nakshatra.name} (Pada ${nakshatra.pada})`,
    `🔮 **Yoga:** ${yoga.name} - ${yoga.meaning}`,
    `⌚ **Karana:** ${karana.name}`,
    '',
    `☉ **Sun:** ${sunRashi.symbol} ${sunRashi.name} (${sunPosition.longitude.toFixed(1)}°)`,
    `🌑 **Moon:** ${rashi.symbol} ${rashi.name} (${moonPosition.longitude.toFixed(1)}°)`,
    '',
    `*The five elements (Panch Anga) are:*`,
    `*Tithi, Nakshatra, Yoga, Karana, and Vara (weekday)*`
  ].join('\n');
}

/**
 * Generate educational explanation of the visualization
 */
export function generateVisualizationGuide(): string {
  return [
    `📚 **Understanding the Ecliptic Belt**`,
    '',
    `The **ecliptic** is the apparent path of the Sun across the sky`,
    `over a year. The Moon and planets also travel near this path.`,
    '',
    `**Reading the Visualization:**`,
    '',
    `🟠 **Top Row (Rashis):** The 12 zodiac signs, each 30° wide`,
    `   Start: ♈ Aries (Mesha) at 0° | End: ♓ Pisces (Meena) at 360°`,
    '',
    `🟡 **Second Row (Nakshatras):** 27 lunar mansions, each 13.33° wide`,
    `   These are the Moon's "stopping places" as it orbits Earth`,
    '',
    `☉🌑 **Middle (Planets):** Current positions of Sun and Moon`,
    `   The Moon moves ~13°/day, Sun moves ~1°/day`,
    '',
    `🌈 **Arc (Tithi):** Shows Moon-Sun separation (determines lunar day)`,
    `   360° separation = complete lunar month (30 tithis)`,
    '',
    `*Click or hover on any element for detailed information!*`
  ].join('\n');
}

// ============================================================================
// Annotation Collection
// ============================================================================

/**
 * Generate all annotations for the current panchangam
 */
export function generateAllAnnotations(
  panchangam: PanchangamElements,
  dimensions: EclipticBeltDimensions
): Annotation[] {
  const { tithi, nakshatra, yoga, karana, sunPosition, moonPosition } = panchangam;

  return [
    generateTithiAnnotation(
      tithi,
      sunPosition.longitude,
      moonPosition.longitude,
      dimensions
    ),
    generateNakshatraAnnotation(
      nakshatra,
      moonPosition.longitude,
      dimensions
    ),
    generateYogaAnnotation(
      yoga,
      sunPosition.longitude,
      moonPosition.longitude,
      dimensions
    ),
    generateKaranaAnnotation(
      karana,
      tithi,
      dimensions
    ),
    generateRashiAnnotation(
      moonPosition.rashi,
      'moon',
      moonPosition.longitude,
      dimensions
    )
  ];
}

// ============================================================================
// Tooltip Content Generators
// ============================================================================

/**
 * Generate short tooltip for hover state
 */
export function getTooltipContent(
  elementType: 'sun' | 'moon' | 'rashi' | 'nakshatra' | 'tithi',
  panchangam: PanchangamElements
): { title: string; subtitle: string } {
  switch (elementType) {
    case 'sun':
      return {
        title: `☉ Sun in ${panchangam.sunRashi.name}`,
        subtitle: `${panchangam.sunPosition.longitude.toFixed(1)}°`
      };

    case 'moon':
      return {
        title: `🌑 Moon in ${panchangam.nakshatra.name}`,
        subtitle: `${panchangam.moonPosition.longitude.toFixed(1)}° • Pada ${panchangam.nakshatra.pada}`
      };

    case 'rashi':
      return {
        title: `${panchangam.rashi.symbol} ${panchangam.rashi.name}`,
        subtitle: `Element: ${panchangam.rashi.element} • Ruler: ${panchangam.rashi.ruler}`
      };

    case 'nakshatra':
      return {
        title: `⭐ ${panchangam.nakshatra.name}`,
        subtitle: `Deity: ${panchangam.nakshatra.deity}`
      };

    case 'tithi':
      return {
        title: `🌙 ${panchangam.tithi.paksha} ${panchangam.tithi.name}`,
        subtitle: `${panchangam.tithi.percentComplete.toFixed(0)}% complete • ${panchangam.tithi.angle.toFixed(0)}° separation`
      };

    default:
      return { title: '', subtitle: '' };
  }
}

// ============================================================================
// Learning Mode Helpers
// ============================================================================

/**
 * Generate step-by-step calculation explanation
 */
export function generateCalculationSteps(
  elementType: 'tithi' | 'yoga' | 'karana',
  sunLong: number,
  moonLong: number
): string[] {
  switch (elementType) {
    case 'tithi':
      const tithiAngle = ((moonLong - sunLong + 360) % 360);
      const tithiNum = Math.floor(tithiAngle / 12) + 1;
      return [
        `Step 1: Get Moon longitude = ${moonLong.toFixed(2)}°`,
        `Step 2: Get Sun longitude = ${sunLong.toFixed(2)}°`,
        `Step 3: Calculate difference = Moon - Sun = ${tithiAngle.toFixed(2)}°`,
        `Step 4: Divide by 12° = ${(tithiAngle / 12).toFixed(2)}`,
        `Step 5: Tithi number = floor(${(tithiAngle / 12).toFixed(2)}) + 1 = ${tithiNum}`
      ];

    case 'yoga':
      const yogaSum = (sunLong + moonLong) % 360;
      const yogaNum = Math.floor(yogaSum / 13.333333) + 1;
      return [
        `Step 1: Get Sun longitude = ${sunLong.toFixed(2)}°`,
        `Step 2: Get Moon longitude = ${moonLong.toFixed(2)}°`,
        `Step 3: Add together = ${yogaSum.toFixed(2)}°`,
        `Step 4: Divide by 13.33° = ${(yogaSum / 13.333333).toFixed(2)}`,
        `Step 5: Yoga number = floor + 1 = ${yogaNum}`
      ];

    case 'karana':
      const karanaAngle = ((moonLong - sunLong + 360) % 360);
      const karanaNum = Math.floor(karanaAngle / 6) + 1;
      return [
        `Step 1: Calculate Moon-Sun difference = ${karanaAngle.toFixed(2)}°`,
        `Step 2: Divide by 6° = ${(karanaAngle / 6).toFixed(2)}`,
        `Step 3: Karana number = floor + 1 = ${karanaNum}`,
        `Step 4: Each tithi (12°) has 2 karanas (6° each)`
      ];

    default:
      return [];
  }
}
