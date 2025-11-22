import {useEffect, useMemo, useState} from 'react';
import {type BudgetDetailed, budgetsApi, categoriesApi, type Category} from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Table, Thead, Tbody, Tfoot, Tr, Th, Td } from '@/components/ui/table';
import { cn } from '@/lib/utils';
import { formatCurrency } from '@/lib/utils';
import { Pencil, Trash2 } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';

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
    const [toDelete, setToDelete] = useState<BudgetDetailed | null>(null);

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

    if (budgets.length === 0) return (
        <Card>
            <CardContent className="py-10 text-center">
                <div className="mb-3 text-lg font-medium">Бюджеты не найдены</div>
                <div className="text-sm text-gray-600">Создайте первый бюджет с нужными параметрами периода.</div>
            </CardContent>
        </Card>
    );

    return (
        <div className="bg-white shadow sm:rounded-lg">
            <div className="w-full overflow-x-auto">
            <Table className="min-w-full">
                <Thead>
                <Tr>
                    <Th className="text-left">Категория</Th>
                    <Th className="text-right">Бюджет</Th>
                    <Th className="text-right">Потрачено</Th>
                    <Th className="text-right">Остаток</Th>
                    <Th />
                </Tr>
                </Thead>
                <Tbody>
                {budgets.map((b) => {
                    const catName = categoryMap.get(b.category_id) || `#${b.category_id}`;
                    const spent = +b.spent;
                    const budgeted = +b.budgeted;
                    const remaining = +b.remaining;
                    const spentPctDisplay = budgeted > 0 ? Math.round((-spent / budgeted) * 100) : 0;
                    const spentPct = Math.max(0, Math.min(100, spentPctDisplay));
                    const remainingPct = Math.max(0, 100 - spentPct);
                    const overPct = spentPctDisplay > 100 ? Math.round(spentPctDisplay - 100) : 0;
                    return (
                        <Tr key={b.id}>
                            <Td className="whitespace-nowrap text-sm text-gray-900">{catName}</Td>
                            <Td className="whitespace-nowrap text-sm text-right">{formatCurrency(budgeted)}</Td>
                            <Td className="whitespace-nowrap text-sm text-right">
                              <div className="flex flex-col items-end">
                                <div>{formatCurrency(spent)}</div>
                                <div className="text-xs text-gray-600">Потр: {Math.max(0, spentPctDisplay)}%</div>
                              </div>
                            </Td>
                            <Td className={cn('whitespace-nowrap text-sm text-right', remaining < 0 ? 'text-red-600' : '')}>
                              <div className="flex flex-col items-end">
                                <div>{formatCurrency(remaining)}</div>
                                <div className="text-xs text-gray-600">Ост: {remainingPct}%{overPct > 0 ? ` (+${overPct}% переп.)` : ''}</div>
                              </div>
                            </Td>
                            <Td className="whitespace-nowrap text-sm">
                                <div className="w-44">
                                    <div className="h-3 w-full overflow-hidden rounded bg-gray-200 flex">
                                        <div className={cn('h-3', spentPctDisplay > 100 ? 'bg-red-500' : 'bg-green-500')} style={{ width: `${spentPct}%` }} aria-label={`Потр. ${spentPct}%`} />
                                        {remainingPct > 0 && (
                                          <div
                                            className="h-3 bg-gray-300"
                                            style={{ width: `${remainingPct}%` }}
                                            aria-label={`Ост. ${remainingPct}%`}
                                          />
                                        )}
                                    </div>
                                </div>
                            </Td>
                            <Td className="whitespace-nowrap text-right text-sm font-medium space-x-1">
                                <Button variant="ghost" size="icon" aria-label="Редактировать" onClick={() => onEdit(b)}>
                                    <Pencil className="h-4 w-4" />
                                </Button>
                                <Button variant="ghost" size="icon" aria-label="Удалить" onClick={() => setToDelete(b)}>
                                    <Trash2 className="h-4 w-4 text-red-600" />
                                </Button>
                            </Td>
                        </Tr>
                    );
                })}
                </Tbody>
                <Tfoot>
                  {(() => {
                    const totals = budgets.reduce(
                      (acc, b) => {
                        acc.budgeted += Number(b.budgeted) || 0;
                        acc.spent += Number(b.spent) || 0;
                        acc.remaining += Number(b.remaining) || 0;
                        return acc;
                      },
                      { budgeted: 0, spent: 0, remaining: 0 }
                    );
                    return (
                      <Tr>
                        <Th className="text-left">Итого</Th>
                        <Th className="text-right">{formatCurrency(totals.budgeted)}</Th>
                        <Th className="text-right">{formatCurrency(totals.spent)}</Th>
                        <Th className={cn('text-right', totals.remaining < 0 ? 'text-red-600' : '')}>{formatCurrency(totals.remaining)}</Th>
                        <Th />
                      </Tr>
                    );
                  })()}
                </Tfoot>
            </Table>
            </div>

            <Dialog open={!!toDelete} onOpenChange={(o) => !o && setToDelete(null)}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Удалить бюджет?</DialogTitle>
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

