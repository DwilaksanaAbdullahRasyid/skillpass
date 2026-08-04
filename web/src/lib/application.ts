import { api } from '@/lib/api';
import type { ApplicationMessage, ApplicationResult } from './api-types';

// ApplicationStatus is a client-side narrowing of the server's free-form
// status string; used to drive the kanban columns.
export type ApplicationStatus = 'applied' | 'reviewed' | 'interviewed' | 'offered' | 'rejected';

// Application keeps the historical name used across the UI; it is the
// generated server response type.
export type Application = ApplicationResult;
export type { ApplicationMessage };

export async function applyToJob(jobId: string): Promise<Application> {
  return api<Application>(`/jobs/${jobId}/apply`, { method: 'POST' });
}

export async function getMyApplications(): Promise<Application[]> {
  return api<Application[]>('/applications/me');
}

export async function getApplicationMessages(applicationId: string): Promise<ApplicationMessage[]> {
  return api<ApplicationMessage[]>(`/applications/${applicationId}/messages`);
}

export async function addApplicationMessage(applicationId: string, body: string): Promise<ApplicationMessage> {
  return api<ApplicationMessage>(`/applications/${applicationId}/messages`, {
    method: 'POST',
    body: { body },
  });
}

export interface ScheduleInterviewPayload {
  scheduledAt: string; // RFC3339 / ISO date-time
  mode: 'onsite' | 'online';
  location?: string;
  meetingLink?: string;
  interviewer?: string;
  notes?: string;
}

// scheduleInterview records an interview appointment and moves the application
// to "interviewed", notifying + emailing the candidate.
export async function scheduleInterview(
  applicationId: string,
  payload: ScheduleInterviewPayload,
): Promise<Application> {
  return api<Application>(`/applications/${applicationId}/interview`, {
    method: 'POST',
    body: payload,
  });
}
