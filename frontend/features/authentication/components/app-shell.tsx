"use client";

import type {
  TenantSummary,
  UserProfile,
} from "@/features/authentication/schemas/auth";

type AppShellProps = {
  tenant: TenantSummary;
  user: UserProfile;
  onSignOut: () => void;
};

const moduleCards = [
  {
    title: "Cases",
    body: "Open, assign, diagnose, close.",
  },
  {
    title: "Images",
    body: "Upload, review, and confirm scans.",
  },
  {
    title: "Audit",
    body: "Track every state change.",
  },
  {
    title: "Admin",
    body: "Manage members and tenant access.",
  },
];

export function AppShell({ tenant, user, onSignOut }: AppShellProps) {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-50">
      <header className="border-b border-white/10 bg-slate-950/95 backdrop-blur">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4">
          <div>
            <p className="text-xs uppercase tracking-[0.35em] text-sky-300">
              MedVault
            </p>
            <h1 className="mt-1 text-xl font-semibold">Clinical workspace</h1>
          </div>
          <div className="flex items-center gap-3 text-sm">
            <div className="rounded-full border border-white/10 bg-white/5 px-3 py-1">
              {tenant.tenantName}
            </div>
            <div className="rounded-full border border-white/10 bg-white/5 px-3 py-1 capitalize">
              {user.role}
            </div>
            <button
              type="button"
              onClick={onSignOut}
              className="rounded-full border border-white/10 px-3 py-1 text-slate-200 transition hover:border-white/30 hover:bg-white/10"
            >
              Sign out
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto grid max-w-6xl gap-6 px-6 py-8 lg:grid-cols-[280px_minmax(0,1fr)]">
        <aside className="rounded-3xl border border-white/10 bg-white/5 p-5">
          <p className="text-xs uppercase tracking-[0.3em] text-sky-300">
            Active tenant
          </p>
          <h2 className="mt-3 text-2xl font-semibold">{tenant.tenantName}</h2>
          <p className="mt-2 text-sm text-slate-300">{user.email}</p>

          <div className="mt-6 grid gap-3 text-sm">
            <div className="rounded-2xl bg-white/5 p-3">
              <div className="text-slate-400">Role</div>
              <div className="mt-1 font-medium capitalize">{user.role}</div>
            </div>
            <div className="rounded-2xl bg-white/5 p-3">
              <div className="text-slate-400">Tenant ID</div>
              <div className="mt-1 break-all font-mono text-xs text-slate-200">
                {tenant.tenantId}
              </div>
            </div>
            <div className="rounded-2xl bg-white/5 p-3">
              <div className="text-slate-400">User status</div>
              <div className="mt-1 font-medium capitalize">{user.status}</div>
            </div>
          </div>
        </aside>

        <section className="grid gap-4 md:grid-cols-2">
          {moduleCards.map((module) => (
            <article
              key={module.title}
              className="rounded-3xl border border-white/10 bg-white/5 p-5"
            >
              <p className="text-xs uppercase tracking-[0.3em] text-sky-300">
                Module
              </p>
              <h2 className="mt-3 text-2xl font-semibold">{module.title}</h2>
              <p className="mt-2 text-sm leading-6 text-slate-300">
                {module.body}
              </p>
            </article>
          ))}

          <article className="rounded-3xl border border-emerald-400/20 bg-emerald-400/10 p-5 md:col-span-2">
            <p className="text-xs uppercase tracking-[0.3em] text-emerald-200">
              Ready
            </p>
            <h2 className="mt-3 text-2xl font-semibold">Auth complete</h2>
            <p className="mt-2 max-w-2xl text-sm leading-6 text-emerald-50/80">
              Session established for {user.email}. The shell is live and ready
              for the patient, doctor, and admin slices.
            </p>
          </article>
        </section>
      </main>
    </div>
  );
}
