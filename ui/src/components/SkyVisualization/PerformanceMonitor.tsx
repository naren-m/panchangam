import React, { useEffect, useState, useRef } from 'react';

interface PerformanceMetrics {
  fps: number;
  frameTime: number;
  memoryUsage: number;
  objectCount: number;
  triangleCount: number;
  drawCalls: number;
}

interface PerformanceMonitorProps {
  enabled: boolean;
  targetFPS?: number;
  onPerformanceIssue?: (metrics: PerformanceMetrics) => void;
  className?: string;
}

export const PerformanceMonitor: React.FC<PerformanceMonitorProps> = ({
  enabled,
  targetFPS = 60,
  onPerformanceIssue,
  className = '',
}) => {
  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    fps: 0,
    frameTime: 0,
    memoryUsage: 0,
    objectCount: 0,
    triangleCount: 0,
    drawCalls: 0,
  });

  const frameCountRef = useRef(0);
  const lastTimeRef = useRef(performance.now());
  const fpsHistoryRef = useRef<number[]>([]);

  useEffect(() => {
    if (!enabled) return;

    let animationFrameId: number;
    const maxHistoryLength = 60; // Keep 1 second of history at 60fps

    const measurePerformance = () => {
      const now = performance.now();
      const deltaTime = now - lastTimeRef.current;

      frameCountRef.current++;

      // Update FPS every second
      if (deltaTime >= 1000) {
        const currentFPS = Math.round((frameCountRef.current * 1000) / deltaTime);
        const avgFrameTime = deltaTime / frameCountRef.current;

        // Add to history
        fpsHistoryRef.current.push(currentFPS);
        if (fpsHistoryRef.current.length > maxHistoryLength) {
          fpsHistoryRef.current.shift();
        }

        // Get memory usage if available
        let memoryUsage = 0;
        if ((performance as any).memory) {
          const memory = (performance as any).memory;
          memoryUsage = Math.round(memory.usedJSHeapSize / 1048576); // Convert to MB
        }

        const newMetrics: PerformanceMetrics = {
          fps: currentFPS,
          frameTime: avgFrameTime,
          memoryUsage,
          objectCount: 0, // To be set by external systems
          triangleCount: 0,
          drawCalls: 0,
        };

        setMetrics(newMetrics);

        // Check for performance issues
        if (currentFPS < targetFPS * 0.8 && onPerformanceIssue) {
          // FPS dropped below 80% of target
          onPerformanceIssue(newMetrics);
        }

        frameCountRef.current = 0;
        lastTimeRef.current = now;
      }

      animationFrameId = requestAnimationFrame(measurePerformance);
    };

    animationFrameId = requestAnimationFrame(measurePerformance);

    return () => {
      if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
      }
    };
  }, [enabled, targetFPS, onPerformanceIssue]);

  if (!enabled) return null;

  const getColorForFPS = (fps: number): string => {
    if (fps >= targetFPS * 0.9) return 'text-green-400';
    if (fps >= targetFPS * 0.7) return 'text-yellow-400';
    return 'text-red-400';
  };

  const getColorForMemory = (mb: number): string => {
    if (mb < 50) return 'text-green-400';
    if (mb < 100) return 'text-yellow-400';
    return 'text-red-400';
  };

  return (
    <div
      className={`bg-black bg-opacity-70 rounded px-3 py-2 font-mono text-xs ${className}`}
    >
      <div className="space-y-1">
        <div className="flex justify-between">
          <span className="text-gray-400">FPS:</span>
          <span className={getColorForFPS(metrics.fps)}>{metrics.fps}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-400">Frame:</span>
          <span className="text-white">{metrics.frameTime.toFixed(2)}ms</span>
        </div>
        {metrics.memoryUsage > 0 && (
          <div className="flex justify-between">
            <span className="text-gray-400">Memory:</span>
            <span className={getColorForMemory(metrics.memoryUsage)}>
              {metrics.memoryUsage}MB
            </span>
          </div>
        )}
      </div>
    </div>
  );
};

