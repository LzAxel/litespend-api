import {format} from 'date-fns';
import {ru} from 'date-fns/locale/ru';
import {type Category, type Transaction} from '@/lib/api';

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
        return <div className="text-center py-8 text-gray-500">Нет транзакций</div>;
    }

    return (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
                {transactions.map((transaction) => (
                    <li key={transaction.id}>
                        <div className="px-4 py-4 sm:px-6">
                            <div className="flex items-center justify-between">
                                <div className="flex-1 min-w-0">
                                    <div className="flex items-center justify-between">
                                        <p className="text-sm font-medium text-gray-900 truncate">
                                            {transaction.description}
                                        </p>
                                    </div>
                                    <div className="mt-2 flex items-center text-sm text-gray-500">
                                        <span>{getCategoryName(transaction.category_id)}</span>
                                        <span className="mx-2">•</span>
                                        <span>
                      {format(new Date(transaction.date), 'dd.MM.yyyy HH:mm', {
                          locale: ru,
                      })}
                    </span>
                                    </div>
                                </div>
                                <div className="">
                                    <div
                                        className={`ml-4 font-semibold ${
                                            transaction.type === 0 ? 'text-green-600' : 'text-red-600'
                                        }`}
                                    >
                                        {transaction.type === 0 ? '+' : '-'}
                                        {parseFloat(transaction.amount).toFixed(2)}
                                    </div>
                                </div>
                                <div className="ml-4 flex-shrink-0 flex space-x-2">
                                    <button
                                        onClick={() => onEdit(transaction)}
                                        className="text-blue-600 hover:text-blue-900"
                                    >
                                        Редактировать
                                    </button>
                                    <button
                                        onClick={() => onDelete(transaction.id)}
                                        className="text-red-600 hover:text-red-900"
                                    >
                                        Удалить
                                    </button>
                                </div>
                            </div>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );
}

