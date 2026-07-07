const features = [
	{
		icon: "🏢",
		title: "Multi-Tenant Isolation",
		body: "Every query scoped to tenant. No cross-tenant data leakage — enforced at database, application, and API layers.",
	},
	{
		icon: "🔐",
		title: "Role-Based Access",
		body: "Patient, Doctor, Administrator — each role sees only what it needs. Minimum Necessary Standard enforced.",
	},
	{
		icon: "📋",
		title: "Complete Audit Trail",
		body: "Every state-changing operation logged. Structured JSON with 6-year retention for HIPAA §164.530(j).",
	},
	{
		icon: "🛡️",
		title: "End-to-End Encryption",
		body: "AES-256 at rest, TLS 1.2+ in transit. KMS-managed keys for medical images. Zero plaintext PHI.",
	},
];

export function Features() {
	return (
		<section id="features" className="bg-slate-950 py-24">
			<div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
				<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
					Platform Capabilities
				</p>
				<h2 className="mt-4 text-3xl font-semibold text-white sm:text-4xl">
					Built for healthcare from day one.
				</h2>
				<p className="mt-4 max-w-xl text-slate-400">
					Every architectural decision maps to a specific HIPAA requirement. No
					afterthought security.
				</p>

				<div className="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
					{features.map((f) => (
						<article
							key={f.title}
							className="rounded-3xl border border-white/10 bg-white/5 p-6"
						>
							<span className="text-3xl">{f.icon}</span>
							<h3 className="mt-4 text-lg font-semibold text-white">
								{f.title}
							</h3>
							<p className="mt-2 text-sm leading-6 text-slate-300">{f.body}</p>
						</article>
					))}
				</div>
			</div>
		</section>
	);
}
