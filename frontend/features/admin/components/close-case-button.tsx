"use client";

import { Button } from "@/shared/components/button";

import { useCloseCase } from "../hooks/use-all-cases";

type CloseCaseButtonProps = {
	caseId: string;
};

export function CloseCaseButton({ caseId }: CloseCaseButtonProps) {
	const closeMutation = useCloseCase(caseId);

	return (
		<Button
			variant="danger"
			busy={closeMutation.isPending}
			onClick={() => {
				if (window.confirm("Are you sure you want to close this case?")) {
					closeMutation.mutate();
				}
			}}
		>
			Close Case
		</Button>
	);
}
