// @vitest-environment happy-dom

import { screen, waitFor } from "@testing-library/react";
import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";
import { apiBase, renderWith, ts, uuid1 } from "../../../test/setup";
import { ImageUpload } from "./image-upload";
import { maxImageUploadSizeBytes } from "../schemas/cases";

const server = setupServer();
beforeAll(() => server.listen({ onUnhandledRequest: "bypass" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

function createTestFile(name = "scan.jpg", type = "image/jpeg") {
	return new File(["dummy"], name, { type });
}

describe("ImageUpload", () => {
	it("uploads a file through the 3-step flow", async () => {
		let uploadUrlCalled = false;
		let s3PutCalled = false;
		let confirmCalled = false;

		server.use(
			http.post(
				`${apiBase}/cases/${uuid1}/images/upload-url`,
				async ({ request }) => {
					const body = (await request.json()) as { file_size: number };
					expect(body.file_size).toBeGreaterThan(0);
					uploadUrlCalled = true;
					return HttpResponse.json({
						data: {
							upload_url: "https://s3.example.com/upload",
							s3_key: "cases/1/scan.jpg",
							expires_in: 300,
						},
					});
				},
			),
			http.put("https://s3.example.com/upload", () => {
				s3PutCalled = true;
				return new HttpResponse(null, { status: 200 });
			}),
			http.post(`${apiBase}/cases/${uuid1}/images`, async ({ request }) => {
				const body = (await request.json()) as {
					s3_key: string;
					file_name: string;
					content_type: string;
				};
				expect(body.s3_key).toBe("cases/1/scan.jpg");
				expect(body.file_name).toBe("scan.jpg");
				confirmCalled = true;

				return HttpResponse.json({
					id: uuid1,
					case_id: uuid1,
					file_name: "scan.jpg",
					content_type: "image/jpeg",
					uploaded_at: ts,
				});
			}),
		);

		const { user } = renderWith(<ImageUpload caseId={uuid1} />);

		const file = createTestFile();
		const input = screen.getByLabelText(/upload image file/i);
		await user.upload(input, file);

		await waitFor(() => {
			expect(uploadUrlCalled).toBe(true);
			expect(s3PutCalled).toBe(true);
			expect(confirmCalled).toBe(true);
		});

		await waitFor(() => {
			expect(screen.getByText(/uploaded successfully/i)).toBeInTheDocument();
		});
	});

	it("shows error when file exceeds 50mb", async () => {
		const { user } = renderWith(<ImageUpload caseId={uuid1} />);

		const file = new File(
			[new Uint8Array(maxImageUploadSizeBytes + 1)],
			"big.jpg",
			{
				type: "image/jpeg",
			},
		);
		const input = screen.getByLabelText(/upload image file/i);
		await user.upload(input, file);

		await waitFor(() => {
			expect(screen.getByText(/file is too large/i)).toBeInTheDocument();
		});
	});

	it("shows error when upload URL request fails", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/images/upload-url`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		const { user } = renderWith(<ImageUpload caseId={uuid1} />);

		const file = createTestFile();
		const input = screen.getByLabelText(/upload image file/i);
		await user.upload(input, file);

		await waitFor(() => {
			expect(
				screen.getByText(/unable to request upload url/i),
			).toBeInTheDocument();
		});
	});

	it("resets error state when dismiss is clicked", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/images/upload-url`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		const { user } = renderWith(<ImageUpload caseId={uuid1} />);

		const file = createTestFile();
		const input = screen.getByLabelText(/upload image file/i);
		await user.upload(input, file);

		await waitFor(() => {
			expect(
				screen.getByText(/unable to request upload url/i),
			).toBeInTheDocument();
		});

		await user.click(screen.getByRole("button", { name: /dismiss/i }));

		expect(
			screen.queryByText(/unable to request upload url/i),
		).not.toBeInTheDocument();
	});

	it("disables upload button while uploading", async () => {
		let resolveUpload = () => {};
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/images/upload-url`, () => {
				return HttpResponse.json({
					data: {
						upload_url: "https://s3.example.com/upload",
						s3_key: "cases/1/scan.jpg",
						expires_in: 300,
					},
				});
			}),
			http.put("https://s3.example.com/upload", () => {
				return new Promise((resolve) => {
					resolveUpload = () =>
						resolve(new HttpResponse(null, { status: 200 }));
				});
			}),
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

		const { user } = renderWith(<ImageUpload caseId={uuid1} />);

		const file = createTestFile();
		const input = screen.getByLabelText(/upload image file/i);
		await user.upload(input, file);

		await waitFor(() => {
			expect(screen.getByRole("button", { name: /uploading/i })).toBeDisabled();
		});

		resolveUpload();

		await waitFor(() => {
			expect(screen.getByText(/uploaded successfully/i)).toBeInTheDocument();
		});
	});
});
