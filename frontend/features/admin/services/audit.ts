import type { components } from "@/generated/api";

import { apiClient } from "@/infrastructure/api/client";

import { auditLogSchema } from "../schemas/audit";

type AuditLogRaw = components["schemas"]["AuditLogResponse"];

function requireData<T>(value: T | null | undefined, message: string): T {
	if (value == null) {
		throw new Error(message);
	}

	return value;
}

function normalizeAuditLog(raw: AuditLogRaw) {
	return auditLogSchema.parse({
		id: requireData(raw.id, "Missing log id"),
		tenantId: requireData(raw.tenant_id, "Missing tenant id"),
		userId: requireData(raw.user_id, "Missing user id"),
		action: requireData(raw.action, "Missing action"),
		resourceType: requireData(raw.resource_type, "Missing resource type"),
		resourceId: requireData(raw.resource_id, "Missing resource id"),
		ipAddress: raw.ip_address,
		metadata: raw.metadata ?? null,
		createdAt: requireData(raw.created_at, "Missing created at"),
	});
}

export async function listAuditLogs(params: {
	page?: number;
	pageSize?: number;
	resourceType?: string;
	resourceId?: string;
}) {
	const response = await apiClient.GET("/audit-logs", {
		params: {
			query: {
				page: params.page,
				page_size: params.pageSize,
				resource_type: params.resourceType,
				resource_id: params.resourceId,
			},
		},
	});

	if (response.error) {
		throw new Error("Unable to load audit logs.");
	}

	const data = requireData(response.data, "Unable to load audit logs.");

	return {
		data: (data.data ?? []).map(normalizeAuditLog),
		meta: data.meta,
	};
}
