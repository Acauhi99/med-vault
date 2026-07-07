"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { MemberList } from "@/features/admin/components/member-list";
import { ReactivateTenantForm } from "@/features/admin/components/reactivate-tenant-form";
import { clearAuthSession } from "@/infrastructure/auth/session-store";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import { RouteLoading } from "@/shared/components/route-loading";
import { Sidebar } from "@/shared/components/sidebar";

export default function MembersPage() {
	const session = useAuthSession();
	const queryClient = useQueryClient();
	const router = useRouter();

	function handleSignOut() {
		queryClient.clear();
		clearAuthSession();
		router.push("/");
	}

	if (!session.user || !session.activeTenant) {
		return <RouteLoading title="Members" />;
	}

	if (session.user.role !== "administrator") {
		return (
			<div className="flex min-h-screen items-center justify-center bg-slate-950 text-slate-50">
				<p className="text-red-400">Access denied. Administrators only.</p>
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-slate-950 text-slate-50">
			<header className="border-b border-white/10 bg-slate-950/95 backdrop-blur">
				<div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4">
					<div>
						<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
							MedVault
						</p>
						<h1 className="mt-1 text-xl font-semibold">Tenant Members</h1>
					</div>
				</div>
			</header>

			<main className="mx-auto grid max-w-6xl gap-6 px-6 py-8 lg:grid-cols-[280px_minmax(0,1fr)]">
				<Sidebar
					user={session.user}
					tenants={session.tenants}
					activeTenant={session.activeTenant}
					accessToken={session.accessToken ?? ""}
					onSignOut={handleSignOut}
				/>
				<section className="space-y-6">
					<div className="rounded-3xl border border-white/10 bg-white/[0.03] p-6">
						<h2 className="text-lg font-semibold text-white">
							Tenant operations
						</h2>
						<p className="mt-1 text-sm text-slate-400">
							Reinstate a suspended tenant by UUID.
						</p>
						<div className="mt-4">
							<ReactivateTenantForm />
						</div>
					</div>
					<MemberList />
				</section>
			</main>
		</div>
	);
}
