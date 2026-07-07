const pillars = [
	{
		icon: "🧩",
		title: "Domain-Driven Design + CQRS",
		body: "Bounded contexts (Identity, Clinical, Imaging, Audit) with command-query separation. Clean architecture layers from domain to infrastructure.",
	},
	{
		icon: "☁️",
		title: "AWS ECS Fargate + RDS",
		body: "Serverless containers with managed PostgreSQL. Private subnets, WAF protection, KMS encryption. Infrastructure as Code via Terraform.",
	},
	{
		icon: "📐",
		title: "OpenAPI Design-First",
		body: "API contracts defined before implementation. Type-safe code generation for both Go backend and TypeScript frontend. Never manually write HTTP schemas.",
	},
];

export function Architecture() {
	return (
		<section id="architecture" className="bg-slate-950 py-24">
			<div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
				<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
					Architecture
				</p>
				<h2 className="mt-4 text-3xl font-semibold text-white sm:text-4xl">
					Enterprise-grade, demonstrated.
				</h2>
				<p className="mt-4 max-w-xl text-slate-400">
					A modular monolith that proves compliance and architecture decisions
					can coexist without over-engineering.
				</p>

				<div className="mt-12 grid gap-4 sm:grid-cols-3">
					{pillars.map((p) => (
						<article
							key={p.title}
							className="rounded-3xl border border-white/10 bg-white/5 p-6"
						>
							<span className="text-3xl">{p.icon}</span>
							<h3 className="mt-4 text-lg font-semibold text-white">
								{p.title}
							</h3>
							<p className="mt-2 text-sm leading-6 text-slate-300">{p.body}</p>
						</article>
					))}
				</div>
			</div>
		</section>
	);
}
