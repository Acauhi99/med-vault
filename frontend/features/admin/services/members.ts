import { apiClient } from "@/infrastructure/api/client";
import { requireData } from "@/shared/utils/require-data";

import {
	addMemberSchema,
	reactivateTenantSchema,
	tenantMemberSchema,
} from "../schemas/members";

type MemberRaw = {
	user_id?: string;
	tenant_id?: string;
	role?: string;
	name?: string;
};

type AddMemberRaw = {
	user_id?: string;
	tenant_id?: string;
	role?: string;
	name?: string;
};

export async function listMembers(tenantId: string) {
	const response = await apiClient.GET(
		"/tenants/{tenant_id}/members" as never,
		{
			params: { path: { tenant_id: tenantId } },
		} as never,
	);

	if (response.error) {
		throw new Error("Unable to load members.");
	}

	const data = requireData(
		response.data as { data?: MemberRaw[] } | null,
		"Unable to load members.",
	);
	const members = data.data ?? [];

	return members.map((m) =>
		tenantMemberSchema.parse({
			userId: requireData(m.user_id, "Missing user id"),
			tenantId: requireData(m.tenant_id, "Missing tenant id"),
			role: requireData(m.role, "Missing role"),
			name: requireData(m.name, "Missing name"),
		}),
	);
}

export async function addMember(
	tenantId: string,
	input: { userId: string; role: "patient" | "doctor" | "administrator" },
) {
	const body = addMemberSchema.parse(input);
	const response = await apiClient.POST(
		"/tenants/{tenant_id}/members" as never,
		{
			params: { path: { tenant_id: tenantId } },
			body: {
				user_id: body.userId,
				role: body.role,
			},
		} as never,
	);

	if (response.error) {
		throw new Error("Unable to add member.");
	}

	const data = requireData(
		response.data as { data?: AddMemberRaw } | null,
		"Unable to add member.",
	);
	const m = requireData(data.data, "Unable to add member.");

	return tenantMemberSchema.parse({
		userId: requireData(m.user_id, "Missing user id"),
		tenantId: requireData(m.tenant_id, "Missing tenant id"),
		role: requireData(m.role, "Missing role"),
		name: requireData(m.name, "Missing name"),
	});
}

export async function removeMember(tenantId: string, userId: string) {
	const response = await apiClient.DELETE(
		"/tenants/{tenant_id}/members/{user_id}" as never,
		{
			params: {
				path: { tenant_id: tenantId, user_id: userId },
			},
		} as never,
	);

	if (response.error) {
		throw new Error("Unable to remove member.");
	}
}

export async function reactivateTenant(input: { tenantId: string }) {
	const body = reactivateTenantSchema.parse(input);
	const response = await apiClient.POST(
		"/tenants/{tenant_id}/reactivate" as never,
		{
			params: { path: { tenant_id: body.tenantId } },
		} as never,
	);

	if (response.error) {
		throw new Error("Unable to reactivate tenant.");
	}
}
