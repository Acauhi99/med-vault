"use client";

import { useState } from "react";

import { Button } from "@/shared/components/button";
import { Dialog } from "@/shared/components/dialog";

import { useAddMember } from "../hooks/use-members";

const ROLES = ["patient", "doctor", "administrator"] as const;

export function AddMemberForm() {
	const [open, setOpen] = useState(false);
	const [userId, setUserId] = useState("");
	const [role, setRole] = useState<"patient" | "doctor" | "administrator">(
		"patient",
	);
	const addMutation = useAddMember();

	function handleSubmit(e: React.FormEvent) {
		e.preventDefault();
		addMutation.mutate(
			{ userId, role },
			{
				onSuccess: () => {
					setUserId("");
					setRole("patient");
					setOpen(false);
				},
			},
		);
	}

	return (
		<>
			<Button onClick={() => setOpen(true)}>Add Member</Button>
			<Dialog
				open={open}
				onClose={() => setOpen(false)}
				title="Add Tenant Member"
			>
				<form onSubmit={handleSubmit} className="space-y-4">
					<div>
						<label
							htmlFor="member-user-id"
							className="block text-sm text-slate-400"
						>
							User ID
						</label>
						<input
							id="member-user-id"
							type="text"
							value={userId}
							onChange={(e) => setUserId(e.target.value)}
							placeholder="Enter user UUID"
							className="mt-1 w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
						/>
					</div>
					<div>
						<label
							htmlFor="member-role"
							className="block text-sm text-slate-400"
						>
							Role
						</label>
						<select
							id="member-role"
							value={role}
							onChange={(e) => setRole(e.target.value as typeof role)}
							className="mt-1 w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
						>
							{ROLES.map((r) => (
								<option key={r} value={r}>
									{r}
								</option>
							))}
						</select>
					</div>
					{addMutation.error && (
						<p className="text-sm text-red-400">{addMutation.error.message}</p>
					)}
					<div className="flex justify-end gap-2">
						<Button variant="ghost" onClick={() => setOpen(false)}>
							Cancel
						</Button>
						<Button
							type="submit"
							disabled={!userId}
							busy={addMutation.isPending}
						>
							Add
						</Button>
					</div>
				</form>
			</Dialog>
		</>
	);
}
