import React from 'react';
import { getCurrentMonthDates, getWeekdayNames, isSameMonth } from '../../utils/dateHelpers';
import { DateCell } from './DateCell';
import { PanchangamData, Settings } from '../../types/panchangam';

interface CalendarGridProps {
  year: number;
  month: number;
  panchangamData: Record<string, PanchangamData>;
  settings: Settings;
  onDateClick: (date: Date) => void;
}

export const CalendarGrid: React.FC<CalendarGridProps> = ({
  year,
  month,
  panchangamData,
  settings,
  onDateClick
}) => {
  const dates = getCurrentMonthDates(year, month);
  const weekdays = getWeekdayNames(settings.locale);

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200">
      {/* Weekday headers */}
      <div className="grid grid-cols-7 bg-gradient-to-r from-orange-400 to-orange-500">
        {weekdays.map((day, index) => (
          <div key={index} className="p-3 text-center">
            <span className="text-white font-semibold text-sm md:text-base">
              {day}
            </span>
          </div>
        ))}
      </div>

      {/* Calendar grid */}
      <div className="grid grid-cols-7 gap-0">
        {dates.map((date, index) => {
          const dateStr = date.toISOString().split('T')[0];
          const data = panchangamData[dateStr];
          const isCurrentMonth = isSameMonth(date, month, year);

          return (
            <DateCell
              key={index}
              date={date}
              data={data}
              isCurrentMonth={isCurrentMonth}
              settings={settings}
              onClick={() => onDateClick(date)}
            />
          );
        })}
      </div>
    </div>
  );
};