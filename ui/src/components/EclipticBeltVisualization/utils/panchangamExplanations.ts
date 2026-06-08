import type { KaranaInfo, TithiInfo, YogaInfo } from '../types/eclipticBelt';

export function getTithiExplanation(tithi: TithiInfo): string {
  const pakshaName = tithi.paksha === 'Shukla' ? 'bright (waxing)' : 'dark (waning)';
  return `${tithi.name} is tithi ${tithi.number} of the ${pakshaName} fortnight. ` +
    `The Moon is ${tithi.angle.toFixed(1)} degrees ahead of the Sun. ` +
    `Deity: ${tithi.deity}. ${tithi.percentComplete.toFixed(0)}% of this tithi has elapsed.`;
}

export function getYogaExplanation(yoga: YogaInfo): string {
  return `${yoga.name} yoga (#${yoga.number}) means "${yoga.meaning}". ` +
    `It is ${yoga.nature.toLowerCase()}. ` +
    `Calculated from Sun + Moon = ${yoga.combinedLongitude.toFixed(1)} degrees.`;
}

export function getKaranaExplanation(karana: KaranaInfo): string {
  return `${karana.name} is a ${karana.type.toLowerCase()} karana (#${karana.number} in this lunar month). ` +
    `It is considered ${karana.nature.toLowerCase()} for activities.`;
}
