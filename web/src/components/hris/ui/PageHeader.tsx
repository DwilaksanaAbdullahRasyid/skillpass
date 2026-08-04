import type { ReactNode } from 'react';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  /** Right-aligned actions (buttons, links). */
  actions?: ReactNode;
}

/**
 * Talenta-style page header: large title on the left, action buttons on the
 * right, an optional muted subtitle underneath.
 */
export function PageHeader({ title, subtitle, actions }: PageHeaderProps) {
  return (
    <div className="flex flex-wrap items-start justify-between gap-3 mb-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight text-base-content">{title}</h1>
        {subtitle && <p className="text-sm text-base-content/60 mt-1">{subtitle}</p>}
      </div>
      {actions && <div className="flex flex-wrap items-center gap-2">{actions}</div>}
    </div>
  );
}
