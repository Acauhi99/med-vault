"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { CreateCaseForm } from "@/features/patients/components/create-case-form";
import { clearAuthSession } from "@/infrastructure/auth/session-store";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import { Sidebar } from "@/shared/components/sidebar";

export default function NewCasePage() {
	const session = useAuthSession();
	const queryClient = useQueryClient();
	const router = useRouter();

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
						<h1 className="mt-1 text-xl font-semibold">New Case</h1>
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
					<CreateCaseForm
						onCreated={() => router.push("/cases/")}
						onCancel={() => router.push("/cases/")}
					/>
				</section>
			</main>
		</div>
	);
}
