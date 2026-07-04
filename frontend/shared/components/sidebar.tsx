"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import type { UserProfile } from "@/features/authentication/schemas/auth";

type NavItem = {
	href: string;
	label: string;
	roles: UserProfile["role"][];
};

const navItems: NavItem[] = [
	{
		href: "/cases/",
		label: "Cases",
		roles: ["patient", "doctor", "administrator"],
	},
	{ href: "/members/", label: "Members", roles: ["administrator"] },
	{ href: "/audit/", label: "Audit Logs", roles: ["administrator"] },
];

export function Sidebar({
	user,
	onSignOut,
}: {
	user: UserProfile;
	onSignOut: () => void;
}) {
	const pathname = usePathname();
	const items = navItems.filter((item) => item.roles.includes(user.role));

	return (
		<aside className="flex flex-col rounded-3xl border border-white/10 bg-white/5 p-5">
			<div className="mb-6">
				<p className="text-xs uppercase tracking-[0.3em] text-sky-300">
					Active tenant
				</p>
				<p className="mt-2 text-sm text-slate-300">{user.email}</p>
				<div className="mt-2 flex gap-2">
					<span className="inline-flex items-center rounded-full border border-sky-400/20 bg-sky-400/10 px-2.5 py-0.5 text-xs font-medium text-sky-300 capitalize">
						{user.role}
					</span>
					<span className="inline-flex items-center rounded-full border border-white/10 bg-white/5 px-2.5 py-0.5 text-xs font-medium text-slate-400 capitalize">
						{user.status}
					</span>
				</div>
			</div>

			<nav className="flex-1 space-y-1">
				{items.map((item) => {
					const active =
						pathname === item.href || pathname.startsWith(item.href);
					return (
						<Link
							key={item.href}
							href={item.href}
							className={`block rounded-2xl px-4 py-2.5 text-sm font-medium transition ${
								active
									? "bg-white/10 text-white"
									: "text-slate-400 hover:bg-white/5 hover:text-white"
							}`}
						>
							{item.label}
						</Link>
					);
				})}
			</nav>

			<button
				type="button"
				onClick={onSignOut}
				className="mt-4 rounded-2xl border border-white/10 px-4 py-2.5 text-sm text-slate-400 transition hover:border-white/30 hover:bg-white/5 hover:text-white"
			>
				Sign out
			</button>
		</aside>
	);
}
