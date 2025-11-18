import React, { useState } from 'react';
import { ChevronDown, ChevronRight, Info } from 'lucide-react';

interface ProgressiveDisclosurePanelProps {
  title: string;
  summary: React.ReactNode;
  details: React.ReactNode;
  defaultExpanded?: boolean;
  icon?: React.ReactNode;
  level?: 'primary' | 'secondary' | 'tertiary';
  ariaLabel?: string;
}

/**
 * ProgressiveDisclosurePanel - A component that implements progressive disclosure pattern
 * Shows summary information first, with option to drill down into details
 *
 * @component
 * @example
 * <ProgressiveDisclosurePanel
 *   title="Tithi"
 *   summary={<div>Shukla Paksha Pratipada</div>}
 *   details={<div>Detailed tithi information...</div>}
 * />
 */
export const ProgressiveDisclosurePanel: React.FC<ProgressiveDisclosurePanelProps> = ({
  title,
  summary,
  details,
  defaultExpanded = false,
  icon,
  level = 'primary',
  ariaLabel,
}) => {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);

  const handleToggle = () => {
    setIsExpanded(!isExpanded);
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleToggle();
    }
  };

  const levelStyles = {
    primary: 'bg-white border-2 border-indigo-200 shadow-md',
    secondary: 'bg-gray-50 border border-gray-300',
    tertiary: 'bg-gray-100 border border-gray-200',
  };

  const headerStyles = {
    primary: 'text-lg font-semibold text-indigo-900',
    secondary: 'text-base font-medium text-gray-800',
    tertiary: 'text-sm font-normal text-gray-700',
  };

  return (
    <div
      className={`rounded-lg overflow-hidden transition-all duration-200 ${levelStyles[level]}`}
      role="region"
      aria-label={ariaLabel || title}
    >
      {/* Header with Summary */}
      <button
        className="w-full px-4 py-3 flex items-center justify-between cursor-pointer hover:bg-opacity-90 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        onClick={handleToggle}
        onKeyDown={handleKeyDown}
        aria-expanded={isExpanded}
        aria-controls={`disclosure-content-${title}`}
        type="button"
      >
        <div className="flex items-center gap-3 flex-1">
          {icon && <div className="flex-shrink-0">{icon}</div>}
          <div className="flex-1 text-left">
            <h3 className={headerStyles[level]}>{title}</h3>
            {!isExpanded && (
              <div className="mt-1 text-sm text-gray-600">{summary}</div>
            )}
          </div>
        </div>
        <div className="flex-shrink-0 ml-3">
          {isExpanded ? (
            <ChevronDown className="w-5 h-5 text-gray-600" aria-hidden="true" />
          ) : (
            <ChevronRight className="w-5 h-5 text-gray-600" aria-hidden="true" />
          )}
        </div>
      </button>

      {/* Expanded Details */}
      {isExpanded && (
        <div
          id={`disclosure-content-${title}`}
          className="px-4 py-3 border-t border-gray-200 animate-fadeIn"
          role="region"
          aria-live="polite"
        >
          <div className="mb-3 text-sm font-medium text-gray-700">{summary}</div>
          <div className="text-sm text-gray-600">{details}</div>
        </div>
      )}
    </div>
  );
};

interface ProgressiveDisclosureGroupProps {
  children: React.ReactNode;
  title?: string;
  description?: string;
  allowMultiple?: boolean;
  className?: string;
}

/**
 * ProgressiveDisclosureGroup - Container for multiple disclosure panels
 * Can be configured to allow only one panel open at a time (accordion) or multiple
 *
 * @component
 */
export const ProgressiveDisclosureGroup: React.FC<ProgressiveDisclosureGroupProps> = ({
  children,
  title,
  description,
  allowMultiple = true,
  className = '',
}) => {
  return (
    <div className={`space-y-3 ${className}`} role="group" aria-label={title}>
      {title && (
        <div className="mb-4">
          <h2 className="text-xl font-bold text-gray-900">{title}</h2>
          {description && (
            <p className="mt-1 text-sm text-gray-600">{description}</p>
          )}
        </div>
      )}
      <div className={allowMultiple ? 'space-y-3' : 'space-y-2'}>
        {children}
      </div>
    </div>
  );
};

interface SummaryDetailViewProps {
  data: {
    label: string;
    summary: string | React.ReactNode;
    details?: string | React.ReactNode;
    icon?: React.ReactNode;
  }[];
  showIcons?: boolean;
  responsive?: boolean;
}

/**
 * SummaryDetailView - Displays data in a summary-first format with optional drill-down
 * Optimized for responsive design
 *
 * @component
 */
export const SummaryDetailView: React.FC<SummaryDetailViewProps> = ({
  data,
  showIcons = true,
  responsive = true,
}) => {
  const gridClass = responsive
    ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3'
    : 'space-y-3';

  return (
    <div className={gridClass}>
      {data.map((item, index) => (
        <ProgressiveDisclosurePanel
          key={index}
          title={item.label}
          summary={item.summary}
          details={item.details || item.summary}
          icon={showIcons ? item.icon : undefined}
          level="secondary"
          ariaLabel={`${item.label} information`}
        />
      ))}
    </div>
  );
};

// Accessibility helper component
export const ScreenReaderOnly: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <span className="sr-only">
      {children}
    </span>
  );
};

// Export types for use in other components
export type { ProgressiveDisclosurePanelProps, ProgressiveDisclosureGroupProps, SummaryDetailViewProps };
