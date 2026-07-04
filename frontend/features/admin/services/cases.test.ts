import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import { assignDoctor, closeCase, getCase, listCases } from "./cases";

const apiBase = "http://localhost:8080";
const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

const uuid1 = "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001";
const uuid2 = "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001";
const uuid3 = "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001";
const ts = "2026-07-04T15:00:00Z";

describe("admin cases service", () => {
	it("lists cases", async () => {
		server.use(
			http.get(`${apiBase}/cases`, () => {
				return HttpResponse.json({
					data: [
						{
							id: uuid1,
							patient_id: uuid2,
							doctor_id: null,
							status: "open",
							created_at: ts,
						},
					],
					meta: { page: 1, page_size: 20, total: 1 },
				});
			}),
		);

		const result = await listCases({ page: 1, pageSize: 20 });
		expect(result.data).toHaveLength(1);
		expect(result.data[0].status).toBe("open");
	});

	it("gets case detail", async () => {
		server.use(
			http.get(`${apiBase}/cases/${uuid1}`, () => {
				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid3,
						patient_id: uuid2,
						doctor_id: uuid3,
						status: "assigned",
						symptoms: [],
						diagnosis: null,
						created_at: ts,
						updated_at: ts,
						closed_at: null,
					},
				});
			}),
		);

		const result = await getCase(uuid1);
		expect(result.status).toBe("assigned");
	});

	it("assigns a doctor to a case", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/assign`, async ({ request }) => {
				const body = (await request.json()) as { doctor_id: string };
				expect(body.doctor_id).toBe(uuid3);

				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid3,
						patient_id: uuid2,
						doctor_id: uuid3,
						status: "assigned",
						symptoms: [],
						diagnosis: null,
						created_at: ts,
						updated_at: ts,
						closed_at: null,
					},
				});
			}),
		);

		const result = await assignDoctor(uuid1, { doctorId: uuid3 });
		expect(result.doctorId).toBe(uuid3);
		expect(result.status).toBe("assigned");
	});

	it("closes a case", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/close`, () => {
				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid3,
						patient_id: uuid2,
						doctor_id: uuid3,
						status: "closed",
						symptoms: [],
						diagnosis: null,
						created_at: ts,
						updated_at: ts,
						closed_at: ts,
					},
				});
			}),
		);

		const result = await closeCase(uuid1);
		expect(result.status).toBe("closed");
		expect(result.closedAt).toBe(ts);
	});

	it("throws on API error", async () => {
		server.use(
			http.get(`${apiBase}/cases`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		await expect(listCases({})).rejects.toThrow("Unable to load cases.");
	});
});
