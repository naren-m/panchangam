import React, { useRef, useEffect, useMemo, useState } from 'react';
import {
  BackSide,
  BufferGeometry,
  Color,
  Material,
  Mesh,
  MeshBasicMaterial,
  MeshPhongMaterial,
  Object3D,
  PerspectiveCamera,
  Scene,
  SphereGeometry,
  WebGLRenderer,
} from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import { 
  SkySphereConfig, 
  CelestialObject, 
  RenderOptions,
  Observer,
  TimeConfig 
} from '../../types/skyVisualization';
import NakshatraVisualization from './NakshatraVisualization';
import {
  addCelestialEquator,
  addCoordinateGrids,
  addEclipticLine,
  addHorizonLine,
} from './skySphereGuides';
import ZodiacVisualization from './ZodiacVisualization';

interface SkySphereProps {
  config?: Partial<SkySphereConfig>;
  celestialObjects?: CelestialObject[];
  observer: Observer;
  timeConfig?: TimeConfig;
  renderOptions?: Partial<RenderOptions>;
  currentNakshatra?: number; // 1-27, current nakshatra to highlight
  currentRashi?: number; // 1-12, current zodiac sign to highlight
  onError?: (error: Error) => void;
  className?: string;
}

const defaultConfig: SkySphereConfig = {
  radius: 100,
  segments: 64,
  rings: 64,
  projection: 'stereographic',
  coordinateSystem: 'equatorial',
  renderMode: 'webgl'
};

const defaultRenderOptions: RenderOptions = {
  showGrid: true,
  showConstellations: true,
  showNakshatras: true,
  showPlanets: true,
  showStars: true,
  showLabels: false,
  showZodiac: true,
  showEcliptic: true,
  showEquator: true,
  showHorizon: true,
  starMagnitudeLimit: 6.0,
  labelMinZoom: 2.0
};

