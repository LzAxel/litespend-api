import { createFileRoute } from '@tanstack/react-router';
import { PrescribedExpansesPage } from '@/pages/PrescribedExpansesPage';

export const Route = createFileRoute('/prescribed-expanses')({
  component: PrescribedExpansesPage,
});
