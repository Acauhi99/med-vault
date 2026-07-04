import type { ReactNode } from "react";

type BadgeVariant = "default" | "success" | "warning" | "danger" | "info";

const variantStyles: Record<BadgeVariant, string> = {
	default: "border-white/10 bg-white/5 text-slate-300",
	success: "border-emerald-400/20 bg-emerald-400/10 text-emerald-300",
	warning: "border-amber-400/20 bg-amber-400/10 text-amber-300",
	danger: "border-red-400/20 bg-red-400/10 text-red-300",
	info: "border-sky-400/20 bg-sky-400/10 text-sky-300",
};

type BadgeProps = {
	children: ReactNode;
	variant?: BadgeVariant;
};

export function Badge({ children, variant = "default" }: BadgeProps) {
	return (
		<span
			className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${variantStyles[variant]}`}
		>
			{children}
		</span>
	);
}

export function StatusBadge({ status }: { status: string }) {
	const variantMap: Record<string, BadgeVariant> = {
		open: "info",
		assigned: "warning",
		diagnosed: "success",
		closed: "default",
		active: "success",
		inactive: "danger",
		patient: "info",
		doctor: "warning",
		administrator: "success",
	};

	return <Badge variant={variantMap[status] ?? "default"}>{status}</Badge>;
}
