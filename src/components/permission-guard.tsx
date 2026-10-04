import React from 'react'
import { useAuthStore } from '@/stores/auth-store'

export interface PermissionGuardProps {
  /**
   * Single permission code required (e.g. 'apps:create')
   */
  permission?: string
  /**
   * Array of permission codes required
   */
  permissions?: string[]
  /**
   * If true, user must have all specified permissions. If false (default), having any permission is enough.
   */
  requireAll?: boolean
  /**
   * Single role required (e.g. 'admin')
   */
  role?: string
  /**
   * Array of roles (any match grants access)
   */
  roles?: string[]
  /**
   * Optional fallback content when access is denied
   */
  fallback?: React.ReactNode
  /**
   * Children components rendered when access is granted
   */
  children: React.ReactNode
}

/**
 * PermissionGuard conditionally renders child elements based on user roles and granular permissions.
 */
export function PermissionGuard({
  permission,
  permissions,
  requireAll = false,
  role,
  roles,
  fallback = null,
  children,
}: PermissionGuardProps) {
  const { hasPermission, hasAnyPermission, hasAllPermissions, hasRole, hasAnyRole } =
    useAuthStore((state) => state.auth)

  // 1. Role checks
  if (role && !hasRole(role)) {
    return <>{fallback}</>
  }

  if (roles && roles.length > 0 && !hasAnyRole(roles)) {
    return <>{fallback}</>
  }

  // 2. Single permission check
  if (permission && !hasPermission(permission)) {
    return <>{fallback}</>
  }

  // 3. Multiple permissions check
  if (permissions && permissions.length > 0) {
    const hasAccess = requireAll
      ? hasAllPermissions(permissions)
      : hasAnyPermission(permissions)
    if (!hasAccess) {
      return <>{fallback}</>
    }
  }

  return <>{children}</>
}
