/**
 * Type definitions for the Celestial Chart (2D Geocentric Circular Visualization)
 *
 * The chart displays Earth at center with concentric rings showing:
 * - Outer: 12 Rashis (30° each)
 * - Middle: 27 Nakshatras (13°20' each)
 * - Inner: 108 Padas (3°20' each, 4 per Nakshatra)
 * - Sun and Moon positions with Tithi arc
 */

import type {
  TithiInfo,
  NakshatraInfo,
  RashiInfo,
  PanchangamElements
} from '../EclipticBeltVisualization/types/eclipticBelt';

// Re-export for convenience
export type { TithiInfo, NakshatraInfo, RashiInfo, PanchangamElements };

// ============================================================================
// Chart Geometry Types
// ============================================================================

/**
 * Point in 2D Cartesian coordinates
 */
export interface Point {
  x: number;
  y: number;
}

/**
 * Arc segment for ring visualization
 */
export interface ArcSegment {
  id: string;
  startAngle: number;  // in degrees, 0° = East (right), counter-clockwise
  endAngle: number;
  innerRadius: number;
  outerRadius: number;
  label: string;
  color: string;
}

/**
 * Ring configuration
 */
export interface RingConfig {
  innerRadius: number;
  outerRadius: number;
  segments: number;
  labelOffset?: number;
}

// ============================================================================
// Chart Dimension Types
// ============================================================================

/**
 * Chart dimensions and layout configuration
 */
export interface ChartDimensions {
  size: number;           // Total chart size (width = height)
  center: Point;          // Center point of the chart
  earthRadius: number;    // Radius of Earth at center
  rings: {
    pada: RingConfig;
    nakshatra: RingConfig;
    rashi: RingConfig;
  };
  celestialOrbit: number; // Radius where Sun/Moon orbit
  margin: number;         // Outer margin
}

/**
 * Default chart dimensions factory
 */
export const createDefaultDimensions = (size: number): ChartDimensions => {
  const center = { x: size / 2, y: size / 2 };
  const usableRadius = (size / 2) - 40; // 40px margin

  return {
    size,
    center,
    earthRadius: usableRadius * 0.08,
    rings: {
      pada: {
        innerRadius: usableRadius * 0.20,
        outerRadius: usableRadius * 0.35,
        segments: 108,
      },
      nakshatra: {
        innerRadius: usableRadius * 0.35,
        outerRadius: usableRadius * 0.60,
        segments: 27,
        labelOffset: usableRadius * 0.475,
      },
      rashi: {
        innerRadius: usableRadius * 0.60,
        outerRadius: usableRadius * 0.85,
        segments: 12,
        labelOffset: usableRadius * 0.725,
      },
    },
    celestialOrbit: usableRadius * 0.15,
    margin: 40,
  };
};

// ============================================================================
// Celestial Body Types
// ============================================================================

/**
 * Celestial body position and display info
 */
export interface CelestialBodyInfo {
  type: 'sun' | 'moon';
  longitude: number;      // 0-360 degrees
  position: Point;        // Screen coordinates
  symbol: string;
  color: string;
  size: number;
  label: string;
}

/**
 * Tithi arc connecting Sun and Moon
 */
export interface TithiArcInfo {
  sunAngle: number;
  moonAngle: number;
  radius: number;
  tithi: TithiInfo;
  arcPath: string;        // SVG path data
}

// ============================================================================
// Interaction Types
// ============================================================================

/**
 * Hover/selection state for chart elements
 */
export interface ChartInteractionState {
  hoveredElement: HoveredElement | null;
  selectedElement: SelectedElement | null;
}

/**
 * Hovered element details
 */
export interface HoveredElement {
  type: 'rashi' | 'nakshatra' | 'pada' | 'sun' | 'moon' | 'tithi';
  id: string;
  data: RashiInfo | NakshatraInfo | PadaInfo | CelestialBodyInfo | TithiInfo;
  position: Point;  // For tooltip positioning
}

/**
 * Selected element details
 */
