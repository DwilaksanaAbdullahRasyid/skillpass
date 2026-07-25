import { Outlet } from 'react-router-dom';
import HRISSidebar from '@/components/hris/HRISSidebar';
import HRISTopbar from '@/components/hris/HRISTopbar';

export default function HRISLayout() {
  return (
    // Scope the Talenta theme to the HRIS console only; the rest of the
    // marketplace app keeps its winter/dark themes.
    <div data-theme="talenta" className="flex min-h-[calc(100vh-4rem)] flex-col bg-base-200 text-base-content">
      <HRISTopbar />
      <div className="flex flex-1 overflow-hidden">
        <HRISSidebar />
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
