import * as React from 'react';
import { cn } from '@/lib/utils';

export interface ProgressProps extends React.HTMLAttributes<HTMLDivElement> {
  value?: number;
}

function Progress({ className, value = 0, ...props }: ProgressProps) {
  const clamped = Math.max(0, Math.min(100, value));
  return (
    <div className={cn('relative h-3 w-40 overflow-hidden rounded bg-gray-200', className)} {...props}>
      <div className="h-3 rounded bg-green-500" style={{ width: `${clamped}%` }} />
    </div>
  );
}

export { Progress };
