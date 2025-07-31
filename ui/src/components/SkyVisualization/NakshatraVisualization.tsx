import React, { useMemo } from 'react';
import * as THREE from 'three';
import { NakshatraVisualization as NakshatraData } from '../../types/skyVisualization';

interface NakshatraProps {
  scene: THREE.Scene;
  radius: number;
  showLabels: boolean;
  showBoundaries: boolean;
  currentNakshatra?: number; // 1-27, highlight current nakshatra
}

// Nakshatra data with astronomical information
const NAKSHATRA_INFO = [
  { id: 1, name: 'Ashwini', deity: 'Ashwini Kumaras', symbol: 'Horse Head', color: '#FF6B6B' },
  { id: 2, name: 'Bharani', deity: 'Yama', symbol: 'Yoni', color: '#4ECDC4' },
  { id: 3, name: 'Krittika', deity: 'Agni', symbol: 'Razor/Knife', color: '#45B7D1' },
  { id: 4, name: 'Rohini', deity: 'Brahma', symbol: 'Cart/Chariot', color: '#96CEB4' },
  { id: 5, name: 'Mrigashira', deity: 'Soma', symbol: 'Deer Head', color: '#FECA57' },
  { id: 6, name: 'Ardra', deity: 'Rudra', symbol: 'Teardrop', color: '#FF9FF3' },
  { id: 7, name: 'Punarvasu', deity: 'Aditi', symbol: 'Bow & Quiver', color: '#54A0FF' },
  { id: 8, name: 'Pushya', deity: 'Brihaspati', symbol: 'Cow Udder', color: '#5F27CD' },
  { id: 9, name: 'Ashlesha', deity: 'Nagas', symbol: 'Serpent', color: '#00D2D3' },
  { id: 10, name: 'Magha', deity: 'Pitrs', symbol: 'Throne', color: '#FF6348' },
  { id: 11, name: 'Purva Phalguni', deity: 'Bhaga', symbol: 'Front Bed Legs', color: '#FF9F43' },
  { id: 12, name: 'Uttara Phalguni', deity: 'Aryaman', symbol: 'Back Bed Legs', color: '#70A1FF' },
  { id: 13, name: 'Hasta', deity: 'Savitar', symbol: 'Hand', color: '#7BED9F' },
  { id: 14, name: 'Chitra', deity: 'Tvashtar', symbol: 'Bright Jewel', color: '#FF6B9D' },
  { id: 15, name: 'Swati', deity: 'Vayu', symbol: 'Young Plant', color: '#C7ECEE' },
  { id: 16, name: 'Vishakha', deity: 'Indra-Agni', symbol: 'Triumphal Arch', color: '#FFA502' },
  { id: 17, name: 'Anuradha', deity: 'Mitra', symbol: 'Lotus', color: '#3742FA' },
  { id: 18, name: 'Jyeshtha', deity: 'Indra', symbol: 'Circular Amulet', color: '#2ED573' },
  { id: 19, name: 'Mula', deity: 'Nirriti', symbol: 'Bunch of Roots', color: '#FF4757' },
  { id: 20, name: 'Purva Ashadha', deity: 'Apas', symbol: 'Elephant Tusk', color: '#FFA726' },
  { id: 21, name: 'Uttara Ashadha', deity: 'Vishve Devas', symbol: 'Elephant Tusk', color: '#42A5F5' },
  { id: 22, name: 'Shravana', deity: 'Vishnu', symbol: 'Ear/Footprints', color: '#66BB6A' },
  { id: 23, name: 'Dhanishta', deity: 'Vasus', symbol: 'Drum', color: '#AB47BC' },
  { id: 24, name: 'Shatabhisha', deity: 'Varuna', symbol: 'Empty Circle', color: '#26C6DA' },
  { id: 25, name: 'Purva Bhadrapada', deity: 'Aja Ekapada', symbol: 'Funeral Cot Front', color: '#EC407A' },
  { id: 26, name: 'Uttara Bhadrapada', deity: 'Ahir Budhnya', symbol: 'Funeral Cot Back', color: '#8BC34A' },
  { id: 27, name: 'Revati', deity: 'Pushan', symbol: 'Fish/Pair of Fish', color: '#FFB74D' }
];

