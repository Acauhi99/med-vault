const companies = [
	"Epic Systems",
	"Oracle Health",
	"Kaiser Permanente",
	"Teladoc Health",
	"Allscripts",
	"athenahealth",
	"eClinicalWorks",
	"NextGen Healthcare",
	"Meditech",
	"GE HealthCare",
];

export function TrustedBy() {
	const doubled = [
		...companies.map((c) => `a-${c}`),
		...companies.map((c) => `b-${c}`),
	];

	return (
		<section className="bg-slate-950 py-24 overflow-hidden">
			<div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
				<p className="text-center text-xs uppercase tracking-[0.35em] text-sky-300">
					Trusted by industry leaders
				</p>
				<p className="mt-3 text-center text-sm text-slate-400">
					HIPAA compliance is table stakes for these organizations.
				</p>
			</div>

			<div className="mt-12 relative">
				<div className="absolute left-0 top-0 z-10 h-full w-32 bg-gradient-to-r from-slate-950 to-transparent" />
				<div className="absolute right-0 top-0 z-10 h-full w-32 bg-gradient-to-l from-slate-950 to-transparent" />

				<div className="flex w-max gap-12 animate-scroll">
					{doubled.map((key) => {
						const name = key.slice(2);
						return (
							<div
								key={key}
								className="flex h-16 w-44 shrink-0 items-center justify-center rounded-2xl border border-white/10 bg-white/5 px-6"
							>
								<span className="text-sm font-medium text-slate-300 whitespace-nowrap">
									{name}
								</span>
							</div>
						);
					})}
				</div>
			</div>
		</section>
	);
}
