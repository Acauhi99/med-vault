import Link from "next/link";

export function Hero() {
	return (
		<section className="relative flex min-h-screen items-center overflow-hidden bg-slate-950 pt-16">
			<div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-sky-900/20 via-slate-950 to-slate-950" />

			<div className="relative mx-auto max-w-6xl px-4 py-24 sm:px-6 lg:px-8">
				<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
					HIPAA-Compliant Platform
				</p>
				<h1 className="mt-6 max-w-3xl text-5xl font-semibold leading-tight text-white sm:text-6xl lg:text-7xl">
					Healthcare data,
					<br />
					secured by design.
				</h1>
				<p className="mt-6 max-w-xl text-lg leading-7 text-slate-300">
					A multi-tenant healthcare workspace built from the ground up with
					HIPAA compliance, end-to-end encryption, and complete audit trails.
					Enterprise-grade security, demonstrated as a PoC.
				</p>

				<div className="mt-10 flex flex-wrap gap-4">
					<Link
						href="/login"
						className="inline-flex h-12 items-center justify-center rounded-2xl bg-sky-400 px-6 text-sm font-semibold text-slate-950 transition hover:bg-sky-300"
					>
						Get Started
					</Link>
					<a
						href="#architecture"
						className="inline-flex h-12 items-center justify-center rounded-2xl border border-white/10 px-6 text-sm font-medium text-slate-300 transition hover:border-white/30 hover:bg-white/5 hover:text-white"
					>
						View Architecture
					</a>
				</div>

				<div className="mt-16 grid max-w-lg grid-cols-3 gap-8 border-t border-white/10 pt-10">
					<div>
						<p className="text-2xl font-semibold text-white">AES-256</p>
						<p className="mt-1 text-sm text-slate-400">Encryption at rest</p>
					</div>
					<div>
						<p className="text-2xl font-semibold text-white">TLS 1.2+</p>
						<p className="mt-1 text-sm text-slate-400">Encryption in transit</p>
					</div>
					<div>
						<p className="text-2xl font-semibold text-white">6 years</p>
						<p className="mt-1 text-sm text-slate-400">Audit retention</p>
					</div>
				</div>
			</div>
		</section>
	);
}
