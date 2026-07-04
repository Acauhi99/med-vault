"use client";

import type { ReactNode } from "react";

type TableProps = {
	children: ReactNode;
	className?: string;
};

export function Table({ children, className = "" }: TableProps) {
	return (
		<div
			className={`overflow-x-auto rounded-3xl border border-white/10 ${className}`}
		>
			<table className="w-full text-left text-sm">{children}</table>
		</div>
	);
}

export function TableHead({ children }: { children: ReactNode }) {
	return (
		<thead className="border-b border-white/10 bg-white/[0.03]">
			<tr>{children}</tr>
		</thead>
	);
}

export function TableBody({ children }: { children: ReactNode }) {
	return <tbody className="divide-y divide-white/5">{children}</tbody>;
}

export function TableRow({ children }: { children: ReactNode }) {
	return <tr className="transition hover:bg-white/[0.02]">{children}</tr>;
}

export function Th({
	children,
	className = "",
}: {
	children: ReactNode;
	className?: string;
}) {
	return (
		<th
			className={`px-4 py-3 text-xs font-medium uppercase tracking-wider text-slate-400 ${className}`}
		>
			{children}
		</th>
	);
}

export function Td({
	children,
	className = "",
}: {
	children: ReactNode;
	className?: string;
}) {
	return (
		<td className={`px-4 py-3 text-slate-200 ${className}`}>{children}</td>
	);
}
