import { createFileRoute } from '@tanstack/react-router';
import { StatisticsPage } from '@/pages/StatisticsPage';

export const Route = createFileRoute('/statistics')({
  component: StatisticsPage,
});

