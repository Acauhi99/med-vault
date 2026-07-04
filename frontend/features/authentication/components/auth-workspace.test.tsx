// @vitest-environment happy-dom

import "../../../test/setup";

import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { AuthSession } from "@/features/authentication/schemas/auth";

const mockUseAuthSession = vi.fn();
const mockUseMutation = vi.fn();
const mockUseQuery = vi.fn();
const mockUseQueryClient = vi.fn();
const mockUseInactivityLogoff = vi.fn();
const mockUseTokenRefresh = vi.fn();
const mockClearAuthSession = vi.fn();
const mockUpdateAuthSession = vi.fn();
const mockLogout = vi.fn();

let session: AuthSession;
let queryClient: { cancelQueries: () => Promise<void>; clear: () => void };
let tokenRefresh: { start: () => void; stop: () => void };

const emptySession = {
  accessToken: null,
  refreshToken: null,
  tenants: [],
  activeTenant: null,
  user: null,
} satisfies AuthSession;

const authedSession = {
  accessToken: "token",
  refreshToken: "refresh",
  tenants: [
    {
      tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
      tenantName: "Test Clinic",
      role: "administrator" as const,
    },
  ],
  activeTenant: {
    tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
    tenantName: "Test Clinic",
    role: "administrator" as const,
  },
  user: {
    id: "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001",
    tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
    email: "admin@example.com",
    role: "administrator",
    status: "active",
    createdAt: "2026-07-04T15:00:00Z",
  },
} satisfies AuthSession;

vi.mock("@tanstack/react-query", () => ({
  useMutation: () => mockUseMutation(),
  useQuery: () => mockUseQuery(),
  useQueryClient: () => mockUseQueryClient(),
}));

vi.mock("@/infrastructure/auth/use-auth-session", () => ({
  useAuthSession: () => mockUseAuthSession(),
}));

vi.mock("@/infrastructure/auth/use-inactivity-logoff", () => ({
  useInactivityLogoff: () => mockUseInactivityLogoff(),
}));

vi.mock("@/infrastructure/auth/use-token-refresh", () => ({
  useTokenRefresh: () => mockUseTokenRefresh(),
}));

vi.mock("@/infrastructure/auth/session-store", () => ({
  clearAuthSession: () => mockClearAuthSession(),
  updateAuthSession: (patch: Partial<AuthSession>) => mockUpdateAuthSession(patch),
}));

vi.mock("../services/auth", () => ({
  getCurrentUser: () => Promise.resolve(null),
  login: () => Promise.resolve(null),
  logout: (refreshToken: string) => mockLogout(refreshToken),
  register: () => Promise.resolve(null),
  selectTenant: () => Promise.resolve(null),
}));

vi.mock("./app-shell", () => ({
  AppShell: () => <div data-testid="app-shell" />,
}));

import { AuthWorkspace } from "./auth-workspace";

beforeEach(() => {
  session = emptySession;
  queryClient = {
    cancelQueries: vi.fn().mockResolvedValue(undefined),
    clear: vi.fn(),
  };
  tokenRefresh = {
    start: vi.fn(),
    stop: vi.fn(),
  };

  mockUseAuthSession.mockImplementation(() => session);
  mockUseMutation.mockReturnValue({
    isPending: false,
    mutate: vi.fn(),
    mutateAsync: vi.fn(),
    error: null,
  });
  mockUseQuery.mockReturnValue({
    data: null,
    error: null,
    isError: false,
    isFetching: false,
  });
  mockUseQueryClient.mockReturnValue(queryClient);
  mockUseInactivityLogoff.mockReturnValue(undefined);
  mockUseTokenRefresh.mockReturnValue(tokenRefresh);
  mockClearAuthSession.mockReset();
  mockUpdateAuthSession.mockReset();
  mockLogout.mockResolvedValue(undefined);
});

describe("AuthWorkspace", () => {
  it("renders auth tabs and switches to register", async () => {
    const user = userEvent.setup();

    render(<AuthWorkspace />);

    expect(
      screen.getByRole("button", { name: /^sign in$/i, pressed: true }),
    ).toHaveAttribute("aria-pressed", "true");

    await user.click(screen.getByRole("button", { name: /^register$/i }));

    expect(
      screen.getByRole("button", { name: /^register$/i, pressed: true }),
    ).toHaveAttribute("aria-pressed", "true");
    expect(
      screen.getByRole("button", { name: /create account/i }),
    ).toBeInTheDocument();
  });

  it("shows empty tenant access and signs out", async () => {
    const user = userEvent.setup();

    session = {
      accessToken: "token",
      refreshToken: "refresh",
      tenants: [],
      activeTenant: null,
      user: null,
    } satisfies AuthSession;

    render(<AuthWorkspace />);

    expect(screen.getByText(/no tenant access yet/i)).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /sign out/i }));

    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalledWith("refresh");
      expect(tokenRefresh.stop).toHaveBeenCalledTimes(1);
      expect(queryClient.cancelQueries).toHaveBeenCalledWith({
        queryKey: ["current-user"],
      });
      expect(queryClient.clear).toHaveBeenCalledTimes(1);
      expect(mockClearAuthSession).toHaveBeenCalledTimes(1);
    });
  });

  it("renders the app shell when the session is ready", () => {
    session = authedSession;

    render(<AuthWorkspace />);

    expect(screen.getByTestId("app-shell")).toBeInTheDocument();
  });
});
