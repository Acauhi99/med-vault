"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { HTMLInputTypeAttribute, ReactNode, SubmitEvent } from "react";
import { useEffect, useRef, useState } from "react";

import {
	clearAuthSession,
	updateAuthSession,
} from "@/infrastructure/auth/session-store";
import { useAuthSession } from "@/infrastructure/auth/use-auth-session";
import {
	loginInputSchema,
	registerInputSchema,
	type TenantSummary,
	tenantSelectionSchema,
} from "../schemas/auth";
import {
	getCurrentUser,
	login,
	register,
	selectTenant,
} from "../services/auth";
import { AppShell } from "./app-shell";

type AuthMode = "login" | "register";

const emptyLoginForm = {
	email: "",
	password: "",
};

const emptyRegisterForm = {
	email: "",
	password: "",
};

export function AuthWorkspace() {
	const session = useAuthSession();
	const queryClient = useQueryClient();
	const authEpoch = useRef(0);
	const [mode, setMode] = useState<AuthMode>("login");
	const [loginForm, setLoginForm] = useState(emptyLoginForm);
	const [registerForm, setRegisterForm] = useState(emptyRegisterForm);
	const [tenantId, setTenantId] = useState("");
	const [notice, setNotice] = useState<string | null>(null);

	const loginMutation = useMutation({ mutationFn: login });
	const registerMutation = useMutation({ mutationFn: register });
	const tenantMutation = useMutation({ mutationFn: selectTenant });

	const currentUserQuery = useQuery({
		queryKey: [
			"current-user",
			session.accessToken,
			session.activeTenant?.tenantId ?? "pending",
		],
		queryFn: () => getCurrentUser(session.accessToken ?? ""),
		enabled: Boolean(session.accessToken && session.activeTenant),
	});

	useEffect(() => {
		if (
			currentUserQuery.data &&
			session.accessToken &&
			session.activeTenant?.tenantId === currentUserQuery.data.tenantId
		) {
			updateAuthSession({ user: currentUserQuery.data });
		}
	}, [
		currentUserQuery.data,
		session.accessToken,
		session.activeTenant?.tenantId,
	]);

	const isBusy =
		loginMutation.isPending ||
		registerMutation.isPending ||
		tenantMutation.isPending ||
		currentUserQuery.isFetching;

	async function pickTenant(
		selectedTenantId: string,
		accessToken: string,
		tenants: TenantSummary[],
	) {
		const epoch = authEpoch.current;
		const parsed = tenantSelectionSchema.parse({ tenantId: selectedTenantId });
		const result = await tenantMutation.mutateAsync({
			accessToken,
			tenantId: parsed.tenantId,
		});
		if (epoch !== authEpoch.current) {
			return;
		}
		const activeTenant = tenants.find(
			(tenant) => tenant.tenantId === parsed.tenantId,
		);

		updateAuthSession({
			accessToken: result.accessToken,
			refreshToken: result.refreshToken,
			tenants,
			activeTenant: activeTenant ?? null,
			user: null,
		});
	}

	async function handleLogin(event: SubmitEvent<HTMLFormElement>) {
		event.preventDefault();
		setNotice(null);
		const epoch = authEpoch.current;

		try {
			const input = loginInputSchema.parse(loginForm);
			const result = await loginMutation.mutateAsync(input);
			if (epoch !== authEpoch.current) {
				return;
			}

			updateAuthSession({
				accessToken: result.accessToken,
				refreshToken: null,
				tenants: result.tenants,
				activeTenant: null,
				user: null,
			});

			setTenantId(result.tenants[0]?.tenantId ?? "");

			if (result.tenants.length === 0) {
				setNotice("No tenant memberships found. Ask an admin to add you.");
				return;
			}

			if (result.tenants.length === 1) {
				await pickTenant(
					result.tenants[0].tenantId,
					result.accessToken,
					result.tenants,
				);
			} else {
				setNotice("Choose a tenant to continue.");
			}

			setLoginForm({ ...emptyLoginForm, email: input.email });
		} catch (error) {
			setNotice(error instanceof Error ? error.message : "Unable to sign in.");
		}
	}

	async function handleRegister(event: SubmitEvent<HTMLFormElement>) {
		event.preventDefault();
		setNotice(null);
		const epoch = authEpoch.current;

		try {
			const input = registerInputSchema.parse(registerForm);
			const result = await registerMutation.mutateAsync(input);
			if (epoch !== authEpoch.current) {
				return;
			}
			setNotice(`Account created for ${result.email}. Sign in to continue.`);
			setMode("login");
			setRegisterForm(emptyRegisterForm);
			setLoginForm((current) => ({ ...current, email: result.email }));
		} catch (error) {
			setNotice(error instanceof Error ? error.message : "Unable to register.");
		}
	}

	async function handleTenantSubmit(event: SubmitEvent<HTMLFormElement>) {
		event.preventDefault();
		setNotice(null);
		const epoch = authEpoch.current;

		if (!session.accessToken) {
			setNotice("Sign in again to continue.");
			return;
		}

		try {
			await pickTenant(tenantId, session.accessToken, session.tenants);
			if (epoch !== authEpoch.current) {
				return;
			}
		} catch (error) {
			setNotice(
				error instanceof Error ? error.message : "Unable to select tenant.",
			);
		}
	}

	async function handleSignOut() {
		authEpoch.current += 1;
		await queryClient.cancelQueries({ queryKey: ["current-user"] });
		queryClient.clear();
		clearAuthSession();
		setLoginForm(emptyLoginForm);
		setRegisterForm(emptyRegisterForm);
		setTenantId("");
		setMode("login");
		setNotice("Signed out.");
	}

	if (session.accessToken && session.activeTenant) {
		if (session.user) {
			return <AppShell />;
		}

		if (currentUserQuery.isError) {
			return (
				<WorkspaceFrame>
					<EmptyState
						title="Workspace failed to load"
						body={
							currentUserQuery.error instanceof Error
								? currentUserQuery.error.message
								: "Unable to load the current workspace."
						}
						actionLabel="Sign out"
						onAction={handleSignOut}
					/>
				</WorkspaceFrame>
			);
		}

		return (
			<WorkspaceFrame>
				<LoadingState
					title="Loading workspace"
					body="Checking the current tenant and loading user context."
				/>
			</WorkspaceFrame>
		);
	}

	if (session.accessToken && !session.activeTenant) {
		if (session.tenants.length === 0) {
			return (
				<WorkspaceFrame>
					<EmptyState
						title="No tenant access yet"
						body="This user is registered, but no tenant membership exists. Sign out and ask an administrator to add the account to a tenant."
						actionLabel="Sign out"
						onAction={handleSignOut}
					/>
				</WorkspaceFrame>
			);
		}

		return (
			<WorkspaceFrame>
				<TenantPicker
					busy={isBusy}
					currentTenantId={tenantId}
					notice={notice}
					tenants={session.tenants}
					onChange={setTenantId}
					onSignOut={handleSignOut}
					onSubmit={handleTenantSubmit}
				/>
			</WorkspaceFrame>
		);
	}

	return (
		<WorkspaceFrame>
			<AuthCard
				busy={isBusy}
				loginForm={loginForm}
				mode={mode}
				notice={notice}
				onLoginFormChange={setLoginForm}
				onLoginSubmit={handleLogin}
				onModeChange={setMode}
				onRegisterFormChange={setRegisterForm}
				onRegisterSubmit={handleRegister}
				registerForm={registerForm}
			/>
		</WorkspaceFrame>
	);
}

