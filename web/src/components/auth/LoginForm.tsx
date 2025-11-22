import { useState } from 'react';
import { Link } from '@tanstack/react-router';
import { authApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';

interface LoginFormProps {
  onSuccess: () => void;
}

export function LoginForm({ onSuccess }: LoginFormProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const setAuthenticated = useAuthStore((state) => state.setAuthenticated);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await authApi.login({ username, password });
      setAuthenticated(true);
      onSuccess();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка входа');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
      {error && (
        <div className="rounded-md border p-3 text-sm bg-[rgb(var(--muted))] border-[rgb(var(--destructive))] text-[rgb(var(--destructive))]">{error}</div>
      )}
      <div className="space-y-4">
        <div>
          <Label htmlFor="username">Имя пользователя</Label>
          <Input
            id="username"
            name="username"
            type="text"
            required
            placeholder="Имя пользователя"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div>
          <Label htmlFor="password">Пароль</Label>
          <Input
            id="password"
            name="password"
            type="password"
            required
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
      </div>

      <div>
        <Button type="submit" disabled={loading} className="w-full">
          {loading ? 'Вход...' : 'Войти'}
        </Button>
      </div>

      <div className="text-center">
        <Link to="/register" className="font-medium text-[rgb(var(--primary))] hover:opacity-90">
          Нет аккаунта? Зарегистрироваться
        </Link>
      </div>
    </form>
  );
}

