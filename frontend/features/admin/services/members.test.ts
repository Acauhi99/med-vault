import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import { addMember, listMembers, removeMember } from "./members";

const apiBase = "http://localhost:8080";
const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

const tenantId = "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001";
const userId = "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001";

describe("admin members service", () => {
	it("lists tenant members", async () => {
		server.use(
			http.get(`${apiBase}/tenants/${tenantId}/members`, () => {
				return HttpResponse.json({
					data: [
						{
							user_id: userId,
							tenant_id: tenantId,
							role: "doctor",
							name: "Dr. Smith",
						},
					],
				});
			}),
		);

		const result = await listMembers(tenantId);
		expect(result).toHaveLength(1);
		expect(result[0].role).toBe("doctor");
		expect(result[0].name).toBe("Dr. Smith");
	});

	it("adds a member to a tenant", async () => {
		server.use(
			http.post(
				`${apiBase}/tenants/${tenantId}/members`,
				async ({ request }) => {
					const body = (await request.json()) as {
						user_id: string;
						role: string;
					};
					expect(body.user_id).toBe(userId);
					expect(body.role).toBe("patient");

					return HttpResponse.json({
						data: {
							user_id: userId,
							tenant_id: tenantId,
							role: "patient",
							name: "Jane Doe",
						},
					});
				},
			),
		);

		const result = await addMember(tenantId, {
			userId,
			role: "patient",
		});
		expect(result.name).toBe("Jane Doe");
		expect(result.role).toBe("patient");
	});

	it("removes a member from a tenant", async () => {
		let called = false;
		server.use(
			http.delete(`${apiBase}/tenants/${tenantId}/members/${userId}`, () => {
				called = true;
				return HttpResponse.json({ data: {} });
			}),
		);

		await removeMember(tenantId, userId);
		expect(called).toBe(true);
	});

	it("throws on API error", async () => {
		server.use(
			http.get(`${apiBase}/tenants/${tenantId}/members`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		await expect(listMembers(tenantId)).rejects.toThrow(
			"Unable to load members.",
		);
	});
});
