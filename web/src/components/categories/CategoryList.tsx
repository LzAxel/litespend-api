import { type Category } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Pencil, Trash2 } from "lucide-react";

interface CategoryListProps {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (id: number) => void;
}

export function CategoryList({
  categories,
  onEdit,
  onDelete,
}: CategoryListProps) {
  if (categories.length === 0) {
    return (
      <div className="text-center py-8 text-[rgb(var(--muted-foreground))]">
        Нет категорий
      </div>
    );
  }

  const groupedCategories = categories.reduce(
    (acc, category) => {
      if (!acc[category.group_name]) {
        acc[category.group_name] = [];
      }
      acc[category.group_name]!.push(category);
      return acc;
    },
    {} as Record<string, Category[]>,
  );

  return (
    <div className="bg-[rgb(var(--card))] shadow overflow-hidden sm:rounded-md">
      <ul className="divide-y divide-[rgb(var(--border))]">
        {Object.entries(groupedCategories).map(
          ([categoryGroup, categories]) => (
            <div key={categoryGroup} className="flex flex-col">
              <div className="px-3 py-2 bg-[rgb(var(--muted))]">{categoryGroup || "Без группы"}</div>
              {categories.map((category) => (
                <li key={category.id}>
                  <div className="px-2 py-2 sm:px-6">
                    <div className="flex items-center justify-between gap-3">
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-medium text-[rgb(var(--app-fg))] truncate">
                          {category.name}
                        </p>
                      </div>
                      <div className="flex-shrink-0 flex items-center gap-1 sm:gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          aria-label="Редактировать"
                          onClick={() => onEdit(category)}
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          aria-label="Удалить"
                          onClick={() => onDelete(category.id)}
                        >
                          <Trash2 className="h-4 w-4 text-[rgb(var(--destructive))]" />
                        </Button>
                      </div>
                    </div>
                  </div>
                </li>
              ))}
            </div>
          ),
        )}
      </ul>
    </div>
  );
}
