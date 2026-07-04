import type { components } from "@/generated/api";

import { apiClient } from "@/infrastructure/api/client";

import {
  loginInputSchema,
  loginResponseSchema,
  registerInputSchema,
  registrationResultSchema,
  tenantSelectionSchema,
  tokenResponseSchema,
  type LoginInput,
  type RegisterInput,
  type TenantSelectionInput,
  type UserProfile,
  userProfileSchema,
} from "../schemas/auth";

type LoginResponseData = NonNullable<
  NonNullable<components["schemas"]["LoginResponse"]["data"]>
>;

type LoginTenant = NonNullable<LoginResponseData["tenants"]>[number];

type UserResponseData = NonNullable<components["schemas"]["UserResponse"]["data"]>;

function requireData<T>(value: T | null | undefined, message: string): T {
  if (value == null) {
    throw new Error(message);
  }

  return value;
}

function authHeaders(accessToken: string) {
  return {
    Authorization: `Bearer ${accessToken}`,
  };
}

function normalizeTenant(
  tenant: LoginTenant,
) {
  return {
    tenantId: requireData(tenant.tenant_id, "Missing tenant id"),
    tenantName: requireData(tenant.tenant_name, "Missing tenant name"),
    role: requireData(tenant.role, "Missing tenant role"),
  };
}

function normalizeUser(
  user: UserResponseData,
): UserProfile {
  return userProfileSchema.parse({
    id: requireData(user.id, "Missing user id"),
    tenantId: requireData(user.tenant_id, "Missing tenant id"),
    email: requireData(user.email, "Missing email"),
    role: requireData(user.role, "Missing role"),
    status: requireData(user.status, "Missing status"),
    createdAt: requireData(user.created_at, "Missing created at"),
  });
}

export async function login(input: LoginInput) {
  const body = loginInputSchema.parse(input);
  const response = await apiClient.POST("/auth/login", {
    body,
  });

  if (response.error) {
    throw new Error("Unable to sign in. Check your email and password.");
  }

  const data = requireData(response.data?.data, "Unable to sign in.");
  const tenants = data.tenants ?? [];

  return loginResponseSchema.parse({
    accessToken: requireData(data.access_token, "Missing access token"),
    tenants: tenants.map(normalizeTenant),
  });
}

export async function register(input: RegisterInput) {
  const body = registerInputSchema.parse(input);
  const response = await apiClient.POST("/auth/register", {
    body,
  });

  if (response.error) {
    throw new Error("Unable to create the account.");
  }

  const data = requireData(response.data?.data, "Unable to create the account.");

  return registrationResultSchema.parse({
    id: requireData(data.id, "Missing user id"),
    email: requireData(data.email, "Missing email"),
    status: requireData(data.status, "Missing status"),
    createdAt: requireData(data.created_at, "Missing created at"),
  });
}

export async function selectTenant(
  input: TenantSelectionInput & { accessToken: string },
) {
  const body = tenantSelectionSchema.parse(input);
  const response = await apiClient.POST("/auth/select-tenant", {
    body: {
      tenant_id: body.tenantId,
    },
    headers: authHeaders(input.accessToken),
  });

  if (response.error) {
    throw new Error("Unable to select tenant.");
  }

  const data = requireData(response.data?.data, "Unable to select tenant.");

  return tokenResponseSchema.parse({
    accessToken: requireData(data.access_token, "Missing access token"),
    refreshToken: requireData(data.refresh_token, "Missing refresh token"),
    expiresIn: requireData(data.expires_in, "Missing token expiry"),
  });
}

export async function refreshSession(refreshToken: string) {
  const response = await apiClient.POST("/auth/refresh", {
    body: {
      refresh_token: refreshToken,
    },
  });

  if (response.error) {
    throw new Error("Unable to refresh the session.");
  }

  const data = requireData(response.data?.data, "Unable to refresh the session.");

  return tokenResponseSchema.parse({
    accessToken: requireData(data.access_token, "Missing access token"),
    refreshToken: requireData(data.refresh_token, "Missing refresh token"),
    expiresIn: requireData(data.expires_in, "Missing token expiry"),
  });
}

export async function getCurrentUser(accessToken: string) {
  const response = await apiClient.GET("/users/me", {
    headers: authHeaders(accessToken),
  });

  if (response.error) {
    throw new Error("Unable to load the current user.");
  }

  const data = requireData(response.data?.data, "Unable to load the current user.");

  return normalizeUser(data);
}
