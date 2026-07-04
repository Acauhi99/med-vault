import { z } from "zod";

const caseStatusSchema = z.enum(["open", "assigned", "diagnosed", "closed"]);

export const caseSummarySchema = z.object({
	id: z.uuid(),
	patientId: z.uuid(),
	doctorId: z.uuid().nullable().optional(),
	status: caseStatusSchema,
	createdAt: z.iso.datetime(),
});

export const assignDoctorSchema = z.object({
	doctorId: z.uuid("Select a doctor"),
});

export const caseDetailSchema = z.object({
	id: z.uuid(),
	tenantId: z.uuid(),
	patientId: z.uuid(),
	doctorId: z.uuid().nullable().optional(),
	status: caseStatusSchema,
	symptoms: z.array(
		z.object({
			id: z.uuid(),
			description: z.string().optional(),
			severity: z.enum(["low", "medium", "high", "critical"]).optional(),
			reportedAt: z.iso.datetime(),
		}),
	),
	diagnosis: z
		.object({
			id: z.uuid(),
			caseId: z.uuid(),
			doctorId: z.uuid(),
			notes: z.string().optional(),
			writtenAt: z.iso.datetime(),
		})
		.nullable(),
	createdAt: z.iso.datetime(),
	updatedAt: z.iso.datetime(),
	closedAt: z.iso.datetime().nullable(),
});

export type CaseStatus = z.infer<typeof caseStatusSchema>;
export type CaseSummary = z.infer<typeof caseSummarySchema>;
export type AssignDoctorInput = z.infer<typeof assignDoctorSchema>;
export type CaseDetail = z.infer<typeof caseDetailSchema>;
