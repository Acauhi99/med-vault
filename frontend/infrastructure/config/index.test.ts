import { afterEach, describe, expect, it, vi } from "vitest";

import { getApiBaseUrl } from "./index";

afterEach(() => {
	vi.unstubAllEnvs();
});

describe("config", () => {
	it("appends the API prefix once", () => {
		vi.stubEnv("NEXT_PUBLIC_API_URL", "https://api.med-vault.space");

		expect(getApiBaseUrl()).toBe("https://api.med-vault.space/api/v1");
	});

	it("keeps an existing API prefix", () => {
		vi.stubEnv("NEXT_PUBLIC_API_URL", "https://api.med-vault.space/api/v1");

		expect(getApiBaseUrl()).toBe("https://api.med-vault.space/api/v1");
	});
});
