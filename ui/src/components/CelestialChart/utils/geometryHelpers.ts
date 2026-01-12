/**
 * Geometry helper functions for the Celestial Chart
 *
 * Handles polar-to-Cartesian conversions and SVG path generation.
 * Uses astronomical convention: 0° = Aries (East), counter-clockwise
 */

import type { Point, ChartDimensions } from '../types';

// ============================================================================
// Coordinate Conversion
// ============================================================================

/**
 * Convert degrees to radians
 */
export const degreesToRadians = (degrees: number): number => {
  return (degrees * Math.PI) / 180;
};

/**
 * Convert radians to degrees
 */
export const radiansToDegrees = (radians: number): number => {
  return (radians * 180) / Math.PI;
};

/**
 * Normalize angle to 0-360 range
 */
export const normalizeAngle = (degrees: number): number => {
  let normalized = degrees % 360;
  if (normalized < 0) normalized += 360;
  return normalized;
};

/**
 * Convert astronomical longitude to SVG angle
 *
 * Astronomical: 0° = Aries (East on chart), increases counter-clockwise
 * SVG: 0° = East (right), but we draw from top (North) clockwise
 *
 * For traditional representation with Aries at the "start" position:
 * - Rotate so 0° is at top (12 o'clock position)
 * - Draw counter-clockwise to match zodiac order
 */
export const longitudeToSvgAngle = (longitude: number): number => {
  // Start from top (subtract 90°), counter-clockwise (negate)
  // This places Aries at top and zodiac flows counter-clockwise
  return normalizeAngle(90 - longitude);
};

/**
 * Convert polar coordinates to Cartesian (SVG coordinates)
 *
 * @param center - Center point of the chart
 * @param radius - Distance from center
 * @param angleDegrees - Angle in degrees (already converted to SVG convention)
 */
export const polarToCartesian = (
  center: Point,
  radius: number,
  angleDegrees: number
): Point => {
  const angleRadians = degreesToRadians(angleDegrees);
  return {
    x: center.x + radius * Math.cos(angleRadians),
    y: center.y - radius * Math.sin(angleRadians), // SVG Y is inverted
  };
};

/**
 * Convert longitude directly to Cartesian coordinates
 */
export const longitudeToCartesian = (
  center: Point,
  radius: number,
  longitude: number
): Point => {
  const svgAngle = longitudeToSvgAngle(longitude);
  return polarToCartesian(center, radius, svgAngle);
};

// ============================================================================
// SVG Path Generation
// ============================================================================

/**
 * Generate SVG path data for an arc segment (used for ring segments)
 *
 * @param center - Center point
 * @param innerRadius - Inner radius of the arc
 * @param outerRadius - Outer radius of the arc
 * @param startLongitude - Start longitude in degrees
 * @param endLongitude - End longitude in degrees
 */
export const createArcSegmentPath = (
  center: Point,
  innerRadius: number,
  outerRadius: number,
  startLongitude: number,
  endLongitude: number
): string => {
  const startAngle = longitudeToSvgAngle(startLongitude);
  const endAngle = longitudeToSvgAngle(endLongitude);

  // Calculate the four corner points
  const outerStart = polarToCartesian(center, outerRadius, startAngle);
  const outerEnd = polarToCartesian(center, outerRadius, endAngle);
  const innerEnd = polarToCartesian(center, innerRadius, endAngle);
  const innerStart = polarToCartesian(center, innerRadius, startAngle);

  // Determine if we need the large arc flag
  // Since we're going counter-clockwise, check the angular difference
  let angleDiff = startAngle - endAngle;
  if (angleDiff < 0) angleDiff += 360;
  const largeArcFlag = angleDiff > 180 ? 1 : 0;

  // SVG arc sweeps counter-clockwise when sweep-flag is 0
  const sweepFlag = 0;

  return [
    `M ${outerStart.x} ${outerStart.y}`,
    `A ${outerRadius} ${outerRadius} 0 ${largeArcFlag} ${sweepFlag} ${outerEnd.x} ${outerEnd.y}`,
    `L ${innerEnd.x} ${innerEnd.y}`,
    `A ${innerRadius} ${innerRadius} 0 ${largeArcFlag} ${1 - sweepFlag} ${innerStart.x} ${innerStart.y}`,
    'Z'
  ].join(' ');
};

/**
 * Generate SVG path for a simple arc (used for tithi arc)
 */
export const createArcPath = (
  center: Point,
  radius: number,
  startLongitude: number,
  endLongitude: number
): string => {
  const startAngle = longitudeToSvgAngle(startLongitude);
  const endAngle = longitudeToSvgAngle(endLongitude);

  const start = polarToCartesian(center, radius, startAngle);
  const end = polarToCartesian(center, radius, endAngle);

  // Calculate angular difference for large arc flag
  let angleDiff = normalizeAngle(endLongitude - startLongitude);
  const largeArcFlag = angleDiff > 180 ? 1 : 0;

  return [
    `M ${start.x} ${start.y}`,
    `A ${radius} ${radius} 0 ${largeArcFlag} 0 ${end.x} ${end.y}`
  ].join(' ');
};

