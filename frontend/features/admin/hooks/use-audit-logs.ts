import { useQuery } from "@tanstack/react-query";

import { listAuditLogs } from "../services/audit";

export function useAuditLogs(params: {
	page?: number;
	pageSize?: number;
	action?: string;
	userId?: string;
	resourceType?: string;
	resourceId?: string;
}) {
	return useQuery({
		queryKey: ["audit-logs", params],
		queryFn: () => listAuditLogs(params),
	});
}
