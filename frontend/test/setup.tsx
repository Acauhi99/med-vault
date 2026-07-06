import "@testing-library/jest-dom/vitest";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, render } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactNode } from "react";
import { afterEach } from "vitest";

import { updateAuthSession } from "@/infrastructure/auth/session-store";

afterEach(() => {
	cleanup();
});

export const apiBase = "http://localhost:8080/api/v1";

// ── Session helpers ───────────────────────────────────────────────────────────

const defaultSession = {
	accessToken: "test-token",
	refreshToken: "refresh-token",
	tenants: [
		{
			tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
			tenantName: "Test Clinic",
			role: "administrator" as const,
		},
	],
	activeTenant: {
		tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
		tenantName: "Test Clinic",
		role: "administrator" as const,
	},
	user: null,
};

export function seedSession(patch?: Record<string, unknown>) {
	updateAuthSession({ ...defaultSession, ...patch });
}

// ── Render helper ─────────────────────────────────────────────────────────────

function createQueryClient() {
	return new QueryClient({
		defaultOptions: {
			queries: { retry: false, gcTime: 0 },
			mutations: { retry: false },
		},
	});
}

export function renderWith(ui: ReactNode) {
	const queryClient = createQueryClient();
	seedSession();

	return {
		...render(
			<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>,
		),
		user: userEvent.setup(),
	};
}

// ── Shared test data ──────────────────────────────────────────────────────────

export const uuid1 = "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001";
export const uuid2 = "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001";
export const uuid3 = "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001";
export const ts = "2026-07-04T15:00:00Z";
