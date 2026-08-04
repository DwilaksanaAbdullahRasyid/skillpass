import { Bell, Building2, Search } from 'lucide-react';
import { useAuth } from '@/hooks/useAuth';
import { usePermissions } from '@/hooks/usePermissions';

function initials(name?: string): string {
  if (!name) return 'HR';
  return name
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((p) => p[0]?.toUpperCase() ?? '')
    .join('');
}

/**
 * Talenta-style HRIS console top bar: product brand on the left, a global
 * search in the middle, notifications + user/role chip on the right.
 */
export default function HRISTopbar() {
  const { user } = useAuth();
  const { roles } = usePermissions();
  const topRole = roles[0]?.name ?? user?.role ?? 'Member';

  return (
    <header className="flex h-14 shrink-0 items-center gap-4 border-b border-base-300 bg-base-100 px-4">
      <div className="flex items-center gap-2 md:w-56">
        <span className="grid h-8 w-8 place-items-center rounded-lg bg-primary text-primary-content">
          <Building2 className="h-[18px] w-[18px]" />
        </span>
        <span className="text-sm font-semibold tracking-tight text-base-content">HRIS Console</span>
      </div>

      <div className="relative hidden max-w-md flex-1 sm:block">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-base-content/40" />
        <input
          type="search"
          aria-label="Search HRIS"
          placeholder="Search employees, requests…"
          className="input input-sm w-full rounded-lg border-base-300 bg-base-200 pl-9 focus:bg-base-100"
        />
      </div>

      <div className="ml-auto flex items-center gap-2">
        <button
          type="button"
          aria-label="Notifications"
          className="btn btn-ghost btn-sm btn-circle text-base-content/70"
        >
          <Bell className="h-5 w-5" />
        </button>
        <div className="flex items-center gap-2 rounded-lg px-1.5 py-1">
          <span className="grid h-8 w-8 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary">
            {initials(user?.name ?? user?.username)}
          </span>
          <div className="hidden leading-tight sm:block">
            <p className="text-sm font-medium text-base-content">{user?.name ?? user?.username ?? 'User'}</p>
            <p className="text-xs capitalize text-base-content/50">{topRole}</p>
          </div>
        </div>
      </div>
    </header>
  );
}