export interface SelectedElement {
  type: 'rashi' | 'nakshatra' | 'pada' | 'sun' | 'moon' | 'tithi';
  id: string;
  data: RashiInfo | NakshatraInfo | PadaInfo | CelestialBodyInfo | TithiInfo;
}

/**
 * Pada (quarter of a Nakshatra) info
 */
export interface PadaInfo {
  number: number;         // 1-108 (global pada number)
  padaInNakshatra: number; // 1-4 (pada within its nakshatra)
  nakshatra: NakshatraInfo;
  startDegree: number;
  endDegree: number;
  navamsha: string;       // Navamsha sign for this pada
}

// ============================================================================
// Component Props Types
// ============================================================================

/**
 * Props for the main CelestialChart container
 */
export interface CelestialChartProps {
  date: Date;
  latitude: number;
  longitude: number;
  timezone?: string;
  panchangamData?: Record<string, unknown>;
  className?: string;
}

/**
 * Props for CelestialChartSVG renderer
 */
export interface CelestialChartSVGProps {
  dimensions: ChartDimensions;
  panchangam: PanchangamElements;
  interactionState: ChartInteractionState;
  onElementHover: (element: HoveredElement | null) => void;
  onElementSelect: (element: SelectedElement | null) => void;
}

/**
 * Props for ring components
 */
export interface RingProps {
  dimensions: ChartDimensions;
  hoveredId: string | null;
  selectedId: string | null;
  onHover: (id: string | null, position: Point) => void;
  onSelect: (id: string) => void;
}

/**
 * Extended props for rings that need position data
 */
export interface RashiRingProps extends RingProps {
  sunRashi: RashiInfo;
  moonRashi: RashiInfo;
}

export interface NakshatraRingProps extends RingProps {
  currentNakshatra: NakshatraInfo;
}

export interface PadaRingProps extends RingProps {
  currentNakshatra: NakshatraInfo;
  currentPada: number;
}

/**
 * Props for celestial body markers
 */
export interface CelestialBodyProps {
  body: CelestialBodyInfo;
  isHovered: boolean;
  onHover: (hovered: boolean, position: Point) => void;
  onClick: () => void;
}

/**
 * Props for tithi arc
 */
export interface TithiArcProps {
  arcInfo: TithiArcInfo;
  isHovered: boolean;
  onHover: (hovered: boolean, position: Point) => void;
  onClick: () => void;
}

/**
 * Props for tooltip
 */
export interface ChartTooltipProps {
  element: HoveredElement;
  chartDimensions: ChartDimensions;
}

// ============================================================================
// Animation Types
// ============================================================================

/**
 * Animation configuration
 */
export interface AnimationConfig {
  enabled: boolean;
  duration: number;      // milliseconds
  easing: 'linear' | 'ease-in' | 'ease-out' | 'ease-in-out';
}

// ============================================================================
// Constants
// ============================================================================

export const DEGREES_PER_PADA = 3.333333;  // 360 / 108
export const PADAS_PER_NAKSHATRA = 4;
export const TOTAL_PADAS = 108;

/**
 * Navamsha signs for each pada (1-108)
 * Each pada corresponds to a navamsha in the D9 chart
 */
export const NAVAMSHA_SEQUENCE = [
  'Mesha', 'Vrishabha', 'Mithuna', 'Karka', 'Simha', 'Kanya',
  'Tula', 'Vrishchika', 'Dhanu', 'Makara', 'Kumbha', 'Meena'
];

/**
 * Color palette for the chart
 */
export const CHART_COLORS = {
  earth: '#2E86AB',
  sun: '#F9A825',
  moon: '#ECEFF1',
  moonShadow: '#90A4AE',
  tithiArc: '#FF7043',
  rashiRing: {
    fire: '#FF5252',     // Aries, Leo, Sagittarius
    earth: '#8BC34A',    // Taurus, Virgo, Capricorn
    air: '#03A9F4',      // Gemini, Libra, Aquarius
    water: '#7C4DFF',    // Cancer, Scorpio, Pisces
  },
  nakshatraRing: '#FFB74D',
  padaRing: '#B0BEC5',
  hover: '#FFF176',
  selected: '#4DD0E1',
  text: '#37474F',
  textLight: '#78909C',
};
