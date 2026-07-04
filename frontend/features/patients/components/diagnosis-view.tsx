"use client";

import { Card, CardHeader, CardTitle } from "@/shared/components/card";
import type { DiagnosisResponse } from "../schemas/cases";

type DiagnosisViewProps = {
	diagnosis: DiagnosisResponse;
};

export function DiagnosisView({ diagnosis }: DiagnosisViewProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Diagnosis</CardTitle>
			</CardHeader>
			<div className="space-y-3">
				<div className="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
					<p className="text-sm leading-relaxed text-slate-200">
						{diagnosis.notes}
					</p>
				</div>
				<div className="flex items-center gap-4 text-xs text-slate-400">
					<span>By doctor {diagnosis.doctorId.slice(0, 8)}</span>
					<span>&middot;</span>
					<span>{new Date(diagnosis.writtenAt).toLocaleString()}</span>
				</div>
			</div>
		</Card>
	);
}
