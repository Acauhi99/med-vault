import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { assignDoctor, closeCase, getCase, listCases } from "../services/cases";

export function useCaseList(params: { page?: number; pageSize?: number }) {
	return useQuery({
		queryKey: ["cases", params],
		queryFn: () => listCases(params),
	});
}

export function useCaseDetail(id: string | null) {
	return useQuery({
		queryKey: ["cases", id],
		queryFn: () => {
			if (id == null) throw new Error("Missing case id");
			return getCase(id);
		},
		enabled: id != null,
	});
}

export function useAssignDoctor(caseId: string) {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (input: { doctorId: string }) => assignDoctor(caseId, input),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["cases"] });
		},
	});
}

export function useCloseCase(caseId: string) {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: () => closeCase(caseId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["cases"] });
		},
	});
}
