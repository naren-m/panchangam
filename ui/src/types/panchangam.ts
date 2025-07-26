export interface PanchangamData {
  date: string;
  tithi: string;
  nakshatra: string;
  yoga: string;
  karana: string;
  sunrise_time: string;
  sunset_time: string;
  moonrise_time?: string;
  moonset_time?: string;
  events: Event[];
  festivals?: string[];
  vara: string;
  planetary_ruler: string;
}

export interface Event {
  name: string;
  time: string;
  event_type: 
    | 'SUNRISE' | 'SUNSET' 
    | 'MOONRISE' | 'MOONSET' | 'MOON_PHASE'
    | 'RAHU_KALAM' | 'YAMAGANDAM' | 'GULIKA_KALAM' | 'ABHIJIT_MUHURTA'
    | 'MUHURTA' | 'ABHIJIT' | 'BRAHMA_MUHURTA' | 'GODHULI'
    | 'TITHI' | 'NAKSHATRA' | 'YOGA' | 'KARANA' | 'VARA'
    | 'FESTIVAL';
  quality: 'auspicious' | 'inauspicious' | 'neutral';
}

export interface Location {
  name: string;
  latitude: number;
  longitude: number;
  timezone: string;
  region: string;
}

export interface Settings {
  calculation_method: 'Drik' | 'Vakya';
  locale: 'en' | 'hi' | 'ta';
  region: string;
  time_format: '12' | '24';
  location: Location;
}

export interface GetPanchangamRequest {
  date: string;
  latitude: number;
  longitude: number;
  timezone: string;
  region: string;
  calculation_method: string;
  locale: string;
}