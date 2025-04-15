import { useEffect, useState } from "react";

export function useDevice() {
  const [mobile, setMobile] = useState(window.innerWidth <= 720);

  function handleWindowSizeChange() {
    setMobile(window.innerWidth <= 720);
  }

  useEffect(() => {
    window.addEventListener("resize", handleWindowSizeChange);
    return () => {
      window.removeEventListener("resize", handleWindowSizeChange);
    };
  }, []);
  return mobile;
}