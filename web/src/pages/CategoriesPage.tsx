import { useState, useEffect } from "react";
import { categoriesApi, type Category } from "@/lib/api";
import { CategoryForm } from "@/components/categories/CategoryForm";
import { CategoryList } from "@/components/categories/CategoryList";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Plus } from "lucide-react";

export function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [toDelete, setToDelete] = useState<Category | null>(null);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const response = await categoriesApi.getAll();
      setCategories(response.data);
    } catch (error) {
      console.error("Error loading categories:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingCategory(null);
    setShowForm(true);
  };

  const handleEdit = (category: Category) => {
    setEditingCategory(category);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await categoriesApi.delete(id);
      loadCategories();
    } catch (error) {
      console.error("Error deleting category:", error);
    }
  };

  const handleFormClose = () => {
    setShowForm(false);
    setEditingCategory(null);
  };

  const handleFormSuccess = () => {
    handleFormClose();
    loadCategories();
  };

  if (loading) {
    return <div className="text-center py-8">Загрузка...</div>;
  }

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="flex justify-between items-center mb-6 gap-3">
        <h1 className="text-2xl sm:text-3xl font-bold">Категории</h1>
        <Button
          onClick={handleCreate}
          className="flex items-center gap-2"
          aria-label="Добавить категорию"
        >
          <Plus className="h-5 w-5 md:hidden" />
          <span className="hidden md:inline">Добавить категорию</span>
        </Button>
      </div>

      <Dialog open={showForm} onOpenChange={(o) => !o && handleFormClose()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingCategory
                ? "Редактировать категорию"
                : "Добавить категорию"}
            </DialogTitle>
          </DialogHeader>
          <CategoryForm
            groups={categories.reduce((acc, category) => {
              if (!acc.includes(category.group_name) && category.group_name) {
                acc.push(category.group_name);
              }
              return acc;
            }, [])}
            category={editingCategory}
            onClose={handleFormClose}
            onSuccess={handleFormSuccess}
          />
        </DialogContent>
      </Dialog>

      <CategoryList
        categories={categories}
        onEdit={handleEdit}
        onDelete={(id) =>
          setToDelete(categories.find((c) => c.id === id) || null)
        }
      />

      <Dialog open={!!toDelete} onOpenChange={(o) => !o && setToDelete(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Удалить категорию?</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-gray-600">Действие нельзя отменить.</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setToDelete(null)}>
              Отмена
            </Button>
            <Button
              variant="destructive"
              onClick={async () => {
                if (toDelete) {
                  await handleDelete(toDelete.id);
                  setToDelete(null);
                }
              }}
            >
              Удалить
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
