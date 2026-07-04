"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import type { AddSymptomInput, CreateCaseInput } from "../schemas/cases";
import {
	addSymptom,
	createCase,
	getCase,
	listCases,
	listImages,
} from "../services/cases";

export function useCases(params?: { page?: number; pageSize?: number }) {
	return useQuery({
		queryKey: ["cases", params?.page ?? 1, params?.pageSize ?? 20],
		queryFn: () => listCases(params),
	});
}

export function useCase(id: string | null) {
	return useQuery({
		queryKey: ["case", id],
		queryFn: () => getCase(id as string),
		enabled: id !== null,
	});
}

export function useCreateCase() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (symptoms: CreateCaseInput["symptoms"]) => createCase(symptoms),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["cases"] });
		},
	});
}

export function useAddSymptom(caseId: string) {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (input: AddSymptomInput) => addSymptom(caseId, input),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["case", caseId] });
			queryClient.invalidateQueries({ queryKey: ["cases"] });
		},
	});
}

export function useCaseImages(caseId: string | null) {
	return useQuery({
		queryKey: ["case-images", caseId],
		queryFn: () => listImages(caseId as string),
		enabled: caseId !== null,
	});
}
