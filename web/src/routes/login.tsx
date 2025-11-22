import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { LoginForm } from '@/components/auth/LoginForm';

export const Route = createFileRoute('/login')({
  component: LoginPage,
});

function LoginPage() {
  const navigate = useNavigate();

  return (
    <div className="space-y-6">
      <h2 className="text-center text-2xl font-semibold text-[rgb(var(--app-fg))]">
        Вход в систему
      </h2>
      <LoginForm onSuccess={() => navigate({ to: '/transactions' })} />
    </div>
  );
}

