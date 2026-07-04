"use client";

import type { SubmitEvent } from "react";
import { useState } from "react";

import { Button } from "@/shared/components/button";
import { Card, CardHeader, CardTitle } from "@/shared/components/card";

import { useAddSymptom } from "../hooks/use-cases";
import type { Severity } from "../schemas/cases";

type AddSymptomFormProps = {
	caseId: string;
};

export function AddSymptomForm({ caseId }: AddSymptomFormProps) {
	const [description, setDescription] = useState("");
	const [severity, setSeverity] = useState<Severity>("low");
	const [notice, setNotice] = useState<string | null>(null);
	const addMutation = useAddSymptom(caseId);

	async function onSubmit(event: SubmitEvent<HTMLFormElement>) {
		event.preventDefault();
		setNotice(null);

		if (!description.trim()) {
			setNotice("Description is required.");
			return;
		}

		try {
			await addMutation.mutateAsync({
				description: description.trim(),
				severity,
			});
			setDescription("");
			setSeverity("low");
			setNotice("Symptom added.");
		} catch (err) {
			setNotice(err instanceof Error ? err.message : "Unable to add symptom.");
		}
	}

	return (
		<Card>
			<CardHeader>
				<CardTitle>Add Symptom</CardTitle>
			</CardHeader>

			{notice && (
				<p
					role="status"
					className="mb-4 rounded-2xl border border-sky-400/30 bg-sky-400/10 px-4 py-3 text-sm text-sky-100"
				>
					{notice}
				</p>
			)}

			<form className="space-y-4" onSubmit={onSubmit}>
				<label
					className="grid gap-2 text-sm text-slate-200"
					htmlFor="symptom-desc"
				>
					Description
					<input
						id="symptom-desc"
						value={description}
						onChange={(e) => setDescription(e.target.value)}
						placeholder="Describe the symptom"
						className="h-10 rounded-2xl border border-white/10 bg-slate-950 px-4 text-slate-50 placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
					/>
				</label>

				<label
					className="grid gap-2 text-sm text-slate-200"
					htmlFor="symptom-severity"
				>
					Severity
					<select
						id="symptom-severity"
						value={severity}
						onChange={(e) => setSeverity(e.target.value as Severity)}
						className="h-10 rounded-2xl border border-white/10 bg-slate-950 px-4 text-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
					>
						<option value="low">Low</option>
						<option value="medium">Medium</option>
						<option value="high">High</option>
						<option value="critical">Critical</option>
					</select>
				</label>

				<Button type="submit" busy={addMutation.isPending}>
					Add symptom
				</Button>
			</form>
		</Card>
	);
}
