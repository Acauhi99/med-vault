type SkeletonProps = {
	className?: string;
};

export function Skeleton({ className = "" }: SkeletonProps) {
	return (
		<div className={`animate-pulse rounded-2xl bg-white/5 ${className}`} />
	);
}

export function CardSkeleton() {
	return (
		<div className="rounded-3xl border border-white/10 bg-white/5 p-5">
			<Skeleton className="h-4 w-1/3 mb-3" />
			<Skeleton className="h-6 w-1/2 mb-2" />
			<Skeleton className="h-4 w-2/3" />
		</div>
	);
}

export function TableSkeleton({ rows = 5 }: { rows?: number }) {
	return (
		<div className="space-y-3">
			{Array.from({ length: rows }, (_, i) => (
				<div key={`skel-${String(i)}`} className="flex gap-4">
					<Skeleton className="h-4 w-1/4" />
					<Skeleton className="h-4 w-1/4" />
					<Skeleton className="h-4 w-1/4" />
					<Skeleton className="h-4 w-1/4" />
				</div>
			))}
		</div>
	);
}
