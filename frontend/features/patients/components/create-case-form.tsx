"use client";

import type { SubmitEvent } from "react";
import { useState } from "react";

import { Button } from "@/shared/components/button";
import { Card, CardHeader, CardTitle } from "@/shared/components/card";

import { useCreateCase } from "../hooks/use-cases";
import { type AddSymptomInput, createCaseInputSchema } from "../schemas/cases";

type CreateCaseFormProps = {
	onCreated: () => void;
	onCancel: () => void;
};

const emptySymptom: AddSymptomInput = { description: "", severity: "low" };

export function CreateCaseForm({ onCreated, onCancel }: CreateCaseFormProps) {
	const [symptoms, setSymptoms] = useState<AddSymptomInput[]>([
		{ ...emptySymptom },
	]);
	const [notice, setNotice] = useState<string | null>(null);
	const createMutation = useCreateCase();

	function addSymptomRow() {
		setSymptoms((prev) => [...prev, { ...emptySymptom }]);
	}

	function removeSymptomRow(index: number) {
		setSymptoms((prev) => prev.filter((_, i) => i !== index));
	}

	function updateSymptom(
		index: number,
		field: keyof AddSymptomInput,
		value: string,
	) {
		setSymptoms((prev) =>
			prev.map((s, i) => (i === index ? { ...s, [field]: value } : s)),
		);
	}

	async function onSubmit(event: SubmitEvent<HTMLFormElement>) {
		event.preventDefault();
		setNotice(null);

		const parsed = createCaseInputSchema.safeParse({ symptoms });
		if (!parsed.success) {
			setNotice(parsed.error.issues[0]?.message ?? "Fix the form errors.");
			return;
		}

		try {
			await createMutation.mutateAsync(parsed.data.symptoms);
			onCreated();
		} catch (err) {
			setNotice(err instanceof Error ? err.message : "Unable to create case.");
		}
	}

	return (
		<Card>
			<CardHeader>
				<CardTitle>New Case</CardTitle>
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
				{symptoms.map((symptom, index) => (
					<div
						key={`symptom-${String(index)}`}
						className="space-y-3 rounded-2xl border border-white/10 bg-white/[0.03] p-4"
					>
						<div className="flex items-center justify-between">
							<span className="text-xs font-medium uppercase tracking-wider text-slate-400">
								Symptom {index + 1}
							</span>
							{symptoms.length > 1 && (
								<button
									type="button"
									onClick={() => removeSymptomRow(index)}
									className="text-xs text-red-400 transition hover:text-red-300"
								>
									Remove
								</button>
							)}
						</div>

						<label
							className="grid gap-2 text-sm text-slate-200"
							htmlFor={`desc-${String(index)}`}
						>
							Description
							<input
								id={`desc-${String(index)}`}
								value={symptom.description}
								onChange={(e) =>
									updateSymptom(index, "description", e.target.value)
								}
								placeholder="Describe the symptom"
								className="h-10 rounded-2xl border border-white/10 bg-slate-950 px-4 text-slate-50 placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
							/>
						</label>

						<label
							className="grid gap-2 text-sm text-slate-200"
							htmlFor={`severity-${String(index)}`}
						>
							Severity
							<select
								id={`severity-${String(index)}`}
								value={symptom.severity}
								onChange={(e) =>
									updateSymptom(index, "severity", e.target.value)
								}
								className="h-10 rounded-2xl border border-white/10 bg-slate-950 px-4 text-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
							>
								<option value="low">Low</option>
								<option value="medium">Medium</option>
								<option value="high">High</option>
								<option value="critical">Critical</option>
							</select>
						</label>
					</div>
				))}

				<Button type="button" variant="secondary" onClick={addSymptomRow}>
					+ Add symptom
				</Button>

				<div className="flex gap-3 pt-2">
					<Button type="submit" busy={createMutation.isPending}>
						Create case
					</Button>
					<Button type="button" variant="ghost" onClick={onCancel}>
						Cancel
					</Button>
				</div>
			</form>
		</Card>
	);
}
