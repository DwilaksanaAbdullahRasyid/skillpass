import type { LucideIcon } from 'lucide-react';

interface StatCardProps {
  label: string;
  value: string | number;
  hint?: string;
  icon?: LucideIcon;
  /** Accent color for the icon chip. */
  tone?: 'primary' | 'success' | 'warning' | 'info' | 'error' | 'neutral';
}

const toneClasses: Record<NonNullable<StatCardProps['tone']>, string> = {
  primary: 'bg-primary/10 text-primary',
  success: 'bg-success/10 text-success',
  warning: 'bg-warning/10 text-warning',
  info: 'bg-info/10 text-info',
  error: 'bg-error/10 text-error',
  neutral: 'bg-base-content/10 text-base-content',
};

/**
 * Talenta-style summary metric card used at the top of list pages
 * (e.g. Total employees / New hires / Leaving).
 */
export function StatCard({ label, value, hint, icon: Icon, tone = 'primary' }: StatCardProps) {
  return (
    <div className="rounded-box border border-base-300 bg-base-100 p-5">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="text-sm font-medium text-base-content/60">{label}</p>
          <p className="mt-2 text-3xl font-semibold text-base-content">{value}</p>
          {hint && <p className="mt-1 text-xs text-base-content/50">{hint}</p>}
        </div>
        {Icon && (
          <span className={`grid h-10 w-10 shrink-0 place-items-center rounded-lg ${toneClasses[tone]}`}>
            <Icon className="h-5 w-5" />
          </span>
        )}
      </div>
    </div>
  );
}
