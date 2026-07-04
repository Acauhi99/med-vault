"use client";

import { useState } from "react";
import { useWriteDiagnosis } from "../hooks/use-diagnosis";
import type { WriteDiagnosisInput } from "../schemas/cases";

type Props = {
	caseId: string;
};

export function DiagnosisForm({ caseId }: Props) {
	const [notes, setNotes] = useState("");
	const [error, setError] = useState<string | null>(null);
	const diagnosisMutation = useWriteDiagnosis(caseId);

	async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault();
		setError(null);

		try {
			const input: WriteDiagnosisInput = { notes: notes.trim() };
			if (!input.notes) {
				setError("Diagnosis notes are required.");
				return;
			}
			await diagnosisMutation.mutateAsync(input);
			setNotes("");
		} catch (err) {
			setError(
				err instanceof Error ? err.message : "Unable to write diagnosis.",
			);
		}
	}

	if (diagnosisMutation.isSuccess) {
		return (
			<div className="rounded-2xl border border-emerald-400/30 bg-emerald-400/10 p-4 text-sm text-emerald-200">
				Diagnosis submitted successfully.
			</div>
		);
	}

	return (
		<form className="grid gap-4" onSubmit={handleSubmit}>
			<label
				className="grid gap-2 text-sm text-slate-200"
				htmlFor="diagnosis-notes"
			>
				Notes
				<textarea
					id="diagnosis-notes"
					rows={5}
					value={notes}
					onChange={(e) => setNotes(e.target.value)}
					placeholder="Enter your diagnosis notes..."
					className="rounded-2xl border border-white/10 bg-slate-950 px-4 py-3 text-sm text-slate-50 placeholder:text-slate-500 transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
				/>
			</label>

			{error && (
				<p className="rounded-xl border border-red-400/30 bg-red-400/10 px-4 py-2 text-sm text-red-200">
					{error}
				</p>
			)}

			<button
				type="submit"
				disabled={diagnosisMutation.isPending}
				aria-busy={diagnosisMutation.isPending}
				className="inline-flex h-12 items-center justify-center rounded-2xl bg-sky-400 px-4 text-sm font-semibold text-slate-950 transition hover:bg-sky-300 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-300"
			>
				{diagnosisMutation.isPending ? "Submitting..." : "Submit Diagnosis"}
			</button>
		</form>
	);
}
