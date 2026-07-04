"use client";

import { useRef } from "react";

import { Button } from "@/shared/components/button";

import { useImageUpload } from "../hooks/use-image-upload";

type ImageUploadProps = {
	caseId: string;
};

const ACCEPTED_TYPES = ["image/jpeg", "image/png", "image/dicom"] as const;

export function ImageUpload({ caseId }: ImageUploadProps) {
	const { upload, status, error, reset } = useImageUpload(caseId);
	const inputRef = useRef<HTMLInputElement>(null);

	async function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
		const file = event.target.files?.[0];
		if (!file) return;
		await upload(file);
		if (inputRef.current) {
			inputRef.current.value = "";
		}
	}

	return (
		<div className="mb-4 space-y-3">
			<div className="flex items-center gap-3">
				<Button
					variant="secondary"
					onClick={() => inputRef.current?.click()}
					disabled={
						status === "requesting" ||
						status === "uploading" ||
						status === "confirming"
					}
				>
					{status === "requesting" ||
					status === "uploading" ||
					status === "confirming"
						? "Uploading..."
						: "Upload image"}
				</Button>

				<input
					ref={inputRef}
					type="file"
					accept={ACCEPTED_TYPES.join(",")}
					onChange={handleChange}
					className="hidden"
					aria-label="Upload image file"
				/>

				{status === "done" && (
					<span className="text-sm text-emerald-400">
						Uploaded successfully.
					</span>
				)}

				{status === "error" && (
					<>
						<span className="text-sm text-red-400">{error}</span>
						<Button variant="ghost" onClick={reset}>
							Dismiss
						</Button>
					</>
				)}
			</div>

			<p className="text-xs text-slate-500">
				Accepted: JPEG, PNG, DICOM. Max 50MB.
			</p>
		</div>
	);
}
