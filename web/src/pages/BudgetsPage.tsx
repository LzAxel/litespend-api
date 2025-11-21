import { useState } from 'react';
import { BudgetForm } from '@/components/budgets/BudgetForm';
import { BudgetList } from '@/components/budgets/BudgetList';
import type { Budget } from '@/lib/api';

export function BudgetsPage() {
  const now = new Date();
  const [year, setYear] = useState<number>(now.getFullYear());
  const [month, setMonth] = useState<number>(now.getMonth() + 1);
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState<Budget | null>(null);

  const handleCreate = () => {
    setEditing(null);
    setShowForm(true);
  };

  const handleEdit = (b: Budget) => {
    setEditing(b);
    setShowForm(true);
  };

  const closeForm = () => {
    setShowForm(false);
    setEditing(null);
  };

  const onSaved = () => {
    setShowForm(false);
    setEditing(null);
  };

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Бюджеты</h1>
        <div className="flex gap-2 items-end">
          <div>
            <label className="block text-sm font-medium text-gray-700">Год</label>
            <input
              type="number"
              className="mt-1 block w-28 border rounded p-2"
              value={year}
              onChange={(e) => setYear(Number(e.target.value))}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Месяц</label>
            <input
              type="number"
              className="mt-1 block w-24 border rounded p-2"
              min={1}
              max={12}
              value={month}
              onChange={(e) => setMonth(Number(e.target.value))}
            />
          </div>
          <button onClick={handleCreate} className="h-10 px-4 bg-blue-600 text-white rounded hover:bg-blue-700">Добавить</button>
        </div>
      </div>

      {showForm && (
        <BudgetForm
          budget={editing}
          defaultYear={year}
          defaultMonth={month}
          onClose={closeForm}
          onSuccess={onSaved}
        />
      )}

      <BudgetList year={year} month={month} onEdit={handleEdit} onDeleted={() => {}} />
    </div>
  );
}
