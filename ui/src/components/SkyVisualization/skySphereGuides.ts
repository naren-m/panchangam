import {
  BufferGeometry,
  EllipseCurve,
  Line,
  LineBasicMaterial,
  Scene,
  Vector3,
} from 'three';

export function addCoordinateGrids(scene: Scene, radius: number) {
  const gridMaterial = new LineBasicMaterial({
    color: 0x444444,
    transparent: true,
    opacity: 0.3
  });

  for (let lat = -80; lat <= 80; lat += 10) {
    const curve = new EllipseCurve(
      0, 0,
      radius * Math.cos(lat * Math.PI / 180),
      radius * Math.cos(lat * Math.PI / 180),
      0, 2 * Math.PI,
      false,
      0
    );

    const points = curve.getPoints(64);
    const geometry = new BufferGeometry().setFromPoints(points);
    const line = new Line(geometry, gridMaterial);
    line.position.y = radius * Math.sin(lat * Math.PI / 180);
    line.rotation.x = Math.PI / 2;
    scene.add(line);
  }

  for (let lon = 0; lon < 360; lon += 15) {
    const points = [];
    for (let lat = -90; lat <= 90; lat += 5) {
      const phi = (90 - lat) * Math.PI / 180;
      const theta = lon * Math.PI / 180;

      points.push(new Vector3(
        radius * Math.sin(phi) * Math.cos(theta),
        radius * Math.cos(phi),
        radius * Math.sin(phi) * Math.sin(theta)
      ));
    }

    const geometry = new BufferGeometry().setFromPoints(points);
    const line = new Line(geometry, gridMaterial);
    scene.add(line);
  }
}

export function addEclipticLine(scene: Scene, radius: number) {
  const material = new LineBasicMaterial({
    color: 0xffff00,
    linewidth: 2
  });

  const points = [];
  for (let lon = 0; lon <= 360; lon += 5) {
    const theta = lon * Math.PI / 180;
    points.push(new Vector3(
      radius * Math.cos(theta),
      0,
      radius * Math.sin(theta)
    ));
  }

  const geometry = new BufferGeometry().setFromPoints(points);
  const line = new Line(geometry, material);
  scene.add(line);
}

export function addCelestialEquator(scene: Scene, radius: number) {
  const material = new LineBasicMaterial({
    color: 0x00ffff,
    linewidth: 2
  });

  const points = [];
  for (let lon = 0; lon <= 360; lon += 5) {
    const theta = lon * Math.PI / 180;
    points.push(new Vector3(
      radius * Math.cos(theta),
      0,
      radius * Math.sin(theta)
    ));
  }

  const geometry = new BufferGeometry().setFromPoints(points);
  const line = new Line(geometry, material);
  line.rotation.x = 23.44 * Math.PI / 180;
  scene.add(line);
}

export function addHorizonLine(scene: Scene, radius: number, latitude: number) {
  const material = new LineBasicMaterial({
    color: 0x00ff00,
    linewidth: 2,
    transparent: true,
    opacity: 0.7
  });

  const points = [];
  for (let az = 0; az <= 360; az += 5) {
    const theta = az * Math.PI / 180;
    points.push(new Vector3(
      radius * Math.sin(theta),
      0,
      radius * Math.cos(theta)
    ));
  }

  const geometry = new BufferGeometry().setFromPoints(points);
  const line = new Line(geometry, material);
  line.rotation.x = (90 - latitude) * Math.PI / 180;
  scene.add(line);
}
