import { useMemo } from 'react';
import { PanchangamData } from '../types/panchangam';

interface UseDayDetailOptions {
  selectedDate: Date | null;
  panchangamData: Record<string, PanchangamData>;
  loading: boolean;
  error: string | null;
  retry: () => void;
}

export const useDayDetail = ({
  selectedDate,
  panchangamData,
  loading,
  error,
  retry,
}: UseDayDetailOptions) => {
  const dayDetailData = useMemo(() => {
    if (!selectedDate) return null;
    const dateStr = selectedDate.toISOString().split('T')[0];
    return panchangamData[dateStr] || null;
  }, [selectedDate, panchangamData]);

  const hasDataForDate = useMemo(() => {
    if (!selectedDate) return false;
    const dateStr = selectedDate.toISOString().split('T')[0];
    return dateStr in panchangamData;
  }, [selectedDate, panchangamData]);

  // Determine if we're loading specifically for the selected date
  const isLoadingDayDetail = loading && !hasDataForDate;

  // Only show error if we don't have data for the specific date
  const dayDetailError = error && !hasDataForDate ? error : null;

  return {
    data: dayDetailData,
    isLoading: isLoadingDayDetail,
    error: dayDetailError,
    retry,
  };
};