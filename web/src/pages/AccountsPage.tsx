import { useState } from "react";
import { type Account } from "@/lib/api";
import { AccountForm } from "@/components/accounts/AccountForm";
import { AccountList } from "@/components/accounts/AccountList";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";

export function AccountsPage() {
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState<Account | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreate = () => {
    setEditing(null);
    setShowForm(true);
  };

  const handleEdit = (a: Account) => {
    setEditing(a);
    setShowForm(true);
  };

  const closeForm = () => {
    setShowForm(false);
    setEditing(null);
  };

  const onSaved = () => {
    setShowForm(false);
    setEditing(null);
    setRefreshKey((k) => k + 1);
  };

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl sm:text-3xl font-bold">Счета</h1>
        <Button onClick={handleCreate} className="flex items-center gap-2" aria-label="Добавить счёт">
          <Plus className="h-5 w-5 md:hidden" />
          <span className="hidden md:inline">Добавить</span>
        </Button>
      </div>

      <Dialog open={showForm} onOpenChange={(o) => !o && closeForm()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? "Редактировать счёт" : "Создать счёт"}</DialogTitle>
          </DialogHeader>
          <AccountForm account={editing} onClose={closeForm} onSuccess={onSaved} />
        </DialogContent>
      </Dialog>

      <AccountList onEdit={handleEdit} onDeleted={() => setRefreshKey((k) => k + 1)} refreshKey={refreshKey} />
    </div>
  );
}
