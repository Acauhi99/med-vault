"use client";

import { StatusBadge } from "@/shared/components/badge";
import { Button } from "@/shared/components/button";
import { Card, CardHeader, CardTitle } from "@/shared/components/card";
import { PageHeader } from "@/shared/components/page-header";
import { Skeleton } from "@/shared/components/skeleton";
import { EmptyState, ErrorState } from "@/shared/components/states";

import { useCase } from "../hooks/use-cases";
import type { SymptomResponse } from "../schemas/cases";
import { AddSymptomForm } from "./add-symptom-form";
import { DiagnosisView } from "./diagnosis-view";
import { ImageList } from "./image-list";
import { ImageUpload } from "./image-upload";

type CaseDetailProps = {
	caseId: string;
	onBack: () => void;
};

export function CaseDetail({ caseId, onBack }: CaseDetailProps) {
	const {
		data: caseData,
		isLoading,
		isError,
		error,
		refetch,
	} = useCase(caseId);

	return (
		<div className="space-y-6">
			<PageHeader
				title={
					isLoading
						? "Loading case..."
						: `Case ${caseData?.id.slice(0, 8) ?? ""}`
				}
				description={caseData ? `Status: ${caseData.status}` : undefined}
				actions={
					<Button variant="ghost" onClick={onBack}>
						Back to cases
					</Button>
				}
			/>

			{isLoading && (
				<div className="space-y-4">
					<Skeleton className="h-24 w-full" />
					<Skeleton className="h-32 w-full" />
					<Skeleton className="h-24 w-full" />
				</div>
			)}

			{isError && (
				<ErrorState
					title="Failed to load case"
					message={
						error instanceof Error
							? error.message
							: "Unable to load case details."
					}
					onRetry={refetch}
				/>
			)}

			{!isLoading && !isError && caseData && (
				<>
					<div className="grid gap-4 sm:grid-cols-3">
						<Card>
							<CardHeader>
								<CardTitle>Status</CardTitle>
							</CardHeader>
							<StatusBadge status={caseData.status} />
						</Card>
						<Card>
							<CardHeader>
								<CardTitle>Created</CardTitle>
							</CardHeader>
							<p className="text-sm text-slate-300">
								{new Date(caseData.createdAt).toLocaleString()}
							</p>
						</Card>
						<Card>
							<CardHeader>
								<CardTitle>Updated</CardTitle>
							</CardHeader>
							<p className="text-sm text-slate-300">
								{new Date(caseData.updatedAt).toLocaleString()}
							</p>
						</Card>
					</div>

					<Card>
						<CardHeader>
							<CardTitle>Symptoms ({caseData.symptoms.length})</CardTitle>
						</CardHeader>
						{caseData.symptoms.length === 0 ? (
							<p className="text-sm text-slate-400">
								No symptoms reported yet.
							</p>
						) : (
							<div className="space-y-3">
								{caseData.symptoms.map((s: SymptomResponse) => (
									<div
										key={s.id}
										className="flex items-start justify-between gap-4 rounded-2xl border border-white/10 bg-white/[0.03] p-4"
									>
										<div className="space-y-1">
											<p className="text-sm text-slate-200">{s.description}</p>
											<p className="text-xs text-slate-400">
												{new Date(s.reportedAt).toLocaleString()}
											</p>
										</div>
										<StatusBadge status={s.severity} />
									</div>
								))}
							</div>
						)}
					</Card>

					{caseData.diagnosis ? (
						<DiagnosisView diagnosis={caseData.diagnosis} />
					) : (
						<Card>
							<CardHeader>
								<CardTitle>Diagnosis</CardTitle>
							</CardHeader>
							<EmptyState
								title="No diagnosis yet"
								description="A doctor has not yet written a diagnosis for this case."
							/>
						</Card>
					)}

					<AddSymptomForm caseId={caseId} />

					<Card>
						<CardHeader>
							<CardTitle>Images</CardTitle>
						</CardHeader>
						<ImageUpload caseId={caseId} />
						<ImageList caseId={caseId} />
					</Card>
				</>
			)}
		</div>
	);
}
