import { useCallback, useEffect, useRef } from "react";

import {
	clearAuthSession,
	getAuthSession,
	updateAuthSession,
} from "@/infrastructure/auth/session-store";
import { refreshSession } from "@/features/authentication/services/auth";

const REFRESH_BUFFER_MS = 30_000;

export function useTokenRefresh() {
	const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

	const schedule = useCallback(function schedule(expiresInSeconds: number) {
		if (timerRef.current !== null) {
			clearTimeout(timerRef.current);
		}

		const delay = Math.max(expiresInSeconds * 1000 - REFRESH_BUFFER_MS, 5_000);

		timerRef.current = setTimeout(async () => {
			const current = getAuthSession();
			if (!current.refreshToken) return;

			try {
				const result = await refreshSession(current.refreshToken);
				updateAuthSession({
					accessToken: result.accessToken,
					refreshToken: result.refreshToken,
				});
				schedule(result.expiresIn);
			} catch {
				clearAuthSession();
			}
		}, delay);
	}, []);

	const start = useCallback(
		(expiresInSeconds: number) => {
			schedule(expiresInSeconds);
		},
		[schedule],
	);

	const stop = useCallback(() => {
		if (timerRef.current !== null) {
			clearTimeout(timerRef.current);
			timerRef.current = null;
		}
	}, []);

	useEffect(() => {
		return () => {
			if (timerRef.current !== null) {
				clearTimeout(timerRef.current);
			}
		};
	}, []);

	return { start, stop };
}
