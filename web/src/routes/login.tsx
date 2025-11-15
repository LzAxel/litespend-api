import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { LoginForm } from '@/components/auth/LoginForm';

export const Route = createFileRoute('/login')({
  component: LoginPage,
});

function LoginPage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
        <div>
          <h2 className="text-center text-3xl font-extrabold text-gray-900">
            Вход в систему
          </h2>
        </div>
        <LoginForm onSuccess={() => navigate({ to: '/transactions' })} />
      </div>
    </div>
  );
}

