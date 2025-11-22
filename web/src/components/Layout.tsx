import {Link, useRouterState} from '@tanstack/react-router';
import {useAuthStore} from '@/store/auth';
import {authApi} from '@/lib/api';
import {Button} from '@/components/ui/button';
import {useEffect, useMemo, useState} from 'react';
import {Dialog, DialogContent} from '@/components/ui/dialog';
import {BarChart3, Folder, Menu, Moon, ReceiptText, Sun, Upload, Wallet} from 'lucide-react';

interface LayoutProps {
    children: React.ReactNode;
}

export function Layout({children}: LayoutProps) {
    const router = useRouterState();
    const {setAuthenticated} = useAuthStore();
    const [mobileOpen, setMobileOpen] = useState(false);
    const [theme, setTheme] = useState<'light' | 'dark'>(() => {
        if (typeof window === 'undefined') return 'light';
        const saved = localStorage.getItem('theme');
        return (saved === 'dark' || saved === 'light') ? (saved as 'light' | 'dark') : 'light';
    });

    const currentPath = router.location.pathname;
    const isAuthPage = currentPath === '/login' || currentPath === '/register';

    useEffect(() => {
        const root = document.documentElement;
        if (theme === 'dark') {
            root.classList.add('dark');
        } else {
            root.classList.remove('dark');
        }
        localStorage.setItem('theme', theme);
    }, [theme]);

    const toggleTheme = () => setTheme((t) => (t === 'dark' ? 'light' : 'dark'));

    const navItems = useMemo(
        () => [
            {to: '/transactions', label: 'Транзакции', icon: ReceiptText},
            {to: '/categories', label: 'Категории', icon: Folder},
            {to: '/budgets', label: 'Бюджеты', icon: Wallet},
            {to: '/import', label: 'Импорт', icon: Upload},
            {to: '/statistics', label: 'Статистика', icon: BarChart3},
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
        return (
            <div className="min-h-screen bg-[rgb(var(--app-bg))] text-[rgb(var(--app-fg))]">
                <header className="flex items-center justify-end h-14 border-b border-[rgb(var(--border))] bg-[rgb(var(--card))] px-4">
                    <Button
                        variant="outline"
                        className="flex items-center gap-2"
                        onClick={toggleTheme}
                        aria-label="Переключить тему"
                    >
                        {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
                        <span className="text-sm">{theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}</span>
                    </Button>
                </header>
                <main className="min-h-[calc(100vh-3.5rem)] flex items-center justify-center p-4">
                    <div className="w-full max-w-md bg-[rgb(var(--card))] border border-[rgb(var(--border))] rounded-lg p-6">
                        {children}
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[rgb(var(--app-bg))] text-[rgb(var(--app-fg))]">
            <div className="flex min-h-screen">
                <aside
                    className="hidden md:flex md:w-60 md:flex-col md:border-r md:bg-[rgb(var(--card))] md:border-[rgb(var(--border))] md:sticky md:top-0 md:h-screen">
                    <div className="p-4 text-lg font-semibold">LiteSpend</div>
                    <nav className="flex-1 px-2 py-2 space-y-1">
                        {navItems.map((item) => (
                            <Link
                                key={item.to}
                                to={item.to}
                                className="flex items-center gap-2 rounded-md px-3 py-2 text-sm"
                                activeProps={{className: 'bg-[rgb(var(--muted))] text-[rgb(var(--app-fg))]'}}
                                inactiveProps={{className: 'text-[rgb(var(--app-fg))] hover:bg-[rgb(var(--muted))]'}}
                            >
                                <item.icon className="h-4 w-4"/>
                                <span>{item.label}</span>
                            </Link>
                        ))}
                    </nav>
                    <div className="p-4 border-t border-[rgb(var(--border))] space-y-2">
                        <Button variant="outline" className="w-full flex items-center justify-center gap-2"
                                onClick={toggleTheme} aria-label="Переключить тему">
                            {theme === 'dark' ? <Sun className="h-4 w-4"/> : <Moon className="h-4 w-4"/>}
                            <span className="text-sm">{theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}</span>
                        </Button>
                        <Button variant="outline" onClick={handleLogout} className="w-full">Выход</Button>
                    </div>
                </aside>

                <div className="flex-1 flex flex-col">
                    <header
                        className="flex items-center justify-between h-14 border-b border-[rgb(var(--border))] bg-[rgb(var(--card))] px-4 md:px-6 md:hidden">
                        <Button variant="outline" size="icon" onClick={() => setMobileOpen(true)} aria-label="Меню">
                            <Menu className="h-5 w-5"/>
                        </Button>
                        <Button variant="outline" onClick={handleLogout}>Выход</Button>
                    </header>
                    <main className="p-4 sm:p-6 lg:p-8 md:ml-0">{children}</main>
                </div>
            </div>

            <Dialog open={mobileOpen} onOpenChange={(o) => !o && setMobileOpen(false)}>
                <DialogContent className="!p-0 w-[280px] h-full fixed left-0 top-0 rounded-none max-h-none">
                    <div className="flex flex-col h-full">
                        <div className="p-4 text-lg font-semibold border-b border-[rgb(var(--border))]">LiteSpend</div>
                        <nav className="flex-1 px-2 py-2 space-y-1">
                            {navItems.map((item) => (
                                <Link
                                    key={item.to}
                                    to={item.to}
                                    className="flex items-center gap-2 rounded-md px-3 py-2 text-sm"
                                    activeProps={{className: 'bg-[rgb(var(--muted))] text-[rgb(var(--app-fg))]'}}
                                    inactiveProps={{className: 'text-[rgb(var(--app-fg))] hover:bg-[rgb(var(--muted))]'}}
                                    onClick={() => setMobileOpen(false)}
                                >
                                    <item.icon className="h-4 w-4"/>
                                    <span>{item.label}</span>
                                </Link>
                            ))}
                        </nav>
                        <div className="p-4 border-t border-[rgb(var(--border))] space-y-2">
                            <Button variant="outline" className="w-full flex items-center justify-center gap-2"
                                    onClick={toggleTheme} aria-label="Переключить тему">
                                {theme === 'dark' ? <Sun className="h-4 w-4"/> : <Moon className="h-4 w-4"/>}
                                <span className="text-sm">{theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}</span>
                            </Button>
                            <Button variant="outline" onClick={handleLogout} className="w-full">Выход</Button>
                        </div>
                    </div>
                </DialogContent>
            </Dialog>
        </div>
    );
}

