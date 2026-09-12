"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { usePathname, useSearchParams } from "next/navigation";
import { cn } from "@/lib/utils";

const MIN_DURATION = 300;
const SAFETY_TIMEOUT = 10_000;

export function TopProgressBar() {
  const [progress, setProgress] = useState(0);
  const [visible, setVisible] = useState(false);

  const progressRef = useRef(0);
  const activeRef = useRef(false);
  const startedAtRef = useRef(0);
  const rafRef = useRef<number | null>(null);
  const safetyTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const hideTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const pathname = usePathname();
  const searchParams = useSearchParams();
  const urlKey = `${pathname}?${searchParams.toString()}`;
  const lastUrlKeyRef = useRef(urlKey);

  const clearRaf = useCallback(() => {
    if (rafRef.current !== null) {
      cancelAnimationFrame(rafRef.current);
      rafRef.current = null;
    }
  }, []);

  const clearTimers = useCallback(() => {
    if (safetyTimerRef.current !== null) {
      clearTimeout(safetyTimerRef.current);
      safetyTimerRef.current = null;
    }
    if (hideTimerRef.current !== null) {
      clearTimeout(hideTimerRef.current);
      hideTimerRef.current = null;
    }
  }, []);

  const reset = useCallback(() => {
    clearRaf();
    clearTimers();
    activeRef.current = false;
    progressRef.current = 0;
    setProgress(0);
    setVisible(false);
  }, [clearRaf, clearTimers]);

  const tick = useCallback((now: number) => {
    if (!activeRef.current) return;
    const elapsed = now - startedAtRef.current;
    const eased = 88 * (1 - Math.exp(-elapsed / 1800));
    const next = Math.min(
      90,
      Math.max(progressRef.current + 0.1, eased + Math.random() * 1.5),
    );
    progressRef.current = next;
    setProgress(next);
    rafRef.current = requestAnimationFrame(tick);
  }, []);

  const done = useCallback(() => {
    if (!activeRef.current) return;
    const elapsed = performance.now() - startedAtRef.current;
    if (elapsed < MIN_DURATION) {
      reset();
      return;
    }
    clearRaf();
    if (safetyTimerRef.current !== null) {
      clearTimeout(safetyTimerRef.current);
      safetyTimerRef.current = null;
    }
    progressRef.current = 100;
    setProgress(100);
    hideTimerRef.current = setTimeout(() => {
      activeRef.current = false;
      progressRef.current = 0;
      setProgress(0);
      setVisible(false);
    }, 250);
  }, [clearRaf, reset]);

  const start = useCallback(() => {
    if (activeRef.current) return;
    activeRef.current = true;
    startedAtRef.current = performance.now();
    progressRef.current = 2;
    setProgress(2);
    setVisible(true);
    clearRaf();
    clearTimers();
    safetyTimerRef.current = setTimeout(done, SAFETY_TIMEOUT);
    rafRef.current = requestAnimationFrame(tick);
  }, [clearRaf, clearTimers, done, tick]);

  useEffect(() => {
    function onClick(event: MouseEvent) {
      if (
        event.defaultPrevented ||
        event.button !== 0 ||
        event.metaKey ||
        event.ctrlKey ||
        event.shiftKey ||
        event.altKey
      )
        return;
      const target = (event.target as Element | null)?.closest("a");
      if (!target) return;
      const href = target.getAttribute("href");
      if (!href || href.startsWith("#") || target.hasAttribute("download")) {
        return;
      }
      if (target.target === "_blank") return;
      try {
        const url = new URL(href, window.location.href);
        if (url.origin !== window.location.origin) return;
      } catch {
        return;
      }
      start();
    }

    const onPopState = () => start();

    const originalPushState = window.history.pushState.bind(window.history);
    const originalReplaceState = window.history.replaceState.bind(window.history);

    window.history.pushState = (...args: Parameters<typeof originalPushState>) => {
      originalPushState(...args);
      start();
    };
    window.history.replaceState = (
      ...args: Parameters<typeof originalReplaceState>
    ) => {
      originalReplaceState(...args);
      start();
    };

    document.addEventListener("click", onClick, true);
    window.addEventListener("popstate", onPopState);

    return () => {
      document.removeEventListener("click", onClick, true);
      window.removeEventListener("popstate", onPopState);
      window.history.pushState = originalPushState;
      window.history.replaceState = originalReplaceState;
    };
  }, [start]);

  useEffect(() => {
    if (lastUrlKeyRef.current !== urlKey) {
      lastUrlKeyRef.current = urlKey;
      done();
    }
  }, [urlKey, done]);

  return (
    <div
      aria-hidden
      className={cn(
        "pointer-events-none fixed inset-x-0 top-0 z-[100] transition-opacity duration-200",
        visible ? "opacity-100" : "opacity-0",
      )}
    >
      <div className="relative h-[3px]">
        <div
          className="absolute left-0 top-0 h-full overflow-hidden rounded-r-full bg-gradient-to-r from-teal-600 via-teal-400 to-cyan-400 shadow-[0_0_12px_rgba(20,184,166,0.7)]"
          style={{ width: `${progress}%` }}
        >
          <span className="top-progress-sheen" />
        </div>
        <div
          className="absolute left-0 top-0 h-full rounded-r-full bg-teal-400/50 blur-[6px]"
          style={{ width: `${progress}%` }}
        />
      </div>
    </div>
  );
}