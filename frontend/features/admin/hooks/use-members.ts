import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { useAuthSession } from "@/infrastructure/auth/use-auth-session";

import {
	addMember,
	listMembers,
	reactivateTenant,
	removeMember,
} from "../services/members";

export function useMemberList() {
	const session = useAuthSession();
	const tenantId = session.activeTenant?.tenantId ?? null;

	return useQuery({
		queryKey: ["members", tenantId],
		queryFn: () => {
			if (tenantId == null) throw new Error("Missing tenant id");
			return listMembers(tenantId);
		},
		enabled: tenantId != null,
	});
}

export function useAddMember() {
	const session = useAuthSession();
	const queryClient = useQueryClient();
	const tenantId = session.activeTenant?.tenantId ?? "";

	return useMutation({
		mutationFn: (input: {
			userId: string;
			role: "patient" | "doctor" | "administrator";
		}) => addMember(tenantId, input),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["members"] });
		},
	});
}

export function useRemoveMember() {
	const session = useAuthSession();
	const queryClient = useQueryClient();
	const tenantId = session.activeTenant?.tenantId ?? "";

	return useMutation({
		mutationFn: (userId: string) => removeMember(tenantId, userId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["members"] });
		},
	});
}

export function useReactivateTenant() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (tenantId: string) => reactivateTenant({ tenantId }),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["members"] });
		},
	});
}
