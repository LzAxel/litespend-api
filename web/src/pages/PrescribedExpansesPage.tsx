import { useState, useEffect } from 'react';
import {
  prescribedExpansesApi,
  categoriesApi,
  transactionsApi,
  type PrescribedExpanseWithPaymentStatus,
  type Category,
  type Transaction,
} from '@/lib/api';
import { PrescribedExpanseForm } from '@/components/prescribedExpanses/PrescribedExpanseForm';
import { PrescribedExpanseList } from '@/components/prescribedExpanses/PrescribedExpanseList';
import { TransactionForm } from '@/components/transactions/TransactionForm';

export function PrescribedExpansesPage() {
  const [prescribedExpanses, setPrescribedExpanses] = useState<PrescribedExpanseWithPaymentStatus[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [showTransactionForm, setShowTransactionForm] = useState(false);
  const [editingPrescribedExpanse, setEditingPrescribedExpanse] =
    useState<PrescribedExpanseWithPaymentStatus | null>(null);
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [prescribedExpansesResponse, categoriesResponse] = await Promise.all([
        prescribedExpansesApi.getAllWithPaymentStatus(),
        categoriesApi.getAll(),
      ]);
      setPrescribedExpanses(prescribedExpansesResponse.data);
      setCategories(categoriesResponse.data);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingPrescribedExpanse(null);
    setShowForm(true);
  };

  const handleEdit = (prescribedExpanse: PrescribedExpanseWithPaymentStatus) => {
    setEditingPrescribedExpanse(prescribedExpanse);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить обязательную трату?')) return;
    try {
      await prescribedExpansesApi.delete(id);
      loadData();
    } catch (error) {
      console.error('Error deleting prescribed expanse:', error);
    }
  };

  const handleMarkAsPaid = async (id: number) => {
    try {
      await prescribedExpansesApi.markAsPaid(id);
      loadData();
    } catch (error) {
      console.error('Error marking as paid:', error);
      alert('Ошибка при отметке оплаты');
    }
  };

  const handleMarkAsPaidPartial = async (id: number, amount: number) => {
    try {
      await prescribedExpansesApi.markAsPaidPartial(id, amount.toFixed(2));
      loadData();
    } catch (error: any) {
      console.error('Error marking as paid partial:', error);
      alert(error.response?.data?.error || 'Ошибка при частичной оплате');
    }
  };

  const handleFormClose = () => {
    setShowForm(false);
    setEditingPrescribedExpanse(null);
  };

  const handleFormSuccess = () => {
    handleFormClose();
    loadData();
  };

  const handleEditTransaction = async (transactionId: number) => {
    try {
      const response = await transactionsApi.getById(transactionId);
      setEditingTransaction(response.data);
      setShowTransactionForm(true);
    } catch (error) {
      console.error('Error loading transaction:', error);
      alert('Ошибка при загрузке транзакции');
    }
  };

  const handleTransactionFormClose = () => {
    setShowTransactionForm(false);
    setEditingTransaction(null);
  };

  const handleTransactionFormSuccess = () => {
    handleTransactionFormClose();
    loadData(); // Перезагружаем список обязательных трат для обновления статуса оплаты
  };

  if (loading) {
    return <div className="text-center py-8">Загрузка...</div>;
  }

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Обязательные траты</h1>
        <button
          onClick={handleCreate}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          Добавить обязательную трату
        </button>
      </div>

      {showForm && (
        <PrescribedExpanseForm
          prescribedExpanse={editingPrescribedExpanse}
          onClose={handleFormClose}
          onSuccess={handleFormSuccess}
        />
      )}

      {showTransactionForm && editingTransaction && (
        <TransactionForm
          transaction={editingTransaction}
          categories={categories}
          onClose={handleTransactionFormClose}
          onSuccess={handleTransactionFormSuccess}
        />
      )}

      <PrescribedExpanseList
        prescribedExpanses={prescribedExpanses}
        categories={categories}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onMarkAsPaid={handleMarkAsPaid}
        onMarkAsPaidPartial={handleMarkAsPaidPartial}
        onEditTransaction={handleEditTransaction}
      />
    </div>
  );
}

