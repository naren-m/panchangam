export const formatTime = (timeString: string, format: '12' | '24' = '12'): string => {
  const [hours, minutes] = timeString.split(':');
  const hour = parseInt(hours);
  const min = minutes || '00';

  if (format === '24') {
    return `${hours.padStart(2, '0')}:${min}`;
  }

  const ampm = hour >= 12 ? 'PM' : 'AM';
  const displayHour = hour === 0 ? 12 : hour > 12 ? hour - 12 : hour;
  return `${displayHour}:${min} ${ampm}`;
};

export const formatTimeRange = (timeRange: string, format: '12' | '24' = '12'): string => {
  const [start, end] = timeRange.split('-');
  return `${formatTime(start, format)} - ${formatTime(end, format)}`;
};

export const getCurrentMonthDates = (year: number, month: number): Date[] => {
  const firstDay = new Date(year, month, 1);
  const dates: Date[] = [];

  // Add dates from previous month to fill the first week
  const startDate = new Date(firstDay);
  startDate.setDate(startDate.getDate() - firstDay.getDay());

  // Add dates until we fill the calendar grid (42 days = 6 weeks)
  for (let i = 0; i < 42; i++) {
    dates.push(new Date(startDate));
    startDate.setDate(startDate.getDate() + 1);
  }

  return dates;
};

export const isToday = (date: Date): boolean => {
  const today = new Date();
  return date.toDateString() === today.toDateString();
};

export const isSameMonth = (date: Date, month: number, year: number): boolean => {
  return date.getMonth() === month && date.getFullYear() === year;
};

export const formatDateForApi = (date: Date): string => {
  return date.toISOString().split('T')[0];
};

const MILLISECONDS_PER_DAY = 1000 * 60 * 60 * 24;

const getUtcCalendarDay = (date: Date): number => {
  return Date.UTC(date.getFullYear(), date.getMonth(), date.getDate());
};

export const getCalendarDayDifference = (startDate: Date, endDate: Date): number => {
  const dayDifference = getUtcCalendarDay(endDate) - getUtcCalendarDay(startDate);
  return dayDifference / MILLISECONDS_PER_DAY;
};

export const countCalendarDaysInclusive = (startDate: Date, endDate: Date): number => {
  return getCalendarDayDifference(startDate, endDate) + 1;
};

export const getCalendarDayOfYear = (date: Date): number => {
  return getCalendarDayDifference(new Date(date.getFullYear(), 0, 0), date);
};

export const parseApiDate = (dateString: string): Date => {
  const [year, month, day] = dateString.split('-').map(Number);
  return new Date(year, month - 1, day);
};

export const getMonthName = (month: number, locale: string = 'en'): string => {
  const monthNames = {
    en: ['January', 'February', 'March', 'April', 'May', 'June',
         'July', 'August', 'September', 'October', 'November', 'December'],
    hi: ['जनवरी', 'फरवरी', 'मार्च', 'अप्रैल', 'मई', 'जून',
         'जुलाई', 'अगस्त', 'सितंबर', 'अक्टूबर', 'नवंबर', 'दिसंबर'],
    ta: ['ஜனவரி', 'பிப்ரவரி', 'மார்ச்', 'ஏப்ரல்', 'மே', 'ஜூன்',
         'ஜூலை', 'ஆகஸ்ட்', 'செப்டம்பர்', 'அக்டோபர்', 'நவம்பர்', 'டிசம்பர்']
  };

  return monthNames[locale as keyof typeof monthNames]?.[month] || monthNames.en[month];
};

export const getWeekdayNames = (locale: string = 'en'): string[] => {
  const weekdays = {
    en: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
    hi: ['रवि', 'सोम', 'मंगल', 'बुध', 'गुरु', 'शुक्र', 'शनि'],
    ta: ['ஞா', 'தி', 'செ', 'பு', 'வி', 'வெ', 'ச']
  };

  return weekdays[locale as keyof typeof weekdays] || weekdays.en;
};