// Hook for throttling updates to avoid excessive re-renders
export const useThrottle = <T,>(value: T, delay: number): T => {
  const [throttledValue, setThrottledValue] = useState(value);
  const lastRun = useRef(Date.now());

  useEffect(() => {
    const handler = setTimeout(() => {
      if (Date.now() - lastRun.current >= delay) {
        setThrottledValue(value);
        lastRun.current = Date.now();
      }
    }, delay - (Date.now() - lastRun.current));

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return throttledValue;
};

// Hook for detecting device capabilities
export const useDeviceCapabilities = () => {
  const [capabilities, setCapabilities] = useState({
    isMobile: false,
    isLowEnd: false,
    supportsWebGL2: false,
    maxTextureSize: 0,
    devicePixelRatio: 1,
  });

  useEffect(() => {
    const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent
    );

    // Check WebGL2 support
    const canvas = document.createElement('canvas');
    const gl2 = canvas.getContext('webgl2');
    const supportsWebGL2 = !!gl2;

    let maxTextureSize = 0;
    if (gl2) {
      maxTextureSize = gl2.getParameter(gl2.MAX_TEXTURE_SIZE);
    } else {
      const gl = canvas.getContext('webgl');
      if (gl) {
        maxTextureSize = gl.getParameter(gl.MAX_TEXTURE_SIZE);
      }
    }

    // Detect low-end devices
    const isLowEnd =
      isMobile &&
      (navigator.hardwareConcurrency <= 4 ||
        (performance as any).memory?.jsHeapSizeLimit < 1000000000);

    setCapabilities({
      isMobile,
      isLowEnd,
      supportsWebGL2,
      maxTextureSize,
      devicePixelRatio: window.devicePixelRatio || 1,
    });
  }, []);

  return capabilities;
};

// Performance optimization helpers
export const PerformanceUtils = {
  // Object culling based on visibility and distance
  shouldCullObject: (
    objectAltitude: number,
    objectDistance: number,
    maxDistance: number = 50
  ): boolean => {
    // Cull objects below horizon or too far away
    return objectAltitude < -5 || objectDistance > maxDistance;
  },

  // Level of detail based on distance and device capabilities
  getLODLevel: (distance: number, isLowEnd: boolean): 'high' | 'medium' | 'low' => {
    if (isLowEnd) {
      return 'low';
    }
    if (distance < 10) {
      return 'high';
    }
    if (distance < 30) {
      return 'medium';
    }
    return 'low';
  },

  // Geometry detail based on LOD
  getGeometrySegments: (lod: 'high' | 'medium' | 'low'): { segments: number; rings: number } => {
    switch (lod) {
      case 'high':
        return { segments: 32, rings: 32 };
      case 'medium':
        return { segments: 16, rings: 16 };
      case 'low':
        return { segments: 8, rings: 8 };
    }
  },

  // Adaptive quality settings based on FPS
  getAdaptiveQuality: (currentFPS: number, targetFPS: number): number => {
    const ratio = currentFPS / targetFPS;
    if (ratio >= 0.95) return 1.0; // Full quality
    if (ratio >= 0.85) return 0.9;
    if (ratio >= 0.75) return 0.8;
    if (ratio >= 0.65) return 0.7;
    return 0.6; // Minimum quality
  },

  // Batch similar objects for single draw call
  batchObjectsByMaterial: <T extends { color: string; size: number }>(
    objects: T[]
  ): Map<string, T[]> => {
    const batches = new Map<string, T[]>();

    for (const obj of objects) {
      const key = `${obj.color}-${obj.size}`;
      if (!batches.has(key)) {
        batches.set(key, []);
      }
      batches.get(key)!.push(obj);
    }

    return batches;
  },

  // Interpolate planetary positions for smooth motion
  interpolatePosition: (
    start: number,
    end: number,
    progress: number
  ): number => {
    // Use easing for smooth motion
    const eased = progress < 0.5
      ? 2 * progress * progress
      : 1 - Math.pow(-2 * progress + 2, 2) / 2;
    return start + (end - start) * eased;
  },
};

export default PerformanceMonitor;
