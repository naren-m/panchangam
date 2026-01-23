import React, { useState } from 'react';
import { Download, FileText, FileSpreadsheet, Table as TableIcon } from 'lucide-react';

interface TableColumn {
  key: string;
  label: string;
  sortable?: boolean;
  render?: (value: any, row: any) => React.ReactNode;
}

interface TabularViewProps {
  columns: TableColumn[];
  data: any[];
  title?: string;
  exportFormats?: ('csv' | 'json' | 'pdf')[];
  searchable?: boolean;
  sortable?: boolean;
  responsive?: boolean;
}

/**
 * TabularView - Displays data in a sortable, searchable table format
 * with export capabilities
 *
 * @component
 */
export const TabularView: React.FC<TabularViewProps> = ({
  columns,
  data,
  title,
  exportFormats = ['csv', 'json'],
  searchable = true,
  sortable = true,
  responsive = true,
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [sortColumn, setSortColumn] = useState<string | null>(null);
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');

  // Filter data based on search term
  const filteredData = searchTerm
    ? data.filter((row) =>
        Object.values(row).some((value) =>
          String(value).toLowerCase().includes(searchTerm.toLowerCase())
        )
      )
    : data;

  // Sort data
  const sortedData = sortColumn
    ? [...filteredData].sort((a, b) => {
        const aVal = a[sortColumn];
        const bVal = b[sortColumn];

        if (aVal === bVal) return 0;

        const comparison = aVal < bVal ? -1 : 1;
        return sortDirection === 'asc' ? comparison : -comparison;
      })
    : filteredData;

  // Handle sort
  const handleSort = (columnKey: string) => {
    if (!sortable) return;

    if (sortColumn === columnKey) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortColumn(columnKey);
      setSortDirection('asc');
    }
  };

  // Export functions
  const exportToCSV = () => {
    const headers = columns.map((col) => col.label).join(',');
    const rows = sortedData
      .map((row) =>
        columns.map((col) => JSON.stringify(row[col.key] ?? '')).join(',')
      )
      .join('\n');

    const csv = `${headers}\n${rows}`;
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${title || 'data'}.csv`;
    link.click();
    URL.revokeObjectURL(url);
  };

  const exportToJSON = () => {
    const json = JSON.stringify(sortedData, null, 2);
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${title || 'data'}.json`;
    link.click();
    URL.revokeObjectURL(url);
  };

  const exportToPDF = () => {
    // Simple PDF export (in production, use a library like jsPDF)
    alert('PDF export would be implemented with a library like jsPDF');
  };

  const handleExport = (format: 'csv' | 'json' | 'pdf') => {
    switch (format) {
      case 'csv':
        exportToCSV();
        break;
      case 'json':
        exportToJSON();
        break;
      case 'pdf':
        exportToPDF();
        break;
    }
  };

  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      {/* Header */}
      <div className="px-4 py-3 border-b border-gray-200 bg-gray-50">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div>
            {title && (
              <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
            )}
            <p className="text-sm text-gray-600">
              Showing {sortedData.length} of {data.length} entries
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-3">
            {searchable && (
              <input
                type="text"
                placeholder="Search..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                aria-label="Search table data"
              />
            )}

            {exportFormats.length > 0 && (
              <div className="flex gap-2">
                {exportFormats.map((format) => (
                  <button
                    key={format}
                    onClick={() => handleExport(format)}
                    className="px-3 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors flex items-center gap-2"
                    aria-label={`Export as ${format.toUpperCase()}`}
                  >
                    {format === 'csv' && <FileSpreadsheet className="w-4 h-4" />}
                    {format === 'json' && <FileText className="w-4 h-4" />}
                    {format === 'pdf' && <Download className="w-4 h-4" />}
                    <span className="uppercase">{format}</span>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Table */}
      <div className={responsive ? 'overflow-x-auto' : ''}>
        <table className="min-w-full divide-y divide-gray-200" role="table">
          <thead className="bg-gray-100">
            <tr>
              {columns.map((column) => (
                <th
                  key={column.key}
                  className={`px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider ${
                    column.sortable !== false && sortable
                      ? 'cursor-pointer hover:bg-gray-200'
                      : ''
                  }`}
                  onClick={() =>
                    column.sortable !== false && handleSort(column.key)
                  }
                  scope="col"
                >
                  <div className="flex items-center gap-2">
                    {column.label}
                    {sortColumn === column.key && (
                      <span aria-label={`Sorted ${sortDirection}ending`}>
                        {sortDirection === 'asc' ? '↑' : '↓'}
                      </span>
                    )}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {sortedData.length === 0 ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-6 py-4 text-center text-gray-500"
                >
                  No data available
                </td>
              </tr>
            ) : (
              sortedData.map((row, rowIndex) => (
                <tr
                  key={rowIndex}
                  className="hover:bg-gray-50 transition-colors"
                >
                  {columns.map((column) => (
                    <td
                      key={column.key}
                      className="px-6 py-4 whitespace-nowrap text-sm text-gray-900"
                    >
                      {column.render
                        ? column.render(row[column.key], row)
                        : row[column.key]}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

interface DataViewSwitcherProps {
  data: any[];
  columns: TableColumn[];
  title?: string;
  defaultView?: 'table' | 'cards' | 'list';
}

/**
 * DataViewSwitcher - Allows switching between different data presentation formats
 *
 * @component
 */
export const DataViewSwitcher: React.FC<DataViewSwitcherProps> = ({
  data,
  columns,
  title,
  defaultView = 'table',
}) => {
  const [currentView, setCurrentView] = useState(defaultView);

  const renderTableView = () => (
    <TabularView
      columns={columns}
      data={data}
      title={title}
      exportFormats={['csv', 'json']}
    />
  );

  const renderCardsView = () => (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {data.map((item, index) => (
        <div
          key={index}
          className="bg-white rounded-lg shadow p-4 hover:shadow-lg transition-shadow"
        >
          {columns.map((column) => (
            <div key={column.key} className="mb-2">
              <span className="font-semibold text-gray-700">{column.label}:</span>{' '}
              <span className="text-gray-600">
                {column.render ? column.render(item[column.key], item) : item[column.key]}
              </span>
            </div>
          ))}
        </div>
      ))}
    </div>
  );

  const renderListView = () => (
    <div className="space-y-2">
      {data.map((item, index) => (
        <div
          key={index}
          className="bg-white rounded-lg shadow p-4 hover:bg-gray-50 transition-colors"
        >
          <div className="flex flex-wrap gap-4">
            {columns.map((column) => (
              <div key={column.key} className="flex-1 min-w-[200px]">
                <div className="text-xs font-semibold text-gray-500 uppercase">
                  {column.label}
                </div>
                <div className="text-sm text-gray-900 mt-1">
                  {column.render ? column.render(item[column.key], item) : item[column.key]}
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );

  return (
    <div>
      {/* View switcher buttons */}
      <div className="mb-4 flex gap-2 border-b border-gray-200 pb-2">
        <button
          onClick={() => setCurrentView('table')}
          className={`px-4 py-2 rounded-t-lg transition-colors ${
            currentView === 'table'
              ? 'bg-indigo-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          }`}
          aria-label="Table view"
          aria-pressed={currentView === 'table'}
        >
          <TableIcon className="w-5 h-5" />
        </button>
        <button
          onClick={() => setCurrentView('cards')}
          className={`px-4 py-2 rounded-t-lg transition-colors ${
            currentView === 'cards'
              ? 'bg-indigo-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          }`}
          aria-label="Cards view"
          aria-pressed={currentView === 'cards'}
        >
          Cards
        </button>
        <button
          onClick={() => setCurrentView('list')}
          className={`px-4 py-2 rounded-t-lg transition-colors ${
            currentView === 'list'
              ? 'bg-indigo-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          }`}
          aria-label="List view"
          aria-pressed={currentView === 'list'}
        >
          List
        </button>
      </div>

      {/* Render current view */}
      {currentView === 'table' && renderTableView()}
      {currentView === 'cards' && renderCardsView()}
      {currentView === 'list' && renderListView()}
    </div>
  );
};

export default TabularView;
