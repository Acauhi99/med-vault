"use client";

import { StatusBadge } from "@/shared/components/badge";
import { Card, CardHeader, CardTitle } from "@/shared/components/card";
import { PageHeader } from "@/shared/components/page-header";
import type { CaseDetail as CaseDetailType } from "../schemas/cases";
import { DiagnosisForm } from "./diagnosis-form";
import { ImageViewer } from "./image-viewer";

type Props = {
	caseData: CaseDetailType;
	onBack: () => void;
};

export function CaseDetail({ caseData, onBack }: Props) {
	return (
		<div className="min-h-screen bg-slate-950 px-4 py-8 text-slate-50 sm:px-6 lg:px-8">
			<div className="mx-auto max-w-4xl">
				<PageHeader
					title={`Case ${caseData.id.slice(0, 8)}`}
					description={`Patient ${caseData.patientId.slice(0, 8)} · Created ${new Date(caseData.createdAt).toLocaleDateString()}`}
					actions={
						<button
							type="button"
							onClick={onBack}
							className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-300 transition hover:bg-white/10"
						>
							Back to cases
						</button>
					}
				/>

				<div className="mt-6 grid gap-6">
					<Card>
						<CardHeader>
							<CardTitle>Status</CardTitle>
						</CardHeader>
						<div className="px-6 pb-6">
							<StatusBadge status={caseData.status} />
						</div>
					</Card>

					{caseData.symptoms.length > 0 && (
						<Card>
							<CardHeader>
								<CardTitle>Symptoms</CardTitle>
							</CardHeader>
							<div className="space-y-3 px-6 pb-6">
								{caseData.symptoms.map((symptom) => (
									<div
										key={symptom.id}
										className="rounded-xl border border-white/10 bg-white/5 p-4"
									>
										<p className="text-sm text-slate-200">
											{symptom.description}
										</p>
										<div className="mt-2 flex items-center gap-3 text-xs text-slate-400">
											<span className="capitalize">{symptom.severity}</span>
											<span>·</span>
											<span>
												{new Date(symptom.reportedAt).toLocaleDateString()}
											</span>
										</div>
									</div>
								))}
							</div>
						</Card>
					)}

					<Card>
						<CardHeader>
							<CardTitle>Images</CardTitle>
						</CardHeader>
						<div className="px-6 pb-6">
							<ImageViewer caseId={caseData.id} />
						</div>
					</Card>

					{caseData.diagnosis ? (
						<Card>
							<CardHeader>
								<CardTitle>Diagnosis</CardTitle>
							</CardHeader>
							<div className="px-6 pb-6">
								<p className="text-sm leading-6 text-slate-200">
									{caseData.diagnosis.notes}
								</p>
								<p className="mt-3 text-xs text-slate-400">
									Written{" "}
									{new Date(caseData.diagnosis.writtenAt).toLocaleDateString()}
								</p>
							</div>
						</Card>
					) : (
						<Card>
							<CardHeader>
								<CardTitle>Write Diagnosis</CardTitle>
							</CardHeader>
							<div className="px-6 pb-6">
								<DiagnosisForm caseId={caseData.id} />
							</div>
						</Card>
					)}
				</div>
			</div>
		</div>
	);
}
