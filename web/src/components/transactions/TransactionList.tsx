import {format} from 'date-fns';
import {ru} from 'date-fns/locale/ru';
import {type Category, type Transaction} from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Pencil, Trash2 } from 'lucide-react';

interface TransactionListProps {
    transactions: Transaction[];
    categories: Category[];
    onEdit: (transaction: Transaction) => void;
    onDelete: (id: number) => void;
}

export function TransactionList({
                                    transactions,
                                    categories,
                                    onEdit,
                                    onDelete,
                                }: TransactionListProps) {
    const getCategoryName = (categoryId: number) => {
        return categories.find((c) => c.id === categoryId)?.name || 'Неизвестно';
    };

    if (transactions.length === 0) {
        return <div className="text-center py-8 text-[rgb(var(--muted-foreground))]">Нет транзакций</div>;
    }

    return (
        <div className="bg-[rgb(var(--card))] shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-[rgb(var(--border))]">
                {transactions.map((transaction) => (
                    <li key={transaction.id}>
                        <div className="px-4 py-3 sm:px-6">
                            <div className="flex items-center justify-between gap-3">
                                <div className="min-w-0 flex-1">
                                    <p className="text-sm font-medium text-[rgb(var(--app-fg))] truncate">
                                        {transaction.description}
                                    </p>
                                    <div className="mt-1 flex flex-wrap items-center gap-x-2 text-xs text-[rgb(var(--muted-foreground))]">
                                        <span className="truncate max-w-[50vw] sm:max-w-none">{getCategoryName(transaction.category_id)}</span>
                                        <span>•</span>
                                        <span>
                                            {format(new Date(transaction.date), 'dd.MM.yyyy HH:mm', { locale: ru })}
                                        </span>
                                    </div>
                                </div>
                                <div className="text-right whitespace-nowrap">
                                    <div
                                        className={`font-semibold ${+transaction.amount > 0 ? 'text-[rgb(var(--success))]' : 'text-[rgb(var(--destructive))]'}`}
                                    >
                                        {parseFloat(transaction.amount).toFixed(2)}
                                    </div>
                                </div>
                                <div className="flex-shrink-0 flex items-center gap-1 sm:gap-2">
                                    <Button variant="ghost" size="icon" aria-label="Редактировать" onClick={() => onEdit(transaction)}>
                                        <Pencil className="h-4 w-4" />
                                    </Button>
                                    <Button variant="ghost" size="icon" aria-label="Удалить" onClick={() => onDelete(transaction.id)}>
                                        <Trash2 className="h-4 w-4 text-[rgb(var(--destructive))]" />
                                    </Button>
                                </div>
                            </div>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );
}

