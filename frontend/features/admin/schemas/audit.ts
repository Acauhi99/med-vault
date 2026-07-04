import { z } from "zod";

export const auditLogSchema = z.object({
	id: z.uuid(),
	tenantId: z.uuid(),
	userId: z.uuid(),
	action: z.string(),
	resourceType: z.string(),
	resourceId: z.uuid(),
	ipAddress: z.string().nullish(),
	metadata: z.record(z.string(), z.unknown()).nullable(),
	createdAt: z.iso.datetime(),
});

export const auditLogFiltersSchema = z.object({
	resourceType: z.string().optional(),
	resourceId: z.string().optional(),
	page: z.number().optional(),
	pageSize: z.number().optional(),
});

export type AuditLog = z.infer<typeof auditLogSchema>;
export type AuditLogFilters = z.infer<typeof auditLogFiltersSchema>;
