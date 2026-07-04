"use client";

import { useState } from "react";

import { Button } from "@/shared/components/button";
import { Dialog } from "@/shared/components/dialog";

import { useAssignDoctor } from "../hooks/use-all-cases";

type AssignDoctorDialogProps = {
	open: boolean;
	onClose: () => void;
	caseId: string;
};

export function AssignDoctorDialog({
	open,
	onClose,
	caseId,
}: AssignDoctorDialogProps) {
	const [doctorId, setDoctorId] = useState("");
	const assignMutation = useAssignDoctor(caseId);

	function handleSubmit(e: React.FormEvent) {
		e.preventDefault();
		assignMutation.mutate(
			{ doctorId },
			{
				onSuccess: () => {
					setDoctorId("");
					onClose();
				},
			},
		);
	}

	return (
		<Dialog open={open} onClose={onClose} title="Assign Doctor">
			<form onSubmit={handleSubmit} className="space-y-4">
				<div>
					<label htmlFor="doctor-id" className="block text-sm text-slate-400">
						Doctor ID
					</label>
					<input
						id="doctor-id"
						type="text"
						value={doctorId}
						onChange={(e) => setDoctorId(e.target.value)}
						placeholder="Enter doctor UUID"
						className="mt-1 w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
					/>
				</div>
				{assignMutation.error && (
					<p className="text-sm text-red-400">{assignMutation.error.message}</p>
				)}
				<div className="flex justify-end gap-2">
					<Button variant="ghost" onClick={onClose}>
						Cancel
					</Button>
					<Button
						type="submit"
						disabled={!doctorId}
						busy={assignMutation.isPending}
					>
						Assign
					</Button>
				</div>
			</form>
		</Dialog>
	);
}
