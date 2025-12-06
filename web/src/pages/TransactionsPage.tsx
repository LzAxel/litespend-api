import { useState, useEffect, useCallback, useRef } from 'react';
import { transactionsApi, categoriesApi, accountsApi, type Transaction, type Category, type Account, type SortField, type SortOrder } from '@/lib/api';
import { TransactionForm } from '@/components/transactions/TransactionForm';
import { TransactionList } from '@/components/transactions/TransactionList';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Plus } from 'lucide-react';

export function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null);
  const [toDelete, setToDelete] = useState<Transaction | null>(null);
  const [sortBy, setSortBy] = useState<SortField | undefined>(undefined);
  const [sortOrder, setSortOrder] = useState<SortOrder>('desc');
  const [search, setSearch] = useState('');
  const [searchInput, setSearchInput] = useState('');
  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadDictionaries();
  }, []);

  useEffect(() => {
    // Сброс при изменении сортировки или поиска
    setPage(1);
    setTransactions([]);
    setHasMore(true);
    loadData(1, true);
  }, [sortBy, sortOrder, search]);

  const loadDictionaries = async () => {
    try {
      const [catsRes, accsRes] = await Promise.all([
        categoriesApi.getAll(),
        accountsApi.getAll(),
      ]);
      setCategories(catsRes.data);
      setAccounts(accsRes.data.accounts);
    } catch (error) {
      console.error('Error loading dictionaries:', error);
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
      <div className="flex justify-between items-center mb-6 gap-3">
        <h1 className="text-2xl sm:text-3xl font-bold">Транзакции</h1>
        <Button onClick={handleCreate} className="flex items-center gap-2" aria-label="Добавить транзакцию">
          <Plus className="h-5 w-5 md:hidden" />
          <span className="hidden md:inline">Добавить транзакцию</span>
        </Button>
      </div>

      {/* Поиск и сортировка */}
      <div className="mb-6 space-y-4">
        <div className="flex gap-4">
          <div className="flex-1">
            <Input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={handleSearchKeyPress}
              placeholder="Поиск по названию..."
            />
          </div>
          <Button onClick={handleSearch} variant="secondary">Найти</Button>
          {search && (
            <Button
              variant="outline"
              onClick={() => {
                setSearchInput('');
                setSearch('');
              }}
            >
              Сбросить
            </Button>
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

      <Dialog open={showForm} onOpenChange={(o) => !o && handleFormClose()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editingTransaction ? 'Редактировать транзакцию' : 'Добавить транзакцию'}</DialogTitle>
          </DialogHeader>
          <TransactionForm
            transaction={editingTransaction}
            categories={categories}
            accounts={accounts}
            onClose={handleFormClose}
            onSuccess={handleFormSuccess}
          />
        </DialogContent>
      </Dialog>

      <TransactionList
        transactions={transactions}
        categories={categories}
        accounts={accounts}
        onEdit={handleEdit}
        onDelete={(id) => setToDelete(transactions.find((t) => t.id === id) || null)}
      />

      {/* Индикатор загрузки и триггер для бесконечной подгрузки */}
      <div ref={observerTarget} className="h-10 flex items-center justify-center">
        {loadingMore && <div className="text-gray-500">Загрузка...</div>}
        {!hasMore && transactions.length > 0 && (
          <div className="text-gray-500">Все транзакции загружены</div>
        )}
      </div>

      <Dialog open={!!toDelete} onOpenChange={(o) => !o && setToDelete(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Удалить транзакцию?</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-gray-600">Действие нельзя отменить.</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setToDelete(null)}>Отмена</Button>
            <Button
              variant="destructive"
              onClick={async () => {
                if (toDelete) {
                  await handleDelete(toDelete.id);
                  setToDelete(null);
                }
              }}
            >
              Удалить
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
