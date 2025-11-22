import { useEffect, useState } from 'react';
import { budgetsApi, categoriesApi, type BudgetDetailed, type CreateBudgetRequest, type UpdateBudgetRequest, type Category } from '@/lib/api';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

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
    <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div>
        <Label htmlFor="category_id">Категория</Label>
        <select id="category_id" name="category_id" value={form.category_id} onChange={handleChange} className="mt-1 block w-full border rounded p-2">
          <option value="">—</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
      </div>
      <div>
        <Label htmlFor="year">Год</Label>
        <Input type="number" id="year" name="year" value={form.year as number} onChange={handleChange} className="mt-1" />
      </div>
      <div>
        <Label htmlFor="month">Месяц</Label>
        <Input type="number" id="month" name="month" value={form.month as number} onChange={handleChange} className="mt-1" min={1} max={12} />
      </div>
      <div>
        <Label htmlFor="budgeted">Сумма</Label>
        <Input type="number" step="0.01" id="budgeted" name="budgeted" value={form.budgeted} onChange={handleChange} className="mt-1" />
      </div>
      <div className="md:col-span-4 flex flex-col-reverse sm:flex-row gap-2 mt-2">
        <Button type="button" variant="outline" onClick={onClose}>
          Отмена
        </Button>
        <Button type="submit" disabled={submitting}>
          {submitting ? 'Сохранение...' : 'Сохранить'}
        </Button>
      </div>
    </form>
  );
}

