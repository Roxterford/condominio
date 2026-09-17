import { useCallback } from 'react';

export function useSmoothScrollToTop() {
  return useCallback(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);
}
