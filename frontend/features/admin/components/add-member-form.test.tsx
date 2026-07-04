// @vitest-environment happy-dom

import { screen, waitFor } from "@testing-library/react";
import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";
import { apiBase, renderWith, uuid1, uuid3 } from "../../../test/setup";
import { AddMemberForm } from "./add-member-form";

const server = setupServer();
beforeAll(() => server.listen({ onUnhandledRequest: "bypass" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("AddMemberForm", () => {
	it("opens dialog, fills form, and adds member", async () => {
		let added = false;

		server.use(
			http.post(`${apiBase}/tenants/${uuid1}/members`, async ({ request }) => {
				const body = (await request.json()) as {
					user_id: string;
					role: string;
				};
				expect(body.user_id).toBe(uuid3);
				expect(body.role).toBe("doctor");
				added = true;

				return HttpResponse.json({
					data: {
						user_id: uuid3,
						tenant_id: uuid1,
						role: "doctor",
						name: "Dr. New",
					},
				});
			}),
		);

		const { user } = renderWith(<AddMemberForm />);

		await user.click(screen.getByRole("button", { name: /add member/i }));
		expect(screen.getByText("Add Tenant Member")).toBeInTheDocument();

		await user.type(screen.getByPlaceholderText(/enter user uuid/i), uuid3);
		await user.selectOptions(screen.getByDisplayValue("patient"), "doctor");
		await user.click(screen.getByRole("button", { name: /^add$/i }));

		await waitFor(() => {
			expect(added).toBe(true);
		});

		await waitFor(() => {
			expect(screen.queryByText("Add Tenant Member")).not.toBeInTheDocument();
		});
	});

	it("closes dialog on cancel", async () => {
		const { user } = renderWith(<AddMemberForm />);

		await user.click(screen.getByRole("button", { name: /add member/i }));
		expect(screen.getByText("Add Tenant Member")).toBeInTheDocument();

		await user.click(screen.getByRole("button", { name: /cancel/i }));

		await waitFor(() => {
			expect(screen.queryByText("Add Tenant Member")).not.toBeInTheDocument();
		});
	});

	it("disables add button when user ID is empty", async () => {
		const { user } = renderWith(<AddMemberForm />);

		await user.click(screen.getByRole("button", { name: /add member/i }));

		expect(screen.getByRole("button", { name: /^add$/i })).toBeDisabled();
	});

	it("shows API error on failure", async () => {
		server.use(
			http.post(`${apiBase}/tenants/${uuid1}/members`, () => {
				return HttpResponse.json({ error: "user not found" }, { status: 404 });
			}),
		);

		const { user } = renderWith(<AddMemberForm />);

		await user.click(screen.getByRole("button", { name: /add member/i }));
		await user.type(screen.getByPlaceholderText(/enter user uuid/i), uuid3);
		await user.click(screen.getByRole("button", { name: /^add$/i }));

		await waitFor(() => {
			expect(screen.getByText(/Unable to add member/)).toBeInTheDocument();
		});
	});
});
