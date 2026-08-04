import type { LucideIcon } from 'lucide-react';
import { Inbox } from 'lucide-react';
import type { ReactNode } from 'react';

interface EmptyStateProps {
  title: string;
  description?: string;
  icon?: LucideIcon;
  action?: ReactNode;
}

/**
 * Talenta-style illustrated empty state: centered icon, heading, muted
 * description and an optional call-to-action.
 */
export function EmptyState({ title, description, icon: Icon = Inbox, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center px-6 py-16 text-center">
      <div className="grid h-20 w-20 place-items-center rounded-full bg-base-200 text-base-content/40">
        <Icon className="h-9 w-9" />
      </div>
      <h3 className="mt-5 text-base font-semibold text-base-content">{title}</h3>
      {description && <p className="mt-1 max-w-sm text-sm text-base-content/60">{description}</p>}
      {action && <div className="mt-5">{action}</div>}
    </div>
  );
}
