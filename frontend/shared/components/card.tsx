import type { ReactNode } from "react";

type CardProps = {
	children: ReactNode;
	className?: string;
};

export function Card({ children, className = "" }: CardProps) {
	return (
		<div
			className={`rounded-3xl border border-white/10 bg-white/5 p-5 ${className}`}
		>
			{children}
		</div>
	);
}

export function CardHeader({ children, className = "" }: CardProps) {
	return <div className={`mb-4 ${className}`}>{children}</div>;
}

export function CardTitle({ children }: { children: ReactNode }) {
	return <h3 className="text-lg font-semibold text-white">{children}</h3>;
}

export function CardDescription({ children }: { children: ReactNode }) {
	return <p className="mt-1 text-sm text-slate-400">{children}</p>;
}
