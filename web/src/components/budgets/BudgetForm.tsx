import { useEffect, useState } from 'react';
import { budgetsApi, categoriesApi, type BudgetDetailed, type CreateBudgetRequest, type UpdateBudgetRequest, type Category } from '@/lib/api';

interface BudgetFormProps {
  budget?: BudgetDetailed | null;
  defaultYear?: number;
  defaultMonth?: number;
  onClose: () => void;
  onSuccess: () => void;
}

export function BudgetForm({ budget, defaultYear, defaultMonth, onClose, onSuccess }: BudgetFormProps) {
  const [categories, setCategories] = useState<Category[]>([]);
  const [form, setForm] = useState<{
    category_id: number | '';
    year: number | '';
    month: number | '';
    budgeted: string;
  }>({
    category_id: budget?.category_id ?? '',
    year: budget?.year ?? defaultYear ?? new Date().getFullYear(),
    month: budget?.month ?? defaultMonth ?? new Date().getMonth() + 1,
    budgeted: budget?.budgeted ?? '',
  });
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    categoriesApi.getAll().then((res) => setCategories(res.data));
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: name === 'category_id' || name === 'year' || name === 'month' ? Number(value) : value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.category_id || !form.year || !form.month || !form.budgeted) return;
    setSubmitting(true);
    try {
      if (budget) {
        const dto: UpdateBudgetRequest = {
          category_id: form.category_id,
          year: form.year,
          month: form.month,
          budgeted: form.budgeted,
        };
        await budgetsApi.update(budget.id, dto);
      } else {
        const dto: CreateBudgetRequest = {
          category_id: form.category_id,
          year: form.year,
          month: form.month,
          budgeted: form.budgeted,
        };
        await budgetsApi.create(dto);
      }
      onSuccess();
    } catch (err) {
      console.error('Budget submit error', err);
      alert('Ошибка сохранения бюджета');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-white shadow sm:rounded-lg p-4 mb-6">
      <h2 className="text-xl font-semibold mb-4">{budget ? 'Редактировать бюджет' : 'Создать бюджет'}</h2>
      <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700">Категория</label>
          <select name="category_id" value={form.category_id} onChange={handleChange} className="mt-1 block w-full border rounded p-2">
            <option value="">—</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Год</label>
          <input type="number" name="year" value={form.year} onChange={handleChange} className="mt-1 block w-full border rounded p-2" />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Месяц</label>
          <input type="number" name="month" value={form.month} onChange={handleChange} className="mt-1 block w-full border rounded p-2" min={1} max={12} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Сумма</label>
          <input type="number" step="0.01" name="budgeted" value={form.budgeted} onChange={handleChange} className="mt-1 block w-full border rounded p-2" />
        </div>
        <div className="md:col-span-4 flex gap-2">
          <button type="submit" disabled={submitting} className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
            {submitting ? 'Сохранение...' : 'Сохранить'}
          </button>
          <button type="button" onClick={onClose} className="px-4 py-2 border rounded">
            Отмена
          </button>
        </div>
      </form>
    </div>
  );
}
