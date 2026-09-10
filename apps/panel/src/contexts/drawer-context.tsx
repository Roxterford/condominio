'use client'

import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'

type DrawerSide = 'right' | 'left' | 'bottom'

type DrawerOptions = {
	content?: ReactNode
	loader?: () => Promise<ReactNode>
	title?: string
	titleBadge?: ReactNode
	side?: DrawerSide
	size?: number | string
}

type DrawerState = {
	isOpen: boolean
	content: ReactNode | null
	title: string
	titleBadge?: ReactNode
	side: DrawerSide
	size: number | string
	loading: boolean
}

type DrawerContextType = {
	state: DrawerState
	open: (options: DrawerOptions) => void
	close: () => void
}

const INITIAL_STATE: DrawerState = {
	isOpen: false,
	content: null,
	title: '',
	side: 'right',
	size: 400,
	loading: false,
}

const DrawerContext = createContext<DrawerContextType | null>(null)

export function DrawerProvider({ children }: { children: ReactNode }) {
	const [state, setState] = useState<DrawerState>(INITIAL_STATE)

	const open = useCallback((options: DrawerOptions) => {
		const { content, loader, title = '', titleBadge, side = 'right', size = 400 } = options

		if (loader) {
			setState({ isOpen: true, content: null, title, titleBadge, side, size, loading: true })
			loader().then((resolved) => {
				setState((prev) => ({ ...prev, content: resolved, loading: false }))
			})
		} else {
			setState({ isOpen: true, content: content ?? null, title, titleBadge, side, size, loading: false })
		}
	}, [])

	const close = useCallback(() => {
		setState(INITIAL_STATE)
	}, [])

	return (
		<DrawerContext value={{ state, open, close }}>
			{children}
		</DrawerContext>
	)
}

export function useDrawer() {
	const ctx = useContext(DrawerContext)
	if (!ctx) throw new Error('useDrawer debe usarse dentro de un DrawerProvider')
	return ctx
}
