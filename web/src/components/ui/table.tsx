import * as React from 'react';
import { cn } from '@/lib/utils';

function Table({ className, ...props }: React.HTMLAttributes<HTMLTableElement>) {
  return <table className={cn('w-full caption-bottom text-sm', className)} {...props} />;
}

function Thead({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) {
  return <thead className={cn('bg-[rgb(var(--muted))]', className)} {...props} />;
}

function Tbody({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) {
  return <tbody className={cn('bg-[rgb(var(--card))] text-[rgb(var(--app-fg))]', className)} {...props} />;
}

function Tfoot({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) {
  return <tfoot className={cn('bg-[rgb(var(--muted))]', className)} {...props} />;
}

function Tr({ className, ...props }: React.HTMLAttributes<HTMLTableRowElement>) {
  return <tr className={cn('border-b border-[rgb(var(--border))]', className)} {...props} />;
}

function Th({ className, ...props }: React.ThHTMLAttributes<HTMLTableCellElement>) {
  return (
    <th
      className={cn('px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-[rgb(var(--muted-foreground))]', className)}
      {...props}
    />
  );
}

function Td({ className, ...props }: React.TdHTMLAttributes<HTMLTableCellElement>) {
  return <td className={cn('px-6 py-4 align-middle', className)} {...props} />;
}

export { Table, Thead, Tbody, Tfoot, Tr, Th, Td };
