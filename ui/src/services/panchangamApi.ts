import { PanchangamData, GetPanchangamRequest } from '../types/panchangam';

// Mock data with authentic panchangam information
const mockPanchangamData: Record<string, PanchangamData> = {
  "2024-01-15": {
    date: "2024-01-15",
    tithi: "Shukla Panchami",
    nakshatra: "Rohini",
    yoga: "Siddha",
    karana: "Bava",
    sunrise_time: "06:45:30",
    sunset_time: "18:15:45",
    moonrise_time: "10:30:00",
    moonset_time: "23:45:00",
    vara: "Monday",
    planetary_ruler: "Moon",
    festivals: ["Makar Sankranti"],
    events: [
      {
        name: "Brahma Muhurta",
        time: "05:09:30-05:57:30",
        event_type: "BRAHMA_MUHURTA",
        quality: "auspicious"
      },
      {
        name: "Rahu Kalam",
        time: "07:30:00-09:00:00",
        event_type: "RAHU_KALAM",
        quality: "inauspicious"
      },
      {
        name: "Abhijit Muhurta",
        time: "12:00:00-12:48:00",
        event_type: "ABHIJIT",
        quality: "auspicious"
      },
      {
        name: "Godhuli Muhurta",
        time: "17:45:45-18:15:45",
        event_type: "GODHULI",
        quality: "auspicious"
      }
    ]
  },
  "2024-01-16": {
    date: "2024-01-16",
    tithi: "Shukla Shashthi",
    nakshatra: "Mrigashirsha",
    yoga: "Sadhya",
    karana: "Balava",
    sunrise_time: "06:45:15",
    sunset_time: "18:16:30",
    moonrise_time: "11:15:00",
    moonset_time: "00:30:00",
    vara: "Tuesday",
    planetary_ruler: "Mars",
    festivals: [],
    events: [
      {
        name: "Brahma Muhurta",
        time: "05:09:15-05:57:15",
        event_type: "BRAHMA_MUHURTA",
        quality: "auspicious"
      },
      {
        name: "Rahu Kalam",
        time: "15:00:00-16:30:00",
        event_type: "RAHU_KALAM",
        quality: "inauspicious"
      },
      {
        name: "Abhijit Muhurta",
        time: "12:00:45-12:48:45",
        event_type: "ABHIJIT",
        quality: "auspicious"
      }
    ]
  },
  "2024-01-17": {
    date: "2024-01-17",
    tithi: "Shukla Saptami",
    nakshatra: "Ardra",
    yoga: "Shubha",
    karana: "Kaulava",
    sunrise_time: "06:45:00",
    sunset_time: "18:17:15",
    moonrise_time: "12:00:00",
    moonset_time: "01:15:00",
    vara: "Wednesday",
    planetary_ruler: "Mercury",
    festivals: [],
    events: [
      {
        name: "Brahma Muhurta",
        time: "05:09:00-05:57:00",
        event_type: "BRAHMA_MUHURTA",
        quality: "auspicious"
      },
      {
        name: "Rahu Kalam",
        time: "12:00:00-13:30:00",
        event_type: "RAHU_KALAM",
        quality: "inauspicious"
      },
      {
        name: "Abhijit Muhurta",
        time: "12:01:30-12:49:30",
        event_type: "ABHIJIT",
        quality: "auspicious"
      }
    ]
  }
};