function WorkspaceFrame({ children }: { children: ReactNode }) {
	return (
		<div className="min-h-screen bg-slate-950 px-4 py-8 text-slate-50 sm:px-6 lg:px-8">
			<div className="mx-auto flex min-h-[calc(100vh-4rem)] max-w-6xl items-center justify-center">
				{children}
			</div>
		</div>
	);
}

type AuthCardProps = {
	busy: boolean;
	loginForm: { email: string; password: string };
	mode: AuthMode;
	notice: string | null;
	registerForm: { email: string; password: string };
	onLoginFormChange: (next: { email: string; password: string }) => void;
	onLoginSubmit: (event: SubmitEvent<HTMLFormElement>) => Promise<void>;
	onModeChange: (next: AuthMode) => void;
	onRegisterFormChange: (next: { email: string; password: string }) => void;
	onRegisterSubmit: (event: SubmitEvent<HTMLFormElement>) => Promise<void>;
};

function AuthCard({
	busy,
	loginForm,
	mode,
	notice,
	registerForm,
	onLoginFormChange,
	onLoginSubmit,
	onModeChange,
	onRegisterFormChange,
	onRegisterSubmit,
}: AuthCardProps) {
	return (
		<div className="grid w-full gap-6 rounded-[2rem] border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 lg:grid-cols-[1.1fr_0.9fr] lg:p-8">
			<section className="rounded-[1.5rem] bg-slate-950/70 p-6">
				<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
					MedVault
				</p>
				<h1 className="mt-4 text-4xl font-semibold leading-tight text-white">
					Secure care workspace.
				</h1>
				<p className="mt-4 max-w-xl text-sm leading-6 text-slate-300">
					Sign in, select a tenant, and land in the shared shell for patient,
					doctor, and admin work.
				</p>

				<div className="mt-8 grid gap-3 sm:grid-cols-3">
					{[
						["Multi-tenant", "Tenant context stays explicit."],
						["Auth first", "Login, register, select tenant."],
						["Shell ready", "Workspace chrome is live."],
					].map(([title, body]) => (
						<article
							key={title}
							className="rounded-2xl border border-white/10 bg-white/5 p-4"
						>
							<h2 className="text-sm font-semibold text-white">{title}</h2>
							<p className="mt-2 text-xs leading-5 text-slate-300">{body}</p>
						</article>
					))}
				</div>
			</section>

			<section className="rounded-[1.5rem] border border-white/10 bg-slate-950 p-6">
				<div className="flex gap-2 rounded-full bg-white/5 p-1 text-sm">
					<TabButton
						active={mode === "login"}
						onClick={() => onModeChange("login")}
					>
						Sign in
					</TabButton>
					<TabButton
						active={mode === "register"}
						onClick={() => onModeChange("register")}
					>
						Register
					</TabButton>
				</div>

				{notice ? (
					<p
						role="status"
						aria-live="polite"
						aria-atomic="true"
						className="mt-4 rounded-2xl border border-sky-400/30 bg-sky-400/10 px-4 py-3 text-sm text-sky-100"
					>
						{notice}
					</p>
				) : null}

				{mode === "login" ? (
					<form className="mt-6 grid gap-4" onSubmit={onLoginSubmit}>
						<Field
							label="Email"
							name="email"
							type="email"
							autoComplete="email"
							value={loginForm.email}
							onChange={(email) => onLoginFormChange({ ...loginForm, email })}
						/>
						<Field
							label="Password"
							name="password"
							type="password"
							autoComplete="current-password"
							value={loginForm.password}
							onChange={(password) =>
								onLoginFormChange({ ...loginForm, password })
							}
						/>
						<ActionButton busy={busy}>Sign in</ActionButton>
					</form>
				) : (
					<form className="mt-6 grid gap-4" onSubmit={onRegisterSubmit}>
						<Field
							label="Email"
							name="register-email"
							type="email"
							autoComplete="email"
							value={registerForm.email}
							onChange={(email) =>
								onRegisterFormChange({ ...registerForm, email })
							}
						/>
						<Field
							label="Password"
							name="register-password"
							type="password"
							autoComplete="new-password"
							value={registerForm.password}
							onChange={(password) =>
								onRegisterFormChange({ ...registerForm, password })
							}
						/>
						<ActionButton busy={busy}>Create account</ActionButton>
					</form>
				)}
			</section>
		</div>
	);
}

