"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";

import type { ContentType } from "../schemas/cases";
import { maxImageUploadSizeBytes } from "../schemas/cases";
import {
	confirmUpload,
	getDownloadURL,
	requestUploadURL,
	uploadFileToS3,
} from "../services/cases";

type UploadStatus =
	| "idle"
	| "requesting"
	| "uploading"
	| "confirming"
	| "done"
	| "error";

export function useImageUpload(caseId: string) {
	const queryClient = useQueryClient();
	const [status, setStatus] = useState<UploadStatus>("idle");
	const [error, setError] = useState<string | null>(null);

	const confirmMutation = useMutation({
		mutationFn: ({
			s3Key,
			fileName,
			contentType,
		}: {
			s3Key: string;
			fileName: string;
			contentType: ContentType;
		}) => confirmUpload(caseId, { s3Key, fileName, contentType }),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["case-images", caseId] });
			setStatus("done");
		},
		onError: () => {
			setStatus("error");
			setError("Failed to confirm upload.");
		},
	});

	const upload = useCallback(
		async (file: File) => {
			setStatus("requesting");
			setError(null);

			if (file.size > maxImageUploadSizeBytes) {
				setStatus("error");
				setError("File is too large. Max 50MB.");
				return;
			}

			try {
				const contentType = file.type as ContentType;
				const { uploadUrl, s3Key } = await requestUploadURL(
					caseId,
					file.name,
					contentType,
					file.size,
				);

				setStatus("uploading");
				await uploadFileToS3(uploadUrl, file);

				setStatus("confirming");
				await confirmMutation.mutateAsync({
					s3Key,
					fileName: file.name,
					contentType,
				});
			} catch (err) {
				setStatus("error");
				setError(err instanceof Error ? err.message : "Upload failed.");
			}
		},
		[caseId, confirmMutation],
	);

	const reset = useCallback(() => {
		setStatus("idle");
		setError(null);
	}, []);

	return { upload, status, error, reset };
}

export function useDownloadURL() {
	return useMutation({
		mutationFn: (imageId: string) => getDownloadURL(imageId),
	});
}
