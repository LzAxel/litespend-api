import { useState } from 'react';
import { Link } from '@tanstack/react-router';
import { authApi } from '@/lib/api';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';

interface RegisterFormProps {
  onSuccess: () => void;
}

export function RegisterForm({ onSuccess }: RegisterFormProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await authApi.register({ username, password });
      onSuccess();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка регистрации');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
      {error && (
        <div className="rounded-md border p-3 text-sm bg-[rgb(var(--muted))] border-[rgb(var(--destructive))] text-[rgb(var(--destructive))]">
          {error}
        </div>
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
          {loading ? 'Регистрация...' : 'Зарегистрироваться'}
        </Button>
      </div>

      <div className="text-center">
        <Link to="/login" className="font-medium text-[rgb(var(--primary))] hover:opacity-90">
          Уже есть аккаунт? Войти
        </Link>
      </div>
    </form>
  );
}

