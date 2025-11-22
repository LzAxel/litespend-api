import { useState, useEffect } from 'react';
import {transactionsApi, type PeriodType, type CurrentBalanceStatistics} from '@/lib/api';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { formatCurrency } from '@/lib/utils';

export function StatisticsPage() {
  const [balance, setBalance] = useState<CurrentBalanceStatistics>();
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
    const currentDate = new Date();
    try {
      setLoading(true);
      const [balanceRes, periodRes, categoryRes] = await Promise.all([
        transactionsApi.getBalanceStatistics(currentDate.getFullYear(), currentDate.getMonth() + 1),
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

  if (loading || !balance) {
    return <div className="text-center py-8">Загрузка...</div>;
  }

  const freeBalance = parseFloat(balance.total_income) - parseFloat(balance.total_expense);

  return (
    <div className="px-4 py-6 sm:px-0 space-y-6">
      <h1 className="text-2xl sm:text-3xl font-bold">Статистика</h1>

      <div className="flex flex-row gap-3 flex-wrap">
        <div className="flex flex-row gap-3 flex-grow">
          <Card className="w-auto flex-grow">
            <CardContent className="p-4 sm:p-6">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h3 className="text-sm font-medium">Свободные</h3>
                  <p className="text-2xl font-bold mt-2">{formatCurrency(parseFloat(balance.free_to_distribute))}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card className="w-auto flex-grow">
            <CardContent className="p-4 sm:p-6">
              <div>
                <h3 className="text-sm font-medium opacity-30">Остаток</h3>
                <p className="text-2xl font-bold mt-2 opacity-30">{formatCurrency(freeBalance)}</p>
              </div>
            </CardContent>
          </Card>
        </div>
        <Card className="flex-grow">
          <CardContent className="p-4 sm:p-6">
            <h3 className="text-sm font-medium text-[rgb(var(--muted-foreground))]">Всего доходов</h3>
            <p className="text-2xl font-bold text-[rgb(var(--success))] mt-2">{formatCurrency(parseFloat(balance.total_income))}</p>
          </CardContent>
        </Card>
        <Card className="flex-grow">
          <CardContent className="p-4 sm:p-6">
            <h3 className="text-sm font-medium text-[rgb(var(--muted-foreground))]">Всего расходов</h3>
            <p className="text-2xl font-bold text-[rgb(var(--destructive))] mt-2">{formatCurrency(parseFloat(balance.total_expense))}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardContent className="p-4 sm:p-6">
          <div className="mb-4 grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
            <div>
              <Label htmlFor="period">Период</Label>
              <select
                id="period"
                value={period}
                onChange={(e) => setPeriod(e.target.value as PeriodType)}
                className="mt-1 block w-full border rounded p-2 text-sm"
              >
                <option value="day">День</option>
                <option value="week">Неделя</option>
                <option value="month">Месяц</option>
              </select>
            </div>
            <div>
              <Label htmlFor="from">От</Label>
              <Input id="from" type="date" className="mt-1 w-full" value={from} onChange={(e) => setFrom(e.target.value)} />
            </div>
            <div>
              <Label htmlFor="to">До</Label>
              <Input id="to" type="date" className="mt-1 w-full" value={to} onChange={(e) => setTo(e.target.value)} />
            </div>
          </div>

        <h2 className="text-lg sm:text-xl font-bold mb-4">Статистика по периодам</h2>
        {periodStats.length === 0 ? (
          <p className="text-[rgb(var(--muted-foreground))]">Нет данных</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-[rgb(var(--border))]">
              <thead className="bg-[rgb(var(--muted))]">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                    Период
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                    Доходы
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                    Расходы
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                    Баланс
                  </th>
                </tr>
              </thead>
              <tbody className="bg-[rgb(var(--card))] divide-y divide-[rgb(var(--border))]">
                {periodStats.map((item, idx) => (
                  <tr key={idx}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--app-fg))]">
                      {item.period}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--success))]">
                      {parseFloat(item.income).toFixed(2)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--destructive))]">
                      {parseFloat(item.expense).toFixed(2)}
                    </td>
                    <td
                      className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${
                        parseFloat(item.balance) >= 0 ? 'text-[rgb(var(--success))]' : 'text-[rgb(var(--destructive))]'
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

        <h2 className="text-lg sm:text-xl font-bold mb-4 mt-8">Статистика по категориям</h2>
        {categoryStats.length === 0 ? (
          <p className="text-[rgb(var(--muted-foreground))]">Нет данных</p>
        ) : (
          <>
            {/* Доходные категории */}
            <div className="mb-6">
              <h3 className="text-base sm:text-lg font-semibold text-[rgb(var(--success))] mb-3">Доходные категории</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-[rgb(var(--border))]">
                  <thead className="bg-[rgb(var(--muted))]">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                        Период
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                        Категория
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                        Доходы
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-[rgb(var(--card))] divide-y divide-[rgb(var(--border))]">
                    {categoryStats
                      .filter((item) => parseFloat(item.income) > 0)
                      .map((item, idx) => (
                        <tr key={idx}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--app-fg))]">
                            {item.period}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--app-fg))]">
                            {item.category_name}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--success))] font-semibold">
                            {parseFloat(item.income).toFixed(2)}
                          </td>
                        </tr>
                      ))}
                    {categoryStats.filter((item) => parseFloat(item.income) > 0).length === 0 && (
                      <tr>
                        <td colSpan={3} className="px-6 py-4 text-center text-sm text-[rgb(var(--muted-foreground))]">
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
              <h3 className="text-base sm:text-lg font-semibold text-[rgb(var(--destructive))] mb-3">Расходные категории</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-[rgb(var(--border))]">
                  <thead className="bg-[rgb(var(--muted))]">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                        Период
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                        Категория
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-[rgb(var(--muted-foreground))] uppercase tracking-wider">
                        Расходы
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-[rgb(var(--card))] divide-y divide-[rgb(var(--border))]">
                    {categoryStats
                      .filter((item) => parseFloat(item.expense) > 0)
                      .map((item, idx) => (
                        <tr key={idx}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--app-fg))]">
                            {item.period}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--app-fg))]">
                            {item.category_name}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-[rgb(var(--destructive))] font-semibold">
                            {parseFloat(item.expense).toFixed(2)}
                          </td>
                        </tr>
                      ))}
                    {categoryStats.filter((item) => parseFloat(item.expense) > 0).length === 0 && (
                      <tr>
                        <td colSpan={3} className="px-6 py-4 text-center text-sm text-[rgb(var(--muted-foreground))]">
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
        </CardContent>
      </Card>
    </div>
  );
}

