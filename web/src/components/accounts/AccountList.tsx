import { useEffect, useState } from "react";
import { type Account, accountsApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Table, Tbody, Td, Th, Thead, Tr } from "@/components/ui/table";
import { Pencil, Trash2 } from "lucide-react";

interface AccountListProps {
  onEdit: (account: Account) => void;
  onDeleted: () => void;
  refreshKey?: number;
}

export function AccountList({ onEdit, onDeleted, refreshKey }: AccountListProps) {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    load();
  }, [refreshKey]);

  async function load() {
    try {
      setLoading(true);
      const res = await accountsApi.getAll();
      setAccounts(res.data.accounts);
    } catch (e) {
      console.error("Failed to load accounts", e);
    } finally {
      setLoading(false);
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await accountsApi.delete(id);
      await load();
      onDeleted();
    } catch (e) {
      console.error("Delete account error", e);
      alert("Не удалось удалить счёт");
    }
  };

  if (loading) return <div className="text-center py-8">Загрузка...</div>;

  if (accounts.length === 0) {
    return <div className="text-center py-8 text-[rgb(var(--muted-foreground))]">Счета не найдены</div>;
  }

  const typeLabel: Record<string, string> = { cash: "Наличные", bank: "Банк", credit: "Кредит" };

  return (
    <div className="bg-[rgb(var(--card))] shadow rounded-lg overflow-hidden">
      <div className="hidden md:block">
        <Table>
          <Thead>
            <Tr>
              <Th>Название</Th>
              <Th>Тип</Th>
              <Th>Порядок</Th>
              <Th>Архивирован</Th>
              <Th>Баланс</Th>
              <Th>Создан</Th>
              <Th>Обновлён</Th>
              <Th></Th>
            </Tr>
          </Thead>
          <Tbody>
            {accounts.map((a) => (
              <Tr key={a.id}>
                <Td>{a.name}</Td>
                <Td>{typeLabel[a.type] || a.type}</Td>
                <Td>{a.order_num}</Td>
                <Td>{a.is_archived ? "Да" : "Нет"}</Td>
                <Td>{a.balance}</Td>
                <Td>{new Date(a.created_at).toLocaleString()}</Td>
                <Td>{new Date(a.updated_at).toLocaleString()}</Td>
                <Td className="text-right">
                  <div className="flex justify-end gap-1">
                    <Button variant="ghost" size="icon" aria-label="Редактировать" onClick={() => onEdit(a)}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" aria-label="Удалить" onClick={() => handleDelete(a.id)}>
                      <Trash2 className="h-4 w-4 text-[rgb(var(--destructive))]" />
                    </Button>
                  </div>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </div>

      <div className="md:hidden divide-y divide-[rgb(var(--border))]">
        {accounts.map((a) => (
          <div key={a.id} className="p-3">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0 flex-1">
                <div className="text-sm font-medium">{a.name}</div>
                <div className="text-xs text-[rgb(var(--muted-foreground))] mt-1">
                  <div>Тип: {typeLabel[a.type] || a.type}</div>
                  <div>Порядок: {a.order_num}</div>
                  <div>Архивирован: {a.is_archived ? "Да" : "Нет"}</div>
                  <div>Баланс: {a.balance}</div>
                </div>
              </div>
              <div className="flex-shrink-0 flex items-center gap-1">
                <Button variant="ghost" size="icon" aria-label="Редактировать" onClick={() => onEdit(a)}>
                  <Pencil className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="icon" aria-label="Удалить" onClick={() => handleDelete(a.id)}>
                  <Trash2 className="h-4 w-4 text-[rgb(var(--destructive))]" />
                </Button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
