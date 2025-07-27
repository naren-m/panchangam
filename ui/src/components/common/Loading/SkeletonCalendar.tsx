import React from 'react';

interface SkeletonCalendarProps {
  rows?: number;
}

export const SkeletonCalendar: React.FC<SkeletonCalendarProps> = ({ rows = 6 }) => {
  const weekdays = Array(7).fill(null);
  const cells = Array(rows * 7).fill(null);

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 animate-pulse">
      {/* Weekday headers skeleton */}
      <div className="grid grid-cols-7 bg-gradient-to-r from-gray-300 to-gray-400">
        {weekdays.map((_, index) => (
          <div key={index} className="p-3 text-center">
            <div className="h-4 bg-gray-200 rounded mx-auto w-8"></div>
          </div>
        ))}
      </div>

      {/* Calendar grid skeleton */}
      <div className="grid grid-cols-7 gap-0">
        {cells.map((_, index) => (
          <div
            key={index}
            className="border-b border-r border-gray-100 p-2 h-20 md:h-24 lg:h-28"
          >
            {/* Date number skeleton */}
            <div className="h-4 bg-gray-200 rounded w-6 mb-2"></div>
            
            {/* Tithi skeleton */}
            <div className="h-3 bg-gray-100 rounded w-12 mb-1"></div>
            
            {/* Event indicators skeleton */}
            <div className="flex space-x-1">
              <div className="h-2 w-2 bg-gray-100 rounded-full"></div>
              <div className="h-2 w-2 bg-gray-100 rounded-full"></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};