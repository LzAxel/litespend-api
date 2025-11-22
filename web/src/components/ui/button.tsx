import * as React from 'react';
import { cn } from '@/lib/utils';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'secondary' | 'destructive' | 'outline' | 'ghost';
  size?: 'default' | 'sm' | 'lg' | 'icon';
}

const variants: Record<NonNullable<ButtonProps['variant']>, string> = {
  default: 'bg-[rgb(var(--primary))] text-[rgb(var(--primary-foreground))] shadow hover:opacity-90',
  secondary: 'bg-[rgb(var(--muted))] text-[rgb(var(--app-fg))] hover:opacity-90',
  destructive: 'bg-[rgb(var(--destructive))] text-[rgb(var(--destructive-foreground))] shadow-sm hover:opacity-90',
  outline: 'border border-[rgb(var(--border))] bg-[rgb(var(--card))] text-[rgb(var(--app-fg))] hover:bg-[rgb(var(--muted))]',
  ghost: 'text-[rgb(var(--app-fg))] hover:bg-[rgb(var(--muted))]',
};

const sizes: Record<NonNullable<ButtonProps['size']>, string> = {
  default: 'h-10 px-4 py-2',
  sm: 'h-9 rounded-md px-3',
  lg: 'h-11 rounded-md px-8',
  icon: 'h-10 w-10',
};

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'default', size = 'default', ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-[rgb(var(--ring))] focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      />
    );
  }
);
Button.displayName = 'Button';
