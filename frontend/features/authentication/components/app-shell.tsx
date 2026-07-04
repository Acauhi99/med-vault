"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

export function AppShell() {
	const router = useRouter();

	useEffect(() => {
		router.replace("/cases/");
	}, [router]);

	return (
		<div className="min-h-screen bg-slate-950 text-slate-50">
			<header className="border-b border-white/10 bg-slate-950/95 backdrop-blur">
				<div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4">
					<div>
						<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
							MedVault
						</p>
						<h1 className="mt-1 text-xl font-semibold">Redirecting...</h1>
					</div>
				</div>
			</header>
		</div>
	);
}
