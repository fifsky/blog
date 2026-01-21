import { useState, useEffect, useRef, ReactNode } from "react";

interface PageTransitionProps {
  children: ReactNode;
  loading?: boolean;
  className?: string;
}

export function PageTransition({ children, loading = false, className = "" }: PageTransitionProps) {
  const [isVisible, setIsVisible] = useState(true);
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (loading) {
      setIsVisible(false);
    } else {
      // Small delay for smooth transition
      const timer = setTimeout(() => setIsVisible(true), 50);
      return () => clearTimeout(timer);
    }
  }, [loading]);

  return (
    <div
      ref={contentRef}
      className={`transition-opacity duration-300 ease-in-out ${isVisible ? "opacity-100" : "opacity-0"} ${className}`}
    >
      {children}
    </div>
  );
}
