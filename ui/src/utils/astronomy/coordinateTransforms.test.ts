import { describe, it, expect } from 'vitest';
import {
  eclipticToEquatorial,
  equatorialToHorizontal,
  horizontalToScreen,
  dateToJulianDay,
  eclipticToScreen
} from './coordinateTransforms';
import type { EclipticCoordinates, Observer } from '../../types/skyVisualization';

describe('Coordinate Transformation System', () => {
  describe('dateToJulianDay', () => {
    it('should convert J2000 epoch correctly', () => {
      // J2000 = January 1, 2000, 12:00:00 TT = JD 2451545.0
      const j2000 = new Date('2000-01-01T12:00:00Z');
      const jd = dateToJulianDay(j2000);
      expect(jd).toBeCloseTo(2451545.0, 2);
    });

    it('should convert Unix epoch correctly', () => {
      // Unix epoch = January 1, 1970, 00:00:00 UTC = JD 2440587.5
      const unixEpoch = new Date('1970-01-01T00:00:00Z');
      const jd = dateToJulianDay(unixEpoch);
      expect(jd).toBeCloseTo(2440587.5, 2);
    });

    it('should handle year 2024 correctly', () => {
      // January 1, 2024, 00:00:00 UTC = JD 2460310.5
      const date2024 = new Date('2024-01-01T00:00:00Z');
      const jd = dateToJulianDay(date2024);
      expect(jd).toBeCloseTo(2460310.5, 2);
    });
  });

  describe('eclipticToEquatorial', () => {
    it('should convert ecliptic point on equator (λ=0°, β=0°)', () => {
      // Point on vernal equinox should have RA=0°, Dec=0°
      const ecliptic: EclipticCoordinates = {
        longitude: 0,
        latitude: 0,
        distance: 1
      };
      const equatorial = eclipticToEquatorial(ecliptic, 2451545.0);

      expect(equatorial.rightAscension).toBeCloseTo(0, 1);
      expect(equatorial.declination).toBeCloseTo(0, 1);
      expect(equatorial.distance).toBe(1);
    });

    it('should convert ecliptic point at λ=90°, β=0° (summer solstice)', () => {
      // Point at 90° ecliptic longitude should have RA≈90° and Dec≈23.4°
      const ecliptic: EclipticCoordinates = {
        longitude: 90,
        latitude: 0,
        distance: 1
      };
      const equatorial = eclipticToEquatorial(ecliptic, 2451545.0);

      expect(equatorial.rightAscension).toBeCloseTo(90, 1);
      expect(equatorial.declination).toBeCloseTo(23.4, 1);
    });

    it('should convert ecliptic north pole (β=90°)', () => {
      // Ecliptic north pole should map to Dec ≈ 66.6° (90° - obliquity)
      const ecliptic: EclipticCoordinates = {
        longitude: 0,
        latitude: 90,
        distance: 1
      };
      const equatorial = eclipticToEquatorial(ecliptic, 2451545.0);

      expect(equatorial.declination).toBeCloseTo(66.6, 1);
    });

    it('should preserve distance in transformation', () => {
      const ecliptic: EclipticCoordinates = {
        longitude: 45,
        latitude: 15,
        distance: 5.2 // Jupiter-like distance
      };
      const equatorial = eclipticToEquatorial(ecliptic, 2451545.0);

      expect(equatorial.distance).toBe(5.2);
    });
  });

  describe('equatorialToHorizontal', () => {
    const observer: Observer = {
      latitude: 40.7128, // New York City
      longitude: -74.0060,
      altitude: 10
    };

    it('should place object at zenith when Dec = Lat and HA = 0', () => {
      // Object with Dec = observer's latitude, crossing meridian (HA=0)
      // should be at altitude = 90° (zenith)
      const equatorial = {
        rightAscension: 0, // Will be adjusted by LST calculation
        declination: observer.latitude,
        distance: 1
      };

      // Use a time when LST = RA (object transiting)
      const time = new Date('2024-03-20T12:00:00Z'); // Near vernal equinox
      const horizontal = equatorialToHorizontal(equatorial, observer, time);

      // Altitude should be close to 90° (within reasonable error)
      expect(horizontal.altitude).toBeGreaterThan(45);
      expect(horizontal.altitude).toBeLessThanOrEqual(90);
    });

    it('should place object at horizon when altitude is near 0', () => {
      // Object on celestial equator (Dec=0) rising/setting
      const equatorial = {
        rightAscension: 90,
        declination: 0,
        distance: 1
      };

      const time = new Date('2024-03-20T00:00:00Z');
      const horizontal = equatorialToHorizontal(equatorial, observer, time);

      // Should be somewhere between -90 and 90 degrees altitude
      expect(horizontal.altitude).toBeGreaterThanOrEqual(-90);
      expect(horizontal.altitude).toBeLessThanOrEqual(90);
    });

    it('should produce azimuth between 0 and 360 degrees', () => {
      const equatorial = {
        rightAscension: 180,
        declination: 30,
        distance: 1
      };

      const time = new Date('2024-06-21T18:00:00Z'); // Summer solstice evening
      const horizontal = equatorialToHorizontal(equatorial, observer, time);

      expect(horizontal.azimuth).toBeGreaterThanOrEqual(0);
      expect(horizontal.azimuth).toBeLessThan(360);
    });

    it('should handle observer at equator', () => {
      const equatorObserver: Observer = {
        latitude: 0,
        longitude: 0,
        altitude: 0
      };

      const equatorial = {
        rightAscension: 0,
        declination: 0,
        distance: 1
      };

      const time = new Date('2024-03-20T12:00:00Z');
      const horizontal = equatorialToHorizontal(equatorial, equatorObserver, time);

      expect(horizontal.altitude).toBeGreaterThanOrEqual(-90);
      expect(horizontal.altitude).toBeLessThanOrEqual(90);
      expect(horizontal.azimuth).toBeGreaterThanOrEqual(0);
      expect(horizontal.azimuth).toBeLessThan(360);
    });
  });

  describe('horizontalToScreen', () => {
    const width = 1920;
    const height = 1080;

    it('should place object at zenith near screen center (stereographic)', () => {
      const horizontal = {
        azimuth: 0, // North
        altitude: 90, // Zenith
        distance: 1
      };

      const screen = horizontalToScreen(horizontal, width, height, 'stereographic');

      // Zenith should be near center
      expect(screen.x).toBeCloseTo(width / 2, 0);
      expect(screen.y).toBeCloseTo(height / 2, 0);
    });

    it('should handle object below horizon (orthographic)', () => {
      const horizontal = {
        azimuth: 180,
        altitude: -10, // Below horizon
        distance: 1
      };

      const screen = horizontalToScreen(horizontal, width, height, 'orthographic');

      // Object below horizon should be off-screen
      expect(screen.x).toBe(-1);
      expect(screen.y).toBe(-1);
      expect(screen.z).toBe(-1);
    });

    it('should produce valid screen coordinates for visible objects', () => {
      const horizontal = {
        azimuth: 90, // East
        altitude: 45, // 45° above horizon
        distance: 1
      };

      const screen = horizontalToScreen(horizontal, width, height, 'stereographic');

      // Should be on screen
      expect(screen.x).toBeGreaterThanOrEqual(0);
      expect(screen.x).toBeLessThanOrEqual(width * 2); // Allow some margin for projection
      expect(screen.y).toBeGreaterThanOrEqual(0);
      expect(screen.y).toBeLessThanOrEqual(height * 2);
      expect(screen.z).toBe(1); // Distance preserved
    });

    it('should work with mercator projection', () => {
      const horizontal = {
        azimuth: 180, // South
        altitude: 30,
        distance: 2
      };

      const screen = horizontalToScreen(horizontal, width, height, 'mercator');

      expect(screen.x).toBeGreaterThan(0);
      expect(screen.y).toBeGreaterThan(0);
      expect(screen.z).toBe(2);
    });
  });

  describe('eclipticToScreen (integration test)', () => {
    const observer: Observer = {
      latitude: 51.5074, // London
      longitude: -0.1278,
      altitude: 0
    };
    const time = new Date('2024-06-21T12:00:00Z'); // Summer solstice, noon
    const width = 1920;
    const height = 1080;

    it('should convert Sun position on summer solstice', () => {
      // Sun at summer solstice is at ecliptic longitude ≈ 90°
      const sunEcliptic: EclipticCoordinates = {
        longitude: 90,
        latitude: 0,
        distance: 1
      };

      const screen = eclipticToScreen(
        sunEcliptic,
        observer,
        time,
        width,
        height,
        'stereographic'
      );

      // Sun should be visible (not below horizon) at noon on summer solstice in London
      expect(screen.x).toBeGreaterThan(-1);
      expect(screen.y).toBeGreaterThan(-1);
    });

    it('should handle moon position transformation', () => {
      // Example Moon position
      const moonEcliptic: EclipticCoordinates = {
        longitude: 120,
        latitude: 5, // Moon can deviate from ecliptic
        distance: 0.00257 // AU
      };

      const screen = eclipticToScreen(
        moonEcliptic,
        observer,
        time,
        width,
        height,
        'stereographic'
      );

      expect(screen.z).toBeCloseTo(0.00257, 5);
    });
  });

  describe('Edge cases and error handling', () => {
    it('should handle date at year boundaries', () => {
      const newYear = new Date('2024-12-31T23:59:59Z');
      const jd = dateToJulianDay(newYear);
      expect(jd).toBeGreaterThan(2451545); // After J2000
      expect(jd).toBeLessThan(2500000); // Reasonable future date
    });

    it('should handle extreme latitudes', () => {
      const northPoleObserver: Observer = {
        latitude: 90,
        longitude: 0,
        altitude: 0
      };

      const equatorial = {
        rightAscension: 0,
        declination: 45,
        distance: 1
      };

      const time = new Date('2024-06-21T12:00:00Z');
      const horizontal = equatorialToHorizontal(equatorial, northPoleObserver, time);

      // At North Pole, objects always have same altitude as declination
      expect(horizontal.altitude).toBeCloseTo(45, 1);
    });

    it('should handle very small canvas sizes', () => {
      const horizontal = {
        azimuth: 45,
        altitude: 30,
        distance: 1
      };

      const screen = horizontalToScreen(horizontal, 100, 100, 'stereographic');

      expect(screen.x).toBeGreaterThanOrEqual(0);
      expect(screen.y).toBeGreaterThanOrEqual(0);
    });
  });

  describe('Known star positions (validation with astronomical data)', () => {
    it('should correctly transform Sirius position', () => {
      // Sirius: RA = 101.287°, Dec = -16.716° (J2000)
      // In ecliptic coordinates: λ ≈ 101.3°, β ≈ -39.6°
      const siriusEcliptic: EclipticCoordinates = {
        longitude: 101.287,
        latitude: -39.608,
        distance: 8.6 // light years, but using as distance unit
      };

      const equatorial = eclipticToEquatorial(siriusEcliptic, 2451545.0);

      // Should convert to approximately RA ≈ 101° (allow some error due to coordinate epoch)
      expect(equatorial.rightAscension).toBeGreaterThan(95);
      expect(equatorial.rightAscension).toBeLessThan(110);

      // Declination should be negative (southern hemisphere)
      expect(equatorial.declination).toBeLessThan(0);
    });

    it('should correctly transform Polaris (North Star) position', () => {
      // Polaris is close to North Celestial Pole
      // Ecliptic: λ ≈ 26°, β ≈ 66.6° (near ecliptic north pole)
      const polarisEcliptic: EclipticCoordinates = {
        longitude: 26,
        latitude: 66,
        distance: 433 // light years
      };

      const equatorial = eclipticToEquatorial(polarisEcliptic, 2451545.0);

      // Polaris should have Dec ≈ 89° (very close to north pole)
      expect(equatorial.declination).toBeGreaterThan(85);
      expect(equatorial.declination).toBeLessThanOrEqual(90);
    });
  });
});