/**
 * Generate combined path for all pada boundaries (optimization)
 * Returns a single path element instead of 108 individual segments
 */
export const createPadaBoundariesPath = (
  center: Point,
  innerRadius: number,
  outerRadius: number
): string => {
  const paths: string[] = [];
  const degreesPerPada = 360 / 108;

  for (let i = 0; i < 108; i++) {
    const longitude = i * degreesPerPada;
    const angle = longitudeToSvgAngle(longitude);
    const inner = polarToCartesian(center, innerRadius, angle);
    const outer = polarToCartesian(center, outerRadius, angle);
    paths.push(`M ${inner.x} ${inner.y} L ${outer.x} ${outer.y}`);
  }

  return paths.join(' ');
};

/**
 * Generate combined path for nakshatra boundaries
 */
export const createNakshatraBoundariesPath = (
  center: Point,
  innerRadius: number,
  outerRadius: number
): string => {
  const paths: string[] = [];
  const degreesPerNakshatra = 360 / 27;

  for (let i = 0; i < 27; i++) {
    const longitude = i * degreesPerNakshatra;
    const angle = longitudeToSvgAngle(longitude);
    const inner = polarToCartesian(center, innerRadius, angle);
    const outer = polarToCartesian(center, outerRadius, angle);
    paths.push(`M ${inner.x} ${inner.y} L ${outer.x} ${outer.y}`);
  }

  return paths.join(' ');
};

/**
 * Generate combined path for rashi boundaries
 */
export const createRashiBoundariesPath = (
  center: Point,
  innerRadius: number,
  outerRadius: number
): string => {
  const paths: string[] = [];

  for (let i = 0; i < 12; i++) {
    const longitude = i * 30;
    const angle = longitudeToSvgAngle(longitude);
    const inner = polarToCartesian(center, innerRadius, angle);
    const outer = polarToCartesian(center, outerRadius, angle);
    paths.push(`M ${inner.x} ${inner.y} L ${outer.x} ${outer.y}`);
  }

  return paths.join(' ');
};

// ============================================================================
// Hit Detection
// ============================================================================

/**
 * Convert screen coordinates to longitude
 * Used for hit detection to determine which segment was clicked/hovered
 */
export const cartesianToLongitude = (
  point: Point,
  center: Point
): number => {
  const dx = point.x - center.x;
  const dy = center.y - point.y; // Invert Y for standard math convention

  // Calculate angle from center
  let angle = radiansToDegrees(Math.atan2(dy, dx));

  // Convert from SVG angle back to longitude
  // SVG: 0° = East, increases counter-clockwise
  // Longitude: 0° at top (after our rotation), increases counter-clockwise
  let longitude = normalizeAngle(90 - angle);

  return longitude;
};

/**
 * Calculate distance from center point
 */
export const distanceFromCenter = (point: Point, center: Point): number => {
  const dx = point.x - center.x;
  const dy = point.y - center.y;
  return Math.sqrt(dx * dx + dy * dy);
};

/**
 * Determine which ring a point falls into
 */
export const getRingFromPosition = (
  point: Point,
  dimensions: ChartDimensions
): 'pada' | 'nakshatra' | 'rashi' | 'earth' | 'outside' | null => {
  const distance = distanceFromCenter(point, dimensions.center);

  if (distance <= dimensions.earthRadius) {
    return 'earth';
  }
  if (distance >= dimensions.rings.pada.innerRadius &&
      distance <= dimensions.rings.pada.outerRadius) {
    return 'pada';
  }
  if (distance >= dimensions.rings.nakshatra.innerRadius &&
      distance <= dimensions.rings.nakshatra.outerRadius) {
    return 'nakshatra';
  }
  if (distance >= dimensions.rings.rashi.innerRadius &&
      distance <= dimensions.rings.rashi.outerRadius) {
    return 'rashi';
  }
  return 'outside';
};

/**
 * Get segment index from longitude for a given ring
 */
export const getSegmentIndex = (
  longitude: number,
  segmentCount: number
): number => {
  const degreesPerSegment = 360 / segmentCount;
  return Math.floor(normalizeAngle(longitude) / degreesPerSegment);
};

// ============================================================================
// Label Positioning
// ============================================================================

/**
 * Calculate position for a label at the center of a segment
 */
export const getLabelPosition = (
  center: Point,
  radius: number,
  startLongitude: number,
  endLongitude: number
): Point => {
  const midLongitude = startLongitude + (endLongitude - startLongitude) / 2;
  return longitudeToCartesian(center, radius, midLongitude);
};

/**
 * Calculate rotation angle for text along the arc
 * Returns the angle to rotate text so it's readable
 */
export const getLabelRotation = (longitude: number): number => {
  // Text should be horizontal or slightly angled
  const angle = longitudeToSvgAngle(longitude);

  // Keep text readable (not upside down)
  if (angle > 90 && angle < 270) {
    return angle + 180;
  }
  return angle;
};
