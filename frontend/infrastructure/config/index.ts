function normalizeApiBaseUrl(baseUrl: string) {
	const trimmed = baseUrl.replace(/\/$/, "");
	return trimmed.endsWith("/api/v1") ? trimmed : `${trimmed}/api/v1`;
}

export function getApiBaseUrl() {
	return normalizeApiBaseUrl(
		process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
	);
}

export const config = {
	apiBaseUrl: getApiBaseUrl(),
} as const;
