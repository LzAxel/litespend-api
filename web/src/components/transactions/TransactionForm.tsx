import { useState } from 'react';
import { transactionsApi, type Transaction, type Category } from '@/lib/api';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

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
    category_id: transaction?.category_id || 0,
    description: transaction?.description || '',
    amount: transaction?.amount || '',
    date: transaction?.date
      ? new Date(transaction.date).toISOString().slice(0, 16)
      : new Date().toISOString().slice(0, 16),
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Фильтруем категории по типу транзакции
  const filteredCategories = categories.filter((cat) => cat.type === 0 && +formData.amount > 0 || cat.type === 1 && +formData.amount <= 0);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const submitData = {
        ...formData,
        date: new Date(formData.date).toISOString(),
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
    <div>
      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-800 rounded-md text-sm">{error}</div>
      )}
      <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <Label className="mb-1" htmlFor="category">Категория</Label>
          <select
            id="category"
            value={formData.category_id}
            onChange={(e) => setFormData({ ...formData, category_id: Number(e.target.value) })}
            className="mt-1 block w-full border rounded p-2"
            required
            disabled={filteredCategories.length === 0}
          >
            {filteredCategories.length === 0 ? (
              <option value={0}>Нет доступных категорий</option>
            ) : (
              filteredCategories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))
            )}
          </select>
          {filteredCategories.length === 0 && (
            <p className="mt-1 text-sm text-yellow-600">
              Создайте категорию для {+formData.amount > 0 ? 'доходов' : 'расходов'}
            </p>
          )}
        </div>

        <div>
          <Label className="mb-1" htmlFor="description">Описание</Label>
          <Input id="description" type="text" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} />
        </div>

        <div>
          <Label className="mb-1" htmlFor="amount">Сумма</Label>
          <Input id="amount" type="text" value={formData.amount} onChange={(e) => setFormData({ ...formData, amount: e.target.value })} required />
        </div>

        <div>
          <Label className="mb-1" htmlFor="date">Дата и время</Label>
          <Input id="date" type="datetime-local" value={formData.date} onChange={(e) => setFormData({ ...formData, date: e.target.value })} required />
        </div>

        <div className="md:col-span-2 flex flex-col-reverse sm:flex-row justify-end gap-2 mt-2">
          <Button type="button" variant="outline" onClick={onClose}>Отмена</Button>
          <Button type="submit" disabled={loading}>{loading ? 'Сохранение...' : 'Сохранить'}</Button>
        </div>
      </form>
    </div>
  );
}

