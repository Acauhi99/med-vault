"use client";

import { useState } from "react";

import { Badge } from "@/shared/components/badge";
import { Button } from "@/shared/components/button";
import { Card, CardHeader, CardTitle } from "@/shared/components/card";
import { Skeleton } from "@/shared/components/skeleton";
import { ErrorState } from "@/shared/components/states";

import { useCaseDetail } from "../hooks/use-all-cases";
import { AssignDoctorDialog } from "./assign-doctor-dialog";
import { CloseCaseButton } from "./close-case-button";

type CaseDetailProps = {
	caseId: string;
	onBack: () => void;
};

export function CaseDetail({ caseId, onBack }: CaseDetailProps) {
	const [assignOpen, setAssignOpen] = useState(false);
	const { data: caseData, isLoading, error, refetch } = useCaseDetail(caseId);

	if (isLoading) {
		return (
			<div className="space-y-4">
				<Skeleton className="h-8 w-48" />
				<Skeleton className="h-40 w-full" />
			</div>
		);
	}

	if (error || !caseData) {
		return (
			<ErrorState
				title="Failed to load case"
				message={error?.message ?? "Case not found."}
				onRetry={refetch}
			/>
		);
	}

	return (
		<div className="space-y-6">
			<div className="flex items-center gap-4">
				<Button variant="ghost" size="sm" onClick={onBack}>
					<svg
						className="h-4 w-4"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
					>
						<path
							fillRule="evenodd"
							d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
							clipRule="evenodd"
						/>
					</svg>
					Back
				</Button>
				<div>
					<h2 className="text-xl font-semibold text-white">
						Case {caseData.id.slice(0, 8)}
					</h2>
					<Badge>{caseData.status}</Badge>
				</div>
			</div>

			<Card>
				<CardHeader>
					<CardTitle>Details</CardTitle>
				</CardHeader>
				<div className="grid grid-cols-2 gap-4 text-sm">
					<div>
						<span className="text-slate-400">Patient</span>
						<p className="mt-1 font-mono text-white">{caseData.patientId}</p>
					</div>
					<div>
						<span className="text-slate-400">Doctor</span>
						<p className="mt-1 font-mono text-white">
							{caseData.doctorId ?? (
								<span className="text-slate-500">Unassigned</span>
							)}
						</p>
					</div>
					<div>
						<span className="text-slate-400">Created</span>
						<p className="mt-1 text-white">
							{new Date(caseData.createdAt).toLocaleString()}
						</p>
					</div>
					<div>
						<span className="text-slate-400">Updated</span>
						<p className="mt-1 text-white">
							{new Date(caseData.updatedAt).toLocaleString()}
						</p>
					</div>
					{caseData.closedAt && (
						<div>
							<span className="text-slate-400">Closed</span>
							<p className="mt-1 text-white">
								{new Date(caseData.closedAt).toLocaleString()}
							</p>
						</div>
					)}
				</div>
			</Card>

			{caseData.symptoms.length > 0 && (
				<Card>
					<CardHeader>
						<CardTitle>Symptoms</CardTitle>
					</CardHeader>
					<div className="space-y-2">
						{caseData.symptoms.map((s) => (
							<div
								key={s.id}
								className="rounded-2xl border border-white/10 bg-white/[0.03] p-3"
							>
								<div className="flex items-center gap-2">
									{s.severity && <Badge variant="warning">{s.severity}</Badge>}
									<span className="text-xs text-slate-500">
										{new Date(s.reportedAt).toLocaleDateString()}
									</span>
								</div>
								{s.description && (
									<p className="mt-1 text-sm text-slate-400">{s.description}</p>
								)}
							</div>
						))}
					</div>
				</Card>
			)}

			{caseData.diagnosis && (
				<Card>
					<CardHeader>
						<CardTitle>Diagnosis</CardTitle>
					</CardHeader>
					<div className="space-y-2 text-sm">
						<div>
							<span className="text-slate-400">Doctor</span>
							<p className="mt-1 font-mono text-white">
								{caseData.diagnosis.doctorId}
							</p>
						</div>
						{caseData.diagnosis.notes && (
							<div>
								<span className="text-slate-400">Notes</span>
								<p className="mt-1 text-white">{caseData.diagnosis.notes}</p>
							</div>
						)}
					</div>
				</Card>
			)}

			<div className="flex gap-3">
				{!caseData.doctorId && (
					<Button onClick={() => setAssignOpen(true)}>Assign Doctor</Button>
				)}
				{caseData.status === "diagnosed" && <CloseCaseButton caseId={caseId} />}
			</div>

			<AssignDoctorDialog
				open={assignOpen}
				onClose={() => setAssignOpen(false)}
				caseId={caseId}
			/>
		</div>
	);
}
