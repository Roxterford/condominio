'use client'

import { X } from 'lucide-react'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { Skeleton } from '@/components/ui/skeleton'
import { useDrawer } from '@/contexts/drawer-context'

function DrawerSkeleton() {
	return (
		<div className="space-y-4 p-6">
			<Skeleton className="h-6 w-3/4" />
			<Skeleton className="h-4 w-full" />
			<Skeleton className="h-4 w-5/6" />
			<Skeleton className="h-32 w-full" />
			<Skeleton className="h-4 w-2/3" />
		</div>
	)
}

export function DynamicDrawer() {
	const { state, close } = useDrawer()

	const width = typeof state.size === 'number' ? `${state.size}px` : state.size

	return (
		<Sheet open={state.isOpen} onOpenChange={(open) => !open && close()}>
			<SheetContent
				side={state.side}
				showCloseButton={false}
				className="gap-0"
				style={{ maxWidth: state.side !== 'bottom' ? width : undefined }}
			>
				<SheetHeader className="px-8">
					<div className="flex items-center justify-between">
<SheetTitle className="font-semibold text-lg">
	<span className="flex items-center gap-2">
		{state.title || 'Detalle'}
		{state.titleBadge}
	</span>
</SheetTitle>
						<button
							onClick={close}
							className="text-gray-400 hover:text-gray-600 transition-colors"
						>
							<X size={20} />
						</button>
					</div>
				</SheetHeader>

				<div className="flex-1 overflow-y-auto">
					{state.loading ? <DrawerSkeleton /> : state.content}
				</div>
			</SheetContent>
		</Sheet>
	)
}
