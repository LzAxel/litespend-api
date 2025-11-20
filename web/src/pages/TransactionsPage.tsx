import { useState, useEffect, useCallback, useRef } from 'react';
import { transactionsApi, categoriesApi, type Transaction, type Category, type SortField, type SortOrder } from '@/lib/api';
import { TransactionForm } from '@/components/transactions/TransactionForm';
import { TransactionList } from '@/components/transactions/TransactionList';

export function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null);
  const [sortBy, setSortBy] = useState<SortField | undefined>(undefined);
  const [sortOrder, setSortOrder] = useState<SortOrder>('desc');
  const [search, setSearch] = useState('');
  const [searchInput, setSearchInput] = useState('');
  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadCategories();
  }, []);

  useEffect(() => {
    // Сброс при изменении сортировки или поиска
    setPage(1);
    setTransactions([]);
    setHasMore(true);
    loadData(1, true);
  }, [sortBy, sortOrder, search]);

  const loadCategories = async () => {
    try {
      const response = await categoriesApi.getAll();
      setCategories(response.data);
    } catch (error) {
      console.error('Error loading categories:', error);
    }
  };

  const loadData = async (pageNum: number, reset: boolean = false) => {
    try {
      if (reset) {
        setLoading(true);
      } else {
        setLoadingMore(true);
      }
      const response = await transactionsApi.getAll(pageNum, 20, sortBy, sortOrder, search || undefined);
      if (reset) {
        setTransactions(response.data.data);
      } else {
        setTransactions((prev) => [...prev, ...response.data.data]);
      }
      setTotalPages(response.data.meta.total_pages);
      setHasMore(pageNum < response.data.meta.total_pages);
      setPage(pageNum);
    } catch (error) {
      setHasMore(false);
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  };

  const loadMore = useCallback(() => {
    if (!loadingMore && hasMore && !loading) {
      const nextPage = page + 1;
      loadData(nextPage, false);
    }
  }, [page, loadingMore, hasMore, loading]);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loadingMore && !loading) {
          loadMore();
        }
      },
      { threshold: 0.1 }
    );

    const currentTarget = observerTarget.current;
    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
    };
  }, [loadMore, hasMore, loadingMore, loading]);

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
      // Перезагружаем с первой страницы
      setPage(1);
      setTransactions([]);
      setHasMore(true);
      loadData(1, true);
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
    // Перезагружаем с первой страницы
    setPage(1);
    setTransactions([]);
    setHasMore(true);
    loadData(1, true);
  };

  const handleSort = (field: SortField) => {
    if (sortBy === field) {
      // Переключаем направление сортировки
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('desc');
    }
  };

  const handleSearch = () => {
    setSearch(searchInput);
  };

  const handleSearchKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
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

      {/* Поиск и сортировка */}
      <div className="mb-6 space-y-4">
        <div className="flex gap-4">
          <div className="flex-1">
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyPress={handleSearchKeyPress}
              placeholder="Поиск по названию..."
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <button
            onClick={handleSearch}
            className="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
          >
            Найти
          </button>
          {search && (
            <button
              onClick={() => {
                setSearchInput('');
                setSearch('');
              }}
              className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
            >
              Сбросить
            </button>
          )}
        </div>

        <div className="flex gap-2 flex-wrap">
          <span className="text-sm text-gray-600 self-center">Сортировка:</span>
          <button
            onClick={() => handleSort('date')}
            className={`px-3 py-1 rounded-md text-sm ${
              sortBy === 'date'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            По дате {sortBy === 'date' && (sortOrder === 'asc' ? '↑' : '↓')}
          </button>
          <button
            onClick={() => handleSort('description')}
            className={`px-3 py-1 rounded-md text-sm ${
              sortBy === 'description'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            По названию {sortBy === 'description' && (sortOrder === 'asc' ? '↑' : '↓')}
          </button>
          <button
            onClick={() => handleSort('category')}
            className={`px-3 py-1 rounded-md text-sm ${
              sortBy === 'category'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            По категории {sortBy === 'category' && (sortOrder === 'asc' ? '↑' : '↓')}
          </button>
          {sortBy && (
            <button
              onClick={() => {
                setSortBy(undefined);
                setSortOrder('desc');
              }}
              className="px-3 py-1 rounded-md text-sm bg-gray-200 text-gray-700 hover:bg-gray-300"
            >
              Сбросить сортировку
            </button>
          )}
        </div>
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

      {/* Индикатор загрузки и триггер для бесконечной подгрузки */}
      <div ref={observerTarget} className="h-10 flex items-center justify-center">
        {loadingMore && <div className="text-gray-500">Загрузка...</div>}
        {!hasMore && transactions.length > 0 && (
          <div className="text-gray-500">Все транзакции загружены</div>
        )}
      </div>
    </div>
  );
}
