import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import {
	getCaseDetail,
	getImageDownloadUrl,
	listAssignedCases,
	listCaseImages,
	writeDiagnosis,
} from "./cases";

const apiBase = "http://localhost:8080";
const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

const uuid1 = "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001";
const uuid2 = "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001";
const uuid3 = "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001";
const ts = "2026-07-04T15:00:00Z";

describe("doctor cases service", () => {
	it("lists assigned cases", async () => {
		server.use(
			http.get(`${apiBase}/cases`, () => {
				return HttpResponse.json({
					data: [
						{
							id: uuid1,
							patient_id: uuid2,
							doctor_id: uuid3,
							status: "assigned",
							created_at: ts,
						},
					],
					meta: { page: 1, page_size: 20, total: 1 },
				});
			}),
		);

		const result = await listAssignedCases(1, 20);
		expect(result.cases).toHaveLength(1);
		expect(result.cases[0].doctorId).toBe(uuid3);
		expect(result.meta.total).toBe(1);
	});

	it("gets case detail with symptoms and diagnosis", async () => {
		server.use(
			http.get(`${apiBase}/cases/${uuid1}`, () => {
				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid3,
						patient_id: uuid2,
						doctor_id: uuid3,
						status: "diagnosed",
						symptoms: [
							{
								id: uuid1,
								description: "Chest pain",
								severity: "high",
								reported_at: ts,
							},
						],
						diagnosis: {
							id: uuid1,
							case_id: uuid1,
							doctor_id: uuid3,
							notes: "Muscle strain",
							written_at: ts,
						},
						created_at: ts,
						updated_at: ts,
						closed_at: null,
					},
				});
			}),
		);

		const result = await getCaseDetail(uuid1);
		expect(result.status).toBe("diagnosed");
		expect(result.symptoms).toHaveLength(1);
		expect(result.diagnosis?.notes).toBe("Muscle strain");
	});

	it("lists case images", async () => {
		server.use(
			http.get(`${apiBase}/cases/${uuid1}/images`, () => {
				return HttpResponse.json({
					data: [
						{
							id: uuid1,
							case_id: uuid1,
							file_name: "xray.png",
							content_type: "image/png",
							uploaded_at: ts,
						},
					],
				});
			}),
		);

		const result = await listCaseImages(uuid1);
		expect(result).toHaveLength(1);
		expect(result[0].fileName).toBe("xray.png");
	});

	it("gets image download URL", async () => {
		server.use(
			http.get(`${apiBase}/images/${uuid1}/download-url`, () => {
				return HttpResponse.json({
					data: {
						download_url: "https://s3.example.com/dl",
						expires_in: 600,
					},
				});
			}),
		);

		const result = await getImageDownloadUrl(uuid1);
		expect(result.downloadUrl).toBe("https://s3.example.com/dl");
		expect(result.expiresIn).toBe(600);
	});

	it("writes a diagnosis", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/diagnosis`, async ({ request }) => {
				const body = (await request.json()) as { notes: string };
				expect(body.notes).toBe("Benign condition");

				return HttpResponse.json({
					id: uuid1,
					case_id: uuid1,
					doctor_id: uuid3,
					notes: "Benign condition",
					written_at: ts,
				});
			}),
		);

		const result = await writeDiagnosis(uuid1, { notes: "Benign condition" });
		expect(result.notes).toBe("Benign condition");
		expect(result.caseId).toBe(uuid1);
	});

	it("throws on API error", async () => {
		server.use(
			http.get(`${apiBase}/cases`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		await expect(listAssignedCases()).rejects.toThrow("Unable to load cases.");
	});
});
