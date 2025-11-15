import { useState } from 'react';
import { categoriesApi, type Category } from '@/lib/api';

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
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">
            {category ? 'Редактировать категорию' : 'Добавить категорию'}
          </h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            ✕
          </button>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-50 text-red-800 rounded-md text-sm">{error}</div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Тип</label>
            <select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: Number(e.target.value) as 0 | 1 })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
            >
              <option value={0}>Доход</option>
              <option value={1}>Расход</option>
            </select>
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? 'Сохранение...' : 'Сохранить'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

