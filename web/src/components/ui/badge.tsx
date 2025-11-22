import * as React from 'react';
import { cn } from '@/lib/utils';

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'secondary' | 'destructive' | 'outline';
}

const variants: Record<NonNullable<BadgeProps['variant']>, string> = {
  default: 'bg-[rgb(var(--primary))] text-[rgb(var(--primary-foreground))]',
  secondary: 'bg-[rgb(var(--muted))] text-[rgb(var(--app-fg))]',
  destructive: 'bg-[rgb(var(--destructive))] text-[rgb(var(--destructive-foreground))]',
  outline: 'border border-[rgb(var(--border))] text-[rgb(var(--app-fg))]',
};

function Badge({ className, variant = 'default', ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium',
        variants[variant],
        className
      )}
      {...props}
    />
  );
}

export { Badge };
