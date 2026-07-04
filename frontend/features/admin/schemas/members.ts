import { z } from "zod";

const memberRoleSchema = z.enum(["patient", "doctor", "administrator"]);

export const tenantMemberSchema = z.object({
	userId: z.uuid(),
	tenantId: z.uuid(),
	role: memberRoleSchema,
	name: z.string(),
});

export const addMemberSchema = z.object({
	userId: z.uuid("Enter a valid user ID"),
	role: memberRoleSchema,
});

export const paginationSchema = z.object({
	page: z.number(),
	pageSize: z.number(),
	total: z.number(),
});

export type MemberRole = z.infer<typeof memberRoleSchema>;
export type TenantMember = z.infer<typeof tenantMemberSchema>;
export type AddMemberInput = z.infer<typeof addMemberSchema>;
