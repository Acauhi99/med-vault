import type { components } from "@/generated/api";
import { apiClient } from "@/infrastructure/api/client";

import {
	type CaseDetail,
	type CaseImage,
	type CaseSummary,
	caseDetailSchema,
	caseSummarySchema,
	type Diagnosis,
	diagnosisSchema,
	downloadUrlSchema,
	imageSchema,
	type WriteDiagnosisInput,
	writeDiagnosisInputSchema,
} from "../schemas/cases";

type CaseSummaryRaw = components["schemas"]["CaseSummary"];
type CaseResponseData = NonNullable<
	components["schemas"]["CaseResponse"]["data"]
>;
type DiagnosisResponse = components["schemas"]["DiagnosisResponse"];
type ImageResponse = components["schemas"]["ImageResponse"];

function requireData<T>(value: T | null | undefined, message: string): T {
	if (value == null) {
		throw new Error(message);
	}
	return value;
}

function normalizeCaseSummary(raw: CaseSummaryRaw): CaseSummary {
	return caseSummarySchema.parse({
		id: requireData(raw.id, "Missing case id"),
		patientId: requireData(raw.patient_id, "Missing patient id"),
		doctorId: raw.doctor_id ?? null,
		status: requireData(raw.status, "Missing case status"),
		createdAt: requireData(raw.created_at, "Missing created at"),
	});
}

function normalizeCaseDetail(data: CaseResponseData): CaseDetail {
	return caseDetailSchema.parse({
		id: requireData(data.id, "Missing case id"),
		tenantId: requireData(data.tenant_id, "Missing tenant id"),
		patientId: requireData(data.patient_id, "Missing patient id"),
		doctorId: data.doctor_id ?? null,
		status: requireData(data.status, "Missing case status"),
		symptoms: (data.symptoms ?? []).map((s) => ({
			id: requireData(s.id, "Missing symptom id"),
			description: requireData(s.description, "Missing symptom description"),
			severity: requireData(s.severity, "Missing symptom severity"),
			reportedAt: requireData(s.reported_at, "Missing symptom reported at"),
		})),
		diagnosis: data.diagnosis ? normalizeDiagnosis(data.diagnosis) : null,
		createdAt: requireData(data.created_at, "Missing created at"),
		updatedAt: requireData(data.updated_at, "Missing updated at"),
		closedAt: data.closed_at ?? null,
	});
}

function normalizeDiagnosis(raw: DiagnosisResponse): Diagnosis {
	return diagnosisSchema.parse({
		id: requireData(raw.id, "Missing diagnosis id"),
		caseId: requireData(raw.case_id, "Missing case id"),
		doctorId: requireData(raw.doctor_id, "Missing doctor id"),
		notes: requireData(raw.notes, "Missing diagnosis notes"),
		writtenAt: requireData(raw.written_at, "Missing written at"),
	});
}

function normalizeImage(raw: ImageResponse): CaseImage {
	return imageSchema.parse({
		id: requireData(raw.id, "Missing image id"),
		caseId: requireData(raw.case_id, "Missing case id"),
		fileName: requireData(raw.file_name, "Missing file name"),
		contentType: requireData(raw.content_type, "Missing content type"),
		uploadedAt: requireData(raw.uploaded_at, "Missing uploaded at"),
	});
}

export async function listAssignedCases(page = 1, pageSize = 20) {
	const response = await apiClient.GET("/cases", {
		params: { query: { page, page_size: pageSize } },
	});

	if (response.error) {
		throw new Error("Unable to load cases.");
	}

	const data = requireData(response.data, "Unable to load cases.");
	const cases = requireData(data.data, "No cases data").map(
		normalizeCaseSummary,
	);
	const meta = data.meta ?? { page: 1, page_size: pageSize, total: 0 };

	return { cases, meta };
}

export async function getCaseDetail(caseId: string) {
	const response = await apiClient.GET("/cases/{id}", {
		params: { path: { id: caseId } },
	});

	if (response.error) {
		throw new Error("Unable to load case details.");
	}

	const envelope = requireData(response.data, "No case data");
	return normalizeCaseDetail(requireData(envelope.data, "No case data"));
}

export async function listCaseImages(caseId: string) {
	const response = await apiClient.GET("/cases/{id}/images", {
		params: { path: { id: caseId } },
	});

	if (response.error) {
		throw new Error("Unable to load images.");
	}

	const data = requireData(response.data, "No images data");
	return requireData(data.data, "No images").map(normalizeImage);
}

export async function getImageDownloadUrl(imageId: string) {
	const response = await apiClient.GET("/images/{id}/download-url", {
		params: { path: { id: imageId } },
	});

	if (response.error) {
		throw new Error("Unable to get download URL.");
	}

	const data = requireData(response.data, "No download URL data");
	const inner = requireData(data.data, "No download URL");

	return downloadUrlSchema.parse({
		downloadUrl: requireData(inner.download_url, "Missing download URL"),
		expiresIn: requireData(inner.expires_in, "Missing expiry"),
	});
}

export async function writeDiagnosis(
	caseId: string,
	input: WriteDiagnosisInput,
) {
	const body = writeDiagnosisInputSchema.parse(input);
	const response = await apiClient.POST("/cases/{id}/diagnosis", {
		params: { path: { id: caseId } },
		body,
	});

	if (response.error) {
		throw new Error("Unable to write diagnosis.");
	}

	return normalizeDiagnosis(requireData(response.data, "No diagnosis data"));
}