export const SkySphere: React.FC<SkySphereProps> = ({
  config = {},
  celestialObjects = [],
  observer,
  renderOptions = {},
  currentNakshatra,
  currentRashi,
  onError,
  className
}) => {
  const mountRef = useRef<HTMLDivElement>(null);
  const sceneRef = useRef<Scene | null>(null);
  const rendererRef = useRef<WebGLRenderer | null>(null);
  const cameraRef = useRef<PerspectiveCamera | null>(null);
  const controlsRef = useRef<OrbitControls | null>(null);
  const frameIdRef = useRef<number>(0);
  
  const [isWebGLSupported, setIsWebGLSupported] = useState(true);
  const [isLoading, setIsLoading] = useState(true);
  
  const mergedConfig = useMemo(
    () => ({ ...defaultConfig, ...config }),
    [config]
  );
  const mergedRenderOptions = useMemo(
    () => ({ ...defaultRenderOptions, ...renderOptions }),
    [renderOptions]
  );

  // Check WebGL support
  useEffect(() => {
    try {
      const canvas = document.createElement('canvas');
      const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
      if (!gl) {
        setIsWebGLSupported(false);
        onError?.(new Error('WebGL is not supported in this browser'));
      }
    } catch {
      setIsWebGLSupported(false);
      onError?.(new Error('WebGL check failed'));
    }
  }, [onError]);

  // Initialize Three.js scene
  useEffect(() => {
    const mountNode = mountRef.current;
    if (!mountNode || !isWebGLSupported) return;

    const width = mountNode.clientWidth;
    const height = mountNode.clientHeight;

    // Scene setup
    const scene = new Scene();
    scene.background = new Color(0x000000);
    sceneRef.current = scene;

    // Camera setup
    const camera = new PerspectiveCamera(
      75, // FOV
      width / height, // Aspect ratio
      0.1, // Near
      1000 // Far
    );
    camera.position.set(0, 0, 150);
    cameraRef.current = camera;

    // Renderer setup
    const renderer = new WebGLRenderer({
      antialias: true,
      alpha: true
    });
    renderer.setSize(width, height);
    renderer.setPixelRatio(window.devicePixelRatio);
    mountNode.appendChild(renderer.domElement);
    rendererRef.current = renderer;

    // Controls setup
    const controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = true;
    controls.dampingFactor = 0.05;
    controls.enablePan = false;
    controls.minDistance = 50;
    controls.maxDistance = 300;
    controlsRef.current = controls;

    // Create sky sphere
    const sphereGeometry = new SphereGeometry(
      mergedConfig.radius,
      mergedConfig.segments,
      mergedConfig.rings
    );
    
    // Invert the sphere so we see it from inside
    sphereGeometry.scale(-1, 1, 1);
    
    const sphereMaterial = new MeshBasicMaterial({
      color: 0x001122,
      side: BackSide
    });
    
    const skySphere = new Mesh(sphereGeometry, sphereMaterial);
    scene.add(skySphere);

    // Add coordinate grids
    if (mergedRenderOptions.showGrid) {
      addCoordinateGrids(scene, mergedConfig.radius);
    }

    // Add ecliptic line
    if (mergedRenderOptions.showEcliptic) {
      addEclipticLine(scene, mergedConfig.radius);
    }

    // Add celestial equator
    if (mergedRenderOptions.showEquator) {
      addCelestialEquator(scene, mergedConfig.radius);
    }

    // Add horizon line
    if (mergedRenderOptions.showHorizon) {
      addHorizonLine(scene, mergedConfig.radius, observer.latitude);
    }

    setIsLoading(false);

    // Handle window resize
    const handleResize = () => {
      if (!camera || !renderer) return;
      
      const newWidth = mountNode.clientWidth;
      const newHeight = mountNode.clientHeight;
      
      camera.aspect = newWidth / newHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(newWidth, newHeight);
    };

    window.addEventListener('resize', handleResize);

    // Animation loop
    const animate = () => {
      frameIdRef.current = requestAnimationFrame(animate);
      
      if (controls) {
        controls.update();
      }
      
      if (renderer && scene && camera) {
        renderer.render(scene, camera);
      }
    };
    
    animate();

    // Cleanup
    return () => {
      window.removeEventListener('resize', handleResize);
      
      if (frameIdRef.current) {
        cancelAnimationFrame(frameIdRef.current);
      }
      
      if (renderer.domElement.parentElement === mountNode) {
        mountNode.removeChild(renderer.domElement);
      }
      
      renderer?.dispose();
      controls?.dispose();
    };
  }, [isWebGLSupported, mergedConfig, mergedRenderOptions, observer.latitude]);

  // Update celestial objects
  useEffect(() => {
    if (!sceneRef.current || !celestialObjects.length) return;

    // Remove existing celestial objects
    const objectsToRemove: Object3D[] = [];
    sceneRef.current.traverse((child) => {
      if (child.userData.isCelestialObject) {
        objectsToRemove.push(child);
      }
    });
    objectsToRemove.forEach(obj => sceneRef.current!.remove(obj));

    // Add celestial objects
    celestialObjects.forEach(obj => {
      if (obj.coordinates.ecliptic) {
        // Convert ecliptic to 3D position
        const phi = (90 - obj.coordinates.ecliptic.latitude) * Math.PI / 180;
        const theta = obj.coordinates.ecliptic.longitude * Math.PI / 180;
        
        const x = mergedConfig.radius * Math.sin(phi) * Math.cos(theta);
        const y = mergedConfig.radius * Math.cos(phi);
        const z = mergedConfig.radius * Math.sin(phi) * Math.sin(theta);

        // Create object based on type
        let geometry: BufferGeometry;
        let material: Material;
        
        if (obj.type === 'star') {
          const size = obj.size || (6 - (obj.magnitude || 0)) * 0.5;
          geometry = new SphereGeometry(size, 8, 8);
          material = new MeshBasicMaterial({
            color: obj.color || 0xffffff
          });
        } else if (obj.type === 'planet') {
          const size = obj.size || 2;
          geometry = new SphereGeometry(size, 16, 16);
          material = new MeshPhongMaterial({
            color: obj.color || 0xffaa00,
            emissive: obj.color || 0xffaa00,
            emissiveIntensity: 0.3
          });
        } else {
          // Default for other objects
          geometry = new SphereGeometry(1, 8, 8);
          material = new MeshBasicMaterial({
            color: obj.color || 0xffffff
          });
        }

        const mesh = new Mesh(geometry, material);
        mesh.position.set(x, y, z);
        mesh.userData = { isCelestialObject: true, ...obj };
        
        sceneRef.current!.add(mesh);
      }
    });
  }, [celestialObjects, mergedConfig.radius]);

  if (!isWebGLSupported) {
    return (
      <div className={`flex items-center justify-center h-full bg-gray-900 text-white ${className || ''}`}>
        <div className="text-center p-8">
          <h3 className="text-xl font-semibold mb-2">WebGL Not Supported</h3>
          <p className="text-gray-400">
            Your browser doesn't support WebGL, which is required for 3D sky visualization.
            Please try a modern browser like Chrome, Firefox, or Edge.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative w-full h-full ${className || ''}`}>
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-900">
          <div className="text-white">Loading sky visualization...</div>
        </div>
      )}
      <div ref={mountRef} className="w-full h-full" />
      
      {/* Nakshatra Visualization */}
      {sceneRef.current && mergedRenderOptions.showNakshatras && (
        <NakshatraVisualization
          scene={sceneRef.current}
          radius={mergedConfig.radius}
          showLabels={mergedRenderOptions.showLabels}
          showBoundaries={true}
          currentNakshatra={currentNakshatra}
        />
      )}
      
      {/* Zodiac Visualization */}
      {sceneRef.current && mergedRenderOptions.showZodiac && (
        <ZodiacVisualization
          scene={sceneRef.current}
          radius={mergedConfig.radius + 2} // Slightly larger radius than nakshatras
          showLabels={mergedRenderOptions.showLabels}
          showBoundaries={true}
          currentRashi={currentRashi}
        />
      )}
    </div>
  );
};

export default SkySphere;
