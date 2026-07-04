import type { AuthSession } from "@/features/authentication/schemas/auth";

type SessionListener = () => void;

function createEmptySession(): AuthSession {
	return {
		accessToken: null,
		refreshToken: null,
		tenants: [],
		activeTenant: null,
		user: null,
	};
}

let session = createEmptySession();
const listeners = new Set<SessionListener>();

function emit() {
	for (const listener of listeners) {
		listener();
	}
}

export function getAuthSession() {
	return session;
}

export function subscribeAuthSession(listener: SessionListener) {
	listeners.add(listener);

	return () => {
		listeners.delete(listener);
	};
}

export function updateAuthSession(patch: Partial<AuthSession>) {
	session = {
		...session,
		...patch,
	};
	emit();
}

export function clearAuthSession() {
	session = createEmptySession();
	emit();
}
