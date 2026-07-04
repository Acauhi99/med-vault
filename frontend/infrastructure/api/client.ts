import createClient from "openapi-fetch";

import type { paths } from "@/generated/api";

const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export const apiClient = createClient<paths>({
  baseUrl,
  fetch: (...args) => globalThis.fetch(...args),
});
