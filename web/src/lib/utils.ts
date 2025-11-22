import { type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import clsx from 'clsx';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const rubFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
  minimumFractionDigits: 2,
});

export function formatCurrency(value: number | string) {
  const num = typeof value === 'string' ? Number(value) : value;
  if (!isFinite(num)) return String(value);
  return rubFormatter.format(num);
}
