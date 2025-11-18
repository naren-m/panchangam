import { useEffect, useRef, useCallback } from 'react';
import * as THREE from 'three';

interface TouchState {
  touches: Map<number, { x: number; y: number; startX: number; startY: number }>;
  lastDistance: number;
  lastAngle: number;
}

interface MobileTouchControlsProps {
  camera: THREE.Camera;
  canvas: HTMLElement;
  enabled: boolean;
  onPinchZoom?: (delta: number) => void;
  onRotate?: (deltaX: number, deltaY: number) => void;
  onDoubleTap?: (x: number, y: number) => void;
}

export const useMobileTouchControls = ({
  camera,
  canvas,
  enabled,
  onPinchZoom,
  onRotate,
  onDoubleTap,
}: MobileTouchControlsProps) => {
  const touchStateRef = useRef<TouchState>({
    touches: new Map(),
    lastDistance: 0,
    lastAngle: 0,
  });

  const lastTapTimeRef = useRef<number>(0);
  const DOUBLE_TAP_DELAY = 300; // ms

  // Calculate distance between two touches
  const getTouchDistance = useCallback((touch1: Touch, touch2: Touch): number => {
    const dx = touch1.clientX - touch2.clientX;
    const dy = touch1.clientY - touch2.clientY;
    return Math.sqrt(dx * dx + dy * dy);
  }, []);

  // Calculate angle between two touches
  const getTouchAngle = useCallback((touch1: Touch, touch2: Touch): number => {
    return Math.atan2(touch2.clientY - touch1.clientY, touch2.clientX - touch1.clientX);
  }, []);

  // Handle touch start
  const handleTouchStart = useCallback(
    (event: TouchEvent) => {
      if (!enabled) return;

      event.preventDefault();

      const state = touchStateRef.current;
      const touches = Array.from(event.touches);

      // Record all touches
      touches.forEach((touch) => {
        state.touches.set(touch.identifier, {
          x: touch.clientX,
          y: touch.clientY,
          startX: touch.clientX,
          startY: touch.clientY,
        });
      });

      // Initialize pinch state for two-finger gestures
      if (touches.length === 2) {
        state.lastDistance = getTouchDistance(touches[0], touches[1]);
        state.lastAngle = getTouchAngle(touches[0], touches[1]);
      }

      // Check for double tap
      if (touches.length === 1) {
        const now = Date.now();
        if (now - lastTapTimeRef.current < DOUBLE_TAP_DELAY) {
          // Double tap detected
          if (onDoubleTap) {
            const rect = canvas.getBoundingClientRect();
            const x = touches[0].clientX - rect.left;
            const y = touches[0].clientY - rect.top;
            onDoubleTap(x, y);
          }
          lastTapTimeRef.current = 0;
        } else {
          lastTapTimeRef.current = now;
        }
      }
    },
    [enabled, canvas, getTouchDistance, getTouchAngle, onDoubleTap]
  );

  // Handle touch move
  const handleTouchMove = useCallback(
    (event: TouchEvent) => {
      if (!enabled) return;

      event.preventDefault();

      const state = touchStateRef.current;
      const touches = Array.from(event.touches);

      if (touches.length === 1) {
        // Single finger pan/rotate
        const touch = touches[0];
        const stored = state.touches.get(touch.identifier);

        if (stored && onRotate) {
          const deltaX = touch.clientX - stored.x;
          const deltaY = touch.clientY - stored.y;

          // Apply rotation with damping for smoother control
          const dampingFactor = 0.5;
          onRotate(deltaX * dampingFactor, deltaY * dampingFactor);

          // Update stored position
          stored.x = touch.clientX;
          stored.y = touch.clientY;
        }
      } else if (touches.length === 2) {
        // Two-finger pinch zoom and rotate
        const distance = getTouchDistance(touches[0], touches[1]);
        const angle = getTouchAngle(touches[0], touches[1]);

        if (state.lastDistance > 0) {
          // Pinch zoom
          const zoomDelta = (distance - state.lastDistance) * 0.01;
          if (onPinchZoom && Math.abs(zoomDelta) > 0.001) {
            onPinchZoom(zoomDelta);
          }

          // Two-finger rotation
          const angleDelta = angle - state.lastAngle;
          if (onRotate && Math.abs(angleDelta) > 0.01) {
            // Convert angle to rotation
            const rotationSpeed = 50;
            onRotate(Math.cos(angle) * angleDelta * rotationSpeed, Math.sin(angle) * angleDelta * rotationSpeed);
          }
        }

        state.lastDistance = distance;
        state.lastAngle = angle;

        // Update touch positions
        touches.forEach((touch) => {
          const stored = state.touches.get(touch.identifier);
          if (stored) {
            stored.x = touch.clientX;
            stored.y = touch.clientY;
          }
        });
      }
    },
    [enabled, getTouchDistance, getTouchAngle, onPinchZoom, onRotate]
  );

  // Handle touch end
  const handleTouchEnd = useCallback(
    (event: TouchEvent) => {
      if (!enabled) return;

      event.preventDefault();

      const state = touchStateRef.current;
      const endedTouches = Array.from(event.changedTouches);

      // Remove ended touches
      endedTouches.forEach((touch) => {
        state.touches.delete(touch.identifier);
      });

      // Reset pinch state if less than 2 touches remain
      if (event.touches.length < 2) {
        state.lastDistance = 0;
        state.lastAngle = 0;
      }
    },
    [enabled]
  );

  // Attach event listeners
  useEffect(() => {
    if (!enabled || !canvas) return;

    canvas.addEventListener('touchstart', handleTouchStart, { passive: false });
    canvas.addEventListener('touchmove', handleTouchMove, { passive: false });
    canvas.addEventListener('touchend', handleTouchEnd, { passive: false });
    canvas.addEventListener('touchcancel', handleTouchEnd, { passive: false });

    return () => {
      canvas.removeEventListener('touchstart', handleTouchStart);
      canvas.removeEventListener('touchmove', handleTouchMove);
      canvas.removeEventListener('touchend', handleTouchEnd);
      canvas.removeEventListener('touchcancel', handleTouchEnd);
    };
  }, [enabled, canvas, handleTouchStart, handleTouchMove, handleTouchEnd]);
};

