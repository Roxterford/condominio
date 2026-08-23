import { AlertTriangle } from "lucide-react";
import { cn } from "@/lib/utils";
import type { GraphqlErrorDetail } from "@/providers/graphql/errors";

interface GraphqlErrorViewProps {
	errors: GraphqlErrorDetail[];
	className?: string;
}

export function GraphqlErrorView({ errors, className }: GraphqlErrorViewProps) {
	return (
		<div
			className={cn(
				"rounded-lg border border-destructive/40 bg-destructive/5 p-4",
				className,
			)}
			role="alert"
		>
			<div className="flex items-center gap-2 text-destructive">
				<AlertTriangle className="size-5" />
				<h2 className="text-sm font-semibold">
					Error de GraphQL
				</h2>
			</div>

			<ul className="mt-3 space-y-3">
				{errors.map((error, index) => (
					<li
						key={index}
						className="rounded-md border border-border bg-card p-3 text-sm"
					>
						<p className="font-medium text-foreground">
							{error.message}
						</p>

						{error.path && (
							<p className="mt-1 text-xs text-muted-foreground">
								<span className="font-semibold">path:</span>{" "}
								{error.path.join(".")}
							</p>
						)}

						{error.locations?.length && (
							<p className="mt-1 text-xs text-muted-foreground">
								<span className="font-semibold">location:</span>{" "}
								{error.locations
									.map(
										(l) =>
											`${l.line}:${l.column}`,
									)
									.join(", ")}
							</p>
						)}

						{error.extensions && (
							<pre className="mt-2 max-h-60 overflow-auto rounded bg-muted p-2 text-xs text-muted-foreground">
								{JSON.stringify(error.extensions, null, 2)}
							</pre>
						)}
					</li>
				))}
			</ul>
		</div>
	);
}
