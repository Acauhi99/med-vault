"use client";

import { useState } from "react";
import { StatusBadge } from "@/shared/components/badge";
import { PageHeader } from "@/shared/components/page-header";
import { TableSkeleton } from "@/shared/components/skeleton";
import { EmptyState, ErrorState } from "@/shared/components/states";
import { useAssignedCases, useCaseDetail } from "../hooks/use-assigned-cases";
import { CaseDetail } from "./case-detail";

export function CaseList() {
	const [selectedCaseId, setSelectedCaseId] = useState<string | null>(null);
	const { data, isLoading, isError, error, refetch } = useAssignedCases();
	const detailQuery = useCaseDetail(selectedCaseId);

	if (selectedCaseId && detailQuery.data) {
		return (
			<CaseDetail
				caseData={detailQuery.data}
				onBack={() => setSelectedCaseId(null)}
			/>
		);
	}

	if (selectedCaseId && detailQuery.isError) {
		return (
			<div className="min-h-screen bg-slate-950 px-4 py-8 text-slate-50 sm:px-6 lg:px-8">
				<div className="mx-auto max-w-4xl">
					<ErrorState
						title="Failed to load case details"
						message={
							detailQuery.error instanceof Error
								? detailQuery.error.message
								: "Unknown error"
						}
						onRetry={() => detailQuery.refetch()}
						action={
							<button
								type="button"
								onClick={() => setSelectedCaseId(null)}
								className="inline-flex h-10 items-center justify-center rounded-2xl border border-white/10 bg-white/5 px-4 text-sm font-medium text-slate-300 transition hover:bg-white/10"
							>
								Back to cases
							</button>
						}
					/>
				</div>
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-slate-950 px-4 py-8 text-slate-50 sm:px-6 lg:px-8">
			<div className="mx-auto max-w-6xl">
				<PageHeader
					title="Assigned Cases"
					description="Cases assigned to you for diagnosis."
				/>

				{isLoading && <TableSkeleton rows={5} />}

				{isError && (
					<ErrorState
						title="Failed to load cases"
						message={error instanceof Error ? error.message : "Unknown error"}
						onRetry={refetch}
					/>
				)}

				{data && data.cases.length === 0 && (
					<EmptyState
						title="No cases assigned"
						description="You have no cases assigned to you yet."
					/>
				)}

				{data && data.cases.length > 0 && (
					<div className="mt-6 overflow-hidden rounded-2xl border border-white/10 bg-white/5">
						<table className="w-full text-sm">
							<thead>
								<tr className="border-b border-white/10 text-left text-xs uppercase tracking-wider text-slate-400">
									<th className="px-4 py-3">Case ID</th>
									<th className="px-4 py-3">Patient</th>
									<th className="px-4 py-3">Status</th>
									<th className="px-4 py-3">Created</th>
									<th className="px-4 py-3" />
								</tr>
							</thead>
							<tbody className="divide-y divide-white/5">
								{data.cases.map((c) => (
									<tr key={c.id} className="transition hover:bg-white/5">
										<td className="px-4 py-3 font-mono text-xs text-slate-300">
											{c.id.slice(0, 8)}
										</td>
										<td className="px-4 py-3 font-mono text-xs text-slate-300">
											{c.patientId.slice(0, 8)}
										</td>
										<td className="px-4 py-3">
											<StatusBadge status={c.status} />
										</td>
										<td className="px-4 py-3 text-xs text-slate-400">
											{new Date(c.createdAt).toLocaleDateString()}
										</td>
										<td className="px-4 py-3 text-right">
											<button
												type="button"
												onClick={() => setSelectedCaseId(c.id)}
												className="rounded-xl bg-sky-400/10 px-3 py-1.5 text-xs font-medium text-sky-300 transition hover:bg-sky-400/20"
											>
												View
											</button>
										</td>
									</tr>
								))}
							</tbody>
					</table>
				</div>
			)}

				{detailQuery.isLoading && selectedCaseId && (
					<div className="mt-6">
						<TableSkeleton rows={3} />
					</div>
				)}
			</div>
		</div>
	);
}
