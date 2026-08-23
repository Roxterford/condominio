import type { ReactNode } from "react";
import { GraphqlErrorView } from "@/components/graphql-error-view";
import {
	GraphqlError,
	GraphqlErrorDetail,
	handleGraphError,
	isDebugMode,
} from "@/providers/graphql/errors";

interface GraphqlResultLike<T> {
	data: T | null;
	errors?: GraphqlErrorDetail[];
}

export async function renderGraphql<T>(
	result: GraphqlResultLike<T>,
	render: (data: T) => ReactNode,
): Promise<ReactNode> {
	try {
		const data = handleGraphError(result);
		return render(data);
	} catch (error) {
		if (isDebugMode()) {
			const errors =
				error instanceof GraphqlError
					? error.errors
					: [{ message: error instanceof Error ? error.message : String(error) }];
			return <GraphqlErrorView errors={errors} />;
		}
		throw error;
	}
}
