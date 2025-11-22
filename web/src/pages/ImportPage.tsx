import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { importApi, type ExcelColumnMapping, type ExcelFileStructure, type ImportResult } from '@/lib/api';
import { ColumnMappingForm } from '@/components/import/ColumnMappingForm';

export function ImportPage() {
  const [file, setFile] = useState<File | null>(null);
  const [fileStructure, setFileStructure] = useState<ExcelFileStructure | null>(null);
  const [mapping, setMapping] = useState<ExcelColumnMapping>({});
  const [loading, setLoading] = useState(false);
  const [importing, setImporting] = useState(false);
  const [importResult, setImportResult] = useState<ImportResult | null>(null);
  const [error, setError] = useState('');

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (!selectedFile) return;

    setFile(selectedFile);
    setError('');
    setImportResult(null);
    setMapping({});
    setLoading(true);

    try {
      const response = await importApi.parseFile(selectedFile);
      setFileStructure(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при чтении файла');
      setFile(null);
      setFileStructure(null);
    } finally {
      setLoading(false);
    }
  };

  const handleImport = async () => {
    if (!file || !fileStructure) return;

    setImporting(true);
    setError('');
    setImportResult(null);

    try {
      const response = await importApi.importData(file, mapping);
      setImportResult(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при импорте');
    } finally {
      setImporting(false);
    }
  };

  const hasAnyMapping = Object.values(mapping).some((v) => v !== undefined && v !== '');

  return (
    <div className="px-4 py-6 sm:px-0">
      <h1 className="text-3xl font-bold text-[rgb(var(--app-fg))] mb-6">Импорт из Excel</h1>

      <div className="space-y-6">
        {/* Загрузка файла */}
        <div className="bg-[rgb(var(--card))] shadow rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">1. Выберите файл Excel</h2>
          <input
            type="file"
            accept=".xlsx,.xls"
            onChange={handleFileSelect}
            disabled={loading}
            className="block w-full text-sm text-[rgb(var(--muted-foreground))] file:mr-4 file:py-2 file:px-4 file:rounded-md file:border file:border-[rgb(var(--border))] file:text-sm file:font-semibold file:bg-[rgb(var(--muted))] file:text-[rgb(var(--app-fg))] hover:file:opacity-90"
          />
          {loading && <p className="mt-2 text-sm text-[rgb(var(--muted-foreground))]">Обработка файла...</p>}
          {error && <p className="mt-2 text-sm text-[rgb(var(--destructive))]">{error}</p>}
        </div>

        {/* Структура файла */}
        {fileStructure && (
          <div className="bg-[rgb(var(--card))] shadow rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">2. Структура файла</h2>
            <div className="mb-4">
              <p className="text-sm text-[rgb(var(--muted-foreground))]">
                Найдено столбцов: <span className="font-semibold">{fileStructure.columns.length}</span>
              </p>
              <p className="text-sm text-[rgb(var(--muted-foreground))]">
                Найдено строк данных: <span className="font-semibold">{fileStructure.rows}</span>
              </p>
            </div>
            <div className="mt-4">
              <p className="text-sm font-medium text-[rgb(var(--muted-foreground))] mb-2">Столбцы:</p>
              <div className="flex flex-wrap gap-2">
                {fileStructure.columns.map((col, idx) => (
                  <span
                    key={idx}
                    className="px-3 py-1 bg-[rgb(var(--muted))] text-[rgb(var(--app-fg))] rounded-md text-sm"
                  >
                    {col}
                  </span>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Маппинг столбцов */}
        {fileStructure && (
          <div className="bg-[rgb(var(--card))] shadow rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">3. Настройка соответствия столбцов</h2>
            <p className="text-sm text-[rgb(var(--muted-foreground))] mb-4">
              Выберите, какие столбцы из файла соответствуют полям в системе. Все поля опциональны.
            </p>
            <ColumnMappingForm
              fileStructure={fileStructure}
              mapping={mapping}
              onChange={setMapping}
            />
          </div>
        )}

        {/* Кнопка импорта */}
        {fileStructure && (
          <div className="bg-[rgb(var(--card))] shadow rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">4. Импорт данных</h2>
            {!hasAnyMapping && (
              <p className="text-sm text-[rgb(var(--primary))] mb-4">
                Внимание: Не выбрано ни одного столбца для импорта. Выберите хотя бы один столбец для продолжения.
              </p>
            )}
            <Button onClick={handleImport} disabled={importing || !hasAnyMapping} className="px-6">
              {importing ? 'Импорт...' : 'Начать импорт'}
            </Button>
          </div>
        )}

        {/* Результаты импорта */}
        {importResult && (
          <div className="bg-[rgb(var(--card))] shadow rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">Результаты импорта</h2>
            <div className="space-y-2">
              <p className="text-sm">
                <span className="font-semibold">Транзакций создано:</span>{' '}
                {importResult.transactions_created}
              </p>
              <p className="text-sm">
                <span className="font-semibold">Категорий создано:</span>{' '}
                {importResult.categories_created}
              </p>
              <p className="text-sm">
                <span className="font-semibold">Обязательных трат создано:</span>{' '}
                {importResult.prescribed_expanses_created}
              </p>
              {importResult.errors && importResult.errors.length > 0 && (
                <div className="mt-4">
                  <p className="text-sm font-semibold text-[rgb(var(--destructive))] mb-2">Ошибки:</p>
                  <ul className="list-disc list-inside space-y-1">
                    {importResult.errors.map((err, idx) => (
                      <li key={idx} className="text-sm text-[rgb(var(--destructive))]">
                        {err}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

