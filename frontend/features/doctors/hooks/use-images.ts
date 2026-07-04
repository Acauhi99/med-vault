"use client";

import { useQuery } from "@tanstack/react-query";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import { getImageDownloadUrl, listCaseImages } from "../services/cases";

export function useCaseImages(caseId: string | null) {
	const session = useAuthSession();

	return useQuery({
		queryKey: ["cases", caseId, "images", session.activeTenant?.tenantId],
		queryFn: () => {
			if (!caseId) throw new Error("Missing case id");
			return listCaseImages(caseId);
		},
		enabled: Boolean(session.accessToken && session.activeTenant && caseId),
	});
}

export function useImageDownloadUrl(imageId: string | null) {
	const session = useAuthSession();

	return useQuery({
		queryKey: [
			"images",
			imageId,
			"download-url",
			session.activeTenant?.tenantId,
		],
		queryFn: () => {
			if (!imageId) throw new Error("Missing image id");
			return getImageDownloadUrl(imageId);
		},
		enabled: Boolean(session.accessToken && session.activeTenant && imageId),
		staleTime: 4 * 60 * 1000,
	});
}
