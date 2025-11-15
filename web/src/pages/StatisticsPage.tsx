import { useState, useEffect } from 'react';
import { transactionsApi, type PeriodType } from '@/lib/api';

export function StatisticsPage() {
  const [balance, setBalance] = useState({ balance: '0', total_income: '0', total_expense: '0' });
  const [period, setPeriod] = useState<PeriodType>('day');
  const [periodStats, setPeriodStats] = useState<any[]>([]);
  const [categoryStats, setCategoryStats] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');

  useEffect(() => {
    loadStatistics();
  }, [period, from, to]);

  const loadStatistics = async () => {
    try {
      setLoading(true);
      const [balanceRes, periodRes, categoryRes] = await Promise.all([
        transactionsApi.getBalanceStatistics(),
        transactionsApi.getPeriodStatistics(
          period,
          from || undefined,
          to || undefined
        ),
        transactionsApi.getCategoryStatistics(
          period,
          from || undefined,
          to || undefined
        ),
      ]);
      setBalance(balanceRes.data);
      setPeriodStats(periodRes.data.items);
      setCategoryStats(categoryRes.data.items);
    } catch (error) {
      console.error('Error loading statistics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-8">Загрузка...</div>;
  }

  return (
    <div className="px-4 py-6 sm:px-0 space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Статистика</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-sm font-medium text-gray-500">Остаток</h3>
          <p className="text-2xl font-bold text-gray-900 mt-2">
            {parseFloat(balance.balance).toFixed(2)}
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-sm font-medium text-gray-500">Всего доходов</h3>
          <p className="text-2xl font-bold text-green-600 mt-2">
            {parseFloat(balance.total_income).toFixed(2)}
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-sm font-medium text-gray-500">Всего расходов</h3>
          <p className="text-2xl font-bold text-red-600 mt-2">
            {parseFloat(balance.total_expense).toFixed(2)}
          </p>
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <div className="mb-4 flex flex-wrap gap-4 items-end">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Период</label>
            <select
              value={period}
              onChange={(e) => setPeriod(e.target.value as PeriodType)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="day">День</option>
              <option value="week">Неделя</option>
              <option value="month">Месяц</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">От</label>
            <input
              type="date"
              value={from}
              onChange={(e) => setFrom(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">До</label>
            <input
              type="date"
              value={to}
              onChange={(e) => setTo(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <h2 className="text-xl font-bold mb-4">Статистика по периодам</h2>
        {periodStats.length === 0 ? (
          <p className="text-gray-500">Нет данных</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Период
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Доходы
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Расходы
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Баланс
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {periodStats.map((item, idx) => (
                  <tr key={idx}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {item.period}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                      {parseFloat(item.income).toFixed(2)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600">
                      {parseFloat(item.expense).toFixed(2)}
                    </td>
                    <td
                      className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${
                        parseFloat(item.balance) >= 0 ? 'text-green-600' : 'text-red-600'
                      }`}
                    >
                      {parseFloat(item.balance).toFixed(2)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <h2 className="text-xl font-bold mb-4 mt-8">Статистика по категориям</h2>
        {categoryStats.length === 0 ? (
          <p className="text-gray-500">Нет данных</p>
        ) : (
          <>
            {/* Доходные категории */}
            <div className="mb-6">
              <h3 className="text-lg font-semibold text-green-700 mb-3">Доходные категории</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-green-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Период
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Категория
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Доходы
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {categoryStats
                      .filter((item) => parseFloat(item.income) > 0)
                      .map((item, idx) => (
                        <tr key={idx}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {item.period}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {item.category_name}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600 font-semibold">
                            {parseFloat(item.income).toFixed(2)}
                          </td>
                        </tr>
                      ))}
                    {categoryStats.filter((item) => parseFloat(item.income) > 0).length === 0 && (
                      <tr>
                        <td colSpan={3} className="px-6 py-4 text-center text-sm text-gray-500">
                          Нет доходных категорий
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>

            {/* Расходные категории */}
            <div>
              <h3 className="text-lg font-semibold text-red-700 mb-3">Расходные категории</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-red-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Период
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Категория
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Расходы
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {categoryStats
                      .filter((item) => parseFloat(item.expense) > 0)
                      .map((item, idx) => (
                        <tr key={idx}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {item.period}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {item.category_name}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600 font-semibold">
                            {parseFloat(item.expense).toFixed(2)}
                          </td>
                        </tr>
                      ))}
                    {categoryStats.filter((item) => parseFloat(item.expense) > 0).length === 0 && (
                      <tr>
                        <td colSpan={3} className="px-6 py-4 text-center text-sm text-gray-500">
                          Нет расходных категорий
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

