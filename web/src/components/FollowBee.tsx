import { useEffect, useRef, useState, useCallback } from "react";

export function FollowBee() {
  const beeRef = useRef<HTMLDivElement>(null);
  const [isFlying, setIsFlying] = useState(false);
  const targetRef = useRef({ x: 100, y: 100 });
  const positionRef = useRef({ x: 100, y: 100 });
  const animationRef = useRef<number | null>(null);

  const animate = useCallback(() => {
    if (!beeRef.current) return;

    const pos = positionRef.current;
    const target = targetRef.current;
    const dx = target.x - pos.x;
    const dy = target.y - pos.y;
    const distance = Math.sqrt(dx * dx + dy * dy);

    if (distance < 2) {
      pos.x = target.x;
      pos.y = target.y;
      beeRef.current.style.left = `${pos.x}px`;
      beeRef.current.style.top = `${pos.y}px`;
      setIsFlying(false);
      animationRef.current = null;
      return;
    }

    const speed = Math.min(distance * 0.08, 15);
    const vx = (dx / distance) * speed;
    const vy = (dy / distance) * speed;
    pos.x += vx;
    pos.y += vy;

    beeRef.current.style.left = `${pos.x}px`;
    beeRef.current.style.top = `${pos.y}px`;
    beeRef.current.style.transform = dx < 0 ? "scaleX(-1)" : "scaleX(1)";

    animationRef.current = requestAnimationFrame(animate);
  }, []);

  const handleMouseMove = useCallback(
    (e: MouseEvent) => {
      targetRef.current = { x: e.clientX - 20, y: e.clientY - 15 };
      if (!animationRef.current) {
        setIsFlying(true);
        animationRef.current = requestAnimationFrame(animate);
      }
    },
    [animate]
  );

  useEffect(() => {
    window.addEventListener("mousemove", handleMouseMove);
    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [handleMouseMove]);

  return (
    <div
      ref={beeRef}
      className="fixed pointer-events-none z-[9999]"
      style={{
        left: positionRef.current.x,
        top: positionRef.current.y,
      }}
    >
      <svg
        width="40"
        height="32"
        viewBox="0 0 40 32"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        {isFlying ? (
          <>
            <ellipse
              cx="8"
              cy="8"
              rx="6"
              ry="4"
              fill="rgba(255,200,0,0.7)"
              className="origin-[8px_8px] animate-[flap-left_0.08s_ease-in-out_infinite_alternate]"
            />
            <ellipse
              cx="8"
              cy="24"
              rx="6"
              ry="4"
              fill="rgba(255,200,0,0.7)"
              className="origin-[8px_24px] animate-[flap-left_0.08s_ease-in-out_infinite_alternate]"
            />
            <ellipse
              cx="32"
              cy="8"
              rx="6"
              ry="4"
              fill="rgba(255,200,0,0.7)"
              className="origin-[32px_8px] animate-[flap-right_0.08s_ease-in-out_infinite_alternate]"
            />
            <ellipse
              cx="32"
              cy="24"
              rx="6"
              ry="4"
              fill="rgba(255,200,0,0.7)"
              className="origin-[32px_24px] animate-[flap-right_0.08s_ease-in-out_infinite_alternate]"
            />
          </>
        ) : (
          <>
            <ellipse cx="8" cy="12" rx="6" ry="3" fill="rgba(255,200,0,0.6)" />
            <ellipse cx="8" cy="20" rx="6" ry="3" fill="rgba(255,200,0,0.6)" />
            <ellipse cx="32" cy="12" rx="6" ry="3" fill="rgba(255,200,0,0.6)" />
            <ellipse cx="32" cy="20" rx="6" ry="3" fill="rgba(255,200,0,0.6)" />
          </>
        )}
        <ellipse cx="20" cy="16" rx="10" ry="8" fill="#FFD700" />
        <rect x="16" y="10" width="8" height="3" fill="#333" />
        <rect x="16" y="16" width="8" height="3" fill="#333" />
        <ellipse cx="14" cy="14" rx="3" ry="4" fill="#FFD700" />
        <circle cx="13" cy="12" r="1.5" fill="#333" />
        <circle cx="12.5" cy="11.5" r="0.5" fill="white" />
        <path d="M8 10 Q6 8 8 6" stroke="#333" strokeWidth="0.5" fill="none" />
        <path d="M10 9 Q8 7 10 5" stroke="#333" strokeWidth="0.5" fill="none" />
        <path d="M32 10 Q34 8 32 6" stroke="#333" strokeWidth="0.5" fill="none" />
        <path d="M30 9 Q32 7 30 5" stroke="#333" strokeWidth="0.5" fill="none" />
      </svg>
      <style>{`
        @keyframes flap-left {
          from { transform: rotate(-30deg) translateY(-2px); }
          to { transform: rotate(30deg) translateY(2px); }
        }
        @keyframes flap-right {
          from { transform: rotate(30deg) translateY(-2px); }
          to { transform: rotate(-30deg) translateY(2px); }
        }
      `}</style>
    </div>
  );
}
