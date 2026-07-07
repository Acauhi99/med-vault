import Link from "next/link";

const navLinks = [
	{ href: "#features", label: "Features" },
	{ href: "#hipaa", label: "HIPAA" },
	{ href: "#architecture", label: "Architecture" },
];

export function Footer() {
	return (
		<footer className="border-t border-white/10 bg-slate-950">
			<div className="mx-auto max-w-6xl px-4 py-12 sm:px-6 lg:px-8">
				<div className="flex flex-col items-start justify-between gap-8 sm:flex-row sm:items-center">
					<div>
						<Link href="/" className="text-lg font-semibold text-white">
							MedVault
						</Link>
						<p className="mt-2 max-w-sm text-sm text-slate-400">
							A HIPAA compliance PoC demonstrating secure multi-tenant
							architecture on AWS.
						</p>
					</div>

					<nav className="flex gap-6">
						{navLinks.map((link) => (
							<a
								key={link.href}
								href={link.href}
								className="text-sm text-slate-400 transition hover:text-white"
							>
								{link.label}
							</a>
						))}
						<Link
							href="/login"
							className="text-sm text-slate-400 transition hover:text-white"
						>
							Sign in
						</Link>
					</nav>
				</div>

				<div className="mt-8 border-t border-white/10 pt-6">
					<p className="text-xs text-slate-500">
						No real patient data. Built for architecture demonstration purposes.
					</p>
				</div>
			</div>
		</footer>
	);
}
