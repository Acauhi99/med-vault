"use client";

import { useSyncExternalStore } from "react";

import { getAuthSession, subscribeAuthSession } from "./session-store";

export function useAuthSession() {
  return useSyncExternalStore(
    subscribeAuthSession,
    getAuthSession,
    getAuthSession,
  );
}
