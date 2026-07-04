import type { components } from "@/generated/api";

import { apiClient } from "@/infrastructure/api/client";

import {
	type AddSymptomInput,
	type CaseResponse,
	type CaseSummary,
	type ConfirmUploadInput,
	caseResponseSchema,
	caseSummarySchema,
	confirmUploadInputSchema,
	type DownloadURLResponse,
	downloadURLResponseSchema,
	type ImageResponse,
	imageResponseSchema,
	type SymptomResponse,
	symptomResponseSchema,
	type UploadURLResponse,
	uploadURLRequestSchema,
	uploadURLResponseSchema,
} from "../schemas/cases";

function requireData<T>(value: T | null | undefined, message: string): T {
	if (value == null) {
		throw new Error(message);
	}

	return value;
}

type CaseSummaryData = components["schemas"]["CaseSummary"];
type SymptomResponseData = components["schemas"]["SymptomResponse"];
type DiagnosisResponseData = components["schemas"]["DiagnosisResponse"];
type ImageResponseData = components["schemas"]["ImageResponse"];

function normalizeCaseSummary(raw: CaseSummaryData): CaseSummary {
	return caseSummarySchema.parse({
		id: requireData(raw.id, "Missing case id"),
		patientId: requireData(raw.patient_id, "Missing patient id"),
		doctorId: raw.doctor_id ?? null,
		status: requireData(raw.status, "Missing case status"),
		createdAt: requireData(raw.created_at, "Missing created at"),
	});
}

function normalizeSymptom(raw: SymptomResponseData): SymptomResponse {
	return symptomResponseSchema.parse({
		id: requireData(raw.id, "Missing symptom id"),
		description: requireData(raw.description, "Missing description"),
		severity: requireData(raw.severity, "Missing severity"),
		reportedAt: requireData(raw.reported_at, "Missing reported at"),
	});
}

function normalizeDiagnosis(raw: DiagnosisResponseData) {
	return {
		id: requireData(raw.id, "Missing diagnosis id"),
		caseId: requireData(raw.case_id, "Missing case id"),
		doctorId: requireData(raw.doctor_id, "Missing doctor id"),
		notes: requireData(raw.notes, "Missing notes"),
		writtenAt: requireData(raw.written_at, "Missing written at"),
	};
}

function normalizeImage(raw: ImageResponseData): ImageResponse {
	return imageResponseSchema.parse({
		id: requireData(raw.id, "Missing image id"),
		caseId: requireData(raw.case_id, "Missing case id"),
		fileName: requireData(raw.file_name, "Missing file name"),
		contentType: requireData(raw.content_type, "Missing content type"),
		uploadedAt: requireData(raw.uploaded_at, "Missing uploaded at"),
	});
}

function normalizeCaseResponse(
	raw: NonNullable<components["schemas"]["CaseResponse"]["data"]>,
): CaseResponse {
	return caseResponseSchema.parse({
		id: requireData(raw.id, "Missing case id"),
		tenantId: requireData(raw.tenant_id, "Missing tenant id"),
		patientId: requireData(raw.patient_id, "Missing patient id"),
		doctorId: raw.doctor_id ?? null,
		status: requireData(raw.status, "Missing case status"),
		symptoms: (raw.symptoms ?? []).map(normalizeSymptom),
		diagnosis: raw.diagnosis ? normalizeDiagnosis(raw.diagnosis) : null,
		createdAt: requireData(raw.created_at, "Missing created at"),
		updatedAt: requireData(raw.updated_at, "Missing updated at"),
		closedAt: raw.closed_at ?? null,
	});
}

export async function listCases(params?: {
	page?: number;
	pageSize?: number;
}): Promise<{ cases: CaseSummary[]; total: number }> {
	const response = await apiClient.GET("/cases", {
		params: {
			query: {
				page: params?.page,
				page_size: params?.pageSize,
			},
		},
	});

	if (response.error) {
		throw new Error("Unable to load cases.");
	}

	const data = requireData(response.data, "Unable to load cases.");

	return {
		cases: (data.data ?? []).map(normalizeCaseSummary),
		total: data.meta?.total ?? 0,
	};
}