export const NakshatraVisualization: React.FC<NakshatraProps> = ({
  scene,
  radius,
  showLabels,
  showBoundaries,
  currentNakshatra
}) => {
  
  const nakshatraObjects = useMemo(() => {
    const objects: THREE.Object3D[] = [];
    
    // Each Nakshatra spans 13.333... degrees (360/27)
    const nakshatraSpan = 360 / 27;
    
    NAKSHATRA_INFO.forEach((nakshatra, index) => {
      const startLongitude = index * nakshatraSpan;
      const endLongitude = (index + 1) * nakshatraSpan;
      const centerLongitude = (startLongitude + endLongitude) / 2;
      
      const isCurrentNakshatra = currentNakshatra === nakshatra.id;
      
      // Create Nakshatra boundary arcs
      if (showBoundaries) {
        const boundaryMaterial = new THREE.LineBasicMaterial({
          color: isCurrentNakshatra ? 0xffffff : new THREE.Color(nakshatra.color).getHex(),
          linewidth: isCurrentNakshatra ? 3 : 1,
          transparent: true,
          opacity: isCurrentNakshatra ? 1.0 : 0.6
        });
        
        // Create arc for nakshatra boundary
        const arcPoints: THREE.Vector3[] = [];
        const numPoints = 32;
        
        for (let i = 0; i <= numPoints; i++) {
          const longitude = startLongitude + (i / numPoints) * nakshatraSpan;
          const theta = longitude * Math.PI / 180;
          
          // Create arc at different latitudes to show the nakshatra band
          for (let lat = -10; lat <= 10; lat += 20) {
            const phi = (90 - lat) * Math.PI / 180;
            const x = radius * Math.sin(phi) * Math.cos(theta);
            const y = radius * Math.cos(phi);
            const z = radius * Math.sin(phi) * Math.sin(theta);
            arcPoints.push(new THREE.Vector3(x, y, z));
          }
        }
        
        const boundaryGeometry = new THREE.BufferGeometry().setFromPoints(arcPoints);
        const boundaryLine = new THREE.Line(boundaryGeometry, boundaryMaterial);
        boundaryLine.userData = { 
          type: 'nakshatra_boundary', 
          nakshatraId: nakshatra.id,
          name: nakshatra.name
        };
        objects.push(boundaryLine);
      }
      
      // Create Nakshatra markers at center
      const markerGeometry = new THREE.SphereGeometry(
        isCurrentNakshatra ? 0.8 : 0.5, 
        16, 
        16
      );
      const markerMaterial = new THREE.MeshBasicMaterial({
        color: new THREE.Color(nakshatra.color).getHex(),
        emissive: new THREE.Color(nakshatra.color).getHex(),
        emissiveIntensity: isCurrentNakshatra ? 0.8 : 0.4,
        transparent: true,
        opacity: isCurrentNakshatra ? 1.0 : 0.7
      });
      
      const marker = new THREE.Mesh(markerGeometry, markerMaterial);
      
      // Position marker at center of nakshatra
      const theta = centerLongitude * Math.PI / 180;
      const phi = Math.PI / 2; // On ecliptic plane
      
      const x = radius * Math.sin(phi) * Math.cos(theta);
      const y = radius * Math.cos(phi);
      const z = radius * Math.sin(phi) * Math.sin(theta);
      
      marker.position.set(x, y, z);
      marker.userData = {
        type: 'nakshatra_marker',
        nakshatraId: nakshatra.id,
        name: nakshatra.name,
        deity: nakshatra.deity,
        symbol: nakshatra.symbol,
        longitude: centerLongitude
      };
      objects.push(marker);
      
      // Create labels if enabled
      if (showLabels) {
        const canvas = document.createElement('canvas');
        const context = canvas.getContext('2d')!;
        
        canvas.width = 256;
        canvas.height = 64;
        
        context.fillStyle = 'rgba(0, 0, 0, 0.8)';
        context.fillRect(0, 0, canvas.width, canvas.height);
        
        context.fillStyle = isCurrentNakshatra ? '#ffffff' : nakshatra.color;
        context.font = isCurrentNakshatra ? 'bold 16px Arial' : '14px Arial';
        context.textAlign = 'center';
        context.textBaseline = 'middle';
        
        // Draw nakshatra name
        context.fillText(nakshatra.name, canvas.width / 2, canvas.height / 2 - 8);
        
        // Draw nakshatra number
        context.font = '12px Arial';
        context.fillStyle = isCurrentNakshatra ? '#cccccc' : '#999999';
        context.fillText(`${nakshatra.id}`, canvas.width / 2, canvas.height / 2 + 8);
        
        const texture = new THREE.CanvasTexture(canvas);
        const labelMaterial = new THREE.SpriteMaterial({
          map: texture,
          transparent: true,
          alphaTest: 0.1
        });
        
        const label = new THREE.Sprite(labelMaterial);
        label.scale.set(4, 1, 1);
        
        // Position label slightly outside the marker
        const labelRadius = radius + 2;
        const labelX = labelRadius * Math.sin(phi) * Math.cos(theta);
        const labelY = labelRadius * Math.cos(phi);
        const labelZ = labelRadius * Math.sin(phi) * Math.sin(theta);
        
        label.position.set(labelX, labelY, labelZ);
        label.userData = {
          type: 'nakshatra_label',
          nakshatraId: nakshatra.id,
          name: nakshatra.name
        };
        objects.push(label);
      }
    });
    
    return objects;
  }, [radius, showLabels, showBoundaries, currentNakshatra]);
  
  // Add objects to scene
  React.useEffect(() => {
    // Remove existing nakshatra objects
    const objectsToRemove: THREE.Object3D[] = [];
    scene.traverse((child) => {
      if (child.userData.type && child.userData.type.startsWith('nakshatra_')) {
        objectsToRemove.push(child);
      }
    });
    objectsToRemove.forEach(obj => scene.remove(obj));
    
    // Add new objects
    nakshatraObjects.forEach(obj => scene.add(obj));
    
    return () => {
      // Cleanup
      nakshatraObjects.forEach(obj => {
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
  }, [scene, nakshatraObjects]);
  
  return null; // This component only adds objects to the scene
};

export default NakshatraVisualization;