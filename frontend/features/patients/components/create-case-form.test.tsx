// @vitest-environment happy-dom

import { screen, waitFor } from "@testing-library/react";
import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";
import { apiBase, renderWith, ts, uuid1, uuid2 } from "../../../test/setup";
import { CreateCaseForm } from "./create-case-form";

const server = setupServer();
beforeAll(() => server.listen({ onUnhandledRequest: "bypass" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("CreateCaseForm", () => {
	it("creates a case with one symptom and calls onCreated", async () => {
		let created = false;
		const onCreated = () => {
			created = true;
		};
		const onCancel = () => {};

		server.use(
			http.post(`${apiBase}/cases`, async ({ request }) => {
				const body = (await request.json()) as { symptoms: unknown[] };
				expect(body.symptoms).toHaveLength(1);
				expect(body.symptoms[0]).toEqual({
					description: "Persistent headache",
					severity: "high",
				});

				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid2,
						patient_id: uuid2,
						doctor_id: null,
						status: "open",
						symptoms: [],
						diagnosis: null,
						created_at: ts,
						updated_at: ts,
						closed_at: null,
					},
				});
			}),
		);

		const { user } = renderWith(
			<CreateCaseForm onCreated={onCreated} onCancel={onCancel} />,
		);

		await user.type(
			screen.getByPlaceholderText("Describe the symptom"),
			"Persistent headache",
		);
		await user.selectOptions(screen.getByDisplayValue("Low"), "high");
		await user.click(screen.getByRole("button", { name: /create case/i }));

		await waitFor(() => {
			expect(created).toBe(true);
		});
	});

	it("shows validation error when description is empty", async () => {
		const { user } = renderWith(
			<CreateCaseForm onCreated={() => {}} onCancel={() => {}} />,
		);

		await user.click(screen.getByRole("button", { name: /create case/i }));

		await waitFor(() => {
			expect(screen.getByRole("status")).toHaveTextContent(/required/i);
		});
	});

	it("adds and removes symptom rows", async () => {
		const { user } = renderWith(
			<CreateCaseForm onCreated={() => {}} onCancel={() => {}} />,
		);

		expect(screen.getAllByText(/Symptom/)).toHaveLength(1);

		await user.click(screen.getByRole("button", { name: /\+ add symptom/i }));
		expect(screen.getAllByText(/Symptom/)).toHaveLength(2);

		await user.click(screen.getAllByText("Remove")[0]);
		expect(screen.getAllByText(/Symptom/)).toHaveLength(1);
	});

	it("calls onCancel when cancel is clicked", async () => {
		let cancelled = false;
		const { user } = renderWith(
			<CreateCaseForm
				onCreated={() => {}}
				onCancel={() => {
					cancelled = true;
				}}
			/>,
		);

		await user.click(screen.getByRole("button", { name: /cancel/i }));
		expect(cancelled).toBe(true);
	});

	it("shows API error message on failure", async () => {
		server.use(
			http.post(`${apiBase}/cases`, () => {
				return HttpResponse.json({ error: "forbidden" }, { status: 403 });
			}),
		);

		const { user } = renderWith(
			<CreateCaseForm onCreated={() => {}} onCancel={() => {}} />,
		);

		await user.type(
			screen.getByPlaceholderText("Describe the symptom"),
			"Fever",
		);
		await user.click(screen.getByRole("button", { name: /create case/i }));

		await waitFor(() => {
			expect(screen.getByRole("status")).toHaveTextContent(
				/Unable to create case/,
			);
		});
	});
});