export async function getCase(id: string): Promise<CaseResponse> {
	const response = await apiClient.GET("/cases/{id}", {
		params: { path: { id } },
	});

	if (response.error) {
		throw new Error("Unable to load case details.");
	}

	const data = requireData(response.data?.data, "Unable to load case details.");

	return normalizeCaseResponse(data);
}

export async function createCase(
	symptoms: AddSymptomInput[],
): Promise<CaseResponse> {
	const response = await apiClient.POST("/cases", {
		body: { symptoms },
	});

	if (response.error) {
		throw new Error("Unable to create case.");
	}

	const data = requireData(response.data?.data, "Unable to create case.");

	return normalizeCaseResponse(data);
}

export async function addSymptom(
	caseId: string,
	input: AddSymptomInput,
): Promise<SymptomResponse> {
	const response = await apiClient.POST("/cases/{id}/symptoms", {
		params: { path: { id: caseId } },
		body: input,
	});

	if (response.error) {
		throw new Error("Unable to add symptom.");
	}

	const data = requireData(response.data, "Unable to add symptom.");

	return normalizeSymptom(data);
}

export async function requestUploadURL(
	caseId: string,
	fileName: string,
	contentType: "image/jpeg" | "image/png" | "image/dicom",
	fileSize: number,
): Promise<UploadURLResponse> {
	const body = { fileName, contentType, fileSize };
	const parsed = uploadURLRequestSchema.parse(body);
	const response = await apiClient.POST("/cases/{id}/images/upload-url", {
		params: { path: { id: caseId } },
		body: {
			file_name: parsed.fileName,
			content_type: parsed.contentType,
			file_size: parsed.fileSize,
		},
	});

	if (response.error) {
		throw new Error("Unable to request upload URL.");
	}

	const data = requireData(
		response.data?.data,
		"Unable to request upload URL.",
	);

	return uploadURLResponseSchema.parse({
		uploadUrl: requireData(data.upload_url, "Missing upload url"),
		s3Key: requireData(data.s3_key, "Missing s3 key"),
		expiresIn: requireData(data.expires_in, "Missing expires in"),
	});
}

export async function uploadFileToS3(
	uploadUrl: string,
	file: File,
): Promise<void> {
	const response = await fetch(uploadUrl, {
		method: "PUT",
		body: file,
		headers: { "Content-Type": file.type },
	});

	if (!response.ok) {
		throw new Error("Upload to S3 failed.");
	}
}

export async function confirmUpload(
	caseId: string,
	input: ConfirmUploadInput,
): Promise<ImageResponse> {
	const body = confirmUploadInputSchema.parse(input);

	const response = await apiClient.POST("/cases/{id}/images", {
		params: { path: { id: caseId } },
		body: {
			s3_key: body.s3Key,
			file_name: body.fileName,
			content_type: body.contentType,
		},
	});

	if (response.error) {
		throw new Error("Unable to confirm upload.");
	}

	const data = requireData(response.data, "Unable to confirm upload.");

	return normalizeImage(data);
}

export async function listImages(caseId: string): Promise<ImageResponse[]> {
	const response = await apiClient.GET("/cases/{id}/images", {
		params: { path: { id: caseId } },
	});

	if (response.error) {
		throw new Error("Unable to load images.");
	}

	const data = requireData(response.data, "Unable to load images.");

	return (data.data ?? []).map(normalizeImage);
}

export async function getDownloadURL(
	imageId: string,
): Promise<DownloadURLResponse> {
	const response = await apiClient.GET("/images/{id}/download-url", {
		params: { path: { id: imageId } },
	});

	if (response.error) {
		throw new Error("Unable to get download URL.");
	}

	const data = requireData(response.data?.data, "Unable to get download URL.");

	return downloadURLResponseSchema.parse({
		downloadUrl: requireData(data.download_url, "Missing download url"),
		expiresIn: requireData(data.expires_in, "Missing expires in"),
	});
}
