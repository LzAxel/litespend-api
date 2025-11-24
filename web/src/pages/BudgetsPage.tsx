import { useState } from 'react';
import { BudgetForm } from '@/components/budgets/BudgetForm';
import { BudgetList } from '@/components/budgets/BudgetList';
import type { BudgetDetailed } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Plus } from 'lucide-react';

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
      <div className="mb-6">
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <h1 className="text-2xl sm:text-3xl font-bold">Бюджеты</h1>
          <div className="w-full md:w-auto grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div className="flex items-end justify-between gap-2 sm:col-span-2 md:col-span-1 items-center">
              <Button variant="outline" onClick={prevMonth} aria-label="Предыдущий месяц" className="shrink-0">←</Button>
              <div className="flex-1 min-w-0 px-1">
                <Label htmlFor="month" className="text-xs sm:text-sm">Месяц</Label>
                <select
                  id="month"
                  className="mt-1 block w-full border rounded p-2 text-sm"
                  value={month}
                  onChange={(e) => setMonth(Number(e.target.value))}
                >
                  {monthNames.map((m, i) => (
                    <option key={m} value={i + 1}>{m}</option>
                  ))}
                </select>
              </div>
              <div className="w-28 md:w-28 px-1">
                <Label htmlFor="year" className="text-xs sm:text-sm">Год</Label>
                <Input
                  id="year"
                  type="number"
                  className="mt-1 w-full"
                  value={year}
                  onChange={(e) => setYear(Number(e.target.value))}
                />
              </div>
              <Button variant="outline" onClick={nextMonth} aria-label="Следующий месяц" className="shrink-0">→</Button>
            </div>
            <Button onClick={handleCreate} className="w-full span sm:w-auto flex items-center justify-center gap-2" aria-label="Добавить бюджет">
              <Plus className="h-5 w-5 md:hidden" />
              <span className="hidden md:inline">Добавить</span>
            </Button>
          </div>
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

