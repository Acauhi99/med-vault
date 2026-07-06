type RouteLoadingProps = {
	title: string;
	description?: string;
};

export function RouteLoading({
	title,
	description = "Loading workspace...",
}: RouteLoadingProps) {
	return (
		<div className="flex min-h-screen items-center justify-center bg-slate-950 text-slate-50">
			<div className="rounded-3xl border border-white/10 bg-white/5 px-6 py-8 text-center shadow-2xl shadow-black/20">
				<div className="mx-auto h-10 w-10 animate-spin rounded-full border-2 border-sky-400 border-t-transparent" />
				<h1 className="mt-4 text-lg font-semibold">{title}</h1>
				<p className="mt-1 text-sm text-slate-400">{description}</p>
			</div>
		</div>
	);
}
