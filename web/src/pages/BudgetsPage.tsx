import { useState } from 'react';
import { BudgetForm } from '@/components/budgets/BudgetForm';
import { BudgetList } from '@/components/budgets/BudgetList';
import type { BudgetDetailed } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';

export function BudgetsPage() {
  const now = new Date();
  const [year, setYear] = useState<number>(now.getFullYear());
  const [month, setMonth] = useState<number>(now.getMonth() + 1);
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState<BudgetDetailed | null>(null);

  const monthNames = ['Январь','Февраль','Март','Апрель','Май','Июнь','Июль','Август','Сентябрь','Октябрь','Ноябрь','Декабрь'];

  const handleCreate = () => {
    setEditing(null);
    setShowForm(true);
  };

  const handleEdit = (b: BudgetDetailed) => {
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

  const prevMonth = () => {
    const m = month - 1;
    if (m < 1) {
      setMonth(12);
      setYear((y) => y - 1);
    } else {
      setMonth(m);
    }
  };

  const nextMonth = () => {
    const m = month + 1;
    if (m > 12) {
      setMonth(1);
      setYear((y) => y + 1);
    } else {
      setMonth(m);
    }
  };

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Бюджеты</h1>
        <div className="flex flex-wrap gap-2 items-end">
          <div className="flex items-end gap-2">
            <Button variant="outline" onClick={prevMonth} aria-label="Предыдущий месяц">←</Button>
            <div>
              <Label htmlFor="month">Месяц</Label>
              <select
                id="month"
                className="mt-1 block w-40 border rounded p-2"
                value={month}
                onChange={(e) => setMonth(Number(e.target.value))}
              >
                {monthNames.map((m, i) => (
                  <option key={m} value={i + 1}>{m}</option>
                ))}
              </select>
            </div>
            <div>
              <Label htmlFor="year">Год</Label>
              <Input
                id="year"
                type="number"
                className="mt-1 w-28"
                value={year}
                onChange={(e) => setYear(Number(e.target.value))}
              />
            </div>
            <Button variant="outline" onClick={nextMonth} aria-label="Следующий месяц">→</Button>
          </div>
          <Button onClick={handleCreate}>Добавить</Button>
        </div>
      </div>

      <Dialog open={showForm} onOpenChange={(o) => !o && closeForm()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? 'Редактировать бюджет' : 'Создать бюджет'}</DialogTitle>
          </DialogHeader>
          <BudgetForm
            budget={editing}
            defaultYear={year}
            defaultMonth={month}
            onClose={closeForm}
            onSuccess={onSaved}
          />
        </DialogContent>
      </Dialog>

      <BudgetList year={year} month={month} onEdit={handleEdit} onDeleted={() => {}} />
    </div>
  );
}

