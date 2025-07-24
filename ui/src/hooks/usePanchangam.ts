import { useState, useEffect } from 'react';
import { PanchangamData, Settings } from '../types/panchangam';
import { panchangamApi } from '../services/panchangamApi';
import { formatDateForApi } from '../utils/dateHelpers';

export const usePanchangam = (date: Date, settings: Settings) => {
  const [data, setData] = useState<PanchangamData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPanchangam = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await panchangamApi.getPanchangam({
          date: formatDateForApi(date),
          latitude: settings.location.latitude,
          longitude: settings.location.longitude,
          timezone: settings.location.timezone,
          region: settings.region,
          calculation_method: settings.calculation_method,
          locale: settings.locale
        });

        setData(response);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch panchangam data');
      } finally {
        setLoading(false);
      }
    };

    fetchPanchangam();
  }, [date, settings]);

  return { data, loading, error };
};

export const usePanchangamRange = (startDate: Date, endDate: Date, settings: Settings) => {
  const [data, setData] = useState<Record<string, PanchangamData>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPanchangamRange = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await panchangamApi.getPanchangamRange(
          formatDateForApi(startDate),
          formatDateForApi(endDate),
          {
            latitude: settings.location.latitude,
            longitude: settings.location.longitude,
            timezone: settings.location.timezone,
            region: settings.region,
            calculation_method: settings.calculation_method,
            locale: settings.locale
          }
        );

        const dataMap: Record<string, PanchangamData> = {};
        response.forEach(item => {
          dataMap[item.date] = item;
        });

        setData(dataMap);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch panchangam data');
      } finally {
        setLoading(false);
      }
    };

    fetchPanchangamRange();
  }, [startDate, endDate, settings]);

  return { data, loading, error };
};