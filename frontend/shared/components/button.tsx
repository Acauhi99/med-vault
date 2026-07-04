"use client";

import type { ReactNode } from "react";

type ButtonProps = {
	children: ReactNode;
	variant?: "primary" | "secondary" | "ghost" | "danger";
	size?: "sm" | "md" | "lg";
	disabled?: boolean;
	busy?: boolean;
	type?: "button" | "submit";
	onClick?: () => void;
};

const variantStyles = {
	primary:
		"bg-sky-400 text-slate-950 hover:bg-sky-300 disabled:bg-slate-700 disabled:text-slate-300",
	secondary:
		"border border-white/10 bg-white/5 text-slate-200 hover:border-white/30 hover:bg-white/10",
	ghost: "text-slate-300 hover:bg-white/5 hover:text-white",
	danger:
		"bg-red-500/10 text-red-400 border border-red-500/20 hover:bg-red-500/20",
};

const sizeStyles = {
	sm: "h-8 px-3 text-xs",
	md: "h-10 px-4 text-sm",
	lg: "h-12 px-6 text-sm",
};

export function Button({
	children,
	variant = "primary",
	size = "md",
	disabled = false,
	busy = false,
	type = "button",
	onClick,
}: ButtonProps) {
	return (
		<button
			type={type}
			disabled={disabled || busy}
			aria-busy={busy}
			onClick={onClick}
			className={`inline-flex items-center justify-center gap-2 rounded-2xl font-semibold transition disabled:cursor-not-allowed disabled:opacity-60 ${variantStyles[variant]} ${sizeStyles[size]}`}
		>
			{busy && (
				<svg
					className="h-4 w-4 animate-spin"
					viewBox="0 0 24 24"
					fill="none"
					aria-hidden="true"
				>
					<circle
						className="opacity-25"
						cx="12"
						cy="12"
						r="10"
						stroke="currentColor"
						strokeWidth="4"
					/>
					<path
						className="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
					/>
				</svg>
			)}
			{children}
		</button>
	);
}
