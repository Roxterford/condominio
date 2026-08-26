export interface GraphqlErrorDetail {
	message: string;
	locations?: Array<{ line: number; column: number }>;
	path?: string[];
	extensions?: Record<string, unknown>;
}

const GENERIC_MESSAGE = "Algo salió mal. Inténtalo de nuevo más tarde.";

export function isDebugMode(): boolean {
	return process.env.NODE_ENV !== "production";
}

export class GraphqlError extends Error {
	readonly errors: GraphqlErrorDetail[];

	constructor(errors: GraphqlErrorDetail[]) {
		super(
			errors[0]?.message ?? "Error en la petición GraphQL",
		);
		this.name = "GraphqlError";
		this.errors = errors;
	}
}

export function handleGraphError<T>(result: {
	data: T | null;
	errors?: GraphqlErrorDetail[];
}): T {
	if (result.errors?.length) {
		if (isDebugMode()) {
			throw new GraphqlError(result.errors);
		}
		throw new Error(GENERIC_MESSAGE);
	}

	if (!result.data) {
		if (isDebugMode()) {
			throw new GraphqlError([
				{ message: "La petición no retornó datos" },
			]);
		}
		throw new Error(GENERIC_MESSAGE);
	}

	return result.data;
}
