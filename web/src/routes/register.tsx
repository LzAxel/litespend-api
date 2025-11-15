import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { RegisterForm } from '@/components/auth/RegisterForm';

export const Route = createFileRoute('/register')({
  component: RegisterPage,
});

function RegisterPage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
        <div>
          <h2 className="text-center text-3xl font-extrabold text-gray-900">
            Регистрация
          </h2>
        </div>
        <RegisterForm onSuccess={() => navigate({ to: '/login' })} />
      </div>
    </div>
  );
}