type TenantPickerProps = {
	busy: boolean;
	currentTenantId: string;
	notice: string | null;
	tenants: TenantSummary[];
	onChange: (tenantId: string) => void;
	onSignOut: () => void;
	onSubmit: (event: SubmitEvent<HTMLFormElement>) => Promise<void>;
};

function TenantPicker({
	busy,
	currentTenantId,
	notice,
	tenants,
	onChange,
	onSignOut,
	onSubmit,
}: TenantPickerProps) {
	return (
		<div className="w-full max-w-2xl rounded-[2rem] border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30">
			<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
				Choose tenant
			</p>
			<h1 className="mt-4 text-3xl font-semibold text-white">
				Pick the workspace for this session.
			</h1>
			<p className="mt-3 text-sm leading-6 text-slate-300">
				The backend returned {tenants.length} tenant
				{tenants.length === 1 ? "" : "s"}.
			</p>

			{notice ? (
				<p
					role="status"
					aria-live="polite"
					aria-atomic="true"
					className="mt-4 rounded-2xl border border-sky-400/30 bg-sky-400/10 px-4 py-3 text-sm text-sky-100"
				>
					{notice}
				</p>
			) : null}

			<form className="mt-6 grid gap-4" onSubmit={onSubmit}>
				<label
					className="grid gap-2 text-sm text-slate-200"
					htmlFor="tenant-id"
				>
					Tenant
					<select
						id="tenant-id"
						name="tenant-id"
						value={currentTenantId}
						onChange={(event) => onChange(event.target.value)}
						disabled={busy}
						className="h-12 rounded-2xl border border-white/10 bg-slate-950 px-4 text-slate-50 transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950 disabled:cursor-not-allowed disabled:opacity-60"
					>
						{tenants.map((tenant) => (
							<option key={tenant.tenantId} value={tenant.tenantId}>
								{tenant.tenantName} · {tenant.role}
							</option>
						))}
					</select>
				</label>

				<div className="flex flex-wrap gap-3">
					<ActionButton busy={busy}>Enter workspace</ActionButton>
					<button
						type="button"
						onClick={onSignOut}
						disabled={busy}
						className="inline-flex h-12 items-center justify-center rounded-2xl border border-white/10 px-4 text-sm text-slate-200 transition hover:border-white/30 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-60"
					>
						Sign out
					</button>
				</div>
			</form>
		</div>
	);
}

