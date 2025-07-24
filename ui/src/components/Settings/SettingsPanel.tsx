import React from 'react';
import { X, Globe, Calculator, Clock, Map } from 'lucide-react';
import { Settings } from '../../types/panchangam';

interface SettingsPanelProps {
  settings: Settings;
  onSettingsChange: (settings: Settings) => void;
  onClose: () => void;
}

export const SettingsPanel: React.FC<SettingsPanelProps> = ({
  settings,
  onSettingsChange,
  onClose
}) => {
  const handleChange = (key: keyof Settings, value: any) => {
    onSettingsChange({
      ...settings,
      [key]: value
    });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-xl shadow-2xl max-w-md w-full max-h-[80vh] overflow-hidden">
        {/* Header */}
        <div className="bg-gray-800 text-white p-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold">Settings</h2>
            <button
              onClick={onClose}
              className="p-1 hover:bg-gray-700 rounded-full transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Calculation Method */}
          <div>
            <label className="flex items-center space-x-2 text-sm font-medium text-gray-700 mb-3">
              <Calculator className="w-4 h-4" />
              <span>Calculation Method</span>
            </label>
            <div className="space-y-2">
              <label className="flex items-center space-x-3">
                <input
                  type="radio"
                  name="calculation_method"
                  value="Drik"
                  checked={settings.calculation_method === 'Drik'}
                  onChange={(e) => handleChange('calculation_method', e.target.value)}
                  className="text-orange-500 focus:ring-orange-500"
                />
                <div>
                  <div className="font-medium">Drik (Modern)</div>
                  <div className="text-xs text-gray-500">Based on astronomical calculations</div>
                </div>
              </label>
              <label className="flex items-center space-x-3">
                <input
                  type="radio"
                  name="calculation_method"
                  value="Vakya"
                  checked={settings.calculation_method === 'Vakya'}
                  onChange={(e) => handleChange('calculation_method', e.target.value)}
                  className="text-orange-500 focus:ring-orange-500"
                />
                <div>
                  <div className="font-medium">Vakya (Traditional)</div>
                  <div className="text-xs text-gray-500">Based on traditional formulas</div>
                </div>
              </label>
            </div>
          </div>

          {/* Language */}
          <div>
            <label className="flex items-center space-x-2 text-sm font-medium text-gray-700 mb-3">
              <Globe className="w-4 h-4" />
              <span>Language</span>
            </label>
            <select
              value={settings.locale}
              onChange={(e) => handleChange('locale', e.target.value)}
              className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent"
            >
              <option value="en">English</option>
              <option value="hi">हिन्दी (Hindi)</option>
              <option value="ta">தமிழ் (Tamil)</option>
            </select>
          </div>

          {/* Time Format */}
          <div>
            <label className="flex items-center space-x-2 text-sm font-medium text-gray-700 mb-3">
              <Clock className="w-4 h-4" />
              <span>Time Format</span>
            </label>
            <div className="space-y-2">
              <label className="flex items-center space-x-3">
                <input
                  type="radio"
                  name="time_format"
                  value="12"
                  checked={settings.time_format === '12'}
                  onChange={(e) => handleChange('time_format', e.target.value)}
                  className="text-orange-500 focus:ring-orange-500"
                />
                <span>12-hour (AM/PM)</span>
              </label>
              <label className="flex items-center space-x-3">
                <input
                  type="radio"
                  name="time_format"
                  value="24"
                  checked={settings.time_format === '24'}
                  onChange={(e) => handleChange('time_format', e.target.value)}
                  className="text-orange-500 focus:ring-orange-500"
                />
                <span>24-hour</span>
              </label>
            </div>
          </div>

          {/* Region */}
          <div>
            <label className="flex items-center space-x-2 text-sm font-medium text-gray-700 mb-3">
              <Map className="w-4 h-4" />
              <span>Region</span>
            </label>
            <select
              value={settings.region}
              onChange={(e) => handleChange('region', e.target.value)}
              className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent"
            >
              <option value="Tamil Nadu">Tamil Nadu</option>
              <option value="Karnataka">Karnataka</option>
              <option value="Andhra Pradesh">Andhra Pradesh</option>
              <option value="Kerala">Kerala</option>
              <option value="Maharashtra">Maharashtra</option>
              <option value="Gujarat">Gujarat</option>
              <option value="Rajasthan">Rajasthan</option>
              <option value="Uttar Pradesh">Uttar Pradesh</option>
              <option value="West Bengal">West Bengal</option>
              <option value="Delhi">Delhi</option>
            </select>
          </div>
        </div>

        {/* Footer */}
        <div className="bg-gray-50 px-6 py-4 border-t">
          <div className="flex justify-end space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-gray-200 rounded-lg hover:bg-gray-300 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors"
            >
              Save Settings
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};