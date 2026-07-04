"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { CaseDetail as AdminCaseDetail } from "@/features/admin/components/case-detail";
import { CaseList as AdminCaseList } from "@/features/admin/components/case-list";
import { CaseList as DoctorCaseList } from "@/features/doctors/components/case-list";
import { CaseDetail as PatientCaseDetail } from "@/features/patients/components/case-detail";
import { CaseList as PatientCaseList } from "@/features/patients/components/case-list";
import { clearAuthSession } from "@/infrastructure/auth/session-store";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import { Sidebar } from "@/shared/components/sidebar";

export default function CasesPage() {
	const session = useAuthSession();
	const queryClient = useQueryClient();
	const router = useRouter();
	const [selectedCaseId, setSelectedCaseId] = useState<string | null>(null);

	function handleSignOut() {
		queryClient.clear();
		clearAuthSession();
		router.push("/");
	}

	if (!session.user || !session.activeTenant) {
		return (
			<div className="flex min-h-screen items-center justify-center bg-slate-950 text-slate-50">
				<p className="text-slate-400">Loading...</p>
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
						<h1 className="mt-1 text-xl font-semibold">Cases</h1>
					</div>
					<div className="flex items-center gap-3 text-sm">
						<span className="rounded-full border border-white/10 bg-white/5 px-3 py-1">
							{session.activeTenant.tenantName}
						</span>
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
				<section>
					{session.user.role === "patient" &&
						(selectedCaseId ? (
							<PatientCaseDetail
								caseId={selectedCaseId}
								onBack={() => setSelectedCaseId(null)}
							/>
						) : (
							<PatientCaseList
								onSelectCase={setSelectedCaseId}
								onCreateCase={() => router.push("/cases/new/")}
							/>
						))}
					{session.user.role === "doctor" && <DoctorCaseList />}
					{session.user.role === "administrator" &&
						(selectedCaseId ? (
							<AdminCaseDetail
								caseId={selectedCaseId}
								onBack={() => setSelectedCaseId(null)}
							/>
						) : (
							<AdminCaseList onSelectCase={setSelectedCaseId} />
						))}
				</section>
			</main>
		</div>
	);
}
