"use client";

import type { ReactNode } from "react";

type PageHeaderProps = {
	title: string;
	description?: string;
	actions?: ReactNode;
};

export function PageHeader({ title, description, actions }: PageHeaderProps) {
	return (
		<div className="flex flex-wrap items-start justify-between gap-4">
			<div>
				<h1 className="text-2xl font-semibold text-white">{title}</h1>
				{description && (
					<p className="mt-1 text-sm text-slate-400">{description}</p>
				)}
			</div>
			{actions && <div className="flex gap-2">{actions}</div>}
		</div>
	);
}
