import { z } from "zod";

const caseStatusSchema = z.enum(["open", "assigned", "diagnosed", "closed"]);
const severitySchema = z.enum(["low", "medium", "high", "critical"]);

export const caseSummarySchema = z.object({
	id: z.uuid(),
	patientId: z.uuid(),
	doctorId: z.uuid().nullable(),
	status: caseStatusSchema,
	createdAt: z.iso.datetime(),
});

export const symptomSchema = z.object({
	id: z.uuid(),
	description: z.string(),
	severity: severitySchema,
	reportedAt: z.iso.datetime(),
});

export const diagnosisSchema = z.object({
	id: z.uuid(),
	caseId: z.uuid(),
	doctorId: z.uuid(),
	notes: z.string(),
	writtenAt: z.iso.datetime(),
});

export const caseDetailSchema = z.object({
	id: z.uuid(),
	tenantId: z.uuid(),
	patientId: z.uuid(),
	doctorId: z.uuid().nullable(),
	status: caseStatusSchema,
	symptoms: z.array(symptomSchema),
	diagnosis: diagnosisSchema.nullable(),
	createdAt: z.iso.datetime(),
	updatedAt: z.iso.datetime(),
	closedAt: z.iso.datetime().nullable(),
});

export const imageSchema = z.object({
	id: z.uuid(),
	caseId: z.uuid(),
	fileName: z.string(),
	contentType: z.string(),
	uploadedAt: z.iso.datetime(),
});

export const writeDiagnosisInputSchema = z.object({
	notes: z.string().min(1, "Diagnosis notes are required"),
});

export const downloadUrlSchema = z.object({
	downloadUrl: z.string().url(),
	expiresIn: z.number().int().positive(),
});

export type CaseSummary = z.infer<typeof caseSummarySchema>;
export type CaseDetail = z.infer<typeof caseDetailSchema>;
export type Symptom = z.infer<typeof symptomSchema>;
export type Diagnosis = z.infer<typeof diagnosisSchema>;
export type CaseImage = z.infer<typeof imageSchema>;
export type WriteDiagnosisInput = z.infer<typeof writeDiagnosisInputSchema>;
export type DownloadUrl = z.infer<typeof downloadUrlSchema>;
export type CaseStatus = z.infer<typeof caseStatusSchema>;
