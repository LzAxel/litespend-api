import { Link, useRouterState } from '@tanstack/react-router';
import { useAuthStore } from '@/store/auth';
import { authApi } from '@/lib/api';
import { Button } from '@/components/ui/button';

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const router = useRouterState();
  const { isAuthenticated, setAuthenticated } = useAuthStore();

  const currentPath = router.location.pathname;
  const isAuthPage = currentPath === '/login' || currentPath === '/register';

  const handleLogout = async () => {
    try {
      await authApi.logout();
      setAuthenticated(false);
      window.location.href = '/login';
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  if (isAuthPage) {
    return <>{children}</>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white/80 backdrop-blur supports-[backdrop-filter]:bg-white/60 border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex gap-6">
              <Link
                to="/transactions"
                className="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                activeProps={{
                  className: 'border-blue-500 text-gray-900',
                }}
                inactiveProps={{
                  className: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                }}
              >
                Транзакции
              </Link>
              <Link
                to="/categories"
                className="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                activeProps={{
                  className: 'border-blue-500 text-gray-900',
                }}
                inactiveProps={{
                  className: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                }}
              >
                Категории
              </Link>
              <Link
                to="/prescribed-expanses"
                className="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                activeProps={{
                  className: 'border-blue-500 text-gray-900',
                }}
                inactiveProps={{
                  className: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                }}
              >
                Обязательные траты
              </Link>
              <Link
                to="/budgets"
                className="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                activeProps={{
                  className: 'border-blue-500 text-gray-900',
                }}
                inactiveProps={{
                  className: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                }}
              >
                Бюджеты
              </Link>
              <Link
                to="/import"
                className="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                activeProps={{
                  className: 'border-blue-500 text-gray-900',
                }}
                inactiveProps={{
                  className: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                }}
              >
                Импорт
              </Link>
              <Link
                to="/statistics"
                className="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                activeProps={{
                  className: 'border-blue-500 text-gray-900',
                }}
                inactiveProps={{
                  className: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                }}
              >
                Статистика
              </Link>
            </div>
            <div className="flex items-center">
              {isAuthenticated && (
                <Button variant="outline" onClick={handleLogout}>
                  Выход
                </Button>
              )}
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">{children}</main>
    </div>
  );
}

