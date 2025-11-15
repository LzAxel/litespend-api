import { useState, useEffect } from 'react';
import { transactionsApi, categoriesApi, type Transaction, type Category } from '@/lib/api';
import { TransactionForm } from '@/components/transactions/TransactionForm';
import { TransactionList } from '@/components/transactions/TransactionList';

export function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null);

  useEffect(() => {
    loadData();
  }, [page]);

  const loadData = async () => {
    try {
      setLoading(true);
      const [transactionsRes, categoriesRes] = await Promise.all([
        transactionsApi.getAll(page, 10),
        categoriesApi.getAll(),
      ]);
      setTransactions(transactionsRes.data.data);
      setTotalPages(transactionsRes.data.meta.total_pages);
      setCategories(categoriesRes.data);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingTransaction(null);
    setShowForm(true);
  };

  const handleEdit = (transaction: Transaction) => {
    setEditingTransaction(transaction);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить транзакцию?')) return;
    try {
      await transactionsApi.delete(id);
      loadData();
    } catch (error) {
      console.error('Error deleting transaction:', error);
    }
  };

  const handleFormClose = () => {
    setShowForm(false);
    setEditingTransaction(null);
  };

  const handleFormSuccess = () => {
    handleFormClose();
    loadData();
  };

  if (loading && transactions.length === 0) {
    return <div className="text-center py-8">Загрузка...</div>;
  }

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Транзакции</h1>
        <button
          onClick={handleCreate}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          Добавить транзакцию
        </button>
      </div>

      {showForm && (
        <TransactionForm
          transaction={editingTransaction}
          categories={categories}
          onClose={handleFormClose}
          onSuccess={handleFormSuccess}
        />
      )}

      <TransactionList
        transactions={transactions}
        categories={categories}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />

      {totalPages > 1 && (
        <div className="mt-6 flex justify-center space-x-2">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-4 py-2 border rounded-md disabled:opacity-50"
          >
            Назад
          </button>
          <span className="px-4 py-2">
            Страница {page} из {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="px-4 py-2 border rounded-md disabled:opacity-50"
          >
            Вперед
          </button>
        </div>
      )}
    </div>
  );
}

