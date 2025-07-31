import React, { useMemo } from 'react';
import * as THREE from 'three';

interface ZodiacVisualizationProps {
  scene: THREE.Scene;
  radius: number;
  showLabels: boolean;
  showBoundaries: boolean;
  currentRashi?: number; // 1-12, current zodiac sign to highlight
}

// Zodiac (Rashi) data with astronomical information
const ZODIAC_INFO = [
  { id: 1, name: 'Mesha', western: 'Aries', symbol: '♈', element: 'Fire', ruler: 'Mars', color: '#FF6B6B' },
  { id: 2, name: 'Vrishabha', western: 'Taurus', symbol: '♉', element: 'Earth', ruler: 'Venus', color: '#4ECDC4' },
  { id: 3, name: 'Mithuna', western: 'Gemini', symbol: '♊', element: 'Air', ruler: 'Mercury', color: '#45B7D1' },
  { id: 4, name: 'Karka', western: 'Cancer', symbol: '♋', element: 'Water', ruler: 'Moon', color: '#96CEB4' },
  { id: 5, name: 'Simha', western: 'Leo', symbol: '♌', element: 'Fire', ruler: 'Sun', color: '#FECA57' },
  { id: 6, name: 'Kanya', western: 'Virgo', symbol: '♍', element: 'Earth', ruler: 'Mercury', color: '#FF9FF3' },
  { id: 7, name: 'Tula', western: 'Libra', symbol: '♎', element: 'Air', ruler: 'Venus', color: '#54A0FF' },
  { id: 8, name: 'Vrishchika', western: 'Scorpio', symbol: '♏', element: 'Water', ruler: 'Mars', color: '#5F27CD' },
  { id: 9, name: 'Dhanus', western: 'Sagittarius', symbol: '♐', element: 'Fire', ruler: 'Jupiter', color: '#00D2D3' },
  { id: 10, name: 'Makara', western: 'Capricorn', symbol: '♑', element: 'Earth', ruler: 'Saturn', color: '#FF6348' },
  { id: 11, name: 'Kumbha', western: 'Aquarius', symbol: '♒', element: 'Air', ruler: 'Saturn', color: '#FF9F43' },
  { id: 12, name: 'Meena', western: 'Pisces', symbol: '♓', element: 'Water', ruler: 'Jupiter', color: '#70A1FF' }
];

