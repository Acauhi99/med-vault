"use client";

import { ErrorState } from "@/shared/components/states";

type Props = {
	error: Error & { digest?: string };
	reset: () => void;
};

export default function RouteError({ error, reset }: Props) {
	return (
		<div className="flex min-h-screen items-center justify-center bg-slate-950 px-4 text-slate-50">
			<ErrorState
				title="Failed to load audit logs"
				message={error.message}
				onRetry={reset}
			/>
		</div>
	);
}
