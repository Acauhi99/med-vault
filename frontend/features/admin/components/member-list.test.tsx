// @vitest-environment happy-dom

import "../../../test/setup";

import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

const mockUseMemberList = vi.fn();
const mockUseRemoveMember = vi.fn();

vi.mock("../hooks/use-members", () => ({
	useMemberList: () => mockUseMemberList(),
	useRemoveMember: () => mockUseRemoveMember(),
}));

vi.mock("./add-member-form", () => ({
	AddMemberForm: () => <div data-testid="add-member-form" />,
}));

vi.mock("@/shared/components/table", () => ({
	Table: ({ children }: { children: ReactNode }) => <div>{children}</div>,
	TableBody: ({ children }: { children: ReactNode }) => <div>{children}</div>,
	TableHead: ({ children }: { children: ReactNode }) => <div>{children}</div>,
	TableRow: ({ children }: { children: ReactNode }) => <div>{children}</div>,
	Td: ({ children }: { children: ReactNode }) => <div>{children}</div>,
	Th: ({ children }: { children: ReactNode }) => <div>{children}</div>,
}));

import { MemberList } from "./member-list";

beforeEach(() => {
	mockUseMemberList.mockReset();
	mockUseRemoveMember.mockReset();
});

describe("MemberList", () => {
	it("shows the loading state", () => {
		mockUseMemberList.mockReturnValue({
			data: undefined,
			error: null,
			isLoading: true,
			refetch: vi.fn(),
		});
		mockUseRemoveMember.mockReturnValue({
			isPending: false,
			mutate: vi.fn(),
		});

		render(<MemberList />);

		expect(screen.getByText("Members")).toBeInTheDocument();
		expect(screen.queryByTestId("add-member-form")).not.toBeInTheDocument();
		expect(screen.queryByText(/no members/i)).not.toBeInTheDocument();
	});

	it("shows the empty state", () => {
		mockUseMemberList.mockReturnValue({
			data: [],
			error: null,
			isLoading: false,
			refetch: vi.fn(),
		});
		mockUseRemoveMember.mockReturnValue({
			isPending: false,
			mutate: vi.fn(),
		});

		render(<MemberList />);

		expect(screen.getByText(/no members/i)).toBeInTheDocument();
		expect(screen.getByTestId("add-member-form")).toBeInTheDocument();
	});

	it("renders members and removes one", async () => {
		const user = userEvent.setup();
		const mutate = vi.fn();

		mockUseMemberList.mockReturnValue({
			data: [
				{
					userId: "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001",
					name: "Dr. Smith",
					role: "doctor",
				},
			],
			error: null,
			isLoading: false,
			refetch: vi.fn(),
		});
		mockUseRemoveMember.mockReturnValue({
			isPending: false,
			mutate,
		});

		render(<MemberList />);

		expect(screen.getByText("Dr. Smith")).toBeInTheDocument();
		expect(screen.getByText("doctor")).toBeInTheDocument();

		await user.click(screen.getByRole("button", { name: /remove/i }));

		expect(mutate).toHaveBeenCalledWith("9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001");
	});
});
