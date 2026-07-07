// @vitest-environment happy-dom

import "../../../test/setup";

import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

const mockUseReactivateTenant = vi.fn();

vi.mock("../hooks/use-members", () => ({
	useReactivateTenant: () => mockUseReactivateTenant(),
}));

import { ReactivateTenantForm } from "./reactivate-tenant-form";

beforeEach(() => {
	mockUseReactivateTenant.mockReset();
});

describe("ReactivateTenantForm", () => {
	it("opens, submits, and closes the dialog", async () => {
		const user = userEvent.setup();
		const tenantUuid = "550e8400-e29b-41d4-a716-446655440000";
		const mutate = vi.fn(
			(_tenantId: string, options?: { onSuccess?: () => void }) => {
				options?.onSuccess?.();
			},
		);

		mockUseReactivateTenant.mockReturnValue({
			error: null,
			isPending: false,
			mutate,
		});

		render(<ReactivateTenantForm />);

		await user.click(
			screen.getByRole("button", { name: /reactivate tenant/i }),
		);
		expect(
			screen.getByRole("heading", { name: /reactivate tenant/i }),
		).toBeInTheDocument();
		expect(
			screen.getByRole("button", { name: /^reactivate$/i }),
		).toBeDisabled();

		await user.type(
			screen.getByPlaceholderText(/enter tenant uuid/i),
			tenantUuid,
		);
		expect(screen.getByRole("button", { name: /^reactivate$/i })).toBeEnabled();

		await user.click(screen.getByRole("button", { name: /^reactivate$/i }));

		expect(mutate).toHaveBeenCalledWith(
			tenantUuid,
			expect.objectContaining({ onSuccess: expect.any(Function) }),
		);

		await waitFor(() => {
			expect(
				screen.queryByRole("heading", { name: /reactivate tenant/i }),
			).not.toBeInTheDocument();
		});
	});

	it("closes when cancel is clicked", async () => {
		const user = userEvent.setup();

		mockUseReactivateTenant.mockReturnValue({
			error: null,
			isPending: false,
			mutate: vi.fn(),
		});

		render(<ReactivateTenantForm />);

		await user.click(
			screen.getAllByRole("button", { name: /reactivate tenant/i })[0],
		);
		await user.click(screen.getByRole("button", { name: /cancel/i }));

		expect(
			screen.queryByRole("heading", { name: /reactivate tenant/i }),
		).not.toBeInTheDocument();
	});
});
