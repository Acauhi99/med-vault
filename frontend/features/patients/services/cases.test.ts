import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import {
	addSymptom,
	confirmUpload,
	createCase,
	getCase,
	getDownloadURL,
	listCases,
	listImages,
	requestUploadURL,
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

describe("patient cases service", () => {
	it("lists cases and normalizes snake_case to camelCase", async () => {
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
		expect(result.cases).toHaveLength(1);
		expect(result.cases[0]).toEqual({
			id: uuid1,
			patientId: uuid2,
			doctorId: null,
			status: "open",
			createdAt: ts,
		});
		expect(result.total).toBe(1);
	});

	it("gets a single case with symptoms and diagnosis", async () => {
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
								description: "Headache",
								severity: "medium",
								reported_at: ts,
							},
						],
						diagnosis: {
							id: uuid1,
							case_id: uuid1,
							doctor_id: uuid3,
							notes: "Tension headache",
							written_at: ts,
						},
						created_at: ts,
						updated_at: ts,
						closed_at: null,
					},
				});
			}),
		);

		const result = await getCase(uuid1);
		expect(result.id).toBe(uuid1);
		expect(result.status).toBe("diagnosed");
		expect(result.symptoms).toHaveLength(1);
		expect(result.diagnosis).not.toBeNull();
		expect(result.diagnosis?.notes).toBe("Tension headache");
	});

	it("creates a case with symptoms", async () => {
		server.use(
			http.post(`${apiBase}/cases`, async ({ request }) => {
				const body = (await request.json()) as { symptoms: unknown[] };
				expect(body.symptoms).toHaveLength(1);

				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid3,
						patient_id: uuid2,
						doctor_id: null,
						status: "open",
						symptoms: [
							{
								id: uuid1,
								description: "Fever",
								severity: "high",
								reported_at: ts,
							},
						],
						diagnosis: null,
						created_at: ts,
						updated_at: ts,
						closed_at: null,
					},
				});
			}),
		);

		const result = await createCase([
			{ description: "Fever", severity: "high" },
		]);
		expect(result.status).toBe("open");
		expect(result.symptoms).toHaveLength(1);
	});

	it("adds a symptom to an existing case", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/symptoms`, async ({ request }) => {
				const body = (await request.json()) as { description: string };
				expect(body.description).toBe("Nausea");

				return HttpResponse.json({
					id: uuid2,
					description: "Nausea",
					severity: "low",
					reported_at: ts,
				});
			}),
		);

		const result = await addSymptom(uuid1, {
			description: "Nausea",
			severity: "low",
		});
		expect(result.description).toBe("Nausea");
		expect(result.severity).toBe("low");
	});

	it("requests an upload URL", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/images/upload-url`, () => {
				return HttpResponse.json({
					data: {
						upload_url: "https://s3.example.com/upload",
						s3_key: "cases/1/img.jpg",
						expires_in: 300,
					},
				});
			}),
		);

		const result = await requestUploadURL(uuid1, "scan.jpg", "image/jpeg");
		expect(result.uploadUrl).toBe("https://s3.example.com/upload");
		expect(result.s3Key).toBe("cases/1/img.jpg");
		expect(result.expiresIn).toBe(300);
	});

	it("confirms an upload", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/images`, () => {
				return HttpResponse.json({
					id: uuid1,
					case_id: uuid1,
					file_name: "scan.jpg",
					content_type: "image/jpeg",
					uploaded_at: ts,
				});
			}),
		);

		const result = await confirmUpload(uuid1, {
			s3Key: "cases/1/img.jpg",
			fileName: "scan.jpg",
			contentType: "image/jpeg",
		});
		expect(result.fileName).toBe("scan.jpg");
	});

	it("lists images for a case", async () => {
		server.use(
			http.get(`${apiBase}/cases/${uuid1}/images`, () => {
				return HttpResponse.json({
					data: [
						{
							id: uuid1,
							case_id: uuid1,
							file_name: "scan.jpg",
							content_type: "image/jpeg",
							uploaded_at: ts,
						},
					],
				});
			}),
		);

		const result = await listImages(uuid1);
		expect(result).toHaveLength(1);
		expect(result[0].fileName).toBe("scan.jpg");
	});

	it("gets a download URL", async () => {
		server.use(
			http.get(`${apiBase}/images/${uuid1}/download-url`, () => {
				return HttpResponse.json({
					data: {
						download_url: "https://s3.example.com/download",
						expires_in: 300,
					},
				});
			}),
		);

		const result = await getDownloadURL(uuid1);
		expect(result.downloadUrl).toBe("https://s3.example.com/download");
		expect(result.expiresIn).toBe(300);
	});

	it("throws on API error", async () => {
		server.use(
			http.get(`${apiBase}/cases`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		await expect(listCases()).rejects.toThrow("Unable to load cases.");
	});
});
