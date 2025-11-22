import {useEffect, useMemo, useState} from 'react';
import {type BudgetDetailed, budgetsApi, categoriesApi, type Category} from '@/lib/api';

interface BudgetListProps {
    year: number;
    month: number;
    onEdit: (budget: BudgetDetailed) => void;
    onDeleted: () => void;
}

export function BudgetList({year, month, onEdit, onDeleted}: BudgetListProps) {
    const [budgets, setBudgets] = useState<BudgetDetailed[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [year, month]);

    async function loadData() {
        try {
            setLoading(true);
            const [budgetsRes, categoriesRes] = await Promise.all([
                budgetsApi.getByPeriod(year, month),
                categoriesApi.getAll(),
            ]);
            setBudgets(budgetsRes.data);
            setCategories(categoriesRes.data);
        } catch (e) {
            console.error('Failed to load budgets', e);
        } finally {
            setLoading(false);
        }
    }

    const categoryMap = useMemo(() => {
        const map = new Map<number, string>();
        categories.forEach((c) => map.set(c.id, c.name));
        return map;
    }, [categories]);

    const handleDelete = async (id: number) => {
        if (!confirm('Удалить бюджет?')) return;
        try {
            await budgetsApi.delete(id);
            await loadData();
            onDeleted();
        } catch (e) {
            console.error('Delete budget error', e);
            alert('Не удалось удалить бюджет');
        }
    };

    if (loading) return <div className="text-center py-8">Загрузка...</div>;

    if (budgets.length === 0) return <div className="text-center py-8">Бюджеты не найдены</div>;

    return (
        <div className="bg-white shadow sm:rounded-lg overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Категория</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Бюджет</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Потрачено</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Остаток</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">%</th>
                    <th className="px-6 py-3"/>
                </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                {budgets.map((b) => {
                    const catName = categoryMap.get(b.category_id) || `#${b.category_id}`;
                    const spent = +b.spent;
                    const budgeted = +b.budgeted;
                    const remaining = +b.remaining;
                    const percent = budgeted > 0 ? -Math.min(100, Math.round((spent / budgeted) * 100)) : 0;
                    const barColor = percent < 80 ? 'bg-green-500' : percent < 100 ? 'bg-yellow-500' : 'bg-red-500';
                    return (
                        <tr key={b.id}>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{catName}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right">{budgeted.toFixed(2)}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right">{spent.toFixed(2)}</td>
                            <td className={`px-6 py-4 whitespace-nowrap text-sm text-right ${remaining < 0 ? 'text-red-600' : ''}`}>{remaining.toFixed(2)}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm">
                                <div className="w-40 bg-gray-200 rounded h-3">
                                    <div className={`${barColor} h-3 rounded`} style={{width: `${percent}%`}}/>
                                </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                                <button onClick={() => onEdit(b)}
                                        className="text-blue-600 hover:text-blue-900">Редактировать
                                </button>
                                <button onClick={() => handleDelete(b.id)}
                                        className="text-red-600 hover:text-red-900">Удалить
                                </button>
                            </td>
                        </tr>
                    );
                })}
                </tbody>
            </table>
        </div>
    );

}
