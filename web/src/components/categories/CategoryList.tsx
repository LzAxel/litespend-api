import { type Category } from '@/lib/api';

interface CategoryListProps {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (id: number) => void;
}

export function CategoryList({ categories, onEdit, onDelete }: CategoryListProps) {
  if (categories.length === 0) {
    return <div className="text-center py-8 text-gray-500">Нет категорий</div>;
  }

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-md">
      <ul className="divide-y divide-gray-200">
        {categories.map((category) => (
          <li key={category.id}>
            <div className="px-4 py-4 sm:px-6">
              <div className="flex items-center justify-between">
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900">{category.name}</p>
                  <p className="text-sm text-gray-500">
                    {category.type === 0 ? 'Доход' : 'Расход'}
                  </p>
                </div>
                <div className="ml-4 flex-shrink-0 flex space-x-2">
                  <button
                    onClick={() => onEdit(category)}
                    className="text-blue-600 hover:text-blue-900"
                  >
                    Редактировать
                  </button>
                  <button
                    onClick={() => onDelete(category.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    Удалить
                  </button>
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

