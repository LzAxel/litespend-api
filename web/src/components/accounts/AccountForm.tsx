import { useEffect, useMemo, useState } from "react";
import { accountsApi, type Account, type AccountType, type CreateAccountRequest, type UpdateAccountRequest } from "@/lib/api";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface AccountFormProps {
  account?: Account | null;
  onClose: () => void;
  onSuccess: () => void;
}

export function AccountForm({ account, onClose, onSuccess }: AccountFormProps) {
  const [form, setForm] = useState<{
    name: string;
    type: AccountType;
    is_archived: boolean;
    order_num: number | "";
  }>({
    name: account?.name ?? "",
    type: (account?.type as AccountType) ?? "cash",
    is_archived: account?.is_archived ?? false,
    order_num: account?.order_num ?? "",
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const typeLabel = useMemo(() => ({
    cash: "Наличные",
    bank: "Банк",
    credit: "Кредит",
  } as Record<AccountType, string>), []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!form.name || form.order_num === "") return;
    setSubmitting(true);
    try {
      if (account) {
        const dto: UpdateAccountRequest = {
          name: form.name,
          is_archived: form.is_archived,
          order_num: Number(form.order_num),
        };
        await accountsApi.update(account.id, dto);
      } else {
        const dto: CreateAccountRequest = {
          name: form.name,
          type: form.type,
          is_archived: form.is_archived,
          order_num: Number(form.order_num),
        };
        await accountsApi.create(dto);
      }
      onSuccess();
    } catch (err: any) {
      setError(err?.response?.data?.error || "Ошибка сохранения счёта");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-800 rounded-md text-sm">{error}</div>
      )}
      <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <Label htmlFor="name">Название</Label>
          <Input
            id="name"
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            className="mt-1"
            required
          />
        </div>
        <div>
          <Label htmlFor="order_num">Порядок</Label>
          <Input
            id="order_num"
            type="number"
            value={form.order_num as number | ""}
            onChange={(e) => setForm({ ...form, order_num: e.target.value === "" ? "" : Number(e.target.value) })}
            className="mt-1"
            required
          />
        </div>
        {!account && (
          <div>
            <Label>Тип</Label>
            <Select value={form.type} onValueChange={(v) => setForm({ ...form, type: v as AccountType })}>
              <SelectTrigger>
                <SelectValue placeholder="Выберите тип" />
              </SelectTrigger>
              <SelectContent>
                {(["cash", "bank", "credit"] as AccountType[]).map((t) => (
                  <SelectItem key={t} value={t}>{typeLabel[t]}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}
        <div>
          <Label>Архивирован</Label>
          <Select value={String(form.is_archived)} onValueChange={(v) => setForm({ ...form, is_archived: v === "true" })}>
            <SelectTrigger>
              <SelectValue placeholder="Выберите" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="false">Нет</SelectItem>
              <SelectItem value="true">Да</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="md:col-span-2 flex flex-col-reverse sm:flex-row justify-end gap-2 mt-2">
          <Button type="button" variant="outline" onClick={onClose}>Отмена</Button>
          <Button type="submit" disabled={submitting}>{submitting ? "Сохранение..." : "Сохранить"}</Button>
        </div>
      </form>
    </div>
  );
}
