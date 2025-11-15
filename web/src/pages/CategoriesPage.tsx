import { useState, useEffect } from 'react';
import { categoriesApi, type Category } from '@/lib/api';
import { CategoryForm } from '@/components/categories/CategoryForm';
import { CategoryList } from '@/components/categories/CategoryList';

export function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const response = await categoriesApi.getAll();
      setCategories(response.data);
    } catch (error) {
      console.error('Error loading categories:', error);
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
    if (!confirm('Удалить категорию?')) return;
    try {
      await categoriesApi.delete(id);
      loadCategories();
    } catch (error) {
      console.error('Error deleting category:', error);
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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Категории</h1>
        <button
          onClick={handleCreate}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          Добавить категорию
        </button>
      </div>

      {showForm && (
        <CategoryForm
          category={editingCategory}
          onClose={handleFormClose}
          onSuccess={handleFormSuccess}
        />
      )}

      <CategoryList categories={categories} onEdit={handleEdit} onDelete={handleDelete} />
    </div>
  );
}

