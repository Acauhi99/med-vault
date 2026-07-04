import { afterEach, describe, expect, it } from "vitest";

import {
	clearAuthSession,
	getAuthSession,
	subscribeAuthSession,
	updateAuthSession,
} from "./session-store";

const emptySession = {
	accessToken: null,
	refreshToken: null,
	tenants: [],
	activeTenant: null,
	user: null,
};

afterEach(() => {
	clearAuthSession();
});

describe("session store", () => {
	it("updates and notifies listeners", () => {
		let calls = 0;
		const unsubscribe = subscribeAuthSession(() => {
			calls += 1;
		});

		updateAuthSession({ accessToken: "token" });

		expect(calls).toBe(1);
		expect(getAuthSession()).toMatchObject({ accessToken: "token" });

		unsubscribe();
	});

	it("clears the session", () => {
		updateAuthSession({ accessToken: "token", refreshToken: "refresh" });
		clearAuthSession();

		expect(getAuthSession()).toEqual(emptySession);
	});
});