// Mobile-optimized settings
export const getMobileOptimizedSettings = () => {
  const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
    navigator.userAgent
  );

  const isLowEnd =
    isMobile &&
    (navigator.hardwareConcurrency <= 4 ||
      (window.performance as any).memory?.jsHeapSizeLimit < 1000000000);

  return {
    isMobile,
    isLowEnd,
    pixelRatio: Math.min(window.devicePixelRatio, isMobile ? 2 : 3),
    targetFPS: isMobile ? 30 : 60,
    antialias: !isLowEnd,
    shadowMapEnabled: !isLowEnd,
    maxLights: isLowEnd ? 2 : 4,
    segments: isLowEnd ? 16 : 32,
    textureSize: isLowEnd ? 512 : 1024,
    starLimit: isLowEnd ? 1000 : 5000,
    enableBloom: !isLowEnd,
  };
};

// Haptic feedback for mobile devices
export const triggerHapticFeedback = (type: 'light' | 'medium' | 'heavy' = 'light') => {
  if ('vibrate' in navigator) {
    const patterns = {
      light: 10,
      medium: 20,
      heavy: 30,
    };
    navigator.vibrate(patterns[type]);
  }
};

// Touch-friendly UI size checker
export const getTouchFriendlySize = (baseSize: number): number => {
  const MIN_TOUCH_SIZE = 44; // iOS HIG recommendation in pixels
  return Math.max(baseSize, MIN_TOUCH_SIZE);
};

export default useMobileTouchControls;
