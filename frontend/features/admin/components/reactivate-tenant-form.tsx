"use client";

import { useState } from "react";

import { Button } from "@/shared/components/button";
import { Dialog } from "@/shared/components/dialog";

import { useReactivateTenant } from "../hooks/use-members";

export function ReactivateTenantForm() {
	const [open, setOpen] = useState(false);
	const [tenantId, setTenantId] = useState("");
	const reactivateMutation = useReactivateTenant();

	function handleSubmit(event: React.FormEvent) {
		event.preventDefault();
		reactivateMutation.mutate(tenantId, {
			onSuccess: () => {
				setTenantId("");
				setOpen(false);
			},
		});
	}

	return (
		<>
			<Button variant="secondary" onClick={() => setOpen(true)}>
				Reactivate Tenant
			</Button>
			<Dialog
				open={open}
				onClose={() => setOpen(false)}
				title="Reactivate Tenant"
			>
				<form onSubmit={handleSubmit} className="space-y-4">
					<div>
						<label htmlFor="tenant-id" className="block text-sm text-slate-400">
							Tenant ID
						</label>
						<input
							id="tenant-id"
							type="text"
							value={tenantId}
							onChange={(event) => setTenantId(event.target.value)}
							placeholder="Enter tenant UUID"
							className="mt-1 w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white placeholder-slate-500 outline-none transition focus:border-sky-400/50 focus:ring-1 focus:ring-sky-400/50"
						/>
					</div>
					{reactivateMutation.error && (
						<p className="text-sm text-red-400">
							{reactivateMutation.error.message}
						</p>
					)}
					<div className="flex justify-end gap-2">
						<Button variant="ghost" onClick={() => setOpen(false)}>
							Cancel
						</Button>
						<Button
							type="submit"
							disabled={!tenantId}
							busy={reactivateMutation.isPending}
						>
							Reactivate
						</Button>
					</div>
				</form>
			</Dialog>
		</>
	);
}
