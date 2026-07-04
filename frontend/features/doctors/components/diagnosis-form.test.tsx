// @vitest-environment happy-dom

import { screen, waitFor } from "@testing-library/react";
import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";
import { apiBase, renderWith, ts, uuid1, uuid3 } from "../../../test/setup";
import { DiagnosisForm } from "./diagnosis-form";

const server = setupServer();
beforeAll(() => server.listen({ onUnhandledRequest: "bypass" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("DiagnosisForm", () => {
	it("submits diagnosis notes and shows success", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/diagnosis`, async ({ request }) => {
				const body = (await request.json()) as { notes: string };
				expect(body.notes).toBe("Muscle strain, rest recommended");

				return HttpResponse.json({
					id: uuid1,
					case_id: uuid1,
					doctor_id: uuid3,
					notes: "Muscle strain, rest recommended",
					written_at: ts,
				});
			}),
		);

		const { user } = renderWith(<DiagnosisForm caseId={uuid1} />);

		await user.type(
			screen.getByPlaceholderText(/enter your diagnosis/i),
			"Muscle strain, rest recommended",
		);
		await user.click(screen.getByRole("button", { name: /submit diagnosis/i }));

		await waitFor(() => {
			expect(
				screen.getByText(/diagnosis submitted successfully/i),
			).toBeInTheDocument();
		});
	});

	it("shows validation error when notes are empty", async () => {
		const { user } = renderWith(<DiagnosisForm caseId={uuid1} />);

		await user.click(screen.getByRole("button", { name: /submit diagnosis/i }));

		await waitFor(() => {
			expect(screen.getByText(/notes are required/i)).toBeInTheDocument();
		});
	});

	it("shows API error on failure", async () => {
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/diagnosis`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		const { user } = renderWith(<DiagnosisForm caseId={uuid1} />);

		await user.type(
			screen.getByPlaceholderText(/enter your diagnosis/i),
			"Notes",
		);
		await user.click(screen.getByRole("button", { name: /submit diagnosis/i }));

		await waitFor(() => {
			expect(screen.getByText(/Unable to write diagnosis/)).toBeInTheDocument();
		});
	});

	it("disables submit while request is in flight", async () => {
		let resolveRequest = () => {};
		server.use(
			http.post(`${apiBase}/cases/${uuid1}/diagnosis`, () => {
				return new Promise((resolve) => {
					resolveRequest = () =>
						resolve(
							HttpResponse.json({
								id: uuid1,
								case_id: uuid1,
								doctor_id: uuid3,
								notes: "Done",
								written_at: ts,
							}),
						);
				});
			}),
		);

		const { user } = renderWith(<DiagnosisForm caseId={uuid1} />);

		await user.type(
			screen.getByPlaceholderText(/enter your diagnosis/i),
			"Notes",
		);
		await user.click(screen.getByRole("button", { name: /submit diagnosis/i }));

		await waitFor(() => {
			expect(
				screen.getByRole("button", { name: /submitting/i }),
			).toBeDisabled();
		});

		resolveRequest();
	});
});
