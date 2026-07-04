import { useEffect, useRef } from "react";

import { clearAuthSession } from "@/infrastructure/auth/session-store";

const ACTIVITY_EVENTS = [
	"mousedown",
	"keydown",
	"touchstart",
	"scroll",
] as const;
const TIMEOUT_MS = 15 * 60 * 1000;

export function useInactivityLogoff(
	isAuthenticated: boolean,
	onSignOut: () => void,
) {
	const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

	useEffect(() => {
		if (!isAuthenticated) return;

		function resetTimer() {
			if (timerRef.current !== null) {
				clearTimeout(timerRef.current);
			}
			timerRef.current = setTimeout(() => {
				clearAuthSession();
				onSignOut();
			}, TIMEOUT_MS);
		}

		for (const event of ACTIVITY_EVENTS) {
			document.addEventListener(event, resetTimer, { passive: true });
		}
		resetTimer();

		return () => {
			if (timerRef.current !== null) {
				clearTimeout(timerRef.current);
			}
			for (const event of ACTIVITY_EVENTS) {
				document.removeEventListener(event, resetTimer);
			}
		};
	}, [isAuthenticated, onSignOut]);
}
