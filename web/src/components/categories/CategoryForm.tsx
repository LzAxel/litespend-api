import { useState } from 'react';
import { categoriesApi, type Category } from '@/lib/api';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

interface CategoryFormProps {
  category?: Category | null;
  onClose: () => void;
  onSuccess: () => void;
}

export function CategoryForm({ category, onClose, onSuccess }: CategoryFormProps) {
  const [formData, setFormData] = useState({
    name: category?.name || '',
    type: (category?.type ?? 0) as 0 | 1,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (category) {
        await categoriesApi.update(category.id, formData);
      } else {
        await categoriesApi.create(formData);
      }
      onSuccess();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка сохранения');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-800 rounded-md text-sm">{error}</div>
      )}
      <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4">
        <div>
          <Label className="mb-1" htmlFor="name">Название</Label>
          <Input id="name" type="text" value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} required />
        </div>
        <div>
          <Label className="mb-1" htmlFor="type">Тип</Label>
          <select id="type" value={formData.type} onChange={(e) => setFormData({ ...formData, type: Number(e.target.value) as 0 | 1 })} className="mt-1 block w-full border rounded p-2" required>
            <option value={0}>Доход</option>
            <option value={1}>Расход</option>
          </select>
        </div>
        <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 mt-2">
          <Button type="button" variant="outline" onClick={onClose}>Отмена</Button>
          <Button type="submit" disabled={loading}>{loading ? 'Сохранение...' : 'Сохранить'}</Button>
        </div>
      </form>
    </div>
  );
}

