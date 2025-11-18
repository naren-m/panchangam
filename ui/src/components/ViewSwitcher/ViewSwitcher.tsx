import React from 'react';
import { Calendar, Table2, BarChart3 } from 'lucide-react';

export type ViewMode = 'calendar' | 'table' | 'graph';

interface ViewSwitcherProps {
  currentView: ViewMode;
  onViewChange: (view: ViewMode) => void;
  className?: string;
}

export const ViewSwitcher: React.FC<ViewSwitcherProps> = ({
  currentView,
  onViewChange,
  className = ''
}) => {
  const views: Array<{ mode: ViewMode; icon: React.ReactNode; label: string; description: string }> = [
    {
      mode: 'calendar',
      icon: <Calendar className="w-5 h-5" />,
      label: 'Calendar',
      description: 'Traditional calendar grid view'
    },
    {
      mode: 'table',
      icon: <Table2 className="w-5 h-5" />,
      label: 'Table',
      description: 'Detailed tabular data view'
    },
    {
      mode: 'graph',
      icon: <BarChart3 className="w-5 h-5" />,
      label: 'Analytics',
      description: 'Visual charts and graphs'
    }
  ];

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <span className="text-sm font-medium text-gray-700 mr-2 hidden sm:inline">View:</span>
      <div className="inline-flex bg-white rounded-lg shadow-md border border-orange-200 p-1">
        {views.map(({ mode, icon, label, description }) => {
          const isActive = currentView === mode;
          return (
            <button
              key={mode}
              onClick={() => onViewChange(mode)}
              className={`
                flex items-center gap-2 px-4 py-2 rounded-md transition-all duration-200
                ${isActive
                  ? 'bg-gradient-to-r from-orange-400 to-orange-500 text-white shadow-md'
                  : 'text-gray-700 hover:bg-orange-50'
                }
              `}
              title={description}
              aria-label={`Switch to ${label} view`}
              aria-pressed={isActive}
            >
              {icon}
              <span className="text-sm font-medium hidden sm:inline">{label}</span>
            </button>
          );
        })}
      </div>
    </div>
  );
};
