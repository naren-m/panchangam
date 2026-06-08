export function normalizeDegrees(degrees: number): number {
  let normalized = degrees;
  while (normalized < 0) normalized += 360;
  while (normalized >= 360) normalized -= 360;
  return normalized;
}

export function calculateAngularDifference(moonLongitude: number, sunLongitude: number): number {
  const diff = normalizeDegrees(moonLongitude) - normalizeDegrees(sunLongitude);
  return normalizeDegrees(diff);
}
