/**
 * Ecliptic Belt Visualization
 *
 * 2D visualization of the ecliptic (Sun's apparent path through the sky)
 * showing Panchangam elements: Tithi, Nakshatra, Yoga, Karana, and Rashi.
 *
 * Usage:
 * ```tsx
 * import { EclipticBeltContainer } from './components/EclipticBeltVisualization';
 *
 * <EclipticBeltContainer
 *   date={new Date()}
 *   latitude={37.4323}
 *   longitude={-121.9066}
 *   onClose={() => setShowVisualization(false)}
 * />
 * ```
 */

// Main container component
export { EclipticBeltContainer } from './EclipticBeltContainer';
export { default as EclipticBeltContainerDefault } from './EclipticBeltContainer';

// SVG renderer (for custom implementations)
export { EclipticBeltSVG } from './EclipticBeltSVG';
export { default as EclipticBeltSVGDefault } from './EclipticBeltSVG';

// Type exports
export type {
  EclipticBeltContainerProps,
  EclipticBeltSVGProps,
  EclipticBeltDimensions,
  PanchangamElements,
  TithiInfo,
  NakshatraInfo,
  YogaInfo,
  KaranaInfo,
  RashiInfo,
  SunPosition,
  MoonPosition,
  EclipticSegment,
  CelestialMarker,
  RelationshipArc,
  Annotation,
  TimeControlState,
  InteractionState,
} from './types/eclipticBelt';

// Utility exports (for custom calculations)
export {
  calculatePanchangamElements,
  calculateTithi,
  calculateNakshatra,
  calculateYoga,
  calculateKarana,
  calculateRashi,
  normalizeDegrees,
  getTithiDisplayName,
  getTithiExplanation,
  getYogaExplanation,
  getKaranaExplanation,
} from './utils/panchangamCalculator';

export {
  longitudeToX,
  xToLongitude,
  getZoneCenterY,
  getZoneYRange,
  generateRashiSegments,
  generateNakshatraSegments,
  generateLayout,
  createDefaultDimensions,
  interpolateLongitude,
  calculateAnimatedPositions,
} from './utils/eclipticLayout';

export {
  generateTithiAnnotation,
  generateNakshatraAnnotation,
  generateYogaAnnotation,
  generateKaranaAnnotation,
  generateRashiAnnotation,
  generatePanchangamSummary,
  generateVisualizationGuide,
  generateAllAnnotations,
  getTooltipContent,
  generateCalculationSteps,
} from './utils/annotationHelper';
