// @vitest-environment happy-dom

import { screen, waitFor } from "@testing-library/react";
import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";
import {
	apiBase,
	renderWith,
	ts,
	uuid1,
	uuid2,
	uuid3,
} from "../../../test/setup";
import { AssignDoctorDialog } from "./assign-doctor-dialog";

const server = setupServer();
beforeAll(() => server.listen({ onUnhandledRequest: "bypass" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("AssignDoctorDialog", () => {
	it("does not render when closed", () => {
		renderWith(
			<AssignDoctorDialog open={false} onClose={() => {}} caseId={uuid1} />,
		);

		expect(screen.queryByText("Assign Doctor")).not.toBeInTheDocument();
	});

	it("renders form when open and assigns doctor", async () => {
		let assigned = false;
		let closed = false;

		server.use(
			http.post(`${apiBase}/cases/${uuid1}/assign`, async ({ request }) => {
				const body = (await request.json()) as { doctor_id: string };
				expect(body.doctor_id).toBe(uuid3);
				assigned = true;

				return HttpResponse.json({
					data: {
						id: uuid1,
						tenant_id: uuid2,
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

		const { user } = renderWith(
			<AssignDoctorDialog
				open={true}
				onClose={() => {
					closed = true;
				}}
				caseId={uuid1}
			/>,
		);

		expect(screen.getByText("Assign Doctor")).toBeInTheDocument();

		await user.type(screen.getByPlaceholderText(/enter doctor uuid/i), uuid3);
		await user.click(screen.getByRole("button", { name: /assign$/i }));

		await waitFor(() => {
			expect(assigned).toBe(true);
			expect(closed).toBe(true);
		});
	});

	it("closes dialog when cancel is clicked", async () => {
		let closed = false;

		const { user } = renderWith(
			<AssignDoctorDialog
				open={true}
				onClose={() => {
					closed = true;
				}}
				caseId={uuid1}
			/>,
		);

		await user.click(screen.getByRole("button", { name: /cancel/i }));
		expect(closed).toBe(true);
	});

	it("closes dialog when backdrop is clicked", async () => {
		let closed = false;

		const { user } = renderWith(
			<AssignDoctorDialog
				open={true}
				onClose={() => {
					closed = true;
				}}
				caseId={uuid1}
			/>,
		);

		await user.click(screen.getByLabelText("Close dialog"));
		expect(closed).toBe(true);
	});

	it("disables assign button when doctor ID is empty", () => {
		renderWith(
			<AssignDoctorDialog open={true} onClose={() => {}} caseId={uuid1} />,
		);

		expect(screen.getByRole("button", { name: /assign$/i })).toBeDisabled();
	});
});
