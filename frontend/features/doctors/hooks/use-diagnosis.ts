"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import type { WriteDiagnosisInput } from "../schemas/cases";
import { writeDiagnosis } from "../services/cases";

export function useWriteDiagnosis(caseId: string) {
	const session = useAuthSession();
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (input: WriteDiagnosisInput) => writeDiagnosis(caseId, input),
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ["cases", caseId, session.activeTenant?.tenantId],
			});
		},
	});
}
