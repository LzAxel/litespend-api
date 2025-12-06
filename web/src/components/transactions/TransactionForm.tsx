import { useEffect, useState } from 'react';
import { transactionsApi, accountsApi, type Transaction, type Category, type Account } from '@/lib/api';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

interface TransactionFormProps {
  transaction?: Transaction | null;
  categories: Category[];
  accounts?: Account[];
  onClose: () => void;
  onSuccess: () => void;
}

export function TransactionForm({
  transaction,
  categories,
  accounts: accountsProp,
  onClose,
  onSuccess,
}: TransactionFormProps) {
  const [accounts, setAccounts] = useState<Account[]>(accountsProp || []);
  const [formData, setFormData] = useState({
    account_id: transaction?.account_id || 0,
    category_id: transaction?.category_id || 0,
    note: (transaction as any)?.note || '',
    amount: transaction?.amount || '',
    date: transaction?.date
      ? new Date(transaction.date).toISOString().slice(0, 16)
      : new Date().toISOString().slice(0, 16),
    is_cleared: (transaction as any)?.is_cleared ?? false,
    is_approved: (transaction as any)?.is_approved ?? false,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!accountsProp) {
      accountsApi.getAll().then((res) => setAccounts(res.data.accounts)).catch(() => {});
    }
  }, []);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const submitData = {
        account_id: Number(formData.account_id),
        category_id: Number(formData.category_id),
        note: formData.note,
        amount: String(formData.amount),
        date: new Date(formData.date).toISOString(),
        is_cleared: Boolean(formData.is_cleared),
        is_approved: Boolean(formData.is_approved),
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
          <Label className="mb-1" htmlFor="account">Счёт</Label>
          <select
            id="account"
            value={formData.account_id}
            onChange={(e) => setFormData({ ...formData, account_id: Number(e.target.value) })}
            className="mt-1 block w-full border rounded p-2"
            required
          >
            <option value={0}>—</option>
            {accounts.map((acc) => (
              <option key={acc.id} value={acc.id}>{acc.name}</option>
            ))}
          </select>
        </div>
        <div>
          <Label className="mb-1" htmlFor="category">Категория</Label>
          <select
            id="category"
            value={formData.category_id}
            onChange={(e) => setFormData({ ...formData, category_id: Number(e.target.value) })}
            className="mt-1 block w-full border rounded p-2"
            required
          >
            <option value={0}>—</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <Label className="mb-1" htmlFor="note">Примечание</Label>
          <Input id="note" type="text" value={formData.note} onChange={(e) => setFormData({ ...formData, note: e.target.value })} />
        </div>

        <div>
          <Label className="mb-1" htmlFor="amount">Сумма</Label>
          <Input id="amount" type="text" value={formData.amount} onChange={(e) => setFormData({ ...formData, amount: e.target.value })} required />
        </div>

        <div>
          <Label className="mb-1" htmlFor="date">Дата и время</Label>
          <Input id="date" type="datetime-local" value={formData.date} onChange={(e) => setFormData({ ...formData, date: e.target.value })} required />
        </div>

        <div>
          <Label className="mb-1" htmlFor="is_cleared">Проведена</Label>
          <select
            id="is_cleared"
            value={String(formData.is_cleared)}
            onChange={(e) => setFormData({ ...formData, is_cleared: e.target.value === 'true' })}
            className="mt-1 block w-full border rounded p-2"
          >
            <option value="false">Нет</option>
            <option value="true">Да</option>
          </select>
        </div>

        <div>
          <Label className="mb-1" htmlFor="is_approved">Подтверждена</Label>
          <select
            id="is_approved"
            value={String(formData.is_approved)}
            onChange={(e) => setFormData({ ...formData, is_approved: e.target.value === 'true' })}
            className="mt-1 block w-full border rounded p-2"
          >
            <option value="false">Нет</option>
            <option value="true">Да</option>
          </select>
        </div>

        <div className="md:col-span-2 flex flex-col-reverse sm:flex-row justify-end gap-2 mt-2">
          <Button type="button" variant="outline" onClick={onClose}>Отмена</Button>
          <Button type="submit" disabled={loading}>{loading ? 'Сохранение...' : 'Сохранить'}</Button>
        </div>
      </form>
    </div>
  );
}

