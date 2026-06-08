import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';

import { TabularView } from './TabularView';

function readBlobText(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = () => resolve(String(reader.result));
    reader.onerror = () => reject(reader.error);
    reader.readAsText(blob);
  });
}

describe('TabularView', () => {
  let originalCreateObjectURL: typeof URL.createObjectURL;
  let originalRevokeObjectURL: typeof URL.revokeObjectURL;
  let createdBlob: Blob | undefined;

  beforeEach(() => {
    originalCreateObjectURL = URL.createObjectURL;
    originalRevokeObjectURL = URL.revokeObjectURL;
    createdBlob = undefined;

    URL.createObjectURL = vi.fn((blob: Blob) => {
      createdBlob = blob;
      return 'blob:tabular-view-export';
    });
    URL.revokeObjectURL = vi.fn();
    vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => undefined);
  });

  afterEach(() => {
    URL.createObjectURL = originalCreateObjectURL;
    URL.revokeObjectURL = originalRevokeObjectURL;
    vi.restoreAllMocks();
  });

  it('escapes CSV headers and cells before download', async () => {
    render(
      <TabularView
        columns={[
          { key: 'name', label: 'Display, "Name"' },
          { key: 'note', label: 'Note' },
        ]}
        data={[
          { name: 'Abhijit "Best" Muhurta', note: 'value, with comma' },
        ]}
        title="events"
        exportFormats={['csv']}
        searchable={false}
        sortable={false}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: 'Export as CSV' }));

    expect(createdBlob).toBeDefined();
    const csvContent = await readBlobText(createdBlob as Blob);

    expect(csvContent).toContain('"Display, ""Name""",Note');
    expect(csvContent).toContain('"Abhijit ""Best"" Muhurta","value, with comma"');
  });
});
