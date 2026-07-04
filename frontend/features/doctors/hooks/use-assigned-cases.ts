"use client";

import { useQuery } from "@tanstack/react-query";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import { getCaseDetail, listAssignedCases } from "../services/cases";

export function useAssignedCases(page = 1, pageSize = 20) {
	const session = useAuthSession();

	return useQuery({
		queryKey: [
			"cases",
			"assigned",
			session.activeTenant?.tenantId,
			page,
			pageSize,
		],
		queryFn: () => listAssignedCases(page, pageSize),
		enabled: Boolean(session.accessToken && session.activeTenant),
	});
}

export function useCaseDetail(caseId: string | null) {
	const session = useAuthSession();

	return useQuery({
		queryKey: ["cases", caseId, session.activeTenant?.tenantId],
		queryFn: () => {
			if (!caseId) throw new Error("Missing case id");
			return getCaseDetail(caseId);
		},
		enabled: Boolean(session.accessToken && session.activeTenant && caseId),
	});
}
