"use client";

import Link from "next/link";
import { useState } from "react";

const links = [
	{ href: "#features", label: "Features" },
	{ href: "#hipaa", label: "HIPAA" },
	{ href: "#architecture", label: "Architecture" },
];

export function Navbar() {
	const [open, setOpen] = useState(false);

	return (
		<nav className="fixed top-0 z-50 w-full border-b border-white/10 bg-slate-950/80 backdrop-blur-xl">
			<div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
				<Link href="/" className="text-lg font-semibold text-white">
					MedVault
				</Link>

				<div className="hidden items-center gap-8 md:flex">
					{links.map((link) => (
						<a
							key={link.href}
							href={link.href}
							className="text-sm text-slate-400 transition hover:text-white"
						>
							{link.label}
						</a>
					))}
					<Link
						href="/login"
						className="rounded-2xl bg-sky-400 px-5 py-2 text-sm font-semibold text-slate-950 transition hover:bg-sky-300"
					>
						Get Started
					</Link>
				</div>

				<button
					type="button"
					onClick={() => setOpen(!open)}
					className="flex h-10 w-10 items-center justify-center rounded-xl text-slate-400 transition hover:bg-white/10 hover:text-white md:hidden"
					aria-label="Toggle menu"
				>
					<svg
						className="h-5 w-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						{open ? (
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M6 18L18 6M6 6l12 12"
							/>
						) : (
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M4 6h16M4 12h16M4 18h16"
							/>
						)}
					</svg>
				</button>
			</div>

			{open ? (
				<div className="border-t border-white/10 bg-slate-950/95 px-4 pb-4 pt-2 md:hidden">
					{links.map((link) => (
						<a
							key={link.href}
							href={link.href}
							onClick={() => setOpen(false)}
							className="block rounded-xl px-4 py-3 text-sm text-slate-400 transition hover:bg-white/5 hover:text-white"
						>
							{link.label}
						</a>
					))}
					<Link
						href="/login"
						onClick={() => setOpen(false)}
						className="mt-2 block rounded-2xl bg-sky-400 px-4 py-3 text-center text-sm font-semibold text-slate-950 transition hover:bg-sky-300"
					>
						Get Started
					</Link>
				</div>
			) : null}
		</nav>
	);
}
