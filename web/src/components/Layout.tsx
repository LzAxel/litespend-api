import { Link, useRouterState } from '@tanstack/react-router';
import { useAuthStore } from '@/store/auth';
import { authApi } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { useState, useMemo } from 'react';
import { Dialog, DialogContent } from '@/components/ui/dialog';
import { Menu, ReceiptText, Folder, ListChecks, Wallet, Upload, BarChart3 } from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const router = useRouterState();
  const { isAuthenticated, setAuthenticated } = useAuthStore();
  const [mobileOpen, setMobileOpen] = useState(false);

  const currentPath = router.location.pathname;
  const isAuthPage = currentPath === '/login' || currentPath === '/register';

  const navItems = useMemo(
    () => [
      { to: '/transactions', label: 'Транзакции', icon: ReceiptText },
      { to: '/categories', label: 'Категории', icon: Folder },
      { to: '/prescribed-expanses', label: 'Обязательные траты', icon: ListChecks },
      { to: '/budgets', label: 'Бюджеты', icon: Wallet },
      { to: '/import', label: 'Импорт', icon: Upload },
      { to: '/statistics', label: 'Статистика', icon: BarChart3 },
    ],
    []
  );

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
      <div className="flex min-h-screen">
        <aside className="hidden md:flex md:w-60 md:flex-col md:border-r md:bg-white md:sticky md:top-0 md:h-screen">
          <div className="p-4 text-lg font-semibold">LiteSpend</div>
          <nav className="flex-1 px-2 py-2 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.to}
                to={item.to}
                className="flex items-center gap-2 rounded-md px-3 py-2 text-sm"
                activeProps={{ className: 'bg-blue-50 text-blue-700' }}
                inactiveProps={{ className: 'text-gray-700 hover:bg-gray-50' }}
              >
                <item.icon className="h-4 w-4" />
                <span>{item.label}</span>
              </Link>
            ))}
          </nav>
          <div className="p-4 border-t">
            {isAuthenticated && (
              <Button variant="outline" onClick={handleLogout} className="w-full">Выход</Button>
            )}
          </div>
        </aside>

        <div className="flex-1 flex flex-col">
          <header className="flex items-center justify-between h-14 border-b bg-white px-4 md:px-6 md:hidden">
            <Button variant="outline" size="icon" onClick={() => setMobileOpen(true)} aria-label="Меню">
              <Menu className="h-5 w-5" />
            </Button>
            {isAuthenticated && (
              <Button variant="outline" onClick={handleLogout}>Выход</Button>
            )}
          </header>
          <main className="p-4 sm:p-6 lg:p-8 md:ml-0">{children}</main>
        </div>
      </div>

      <Dialog open={mobileOpen} onOpenChange={(o) => !o && setMobileOpen(false)}>
        <DialogContent className="p-0 w-[280px] h-full fixed left-0 top-0 rounded-none">
          <div className="flex flex-col h-full">
            <div className="p-4 text-lg font-semibold border-b">LiteSpend</div>
            <nav className="flex-1 px-2 py-2 space-y-1">
              {navItems.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  className="flex items-center gap-2 rounded-md px-3 py-2 text-sm"
                  activeProps={{ className: 'bg-blue-50 text-blue-700' }}
                  inactiveProps={{ className: 'text-gray-700 hover:bg-gray-50' }}
                  onClick={() => setMobileOpen(false)}
                >
                  <item.icon className="h-4 w-4" />
                  <span>{item.label}</span>
                </Link>
              ))}
            </nav>
            <div className="p-4 border-t">
              {isAuthenticated && (
                <Button variant="outline" onClick={handleLogout} className="w-full">Выход</Button>
              )}
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}

