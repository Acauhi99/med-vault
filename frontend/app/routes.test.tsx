// @vitest-environment happy-dom

import "../test/setup";

import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { AuthSession } from "@/features/authentication/schemas/auth";

const mockUseAuthSession = vi.fn();
const mockUseRouter = vi.fn();
const mockUseQueryClient = vi.fn();

let session: AuthSession;
let routerPush: ReturnType<typeof vi.fn>;
let queryClient: { cancelQueries: () => Promise<void>; clear: () => void };

const tenant = {
  tenantId: "7b4a8cf8-6c5a-4f4b-8b36-4c6f4c59d001",
  tenantName: "Test Clinic",
  role: "administrator" as const,
};

const adminUser = {
  id: "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001",
  tenantId: tenant.tenantId,
  email: "admin@example.com",
  role: "administrator" as const,
  status: "active" as const,
  createdAt: "2026-07-04T15:00:00Z",
};

const patientUser = {
  id: "7d9db2bb-3f33-40a5-98d6-bf3b9f9b0001",
  tenantId: tenant.tenantId,
  email: "patient@example.com",
  role: "patient" as const,
  status: "active" as const,
  createdAt: "2026-07-04T15:00:00Z",
};

const doctorUser = {
  id: "9a7ce2c1-1411-4a63-93f7-0c2e0f1d0001",
  tenantId: tenant.tenantId,
  email: "doctor@example.com",
  role: "doctor" as const,
  status: "active" as const,
  createdAt: "2026-07-04T15:00:00Z",
};

const adminSession = {
  accessToken: "token",
  refreshToken: "refresh",
  tenants: [tenant],
  activeTenant: tenant,
  user: adminUser,
} satisfies AuthSession;

const patientSession = {
  accessToken: "token",
  refreshToken: "refresh",
  tenants: [tenant],
  activeTenant: tenant,
  user: patientUser,
} satisfies AuthSession;

const doctorSession = {
  accessToken: "token",
  refreshToken: "refresh",
  tenants: [tenant],
  activeTenant: tenant,
  user: doctorUser,
} satisfies AuthSession;

vi.mock("@tanstack/react-query", () => ({
  useQueryClient: () => mockUseQueryClient(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => mockUseRouter(),
}));

vi.mock("@/infrastructure/auth/use-auth-session", () => ({
  useAuthSession: () => mockUseAuthSession(),
}));

vi.mock("@/features/authentication/components/auth-workspace", () => ({
  AuthWorkspace: () => <div data-testid="auth-workspace" />,
}));

vi.mock("@/features/admin/components/audit-log-table", () => ({
  AuditLogTable: () => <div data-testid="audit-log-table" />,
}));

vi.mock("@/features/admin/components/member-list", () => ({
  MemberList: () => <div data-testid="member-list" />,
}));

vi.mock("@/features/admin/components/reactivate-tenant-form", () => ({
  ReactivateTenantForm: () => <div data-testid="reactivate-tenant-form" />,
}));

vi.mock("@/shared/components/sidebar", () => ({
  Sidebar: () => <aside data-testid="sidebar" />,
}));

type PatientCaseListProps = {
  onSelectCase: (id: string) => void;
  onCreateCase: () => void;
};

type CaseDetailProps = {
  caseId: string;
  onBack: () => void;
};

type CreateCaseFormProps = {
  onCreated: () => void;
  onCancel: () => void;
};

vi.mock("@/features/patients/components/case-list", () => ({
  CaseList: ({ onSelectCase, onCreateCase }: PatientCaseListProps) => (
    <div data-testid="patient-case-list">
      <button type="button" onClick={() => onSelectCase("case-1")}>Select case</button>
      <button type="button" onClick={onCreateCase}>Create case</button>
    </div>
  ),
}));

vi.mock("@/features/patients/components/case-detail", () => ({
  CaseDetail: ({ caseId, onBack }: CaseDetailProps) => (
    <div data-testid="patient-case-detail">
      <span>{caseId}</span>
      <button type="button" onClick={onBack}>Back</button>
    </div>
  ),
}));

vi.mock("@/features/doctors/components/case-list", () => ({
  CaseList: () => <div data-testid="doctor-case-list" />,
}));

vi.mock("@/features/admin/components/case-list", () => ({
  CaseList: () => <div data-testid="admin-case-list" />,
}));

vi.mock("@/features/admin/components/case-detail", () => ({
  CaseDetail: ({ caseId, onBack }: CaseDetailProps) => (
    <div data-testid="admin-case-detail">
      <span>{caseId}</span>
      <button type="button" onClick={onBack}>Back</button>
    </div>
  ),
}));

vi.mock("@/features/patients/components/create-case-form", () => ({
  CreateCaseForm: ({ onCreated, onCancel }: CreateCaseFormProps) => (
    <div data-testid="create-case-form">
      <button type="button" onClick={onCreated}>Save case</button>
      <button type="button" onClick={onCancel}>Cancel</button>
    </div>
  ),
}));

vi.mock("@/infrastructure/auth/session-store", () => ({
  clearAuthSession: vi.fn(),
}));

import Home from "./page";
import AuditPage from "./audit/page";
import CasesPage from "./cases/page";
import NewCasePage from "./cases/new/page";
import MembersPage from "./members/page";

beforeEach(() => {
  session = adminSession;
  routerPush = vi.fn();
  queryClient = {
    cancelQueries: vi.fn().mockResolvedValue(undefined),
    clear: vi.fn(),
  };

  mockUseAuthSession.mockImplementation(() => session);
  mockUseRouter.mockReturnValue({ push: routerPush });
  mockUseQueryClient.mockReturnValue(queryClient);
});

describe("app routes", () => {
  it("renders the home auth workspace", () => {
    render(<Home />);

    expect(screen.getByTestId("auth-workspace")).toBeInTheDocument();
  });

  it("renders the audit page shell", () => {
    render(<AuditPage />);

    expect(screen.getByText("Audit Logs")).toBeInTheDocument();
    expect(screen.getByTestId("sidebar")).toBeInTheDocument();
    expect(screen.getByTestId("audit-log-table")).toBeInTheDocument();
  });

  it("renders the members page shell", () => {
    render(<MembersPage />);

    expect(screen.getByText("Tenant Members")).toBeInTheDocument();
    expect(screen.getByTestId("sidebar")).toBeInTheDocument();
    expect(screen.getByTestId("reactivate-tenant-form")).toBeInTheDocument();
    expect(screen.getByTestId("member-list")).toBeInTheDocument();
  });

  it("renders the patient cases page and opens case detail", async () => {
    const user = userEvent.setup();

    session = patientSession;

    render(<CasesPage />);

    expect(screen.getByTestId("patient-case-list")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /select case/i }));

    expect(screen.getByTestId("patient-case-detail")).toBeInTheDocument();
    expect(screen.getByText("case-1")).toBeInTheDocument();
  });

  it("renders the doctor cases page", () => {
    session = doctorSession;

    render(<CasesPage />);

    expect(screen.getByTestId("doctor-case-list")).toBeInTheDocument();
  });

  it("renders the administrator cases page", () => {
    render(<CasesPage />);

    expect(screen.getByTestId("admin-case-list")).toBeInTheDocument();
  });

  it("renders the new case page and routes back", async () => {
    const user = userEvent.setup();

    render(<NewCasePage />);

    expect(screen.getByTestId("create-case-form")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /save case/i }));
    await user.click(screen.getByRole("button", { name: /cancel/i }));

    expect(routerPush).toHaveBeenCalledWith("/cases/");
    expect(routerPush).toHaveBeenCalledTimes(2);
  });
});
