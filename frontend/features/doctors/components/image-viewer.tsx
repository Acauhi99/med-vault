"use client";
import { ErrorState } from "@/shared/components/states";
import { useCaseImages, useImageDownloadUrl } from "../hooks/use-images";

type Props = {
	caseId: string;
};

function ImageCard({
	imageId,
	fileName,
}: {
	imageId: string;
	fileName: string;
}) {
	const downloadQuery = useImageDownloadUrl(imageId);

	return (
		<div className="rounded-xl border border-white/10 bg-white/5 p-4">
			<p className="truncate text-sm text-slate-200">{fileName}</p>
			<div className="mt-3">
				{downloadQuery.isLoading && (
					<div className="h-8 w-24 animate-pulse rounded-lg bg-white/10" />
				)}
				{downloadQuery.data && (
					<a
						href={downloadQuery.data.downloadUrl}
						target="_blank"
						rel="noopener noreferrer"
						className="inline-block rounded-lg bg-sky-400/10 px-3 py-1.5 text-xs font-medium text-sky-300 transition hover:bg-sky-400/20"
					>
						Download
					</a>
				)}
				{downloadQuery.isError && (
					<span className="text-xs text-red-400">URL unavailable</span>
				)}
			</div>
		</div>
	);
}

export function ImageViewer({ caseId }: Props) {
	const {
		data: images,
		isLoading,
		isError,
		error,
		refetch,
	} = useCaseImages(caseId);

	if (isLoading) {
		return (
			<div className="grid gap-3 sm:grid-cols-2">
				{[1, 2].map((n) => (
					<div key={n} className="h-24 animate-pulse rounded-xl bg-white/5" />
				))}
			</div>
		);
	}

	if (isError) {
		return (
			<ErrorState
				title="Failed to load images"
				message={error instanceof Error ? error.message : "Unknown error"}
				onRetry={refetch}
			/>
		);
	}

	if (!images || images.length === 0) {
		return (
			<p className="text-sm text-slate-400">
				No images uploaded for this case.
			</p>
		);
	}

	return (
		<div className="grid gap-3 sm:grid-cols-2">
			{images?.map((img) => (
				<ImageCard key={img.id} imageId={img.id} fileName={img.fileName} />
			))}
		</div>
	);
}
