"use client";

import { useState } from "react";

import { Button } from "@/shared/components/button";
import { Skeleton } from "@/shared/components/skeleton";
import { EmptyState, ErrorState } from "@/shared/components/states";

import { useCaseImages } from "../hooks/use-cases";
import { useDownloadURL } from "../hooks/use-image-upload";
import type { ImageResponse } from "../schemas/cases";

type ImageListProps = {
	caseId: string;
};

export function ImageList({ caseId }: ImageListProps) {
	const {
		data: images,
		isLoading,
		isError,
		error,
		refetch,
	} = useCaseImages(caseId);
	const downloadMutation = useDownloadURL();
	const [downloadingId, setDownloadingId] = useState<string | null>(null);

	async function handleDownload(image: ImageResponse) {
		setDownloadingId(image.id);
		try {
			const result = await downloadMutation.mutateAsync(image.id);
			window.open(result.downloadUrl, "_blank");
		} catch {
		} finally {
			setDownloadingId(null);
		}
	}

	if (isLoading) {
		return <Skeleton className="h-20 w-full" />;
	}

	if (isError) {
		return (
			<ErrorState
				title="Failed to load images"
				message={
					error instanceof Error ? error.message : "Unable to load images."
				}
				onRetry={refetch}
			/>
		);
	}

	if (!images || images.length === 0) {
		return (
			<EmptyState
				title="No images"
				description="Upload an image to attach it to this case."
			/>
		);
	}

	return (
		<div className="space-y-2">
			{images.map((image: ImageResponse) => (
				<div
					key={image.id}
					className="flex items-center justify-between gap-4 rounded-2xl border border-white/10 bg-white/[0.03] p-3"
				>
					<div className="min-w-0 space-y-1">
						<p className="truncate text-sm text-slate-200">{image.fileName}</p>
						<p className="text-xs text-slate-400">
							{image.contentType} &middot;{" "}
							{new Date(image.uploadedAt).toLocaleDateString()}
						</p>
					</div>
					<Button
						variant="ghost"
						size="sm"
						disabled={downloadingId === image.id}
						onClick={() => handleDownload(image)}
					>
						{downloadingId === image.id ? "Loading..." : "Download"}
					</Button>
				</div>
			))}
		</div>
	);
}