function EmptyState({
	actionLabel,
	body,
	onAction,
	title,
}: {
	actionLabel: string;
	body: string;
	onAction: () => void;
	title: string;
}) {
	return (
		<div
			role="status"
			aria-busy="true"
			className="w-full max-w-xl rounded-[2rem] border border-white/10 bg-white/5 p-6 text-center shadow-2xl shadow-black/30"
		>
			<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
				Workspace
			</p>
			<h1 className="mt-4 text-3xl font-semibold text-white">{title}</h1>
			<p className="mt-3 text-sm leading-6 text-slate-300">{body}</p>
			<button
				type="button"
				onClick={onAction}
				className="mt-6 inline-flex h-12 items-center justify-center rounded-2xl bg-white px-4 text-sm font-medium text-slate-950 transition hover:bg-slate-200"
			>
				{actionLabel}
			</button>
		</div>
	);
}

function LoadingState({ title, body }: { title: string; body: string }) {
	return (
		<div className="w-full max-w-xl rounded-[2rem] border border-white/10 bg-white/5 p-6 text-center shadow-2xl shadow-black/30">
			<p className="text-xs uppercase tracking-[0.35em] text-sky-300">
				Workspace
			</p>
			<h1 className="mt-4 text-3xl font-semibold text-white">{title}</h1>
			<p className="mt-3 text-sm leading-6 text-slate-300">{body}</p>
		</div>
	);
}

function Field({
	autoComplete,
	label,
	name,
	onChange,
	type,
	value,
}: {
	autoComplete?: string;
	label: string;
	name: string;
	onChange: (value: string) => void;
	type: HTMLInputTypeAttribute;
	value: string;
}) {
	return (
		<label className="grid gap-2 text-sm text-slate-200" htmlFor={name}>
			{label}
			<input
				id={name}
				name={name}
				type={type}
				autoComplete={autoComplete}
				value={value}
				onChange={(event) => onChange(event.target.value)}
				className="h-12 rounded-2xl border border-white/10 bg-slate-950 px-4 text-slate-50 transition placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
			/>
		</label>
	);
}

function TabButton({
	active,
	children,
	onClick,
}: {
	active: boolean;
	children: ReactNode;
	onClick: () => void;
}) {
	return (
		<button
			type="button"
			onClick={onClick}
			aria-pressed={active}
			className={`flex-1 rounded-full px-4 py-2 text-sm font-medium transition ${
				active
					? "bg-sky-400 text-slate-950"
					: "text-slate-300 hover:bg-white/5 hover:text-white"
			}`}
		>
			{children}
		</button>
	);
}

function ActionButton({
	busy,
	children,
}: {
	busy: boolean;
	children: ReactNode;
}) {
	return (
		<button
			type="submit"
			disabled={busy}
			aria-busy={busy}
			className="inline-flex h-12 items-center justify-center rounded-2xl bg-sky-400 px-4 text-sm font-semibold text-slate-950 transition hover:bg-sky-300 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-300"
		>
			{children}
		</button>
	);
}
