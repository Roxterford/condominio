import { useCallback, useState } from 'react';

export interface OverlayProps {
    open?: boolean;
    onOpenChange?: (open: boolean) => void;
}

export function useOverlay(initialState = false) {
    const [isOpen, setIsOpen] = useState(initialState);

    const open = useCallback(() => setIsOpen(true), []);
    const close = useCallback(() => setIsOpen(false), []);
    const toggle = useCallback(() => setIsOpen((prev) => !prev), []);

    return {
        isOpen,
        setIsOpen,
        open,
        close,
        toggle,
        overlayProps: {
            open: isOpen,
            onOpenChange: setIsOpen,
        } as OverlayProps,
    };
}
