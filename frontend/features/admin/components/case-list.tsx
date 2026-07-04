"use client";

import { useState } from "react";
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

import { useCaseList } from "../hooks/use-all-cases";
import type { CaseStatus } from "../schemas/cases";

const STATUS_FILTERS: Array<{ label: string; value: CaseStatus | "all" }> = [
	{ label: "All", value: "all" },
	{ label: "Open", value: "open" },
	{ label: "Assigned", value: "assigned" },
	{ label: "Diagnosed", value: "diagnosed" },
	{ label: "Closed", value: "closed" },
];

type CaseListProps = {
	onSelectCase: (id: string) => void;
};

export function CaseList({ onSelectCase }: CaseListProps) {
	const [page, setPage] = useState(1);
	const [statusFilter, setStatusFilter] = useState<CaseStatus | "all">("all");

	const { data, isLoading, error, refetch } = useCaseList({
		page,
		pageSize: 20,
	});

	if (isLoading) {
		return (
			<div>
				<PageHeader title="Cases" description="All cases across the tenant" />
				<div className="mt-6">
					<TableSkeleton rows={5} />
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div>
				<PageHeader title="Cases" description="All cases across the tenant" />
				<div className="mt-6">
					<ErrorState
						title="Failed to load cases"
						message={error.message}
						onRetry={refetch}
					/>
				</div>
			</div>
		);
	}

	const cases = data?.data ?? [];
	const filtered =
		statusFilter === "all"
			? cases
			: cases.filter((c) => c.status === statusFilter);

	return (
		<div>
			<PageHeader
				title="Cases"
				description="All cases across the tenant"
				actions={
					<div className="flex gap-1">
						{STATUS_FILTERS.map((f) => (
							<Button
								key={f.value}
								variant={statusFilter === f.value ? "primary" : "ghost"}
								size="sm"
								onClick={() => {
									setStatusFilter(f.value);
									setPage(1);
								}}
							>
								{f.label}
							</Button>
						))}
					</div>
				}
			/>

			<div className="mt-6">
				{filtered.length === 0 ? (
					<EmptyState
						title="No cases"
						description={
							statusFilter === "all"
								? "No cases exist yet."
								: `No ${statusFilter} cases.`
						}
					/>
				) : (
					<Table>
						<TableHead>
							<TableRow>
								<Th>Patient</Th>
								<Th>Doctor</Th>
								<Th>Status</Th>
								<Th>Created</Th>
								<Th>&nbsp;</Th>
							</TableRow>
						</TableHead>
						<TableBody>
							{filtered.map((c) => (
								<TableRow key={c.id}>
									<Td>{c.patientId}</Td>
									<Td>
										{c.doctorId ? (
											c.doctorId
										) : (
											<span className="text-slate-500">Unassigned</span>
										)}
									</Td>
									<Td>
										<StatusBadge status={c.status} />
									</Td>
									<Td>{new Date(c.createdAt).toLocaleDateString()}</Td>
									<Td>
										<Button
											variant="ghost"
											size="sm"
											onClick={() => onSelectCase(c.id)}
										>
											View
										</Button>
									</Td>
								</TableRow>
							))}
						</TableBody>
					</Table>
				)}
			</div>

			{data?.meta?.total != null && data.meta.total > 20 && (
				<div className="mt-4 flex items-center justify-between">
					<span className="text-sm text-slate-400">
						Page {data.meta.page ?? 1} of {Math.ceil(data.meta.total / 20)}
					</span>
					<div className="flex gap-2">
						<Button
							variant="secondary"
							size="sm"
							disabled={page <= 1}
							onClick={() => setPage((p) => p - 1)}
						>
							Previous
						</Button>
						<Button
							variant="secondary"
							size="sm"
							disabled={page * 20 >= data.meta.total}
							onClick={() => setPage((p) => p + 1)}
						>
							Next
						</Button>
					</div>
				</div>
			)}
		</div>
	);
}
