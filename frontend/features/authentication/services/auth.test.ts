import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import {
	getCurrentUser,
	login,
	refreshSession,
	register,
	selectTenant,
} from "./auth";

const apiBase = "http://localhost:8080";

const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("auth service", () => {
	it("logs in and normalizes tenants", async () => {
		server.use(
			http.post(`${apiBase}/auth/login`, async ({ request }) => {
				expect(await request.json()).toEqual({
					email: "doctor@example.com",
					password: "password123",
				});

				return HttpResponse.json({
					data: {
						access_token: "login-token",
						tenants: [
							{
								tenant_id: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
								tenant_name: "North Clinic",
								role: "doctor",
							},
						],
					},
				});
			}),
		);

		await expect(
			login({ email: "doctor@example.com", password: "password123" }),
		).resolves.toEqual({
			accessToken: "login-token",
			tenants: [
				{
					tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
					tenantName: "North Clinic",
					role: "doctor",
				},
			],
		});
	});

	it("registers a user", async () => {
		server.use(
			http.post(`${apiBase}/auth/register`, async ({ request }) => {
				expect(await request.json()).toEqual({
					email: "new@example.com",
					password: "password123",
				});

				return HttpResponse.json({
					data: {
						id: "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001",
						email: "new@example.com",
						status: "active",
						created_at: "2026-07-04T15:00:00Z",
					},
				});
			}),
		);

		await expect(
			register({ email: "new@example.com", password: "password123" }),
		).resolves.toEqual({
			id: "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001",
			email: "new@example.com",
			status: "active",
			createdAt: "2026-07-04T15:00:00Z",
		});
	});

	it("selects a tenant with bearer auth", async () => {
		let authorization = "";

		server.use(
			http.post(`${apiBase}/auth/select-tenant`, async ({ request }) => {
				authorization = request.headers.get("authorization") ?? "";

				expect(await request.json()).toEqual({
					tenant_id: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
				});

				return HttpResponse.json({
					data: {
						access_token: "tenant-token",
						refresh_token: "refresh-token",
						expires_in: 3600,
					},
				});
			}),
		);

		await expect(
			selectTenant({
				accessToken: "login-token",
				tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
			}),
		).resolves.toEqual({
			accessToken: "tenant-token",
			refreshToken: "refresh-token",
			expiresIn: 3600,
		});

		expect(authorization).toBe("Bearer login-token");
	});

	it("loads the current user", async () => {
		let authorization = "";

		server.use(
			http.get(`${apiBase}/users/me`, ({ request }) => {
				authorization = request.headers.get("authorization") ?? "";

				return HttpResponse.json({
					data: {
						id: "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001",
						tenant_id: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
						email: "doctor@example.com",
						role: "doctor",
						status: "active",
						created_at: "2026-07-04T15:00:00Z",
					},
				});
			}),
		);

		await expect(getCurrentUser("tenant-token")).resolves.toEqual({
			id: "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001",
			tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
			email: "doctor@example.com",
			role: "doctor",
			status: "active",
			createdAt: "2026-07-04T15:00:00Z",
		});

		expect(authorization).toBe("Bearer tenant-token");
	});

	it("refreshes the session", async () => {
		server.use(
			http.post(`${apiBase}/auth/refresh`, async ({ request }) => {
				expect(await request.json()).toEqual({
					refresh_token: "refresh-token",
				});

				return HttpResponse.json({
					data: {
						access_token: "refreshed-token",
						refresh_token: "refresh-token-2",
						expires_in: 7200,
					},
				});
			}),
		);

		await expect(refreshSession("refresh-token")).resolves.toEqual({
			accessToken: "refreshed-token",
			refreshToken: "refresh-token-2",
			expiresIn: 7200,
		});
	});
});
