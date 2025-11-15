import { useState } from 'react';
import { transactionsApi, type Transaction, type Category } from '@/lib/api';

interface TransactionFormProps {
  transaction?: Transaction | null;
  categories: Category[];
  onClose: () => void;
  onSuccess: () => void;
}

export function TransactionForm({
  transaction,
  categories,
  onClose,
  onSuccess,
}: TransactionFormProps) {
  const [formData, setFormData] = useState({
    category_id: transaction?.category_id || categories[0]?.id || 0,
    description: transaction?.description || '',
    amount: transaction?.amount || '',
    type: (transaction?.type ?? 0) as 0 | 1,
    date_time: transaction?.date_time
      ? new Date(transaction.date_time).toISOString().slice(0, 16)
      : new Date().toISOString().slice(0, 16),
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const submitData = {
        ...formData,
        date_time: new Date(formData.date_time).toISOString(),
      };

      if (transaction) {
        await transactionsApi.update(transaction.id, submitData);
      } else {
        await transactionsApi.create(submitData);
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
            {transaction ? 'Редактировать транзакцию' : 'Добавить транзакцию'}
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
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Тип транзакции
            </label>
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

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Категория</label>
            <select
              value={formData.category_id}
              onChange={(e) =>
                setFormData({ ...formData, category_id: Number(e.target.value) })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
            >
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Описание</label>
            <input
              type="text"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Сумма</label>
            <input
              type="number"
              step="0.01"
              value={formData.amount}
              onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Дата и время</label>
            <input
              type="datetime-local"
              value={formData.date_time}
              onChange={(e) => setFormData({ ...formData, date_time: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              required
            />
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

