import type { components } from "@/generated/api";

import { apiClient } from "@/infrastructure/api/client";

import {
	assignDoctorSchema,
	caseDetailSchema,
	caseSummarySchema,
} from "../schemas/cases";

type CaseSummaryRaw = components["schemas"]["CaseSummary"];
type CaseResponseRaw = NonNullable<
	components["schemas"]["CaseResponse"]["data"]
>;

function requireData<T>(value: T | null | undefined, message: string): T {
	if (value == null) {
		throw new Error(message);
	}

	return value;
}

function normalizeCaseSummary(raw: CaseSummaryRaw) {
	return caseSummarySchema.parse({
		id: requireData(raw.id, "Missing case id"),
		patientId: requireData(raw.patient_id, "Missing patient id"),
		doctorId: raw.doctor_id ?? null,
		status: requireData(raw.status, "Missing status"),
		createdAt: requireData(raw.created_at, "Missing created at"),
	});
}

function normalizeCaseDetail(raw: CaseResponseRaw) {
	return caseDetailSchema.parse({
		id: requireData(raw.id, "Missing case id"),
		tenantId: requireData(raw.tenant_id, "Missing tenant id"),
		patientId: requireData(raw.patient_id, "Missing patient id"),
		doctorId: raw.doctor_id ?? null,
		status: requireData(raw.status, "Missing status"),
		symptoms: (raw.symptoms ?? []).map((s) => ({
			id: requireData(s.id, "Missing symptom id"),
			description: s.description,
			severity: s.severity,
			reportedAt: requireData(s.reported_at, "Missing symptom reported at"),
		})),
		diagnosis: raw.diagnosis
			? {
					id: requireData(raw.diagnosis.id, "Missing diagnosis id"),
					caseId: requireData(
						raw.diagnosis.case_id,
						"Missing diagnosis case id",
					),
					doctorId: requireData(
						raw.diagnosis.doctor_id,
						"Missing diagnosis doctor id",
					),
					notes: raw.diagnosis.notes,
					writtenAt: requireData(
						raw.diagnosis.written_at,
						"Missing diagnosis written at",
					),
				}
			: null,
		createdAt: requireData(raw.created_at, "Missing created at"),
		updatedAt: requireData(raw.updated_at, "Missing updated at"),
		closedAt: raw.closed_at ?? null,
	});
}

export async function listCases(params: { page?: number; pageSize?: number }) {
	const response = await apiClient.GET("/cases", {
		params: {
			query: {
				page: params.page,
				page_size: params.pageSize,
			},
		},
	});

	if (response.error) {
		throw new Error("Unable to load cases.");
	}

	const data = requireData(response.data, "Unable to load cases.");

	return {
		data: (data.data ?? []).map(normalizeCaseSummary),
		meta: data.meta,
	};
}

export async function getCase(id: string) {
	const response = await apiClient.GET("/cases/{id}", {
		params: { path: { id } },
	});

	if (response.error) {
		throw new Error("Unable to load case.");
	}

	const data = requireData(response.data?.data, "Unable to load case.");

	return normalizeCaseDetail(data);
}

export async function assignDoctor(
	caseId: string,
	input: { doctorId: string },
) {
	const body = assignDoctorSchema.parse(input);
	const response = await apiClient.POST("/cases/{id}/assign", {
		params: { path: { id: caseId } },
		body: { doctor_id: body.doctorId },
	});

	if (response.error) {
		throw new Error("Unable to assign doctor.");
	}

	const data = requireData(response.data?.data, "Unable to assign doctor.");

	return normalizeCaseDetail(data);
}

export async function closeCase(caseId: string) {
	const response = await apiClient.POST("/cases/{id}/close", {
		params: { path: { id: caseId } },
	});

	if (response.error) {
		throw new Error("Unable to close case.");
	}

	const data = requireData(response.data?.data, "Unable to close case.");

	return normalizeCaseDetail(data);
}
