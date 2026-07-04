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

import { useMemberList, useRemoveMember } from "../hooks/use-members";
import { AddMemberForm } from "./add-member-form";

export function MemberList() {
	const { data: members, isLoading, error, refetch } = useMemberList();
	const removeMutation = useRemoveMember();

	if (isLoading) {
		return (
			<div>
				<PageHeader title="Members" description="Users in this tenant" />
				<div className="mt-6">
					<TableSkeleton rows={5} />
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div>
				<PageHeader title="Members" description="Users in this tenant" />
				<div className="mt-6">
					<ErrorState
						title="Failed to load members"
						message={error.message}
						onRetry={refetch}
					/>
				</div>
			</div>
		);
	}

	const list = members ?? [];

	return (
		<div>
			<PageHeader
				title="Members"
				description="Users in this tenant"
				actions={<AddMemberForm />}
			/>

			<div className="mt-6">
				{list.length === 0 ? (
					<EmptyState
						title="No members"
						description="No users have been added to this tenant yet."
					/>
				) : (
					<Table>
						<TableHead>
							<TableRow>
								<Th>Name</Th>
								<Th>User ID</Th>
								<Th>Role</Th>
								<Th>&nbsp;</Th>
							</TableRow>
						</TableHead>
						<TableBody>
							{list.map((m: { userId: string; name: string; role: string }) => (
								<TableRow key={m.userId}>
									<Td>{m.name}</Td>
									<Td className="font-mono text-xs">{m.userId}</Td>
									<Td>
										<StatusBadge status={m.role} />
									</Td>
									<Td>
										<Button
											variant="danger"
											size="sm"
											disabled={removeMutation.isPending}
											onClick={() => removeMutation.mutate(m.userId)}
										>
											Remove
										</Button>
									</Td>
								</TableRow>
							))}
						</TableBody>
					</Table>
				)}
			</div>
		</div>
	);
}
