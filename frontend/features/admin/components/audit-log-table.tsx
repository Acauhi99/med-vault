"use client";

import { useState } from "react";

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

import { useAuditLogs } from "../hooks/use-audit-logs";

export function AuditLogTable() {
	const [page, setPage] = useState(1);
	const [action, setAction] = useState("");
	const [userId, setUserId] = useState("");
	const [resourceType, setResourceType] = useState("");
	const [resourceId, setResourceId] = useState("");
	const [filters, setFilters] = useState({
		action: "",
		userId: "",
		resourceType: "",
		resourceId: "",
	});

	const { data, isLoading, error, refetch } = useAuditLogs({
		page,
		pageSize: 20,
		action: filters.action || undefined,
		userId: filters.userId || undefined,
		resourceType: filters.resourceType || undefined,
		resourceId: filters.resourceId || undefined,
	});

	function handleFilter() {
		setFilters({
			action,
			userId,
			resourceType,
			resourceId,
		});
		setPage(1);
	}

	if (isLoading) {
		return (
			<div>
				<PageHeader title="Audit Logs" description="System audit trail" />
				<div className="mt-6">
					<TableSkeleton rows={5} />
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div>
				<PageHeader title="Audit Logs" description="System audit trail" />
				<div className="mt-6">
					<ErrorState
						title="Failed to load audit logs"
						message={error.message}
						onRetry={refetch}
					/>
				</div>
			</div>
		);
	}

	const logs = data?.data ?? [];

	return (
		<div>
			<PageHeader title="Audit Logs" description="System audit trail" />

			<div className="mt-6 grid gap-2 md:grid-cols-2 lg:grid-cols-4">
				<input
					type="text"
					value={action}
					onChange={(e) => setAction(e.target.value)}
					placeholder="Filter by action"
					aria-label="Filter by action"
					className="rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
				/>
				<input
					type="text"
					value={userId}
					onChange={(e) => setUserId(e.target.value)}
					placeholder="Filter by user ID"
					aria-label="Filter by user ID"
					className="rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
				/>
				<input
					type="text"
					value={resourceType}
					onChange={(e) => setResourceType(e.target.value)}
					placeholder="Filter by resource type"
					aria-label="Filter by resource type"
					className="rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
				/>
				<input
					type="text"
					value={resourceId}
					onChange={(e) => setResourceId(e.target.value)}
					placeholder="Filter by resource ID"
					aria-label="Filter by resource ID"
					className="rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
				/>
				<div className="lg:col-span-4">
					<Button variant="secondary" size="sm" onClick={handleFilter}>
						Filter
					</Button>
				</div>
			</div>

			<div className="mt-4">
				{logs.length === 0 ? (
					<EmptyState
						title="No audit logs"
						description="No audit events recorded."
					/>
				) : (
					<Table>
						<TableHead>
							<TableRow>
								<Th>Action</Th>
								<Th>Resource</Th>
								<Th>User</Th>
								<Th>IP</Th>
								<Th>Time</Th>
							</TableRow>
						</TableHead>
						<TableBody>
							{logs.map((log) => (
								<TableRow key={log.id}>
									<Td>{log.action}</Td>
									<Td>
										<span className="text-slate-400">{log.resourceType}</span>
										<span className="ml-2 font-mono text-xs text-slate-500">
											{log.resourceId.slice(0, 8)}
										</span>
									</Td>
									<Td className="font-mono text-xs">{log.userId}</Td>
									<Td className="text-slate-400">{log.ipAddress ?? "—"}</Td>
									<Td>{new Date(log.createdAt).toLocaleString()}</Td>
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
