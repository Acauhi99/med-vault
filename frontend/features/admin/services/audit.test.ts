import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import { listAuditLogs } from "./audit";

const apiBase = "http://localhost:8080";
const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

const uuid1 = "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001";
const uuid2 = "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001";
const uuid3 = "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001";
const ts = "2026-07-04T15:00:00Z";

describe("admin audit service", () => {
	it("lists audit logs", async () => {
		server.use(
			http.get(`${apiBase}/audit-logs`, () => {
				return HttpResponse.json({
					data: [
						{
							id: uuid1,
							tenant_id: uuid2,
							user_id: uuid3,
							action: "case.create",
							resource_type: "case",
							resource_id: uuid1,
							ip_address: "127.0.0.1",
							metadata: { key: "value" },
							created_at: ts,
						},
					],
					meta: { page: 1, page_size: 20, total: 1 },
				});
			}),
		);

		const result = await listAuditLogs({ page: 1, pageSize: 20 });
		expect(result.data).toHaveLength(1);
		expect(result.data[0].action).toBe("case.create");
		expect(result.data[0].resourceType).toBe("case");
		expect(result.data[0].ipAddress).toBe("127.0.0.1");
		expect(result.data[0].metadata).toEqual({ key: "value" });
	});

	it("passes filters as query params", async () => {
		let capturedUrl = "";
		server.use(
			http.get(`${apiBase}/audit-logs`, ({ request }) => {
				capturedUrl = request.url;
				return HttpResponse.json({ data: [], meta: { total: 0 } });
			}),
		);

		await listAuditLogs({
			page: 2,
			pageSize: 10,
			resourceType: "case",
			resourceId: uuid1,
		});

		const url = new URL(capturedUrl);
		expect(url.searchParams.get("page")).toBe("2");
		expect(url.searchParams.get("page_size")).toBe("10");
		expect(url.searchParams.get("resource_type")).toBe("case");
		expect(url.searchParams.get("resource_id")).toBe(uuid1);
	});

	it("handles null metadata", async () => {
		server.use(
			http.get(`${apiBase}/audit-logs`, () => {
				return HttpResponse.json({
					data: [
						{
							id: uuid1,
							tenant_id: uuid2,
							user_id: uuid3,
							action: "user.login",
							resource_type: "user",
							resource_id: uuid3,
							ip_address: null,
							metadata: null,
							created_at: ts,
						},
					],
					meta: { total: 1 },
				});
			}),
		);

		const result = await listAuditLogs({});
		expect(result.data[0].ipAddress).toBeNull();
		expect(result.data[0].metadata).toBeNull();
	});

	it("throws on API error", async () => {
		server.use(
			http.get(`${apiBase}/audit-logs`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		await expect(listAuditLogs({})).rejects.toThrow(
			"Unable to load audit logs.",
		);
	});
});
