import type { ReactNode } from "react";

type EmptyStateProps = {
	title: string;
	description: string;
	action?: ReactNode;
};

export function EmptyState({ title, description, action }: EmptyStateProps) {
	return (
		<div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-white/10 bg-white/[0.02] px-6 py-12 text-center">
			<h3 className="text-lg font-semibold text-white">{title}</h3>
			<p className="mt-2 max-w-sm text-sm text-slate-400">{description}</p>
			{action && <div className="mt-4">{action}</div>}
		</div>
	);
}

type ErrorStateProps = {
	title: string;
	message: string;
	onRetry?: () => void;
	action?: ReactNode;
};

export function ErrorState({
	title,
	message,
	onRetry,
	action,
}: ErrorStateProps) {
	return (
		<div className="flex flex-col items-center justify-center rounded-3xl border border-red-400/20 bg-red-400/5 px-6 py-12 text-center">
			<h3 className="text-lg font-semibold text-white">{title}</h3>
			<p className="mt-2 max-w-sm text-sm text-slate-400">{message}</p>
			{action && <div className="mt-4">{action}</div>}
			{onRetry && (
				<button
					type="button"
					onClick={onRetry}
					className="mt-4 inline-flex h-10 items-center justify-center rounded-2xl bg-white px-4 text-sm font-medium text-slate-950 transition hover:bg-slate-200"
				>
					Try again
				</button>
			)}
		</div>
	);
}
