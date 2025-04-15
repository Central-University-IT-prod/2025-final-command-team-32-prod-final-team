import * as React from "react";
import { useEffect, useRef } from "react";

type IntersectionObserverComponentProps = {
  onIntersect: () => Promise<void>;
};

const IntersectionObserverComponent: React.FC<IntersectionObserverComponentProps> = ({ onIntersect }) => {
  const ref = useRef<HTMLDivElement>(null);
  const isProcessing = useRef(false);

  useEffect(() => {
    const observer = new IntersectionObserver(async (entries) => {
      if (entries[0].isIntersecting && !isProcessing.current) {
        isProcessing.current = true;
        try {
          await onIntersect();
        } finally {
          isProcessing.current = false;
        }
      }
    });

    if (ref.current) {
      observer.observe(ref.current);
    }

    return () => {
      if (ref.current) {
        observer.unobserve(ref.current);
      }
    };
  }, [onIntersect]);

  return <div ref={ref} style={{ width: "100%", height: "10px" }} />;
};

export default IntersectionObserverComponent;
