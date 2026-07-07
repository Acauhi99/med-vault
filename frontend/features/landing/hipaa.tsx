const rules = [
	{
		icon: "👁️",
		title: "Privacy Rule",
		regulation: "45 CFR §164.500–534",
		body: "Controls who can access and disclose Protected Health Information (PHI). Patients have rights to access, amend, and restrict sharing of their data.",
		implementation:
			"Role-based access control at every endpoint. Tenant isolation ensures data never crosses organizational boundaries.",
	},
	{
		icon: "🔒",
		title: "Security Rule",
		regulation: "45 CFR §164.302–318",
		body: "Requires administrative, physical, and technical safeguards for electronic PHI (ePHI). Encryption, access controls, and audit logging are mandatory.",
		implementation:
			"AES-256 encryption at rest, TLS 1.2+ in transit, bcrypt password hashing, and network segmentation via VPC.",
	},
	{
		icon: "🚨",
		title: "Breach Notification",
		regulation: "45 CFR §164.400–414",
		body: "Covered entities must notify affected individuals within 60 days of discovering a breach. HHS and media notification required for large breaches.",
		implementation:
			"Complete audit trail with structured logging. All mutations recorded with timestamps, user identity, and tenant context for forensic analysis.",
	},
	{
		icon: "⚖️",
		title: "Minimum Necessary",
		regulation: "45 CFR §164.502(b)",
		body: "Access and disclosure of PHI must be limited to the minimum necessary for the intended purpose. Each role sees only what it needs.",
		implementation:
			"Three distinct roles (Patient, Doctor, Administrator) with scoped permissions. Patients see own cases, Doctors see assigned cases, Admins see all within tenant.",
	},
];

export function Hipaa() {
	return (
		<section id="hipaa" className="bg-slate-900/50 py-24">
			<div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
				<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
					HIPAA Compliance
				</p>
				<h2 className="mt-4 text-3xl font-semibold text-white sm:text-4xl">
					Built for compliance from day one.
				</h2>
				<p className="mt-4 max-w-xl text-slate-400">
					Every feature in MedVault maps to a specific HIPAA requirement. This
					isn&apos;t security added later — it&apos;s the foundation.
				</p>

				<div className="mt-12 grid gap-4 sm:grid-cols-2">
					{rules.map((r) => (
						<article
							key={r.title}
							className="rounded-3xl border border-white/10 bg-white/5 p-6"
						>
							<div className="flex items-center gap-3">
								<span className="text-2xl">{r.icon}</span>
								<div>
									<h3 className="text-lg font-semibold text-white">
										{r.title}
									</h3>
									<p className="text-xs text-sky-300/80">{r.regulation}</p>
								</div>
							</div>
							<p className="mt-4 text-sm leading-6 text-slate-300">{r.body}</p>
							<div className="mt-4 rounded-2xl border border-sky-400/10 bg-sky-400/5 px-4 py-3">
								<p className="text-xs font-medium text-sky-300">
									How MedVault implements this
								</p>
								<p className="mt-1 text-sm leading-6 text-slate-200">
									{r.implementation}
								</p>
							</div>
						</article>
					))}
				</div>
			</div>
		</section>
	);
}
