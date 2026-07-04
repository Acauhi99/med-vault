import { z } from "zod";

const roleSchema = z.enum(["patient", "doctor", "administrator"]);
const statusSchema = z.enum(["active", "inactive"]);

export const loginInputSchema = z.object({
  email: z.email("Enter a valid email").trim(),
  password: z.string().min(8, "Use at least 8 characters"),
});

export const registerInputSchema = loginInputSchema;

export const tenantSelectionSchema = z.object({
  tenantId: z.uuid("Select a tenant"),
});

export const tenantSummarySchema = z.object({
  tenantId: z.uuid(),
  tenantName: z.string().min(1),
  role: roleSchema,
});

export const loginResponseSchema = z.object({
  accessToken: z.string().min(1),
  tenants: z.array(tenantSummarySchema),
});

export const tokenResponseSchema = z.object({
  accessToken: z.string().min(1),
  refreshToken: z.string().min(1),
  expiresIn: z.number().int().positive(),
});

export const registrationResultSchema = z.object({
  id: z.uuid(),
  email: z.email(),
  status: statusSchema,
  createdAt: z.iso.datetime(),
});

export const userProfileSchema = z.object({
  id: z.uuid(),
  tenantId: z.uuid(),
  email: z.email(),
  role: roleSchema,
  status: statusSchema,
  createdAt: z.iso.datetime(),
});

export type LoginInput = z.infer<typeof loginInputSchema>;
export type RegisterInput = z.infer<typeof registerInputSchema>;
export type TenantSelectionInput = z.infer<typeof tenantSelectionSchema>;
export type TenantSummary = z.infer<typeof tenantSummarySchema>;
export type UserProfile = z.infer<typeof userProfileSchema>;

export type AuthSession = {
  accessToken: string | null;
  refreshToken: string | null;
  tenants: TenantSummary[];
  activeTenant: TenantSummary | null;
  user: UserProfile | null;
};

export type LoginResult = z.infer<typeof loginResponseSchema>;
export type TokenResult = z.infer<typeof tokenResponseSchema>;
export type RegistrationResult = z.infer<typeof registrationResultSchema>;
