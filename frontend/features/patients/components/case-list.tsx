"use client";

import { StatusBadge } from "@/shared/components/badge";
import { Button } from "@/shared/components/button";
import { PageHeader } from "@/shared/components/page-header";
import { TableSkeleton } from "@/shared/components/skeleton";
import { EmptyState, ErrorState } from "@/shared/components/states";
import {
	Table,
	TableBody,
	TableHead,
	TableRow,
	Td,
	Th,
} from "@/shared/components/table";

import { useCases } from "../hooks/use-cases";
import type { CaseSummary } from "../schemas/cases";

type CaseListProps = {
	onSelectCase: (id: string) => void;
	onCreateCase: () => void;
};

export function CaseList({ onSelectCase, onCreateCase }: CaseListProps) {
	const { data, isLoading, isError, error, refetch } = useCases();

	return (
		<div className="space-y-6">
			<PageHeader
				title="My Cases"
				description="View and manage your medical cases."
				actions={
					<Button variant="primary" onClick={onCreateCase}>
						New case
					</Button>
				}
			/>

			{isLoading && <TableSkeleton rows={5} />}

			{isError && (
				<ErrorState
					title="Failed to load cases"
					message={
						error instanceof Error ? error.message : "Unable to load cases."
					}
					onRetry={refetch}
				/>
			)}

			{!isLoading && !isError && data && data.cases.length === 0 && (
				<EmptyState
					title="No cases yet"
					description="Create your first case to get started."
					action={
						<Button variant="primary" onClick={onCreateCase}>
							New case
						</Button>
					}
				/>
			)}

			{!isLoading && !isError && data && data.cases.length > 0 && (
				<Table>
					<TableHead>
						<Th>Case ID</Th>
						<Th>Status</Th>
						<Th>Doctor</Th>
						<Th>Created</Th>
					</TableHead>
					<TableBody>
						{data.cases.map((c: CaseSummary) => (
							<TableRow key={c.id}>
								<Td>
									<button
										type="button"
										onClick={() => onSelectCase(c.id)}
										className="font-medium text-sky-400 transition hover:text-sky-300 hover:underline"
									>
										{c.id.slice(0, 8)}
									</button>
								</Td>
								<Td>
									<StatusBadge status={c.status} />
								</Td>
								<Td className="text-slate-400">
									{c.doctorId ? c.doctorId.slice(0, 8) : "—"}
								</Td>
								<Td className="text-slate-400">
									{new Date(c.createdAt).toLocaleDateString()}
								</Td>
							</TableRow>
						))}
					</TableBody>
				</Table>
			)}
		</div>
	);
}
