import { type ExcelColumnMapping, type ExcelFileStructure } from '@/lib/api';

interface ColumnMappingFormProps {
  fileStructure: ExcelFileStructure;
  mapping: ExcelColumnMapping;
  onChange: (mapping: ExcelColumnMapping) => void;
}

const fieldLabels = {
  transaction_description: 'Описание транзакции',
  transaction_amount: 'Сумма транзакции',
  transaction_type: 'Тип транзакции',
  transaction_date: 'Дата транзакции',
  transaction_category: 'Категория транзакции',
  category_name: 'Название категории',
  category_type: 'Тип категории',
  prescribed_expanse_description: 'Описание обязательной траты',
  prescribed_expanse_amount: 'Сумма обязательной траты',
  prescribed_expanse_frequency: 'Частота обязательной траты',
  prescribed_expanse_date: 'Дата обязательной траты',
  prescribed_expanse_category: 'Категория обязательной траты',
};

const fieldGroups = {
  'Транзакции': [
    'transaction_description',
    'transaction_amount',
    'transaction_type',
    'transaction_date',
    'transaction_category',
  ],
  'Категории': ['category_name', 'category_type'],
  'Обязательные траты': [
    'prescribed_expanse_description',
    'prescribed_expanse_amount',
    'prescribed_expanse_frequency',
    'prescribed_expanse_date',
    'prescribed_expanse_category',
  ],
};

export function ColumnMappingForm({ fileStructure, mapping, onChange }: ColumnMappingFormProps) {
  const handleFieldChange = (field: keyof ExcelColumnMapping, value: string) => {
    const newMapping = { ...mapping };
    if (value === '') {
      delete newMapping[field];
    } else {
      newMapping[field] = value;
    }
    onChange(newMapping);
  };

  return (
    <div className="space-y-6">
      {Object.entries(fieldGroups).map(([groupName, fields]) => (
        <div key={groupName} className="border rounded-lg p-4">
          <h3 className="text-lg font-semibold mb-4">{groupName}</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {fields.map((field) => {
              const fieldKey = field as keyof ExcelColumnMapping;
              return (
                <div key={field}>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {fieldLabels[fieldKey]}
                  </label>
                  <select
                    value={mapping[fieldKey] || ''}
                    onChange={(e) => handleFieldChange(fieldKey, e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="">Не использовать</option>
                    {fileStructure.columns.map((col) => (
                      <option key={col} value={col}>
                        {col}
                      </option>
                    ))}
                  </select>
                </div>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}