// Generate additional mock data for the current month
const generateMockData = (date: string): PanchangamData => {
  const tithis = [
    "Krishna Pratipada", "Krishna Dwitiya", "Krishna Tritiya", "Krishna Chaturthi", "Krishna Panchami",
    "Krishna Shashthi", "Krishna Saptami", "Krishna Ashtami", "Krishna Navami", "Krishna Dashami",
    "Krishna Ekadashi", "Krishna Dwadashi", "Krishna Trayodashi", "Krishna Chaturdashi", "Amavasya",
    "Shukla Pratipada", "Shukla Dwitiya", "Shukla Tritiya", "Shukla Chaturthi", "Shukla Panchami",
    "Shukla Shashthi", "Shukla Saptami", "Shukla Ashtami", "Shukla Navami", "Shukla Dashami",
    "Shukla Ekadashi", "Shukla Dwadashi", "Shukla Trayodashi", "Shukla Chaturdashi", "Purnima"
  ];

  const nakshatras = [
    "Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashirsha", "Ardra", "Punarvasu",
    "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni", "Hasta",
    "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha", "Mula", "Purva Ashadha",
    "Uttara Ashadha", "Shravana", "Dhanishtha", "Shatabhisha", "Purva Bhadrapada",
    "Uttara Bhadrapada", "Revati"
  ];

  const yogas = [
    "Vishkumbha", "Preeti", "Ayushman", "Saubhagya", "Shobhana", "Atiganda", "Sukarman",
    "Dhriti", "Shoola", "Ganda", "Vriddhi", "Dhruva", "Vyaghata", "Harshana", "Vajra",
    "Siddhi", "Vyatipata", "Variyana", "Parigha", "Shiva", "Siddha", "Sadhya", "Shubha",
    "Shukla", "Brahma", "Indra", "Vaidhriti"
  ];

  const karanas = [
    "Bava", "Balava", "Kaulava", "Taitila", "Gara", "Vanija", "Vishti", "Shakuni",
    "Chatushpada", "Naga", "Kimstughna"
  ];

  const weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  const rulers = ["Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"];

  const dateObj = new Date(date);
  const dayIndex = Math.floor(Math.random() * tithis.length);
  const dayOfWeek = dateObj.getDay();

  return {
    date,
    tithi: tithis[dayIndex],
    nakshatra: nakshatras[Math.floor(Math.random() * nakshatras.length)],
    yoga: yogas[Math.floor(Math.random() * yogas.length)],
    karana: karanas[Math.floor(Math.random() * karanas.length)],
    sunrise_time: "06:45:00",
    sunset_time: "18:15:00",
    moonrise_time: `${10 + Math.floor(Math.random() * 6)}:${Math.floor(Math.random() * 60)}:00`,
    moonset_time: `${22 + Math.floor(Math.random() * 4)}:${Math.floor(Math.random() * 60)}:00`,
    vara: weekdays[dayOfWeek],
    planetary_ruler: rulers[dayOfWeek],
    festivals: Math.random() > 0.8 ? ["Festival Day"] : [],
    events: [
      {
        name: "Brahma Muhurta",
        time: "05:09:00-05:57:00",
        event_type: "BRAHMA_MUHURTA",
        quality: "auspicious"
      },
      {
        name: "Rahu Kalam",
        time: `${9 + (dayOfWeek * 1.5)}:00:00-${10 + (dayOfWeek * 1.5)}:30:00`,
        event_type: "RAHU_KALAM",
        quality: "inauspicious"
      },
      {
        name: "Abhijit Muhurta",
        time: "12:00:00-12:48:00",
        event_type: "ABHIJIT",
        quality: "auspicious"
      }
    ]
  };
};

class PanchangamApiService {
  async getPanchangam(request: GetPanchangamRequest): Promise<PanchangamData> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 300));

    // Return mock data or generate it
    return mockPanchangamData[request.date] || generateMockData(request.date);
  }

  async getPanchangamRange(startDate: string, endDate: string, request: Omit<GetPanchangamRequest, 'date'>): Promise<PanchangamData[]> {
    const start = new Date(startDate);
    const end = new Date(endDate);
    const results: PanchangamData[] = [];

    for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
      const dateStr = d.toISOString().split('T')[0];
      const data = await this.getPanchangam({ ...request, date: dateStr });
      results.push(data);
    }

    return results;
  }
}

export const panchangamApi = new PanchangamApiService();