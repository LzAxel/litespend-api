import {useState} from 'react';
import {format} from 'date-fns';
import {ru} from 'date-fns/locale/ru';
import {type Category, type FrequencyType, type PrescribedExpanseWithPaymentStatus} from '@/lib/api';

interface PrescribedExpanseListProps {
    prescribedExpanses: PrescribedExpanseWithPaymentStatus[];
    categories: Category[];
    onEdit: (prescribedExpanse: PrescribedExpanseWithPaymentStatus) => void;
    onDelete: (id: number) => void;
    onMarkAsPaid: (id: number) => void;
    onMarkAsPaidPartial: (id: number, amount: number) => void;
    onEditTransaction: (transactionId: number) => void;
}

const frequencyLabels: Record<FrequencyType, string> = {
    0: 'Ежемесячно',
    1: 'Ежедневно',
    2: 'Еженедельно',
    3: 'Ежеквартально',
};

export function PrescribedExpanseList({
                                          prescribedExpanses,
                                          categories,
                                          onEdit,
                                          onDelete,
                                          onMarkAsPaid,
                                          onMarkAsPaidPartial,
                                          onEditTransaction,
                                      }: PrescribedExpanseListProps) {
    const [showPartialPaymentModal, setShowPartialPaymentModal] = useState(false);
    const [selectedExpanse, setSelectedExpanse] = useState<PrescribedExpanseWithPaymentStatus | null>(null);
    const [partialAmount, setPartialAmount] = useState('');

    const handlePartialPaymentClick = (expanse: PrescribedExpanseWithPaymentStatus) => {
        setSelectedExpanse(expanse);
        setPartialAmount('');
        setShowPartialPaymentModal(true);
    };

    const handlePartialPaymentSubmit = () => {
        if (selectedExpanse && partialAmount) {
            const amount = parseFloat(partialAmount);
            if (amount > 0 && amount <= parseFloat(selectedExpanse.amount)) {
                onMarkAsPaidPartial(selectedExpanse.id, amount);
                setShowPartialPaymentModal(false);
                setSelectedExpanse(null);
                setPartialAmount('');
            }
        }
    };
    const getCategoryName = (categoryId: number) => {
        return categories.find((c) => c.id === categoryId)?.name || 'Неизвестно';
    };

    if (prescribedExpanses.length === 0) {
        return <div className="text-center py-8 text-gray-500">Нет обязательных трат</div>;
    }

    const totalAmount = prescribedExpanses.reduce((sum, pe) => sum + parseFloat(pe.amount), 0);
    const unpaidAmount = prescribedExpanses.reduce((sum, pe) => {
        const paid = parseFloat(pe.paid_amount || '0');
        const total = parseFloat(pe.amount);
        return sum + Math.max(0, total - paid);
    }, 0);

    return (
        <div>
            <div className="mb-4 p-4 bg-blue-50 rounded-lg">
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-gray-600">Всего обязательных трат</p>
                        <p className="text-2xl font-bold text-gray-900">{totalAmount.toFixed(2)} ₽</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-600">Неоплаченные траты</p>
                        <p className="text-2xl font-bold text-red-600">{unpaidAmount.toFixed(2)} ₽</p>
                    </div>
                </div>
            </div>

            <div className="bg-white shadow overflow-hidden sm:rounded-md">
                <ul className="divide-y divide-gray-200">
                    {prescribedExpanses.map((prescribedExpanse) => (
                        <li key={prescribedExpanse.id}>
                            <div className="px-4 py-4 sm:px-6">
                                <div className="flex items-center justify-between">
                                    <div className="w-[500px] min-w-0">
                                        <div className="flex items-center">
                                            <p className="text-sm font-medium text-gray-900 truncate">
                                                {prescribedExpanse.description}
                                            </p>
                                            <div
                                                className={`ml-4 px-2 py-1 rounded text-xs font-semibold ${
                                                    prescribedExpanse.is_paid
                                                        ? 'bg-green-100 text-green-800'
                                                        : 'bg-red-100 text-red-800'
                                                }`}
                                            >
                                                {prescribedExpanse.is_paid ? 'Оплачено' : 'Не оплачено'}
                                            </div>
                                        </div>
                                        <div className="mt-2 flex items-center text-sm text-gray-500">
                                            <span>{getCategoryName(prescribedExpanse.category_id)}</span>
                                            <span className="mx-2">•</span>
                                            <span>{frequencyLabels[prescribedExpanse.frequency]}</span>
                                            <span className="mx-2">•</span>
                                            <span>
                        {format(new Date(prescribedExpanse.date_time), 'dd.MM.yyyy', {
                            locale: ru,
                        })}
                      </span>
                                        </div>
                                        {/* Прогресс оплаты */}
                                        {parseFloat(prescribedExpanse.paid_amount || '0') > 0 && (
                                            <div className="mt-3">
                                                <div
                                                    className="flex items-center justify-between text-xs text-gray-600 mb-1">
                          <span>
                            Оплачено: {parseFloat(prescribedExpanse.paid_amount || '0').toFixed(2)} / {parseFloat(prescribedExpanse.amount).toFixed(2)} ₽
                          </span>
                                                    <span>
                            {Math.round((parseFloat(prescribedExpanse.paid_amount || '0') / parseFloat(prescribedExpanse.amount)) * 100)}%
                          </span>
                                                </div>
                                                <div className="w-full bg-gray-200 rounded-full h-2">
                                                    <div
                                                        className={`h-2 rounded-full ${
                                                            prescribedExpanse.is_paid ? 'bg-green-600' : 'bg-yellow-500'
                                                        }`}
                                                        style={{
                                                            width: `${Math.min(100, (parseFloat(prescribedExpanse.paid_amount || '0') / parseFloat(prescribedExpanse.amount)) * 100)}%`,
                                                        }}
                                                    />
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                    <div className="ml-auto">
                                        <div
                                            className={`font-semibold text-red-600`}
                                        >
                                            -{parseFloat(prescribedExpanse.amount).toFixed(2)} ₽
                                        </div>
                                    </div>
                                    <div className="ml-4 flex-shrink-0 flex flex-col space-x-2 flex-wrap">
                                        <div className="flex gap-2">
                                            {!prescribedExpanse.is_paid && !prescribedExpanse.transaction_id && (
                                                <>
                                                    <button
                                                        onClick={() => onMarkAsPaid(prescribedExpanse.id)}
                                                        className="px-3 py-1 bg-green-600 text-white text-sm rounded-md hover:bg-green-700"
                                                    >
                                                        Полностью
                                                    </button>
                                                    <button
                                                        onClick={() => handlePartialPaymentClick(prescribedExpanse)}
                                                        className="px-3 py-1 bg-yellow-600 text-white text-sm rounded-md hover:bg-yellow-700"
                                                    >
                                                        Частично
                                                    </button>
                                                </>
                                            )}
                                            {prescribedExpanse.transaction_id && parseFloat(prescribedExpanse.paid_amount || '0') > 0 && (
                                                <button
                                                    onClick={() => onEditTransaction(prescribedExpanse.transaction_id!)}
                                                    className="px-3 w-full py-1 bg-blue-600 text-white text-sm rounded-md hover:bg-blue-700"
                                                    title="Редактировать транзакцию оплаты"
                                                >
                                                    Изменить оплату
                                                </button>
                                            )}</div>
                                        <div className="flex gap-2">
                                            <button
                                                onClick={() => onEdit(prescribedExpanse)}
                                                className="text-blue-600 hover:text-blue-900"
                                            >
                                                Редактировать
                                            </button>
                                            <button
                                                onClick={() => onDelete(prescribedExpanse.id)}
                                                className="text-red-600 hover:text-red-900"
                                            >
                                                Удалить
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </li>
                    ))}
                </ul>
            </div>

            {/* Модальное окно для частичной оплаты */}
            {showPartialPaymentModal && selectedExpanse && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
                        <div className="flex justify-between items-center mb-4">
                            <h2 className="text-xl font-bold">Частичная оплата</h2>
                            <button
                                onClick={() => {
                                    setShowPartialPaymentModal(false);
                                    setSelectedExpanse(null);
                                    setPartialAmount('');
                                }}
                                className="text-gray-500 hover:text-gray-700"
                            >
                                ✕
                            </button>
                        </div>
                        <div className="mb-4">
                            <p className="text-sm text-gray-600 mb-2">
                                Обязательная трата: <span className="font-semibold">{selectedExpanse.description}</span>
                            </p>
                            <p className="text-sm text-gray-600 mb-4">
                                Сумма: <span
                                className="font-semibold">{parseFloat(selectedExpanse.amount).toFixed(2)} ₽</span>
                            </p>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Сумма оплаты
                            </label>
                            <input
                                type="number"
                                step="0.01"
                                min="0.01"
                                max={selectedExpanse.amount}
                                value={partialAmount}
                                onChange={(e) => setPartialAmount(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                                placeholder="Введите сумму"
                            />
                            {partialAmount && parseFloat(partialAmount) > parseFloat(selectedExpanse.amount) && (
                                <p className="mt-1 text-sm text-red-600">
                                    Сумма не может быть больше {parseFloat(selectedExpanse.amount).toFixed(2)} ₽
                                </p>
                            )}
                        </div>
                        <div className="flex justify-end space-x-3">
                            <button
                                onClick={() => {
                                    setShowPartialPaymentModal(false);
                                    setSelectedExpanse(null);
                                    setPartialAmount('');
                                }}
                                className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
                            >
                                Отмена
                            </button>
                            <button
                                onClick={handlePartialPaymentSubmit}
                                disabled={
                                    !partialAmount ||
                                    parseFloat(partialAmount) <= 0 ||
                                    parseFloat(partialAmount) > parseFloat(selectedExpanse.amount)
                                }
                                className="px-4 py-2 bg-yellow-600 text-white rounded-md hover:bg-yellow-700 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                Оплатить
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

