import { useState } from "react";
import { categoriesApi, type Category } from "@/lib/api";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

interface CategoryFormProps {
  category?: Category | null;
  groups?: string[];
  onClose: () => void;
  onSuccess: () => void;
}

export function CategoryForm({
  category,
  onClose,
  onSuccess,
  ...props
}: CategoryFormProps) {
  const [formData, setFormData] = useState({
    name: category?.name || "",
    group_name: category?.group_name || "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [groups, setGroups] = useState<string[]>([...(props.groups || [])]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      if (category) {
        await categoriesApi.update(category.id, formData);
      } else {
        await categoriesApi.create(formData);
      }
      onSuccess();
    } catch (err) {
      setError(err.response?.data?.error || "Ошибка сохранения");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-800 rounded-md text-sm">
          {error}
        </div>
      )}
      <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4">
        <div>
          <Label className="mb-1" htmlFor="name">
            Название
          </Label>
          <Input
            id="name"
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
        </div>
        <div>
          <Label className="mb-1" htmlFor="group">
            Группа
          </Label>
          <Select
            onValueChange={(value) =>
              setFormData({ ...formData, group_name: value })
            }
            value={formData.group_name}
          >
            <SelectTrigger>
              <SelectValue placeholder="Выберите группу" />
            </SelectTrigger>
            <SelectContent>
              <Input
                className="mb-2"
                placeholder="Создать группу (Enter)"
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    setGroups([...groups, e.currentTarget.value]);
                  }
                }}
              />
              {groups?.map((group) => (
                <SelectItem key={group} value={group}>
                  {group}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 mt-2">
          <Button type="button" variant="outline" onClick={onClose}>
            Отмена
          </Button>
          <Button type="submit" disabled={loading}>
            {loading ? "Сохранение..." : "Сохранить"}
          </Button>
        </div>
      </form>
    </div>
  );
}