export const ZodiacVisualization: React.FC<ZodiacVisualizationProps> = ({
  scene,
  radius,
  showLabels,
  showBoundaries,
  currentRashi
}) => {
  
  const zodiacObjects = useMemo(() => {
    const objects: THREE.Object3D[] = [];
    
    // Each zodiac sign spans 30 degrees (360/12)
    const rashiSpan = 360 / 12;
    
    ZODIAC_INFO.forEach((rashi, index) => {
      const startLongitude = index * rashiSpan;
      const endLongitude = (index + 1) * rashiSpan;
      const centerLongitude = (startLongitude + endLongitude) / 2;
      
      const isCurrentRashi = currentRashi === rashi.id;
      
      // Create Zodiac boundary arcs
      if (showBoundaries) {
        const boundaryMaterial = new THREE.LineBasicMaterial({
          color: isCurrentRashi ? 0xffffff : new THREE.Color(rashi.color).getHex(),
          linewidth: isCurrentRashi ? 4 : 2,
          transparent: true,
          opacity: isCurrentRashi ? 1.0 : 0.8
        });
        
        // Create zodiac boundary lines - these are meridian lines from pole to pole
        const meridianPoints: THREE.Vector3[] = [];
        
        // Draw meridian lines at start and end of zodiac sign
        [startLongitude, endLongitude].forEach(longitude => {
          for (let lat = -90; lat <= 90; lat += 2) {
            const phi = (90 - lat) * Math.PI / 180;
            const theta = longitude * Math.PI / 180;
            
            const x = radius * Math.sin(phi) * Math.cos(theta);
            const y = radius * Math.cos(phi);
            const z = radius * Math.sin(phi) * Math.sin(theta);
            
            meridianPoints.push(new THREE.Vector3(x, y, z));
          }
          
          // Add separator between start and end line
          if (longitude === startLongitude) {
            meridianPoints.push(new THREE.Vector3(NaN, NaN, NaN)); // Line break
          }
        });
        
        const boundaryGeometry = new THREE.BufferGeometry().setFromPoints(
          meridianPoints.filter(p => !isNaN(p.x))
        );
        const boundaryLine = new THREE.Line(boundaryGeometry, boundaryMaterial);
        boundaryLine.userData = { 
          type: 'zodiac_boundary', 
          rashiId: rashi.id,
          name: rashi.name
        };
        objects.push(boundaryLine);
      }
      
      // Create Zodiac sector fill (optional)
      if (showBoundaries) {
        // Create a subtle sector fill on the ecliptic plane
        const sectorGeometry = new THREE.RingGeometry(
          radius * 0.98, 
          radius * 1.02, 
          Math.floor(startLongitude * 2), 
          Math.floor(rashiSpan * 2)
        );
        
        const sectorMaterial = new THREE.MeshBasicMaterial({
          color: new THREE.Color(rashi.color).getHex(),
          transparent: true,
          opacity: isCurrentRashi ? 0.2 : 0.1,
          side: THREE.DoubleSide
        });
        
        const sector = new THREE.Mesh(sectorGeometry, sectorMaterial);
        sector.rotation.x = Math.PI / 2; // Rotate to lie on ecliptic plane
        sector.userData = {
          type: 'zodiac_sector',
          rashiId: rashi.id,
          name: rashi.name
        };
        objects.push(sector);
      }
      
      // Create Zodiac symbol markers
      const markerGeometry = new THREE.SphereGeometry(
        isCurrentRashi ? 1.2 : 0.8, 
        16, 
        16
      );
      const markerMaterial = new THREE.MeshBasicMaterial({
        color: new THREE.Color(rashi.color).getHex(),
        emissive: new THREE.Color(rashi.color).getHex(),
        emissiveIntensity: isCurrentRashi ? 1.0 : 0.5,
        transparent: true,
        opacity: isCurrentRashi ? 1.0 : 0.8
      });
      
      const marker = new THREE.Mesh(markerGeometry, markerMaterial);
      
      // Position marker at center of zodiac sign on ecliptic
      const theta = centerLongitude * Math.PI / 180;
      const phi = Math.PI / 2; // On ecliptic plane
      
      const x = radius * Math.sin(phi) * Math.cos(theta);
      const y = 0; // On ecliptic plane
      const z = radius * Math.sin(phi) * Math.sin(theta);
      
      marker.position.set(x, y, z);
      marker.userData = {
        type: 'zodiac_marker',
        rashiId: rashi.id,
        name: rashi.name,
        western: rashi.western,
        symbol: rashi.symbol,
        element: rashi.element,
        ruler: rashi.ruler,
        longitude: centerLongitude
      };
      objects.push(marker);
      
      // Create labels if enabled
      if (showLabels) {
        const canvas = document.createElement('canvas');
        const context = canvas.getContext('2d')!;
        
        canvas.width = 320;
        canvas.height = 80;
        
        context.fillStyle = 'rgba(0, 0, 0, 0.8)';
        context.fillRect(0, 0, canvas.width, canvas.height);
        
        context.fillStyle = isCurrentRashi ? '#ffffff' : rashi.color;
        context.font = isCurrentRashi ? 'bold 18px Arial' : '16px Arial';
        context.textAlign = 'center';
        context.textBaseline = 'middle';
        
        // Draw Sanskrit name and symbol
        context.fillText(`${rashi.symbol} ${rashi.name}`, canvas.width / 2, canvas.height / 2 - 12);
        
        // Draw Western name
        context.font = '14px Arial';
        context.fillStyle = isCurrentRashi ? '#cccccc' : '#999999';
        context.fillText(`${rashi.western} (${rashi.element})`, canvas.width / 2, canvas.height / 2 + 8);
        
        const texture = new THREE.CanvasTexture(canvas);
        const labelMaterial = new THREE.SpriteMaterial({
          map: texture,
          transparent: true,
          alphaTest: 0.1
        });
        
        const label = new THREE.Sprite(labelMaterial);
        label.scale.set(6, 1.5, 1);
        
        // Position label slightly outside the marker
        const labelRadius = radius + 4;
        const labelX = labelRadius * Math.sin(phi) * Math.cos(theta);
        const labelY = 2; // Slightly above ecliptic
        const labelZ = labelRadius * Math.sin(phi) * Math.sin(theta);
        
        label.position.set(labelX, labelY, labelZ);
        label.userData = {
          type: 'zodiac_label',
          rashiId: rashi.id,
          name: rashi.name
        };
        objects.push(label);
      }
    });
    
    return objects;
  }, [radius, showLabels, showBoundaries, currentRashi]);
  
  // Add objects to scene
  React.useEffect(() => {
    // Remove existing zodiac objects
    const objectsToRemove: THREE.Object3D[] = [];
    scene.traverse((child) => {
      if (child.userData.type && child.userData.type.startsWith('zodiac_')) {
        objectsToRemove.push(child);
      }
    });
    objectsToRemove.forEach(obj => scene.remove(obj));
    
    // Add new objects
    zodiacObjects.forEach(obj => scene.add(obj));
    
    return () => {
      // Cleanup
      zodiacObjects.forEach(obj => {
        scene.remove(obj);
        if (obj instanceof THREE.Mesh || obj instanceof THREE.Line) {
          obj.geometry.dispose();
          if (Array.isArray(obj.material)) {
            obj.material.forEach(mat => mat.dispose());
          } else {
            obj.material.dispose();
          }
        }
      });
    };
  }, [scene, zodiacObjects]);
  
  return null; // This component only adds objects to the scene
};

export default ZodiacVisualization;