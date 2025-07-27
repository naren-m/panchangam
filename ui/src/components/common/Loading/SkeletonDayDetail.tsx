import React from 'react';

export const SkeletonDayDetail: React.FC = () => {
  return (
    <div className="bg-white rounded-lg p-6 animate-pulse">
      {/* Header skeleton */}
      <div className="border-b border-gray-200 pb-4 mb-6">
        <div className="h-6 bg-gray-200 rounded w-32 mb-2"></div>
        <div className="h-4 bg-gray-100 rounded w-24"></div>
      </div>

      {/* Panchanga elements skeleton */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        {Array(8).fill(null).map((_, index) => (
          <div key={index} className="bg-gray-50 p-3 rounded-lg">
            <div className="h-4 bg-gray-200 rounded w-16 mb-2"></div>
            <div className="h-5 bg-gray-300 rounded w-20"></div>
          </div>
        ))}
      </div>

      {/* Sun and Moon times skeleton */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-orange-50 p-4 rounded-lg">
          <div className="h-5 bg-orange-200 rounded w-24 mb-2"></div>
          <div className="h-6 bg-orange-300 rounded w-16"></div>
        </div>
        <div className="bg-blue-50 p-4 rounded-lg">
          <div className="h-5 bg-blue-200 rounded w-24 mb-2"></div>
          <div className="h-6 bg-blue-300 rounded w-16"></div>
        </div>
      </div>

      {/* Events section skeleton */}
      <div className="mb-6">
        <div className="h-5 bg-gray-200 rounded w-20 mb-3"></div>
        <div className="space-y-2">
          {Array(4).fill(null).map((_, index) => (
            <div key={index} className="flex justify-between items-center p-2 bg-gray-50 rounded">
              <div className="h-4 bg-gray-200 rounded w-32"></div>
              <div className="h-4 bg-gray-100 rounded w-16"></div>
            </div>
          ))}
        </div>
      </div>

      {/* Festivals section skeleton */}
      <div>
        <div className="h-5 bg-gray-200 rounded w-24 mb-3"></div>
        <div className="h-4 bg-gray-100 rounded w-48"></div>
      </div>
    </div>
  );
};