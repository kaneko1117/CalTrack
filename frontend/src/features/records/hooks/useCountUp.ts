import { useState, useEffect } from "react";

type UseCountUpOptions = {
  end: number;
  duration?: number;
  startOnMount?: boolean;
};

/**
 * 数値をカウントアップアニメーションで表示するフック
 */
export function useCountUp({
  end,
  duration = 1000,
  startOnMount = true,
}: UseCountUpOptions) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    if (!startOnMount) return;

    let startTime: number;
    let animationFrame: number;

    const animate = (timestamp: number) => {
      if (!startTime) startTime = timestamp;
      const progress = Math.min((timestamp - startTime) / duration, 1);

      // easeOutQuart イージング
      const eased = 1 - Math.pow(1 - progress, 4);
      setCount(Math.floor(eased * end));

      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate);
      }
    };

    animationFrame = requestAnimationFrame(animate);

    return () => cancelAnimationFrame(animationFrame);
  }, [end, duration, startOnMount]);

  return count;
}
