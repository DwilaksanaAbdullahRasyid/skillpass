import { api } from '@/lib/api';

export interface HrisRole {
  id: string;
  name: string;
  description?: string;
  isSystem: boolean;
}

export interface MyPermissions {
  permissions: string[];
  roles: HrisRole[];
}

export function getMyPermissions(): Promise<MyPermissions> {
  return api<MyPermissions>('/hris/me/permissions');
}

export interface HrisPermission {
  id: string;
  code: string;
  module: string;
  description: string;
}

export function listRoles(): Promise<HrisRole[]> {
  return api<HrisRole[]>('/hris/roles');
}

export function listPermissions(): Promise<HrisPermission[]> {
  return api<HrisPermission[]>('/hris/permissions');
}

export async function getRolePermissionIds(roleId: string): Promise<string[]> {
  const res = await api<{ permissionIds: string[] }>(`/hris/roles/${roleId}/permissions`);
  return res.permissionIds ?? [];
}

export function setRolePermissions(roleId: string, permissionIds: string[]): Promise<void> {
  return api(`/hris/roles/${roleId}/permissions`, {
    method: 'PUT',
    body: { permissionIds },
  });
}

export function createRole(name: string, description?: string): Promise<HrisRole> {
  return api<HrisRole>('/hris/roles', {
    method: 'POST',
    body: { name, description },
  });
}

export function updateRole(roleId: string, name: string, description?: string): Promise<HrisRole> {
  return api<HrisRole>(`/hris/roles/${roleId}`, {
    method: 'PUT',
    body: { name, description },
  });
}

export function deleteRole(roleId: string): Promise<void> {
  return api(`/hris/roles/${roleId}`, { method: 'DELETE' });
}

export function getEmployeeRoles(employeeId: string): Promise<HrisRole[]> {
  return api<HrisRole[]>(`/hris/employees/${employeeId}/roles`);
}

export function assignRole(employeeId: string, roleId: string): Promise<void> {
  return api(`/hris/employees/${employeeId}/roles`, {
    method: 'POST',
    body: { roleId },
  });
}

export function removeRole(employeeId: string, roleId: string): Promise<void> {
  return api(`/hris/employees/${employeeId}/roles/${roleId}`, {
    method: 'DELETE',
  });
}
