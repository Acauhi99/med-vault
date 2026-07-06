import createClient from "openapi-fetch";

import type { paths } from "@/generated/api";
import {
	clearAuthSession,
	getAuthSession,
} from "@/infrastructure/auth/session-store";
import { config } from "@/infrastructure/config";

async function authFetch(
	input: RequestInfo | URL,
	init?: RequestInit,
): Promise<Response> {
	const session = getAuthSession();
	const hasAuth =
		init?.headers &&
		("Authorization" in init.headers ||
			(init.headers instanceof Headers && init.headers.has("Authorization")));

	if (session.accessToken && !hasAuth) {
		const headers = new Headers(init?.headers);
		headers.set("Authorization", `Bearer ${session.accessToken}`);
		const response = await globalThis.fetch(input, { ...init, headers });
		if (response.status === 401) {
			clearAuthSession();
			window.location.href = "/";
		}
		return response;
	}

	const response = await globalThis.fetch(input, init);
	if (response.status === 401) {
		clearAuthSession();
		window.location.href = "/";
	}
	return response;
}

export const apiClient = createClient<paths>({
	baseUrl: config.apiBaseUrl,
	fetch: authFetch,
});
