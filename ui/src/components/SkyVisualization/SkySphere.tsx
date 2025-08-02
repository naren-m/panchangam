import React, { useRef, useEffect, useState } from 'react';
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
import { 
  SkySphereConfig, 
  CelestialObject, 
  RenderOptions,
  Observer,
  TimeConfig 
} from '../../types/skyVisualization';
import { eclipticToScreen } from '../../utils/astronomy/coordinateTransforms';
import NakshatraVisualization from './NakshatraVisualization';
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
  timeConfig = { date: new Date(), speed: 1, paused: false },
  renderOptions = {},
  currentNakshatra,
  currentRashi,
  onError,
  className
}) => {
  const mountRef = useRef<HTMLDivElement>(null);
  const sceneRef = useRef<THREE.Scene | null>(null);
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null);
  const cameraRef = useRef<THREE.PerspectiveCamera | null>(null);
  const controlsRef = useRef<OrbitControls | null>(null);
  const frameIdRef = useRef<number>(0);
  
  const [isWebGLSupported, setIsWebGLSupported] = useState(true);
  const [isLoading, setIsLoading] = useState(true);
  
  const mergedConfig = { ...defaultConfig, ...config };
  const mergedRenderOptions = { ...defaultRenderOptions, ...renderOptions };

  // Check WebGL support
  useEffect(() => {
    try {
      const canvas = document.createElement('canvas');
      const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
      if (!gl) {
        setIsWebGLSupported(false);
        onError?.(new Error('WebGL is not supported in this browser'));
      }
    } catch (e) {
      setIsWebGLSupported(false);
      onError?.(new Error('WebGL check failed'));
    }
  }, [onError]);

  // Initialize Three.js scene
  useEffect(() => {
    if (!mountRef.current || !isWebGLSupported) return;

    const width = mountRef.current.clientWidth;
    const height = mountRef.current.clientHeight;

    // Scene setup
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x000000);
    sceneRef.current = scene;

    // Camera setup
    const camera = new THREE.PerspectiveCamera(
      75, // FOV
      width / height, // Aspect ratio
      0.1, // Near
      1000 // Far
    );
    camera.position.set(0, 0, 150);
    cameraRef.current = camera;

    // Renderer setup
    const renderer = new THREE.WebGLRenderer({ 
      antialias: true,
      alpha: true 
    });
    renderer.setSize(width, height);
    renderer.setPixelRatio(window.devicePixelRatio);
    mountRef.current.appendChild(renderer.domElement);
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
    const sphereGeometry = new THREE.SphereGeometry(
      mergedConfig.radius,
      mergedConfig.segments,
      mergedConfig.rings
    );
    
    // Invert the sphere so we see it from inside
    sphereGeometry.scale(-1, 1, 1);
    
    const sphereMaterial = new THREE.MeshBasicMaterial({
      color: 0x001122,
      side: THREE.BackSide
    });
    
    const skySphere = new THREE.Mesh(sphereGeometry, sphereMaterial);
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
      if (!mountRef.current || !camera || !renderer) return;
      
      const newWidth = mountRef.current.clientWidth;
      const newHeight = mountRef.current.clientHeight;
      
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
      
      if (mountRef.current && renderer) {
        mountRef.current.removeChild(renderer.domElement);
      }
      
      renderer?.dispose();
      controls?.dispose();
    };
  }, [isWebGLSupported, mergedConfig, mergedRenderOptions, observer.latitude]);

  // Update celestial objects
  useEffect(() => {
    if (!sceneRef.current || !celestialObjects.length) return;

    // Remove existing celestial objects
    const objectsToRemove: THREE.Object3D[] = [];
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
        let geometry: THREE.BufferGeometry;
        let material: THREE.Material;
        
        if (obj.type === 'star') {
          const size = obj.size || (6 - (obj.magnitude || 0)) * 0.5;
          geometry = new THREE.SphereGeometry(size, 8, 8);
          material = new THREE.MeshBasicMaterial({
            color: obj.color || 0xffffff,
            emissive: obj.color || 0xffffff,
            emissiveIntensity: 0.8
          });
        } else if (obj.type === 'planet') {
          const size = obj.size || 2;
          geometry = new THREE.SphereGeometry(size, 16, 16);
          material = new THREE.MeshPhongMaterial({
            color: obj.color || 0xffaa00,
            emissive: obj.color || 0xffaa00,
            emissiveIntensity: 0.3
          });
        } else {
          // Default for other objects
          geometry = new THREE.SphereGeometry(1, 8, 8);
          material = new THREE.MeshBasicMaterial({
            color: obj.color || 0xffffff
          });
        }

        const mesh = new THREE.Mesh(geometry, material);
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

// Helper functions to add coordinate elements
function addCoordinateGrids(scene: THREE.Scene, radius: number) {
  const gridMaterial = new THREE.LineBasicMaterial({ 
    color: 0x444444, 
    transparent: true, 
    opacity: 0.3 
  });

  // Add latitude lines
  for (let lat = -80; lat <= 80; lat += 10) {
    const curve = new THREE.EllipseCurve(
      0, 0,
      radius * Math.cos(lat * Math.PI / 180),
      radius * Math.cos(lat * Math.PI / 180),
      0, 2 * Math.PI,
      false,
      0
    );
    
    const points = curve.getPoints(64);
    const geometry = new THREE.BufferGeometry().setFromPoints(points);
    const line = new THREE.Line(geometry, gridMaterial);
    line.position.y = radius * Math.sin(lat * Math.PI / 180);
    line.rotation.x = Math.PI / 2;
    scene.add(line);
  }

  // Add longitude lines
  for (let lon = 0; lon < 360; lon += 15) {
    const points = [];
    for (let lat = -90; lat <= 90; lat += 5) {
      const phi = (90 - lat) * Math.PI / 180;
      const theta = lon * Math.PI / 180;
      
      const x = radius * Math.sin(phi) * Math.cos(theta);
      const y = radius * Math.cos(phi);
      const z = radius * Math.sin(phi) * Math.sin(theta);
      
      points.push(new THREE.Vector3(x, y, z));
    }
    
    const geometry = new THREE.BufferGeometry().setFromPoints(points);
    const line = new THREE.Line(geometry, gridMaterial);
    scene.add(line);
  }
}

function addEclipticLine(scene: THREE.Scene, radius: number) {
  const material = new THREE.LineBasicMaterial({ 
    color: 0xffff00, 
    linewidth: 2 
  });
  
  const points = [];
  for (let lon = 0; lon <= 360; lon += 5) {
    const theta = lon * Math.PI / 180;
    const x = radius * Math.cos(theta);
    const z = radius * Math.sin(theta);
    points.push(new THREE.Vector3(x, 0, z));
  }
  
  const geometry = new THREE.BufferGeometry().setFromPoints(points);
  const line = new THREE.Line(geometry, material);
  scene.add(line);
}

function addCelestialEquator(scene: THREE.Scene, radius: number) {
  const material = new THREE.LineBasicMaterial({ 
    color: 0x00ffff, 
    linewidth: 2 
  });
  
  const points = [];
  for (let lon = 0; lon <= 360; lon += 5) {
    const theta = lon * Math.PI / 180;
    const x = radius * Math.cos(theta);
    const z = radius * Math.sin(theta);
    points.push(new THREE.Vector3(x, 0, z));
  }
  
  const geometry = new THREE.BufferGeometry().setFromPoints(points);
  const line = new THREE.Line(geometry, material);
  // Tilt by Earth's obliquity
  line.rotation.x = 23.44 * Math.PI / 180;
  scene.add(line);
}

function addHorizonLine(scene: THREE.Scene, radius: number, latitude: number) {
  const material = new THREE.LineBasicMaterial({ 
    color: 0x00ff00, 
    linewidth: 2,
    transparent: true,
    opacity: 0.7
  });
  
  const points = [];
  for (let az = 0; az <= 360; az += 5) {
    const theta = az * Math.PI / 180;
    const x = radius * Math.sin(theta);
    const z = radius * Math.cos(theta);
    points.push(new THREE.Vector3(x, 0, z));
  }
  
  const geometry = new THREE.BufferGeometry().setFromPoints(points);
  const line = new THREE.Line(geometry, material);
  // Rotate based on observer latitude
  line.rotation.x = (90 - latitude) * Math.PI / 180;
  scene.add(line);
}

export default SkySphere;