import { z } from "zod";

const caseStatusSchema = z.enum(["open", "assigned", "diagnosed", "closed"]);
export const severitySchema = z.enum(["low", "medium", "high", "critical"]);
export const contentTypeSchema = z.enum([
	"image/jpeg",
	"image/png",
	"image/dicom",
]);

export const addSymptomInputSchema = z.object({
	description: z.string().min(1, "Description is required").trim(),
	severity: severitySchema,
});

export const createCaseInputSchema = z.object({
	symptoms: z.array(addSymptomInputSchema).min(1, "Add at least one symptom"),
});

export const caseSummarySchema = z.object({
	id: z.uuid(),
	patientId: z.uuid(),
	doctorId: z.uuid().nullable(),
	status: caseStatusSchema,
	createdAt: z.iso.datetime(),
});

export const symptomResponseSchema = z.object({
	id: z.uuid(),
	description: z.string(),
	severity: severitySchema,
	reportedAt: z.iso.datetime(),
});

export const diagnosisResponseSchema = z.object({
	id: z.uuid(),
	caseId: z.uuid(),
	doctorId: z.uuid(),
	notes: z.string(),
	writtenAt: z.iso.datetime(),
});

export const caseResponseSchema = z.object({
	id: z.uuid(),
	tenantId: z.uuid(),
	patientId: z.uuid(),
	doctorId: z.uuid().nullable(),
	status: caseStatusSchema,
	symptoms: z.array(symptomResponseSchema),
	diagnosis: diagnosisResponseSchema.nullable(),
	createdAt: z.iso.datetime(),
	updatedAt: z.iso.datetime(),
	closedAt: z.iso.datetime().nullable(),
});

export const imageResponseSchema = z.object({
	id: z.uuid(),
	caseId: z.uuid(),
	fileName: z.string(),
	contentType: z.string(),
	uploadedAt: z.iso.datetime(),
});

export const uploadURLRequestSchema = z.object({
	fileName: z.string().min(1),
	contentType: contentTypeSchema,
});

export const uploadURLResponseSchema = z.object({
	uploadUrl: z.url(),
	s3Key: z.string(),
	expiresIn: z.number().int().positive(),
});

export const confirmUploadInputSchema = z.object({
	s3Key: z.string(),
	fileName: z.string(),
	contentType: contentTypeSchema,
});

export const downloadURLResponseSchema = z.object({
	downloadUrl: z.url(),
	expiresIn: z.number().int().positive(),
});

export type CaseStatus = z.infer<typeof caseStatusSchema>;
export type Severity = z.infer<typeof severitySchema>;
export type ContentType = z.infer<typeof contentTypeSchema>;
export type AddSymptomInput = z.infer<typeof addSymptomInputSchema>;
export type CreateCaseInput = z.infer<typeof createCaseInputSchema>;
export type CaseSummary = z.infer<typeof caseSummarySchema>;
export type SymptomResponse = z.infer<typeof symptomResponseSchema>;
export type DiagnosisResponse = z.infer<typeof diagnosisResponseSchema>;
export type CaseResponse = z.infer<typeof caseResponseSchema>;
export type ImageResponse = z.infer<typeof imageResponseSchema>;
export type UploadURLRequest = z.infer<typeof uploadURLRequestSchema>;
export type UploadURLResponse = z.infer<typeof uploadURLResponseSchema>;
export type ConfirmUploadInput = z.infer<typeof confirmUploadInputSchema>;
export type DownloadURLResponse = z.infer<typeof downloadURLResponseSchema>;
