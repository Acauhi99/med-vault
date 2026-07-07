"use client";

import type { ReactNode } from "react";

type DialogProps = {
	open: boolean;
	onClose: () => void;
	title: string;
	children: ReactNode;
};

export function Dialog({ open, onClose, title, children }: DialogProps) {
	if (!open) return null;

	return (
		<div className="fixed inset-0 z-50 flex items-center justify-center">
			<button
				type="button"
				className="fixed inset-0 bg-black/60"
				onClick={onClose}
				aria-label="Close dialog"
			/>
			<div className="relative mx-4 w-full max-w-lg rounded-3xl border border-white/10 bg-slate-950 p-6 shadow-2xl">
				<div className="mb-4 flex items-center justify-between">
					<h2 className="text-lg font-semibold text-white">{title}</h2>
					<button
						type="button"
						onClick={onClose}
						aria-label="Close dialog"
						className="rounded-xl p-1 text-slate-400 transition hover:bg-white/5 hover:text-white"
					>
						<svg
							className="h-5 w-5"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
						>
							<path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
						</svg>
					</button>
				</div>
				{children}
			</div>
		</div>
	);
}
