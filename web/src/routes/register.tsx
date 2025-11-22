import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { RegisterForm } from '@/components/auth/RegisterForm';

export const Route = createFileRoute('/register')({
  component: RegisterPage,
});

function RegisterPage() {
  const navigate = useNavigate();

  return (
    <div className="space-y-6">
      <h2 className="text-center text-2xl font-semibold text-[rgb(var(--app-fg))]">
        Регистрация
      </h2>
      <RegisterForm onSuccess={() => navigate({ to: '/login' })} />
    </div>
  );
}

