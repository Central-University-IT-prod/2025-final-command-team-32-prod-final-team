import { cloneElement } from "react";

import "./error.css";

export function ErrorWrapper(props) {
  const { children, error, ...rest } = props;
  return <>
    {cloneElement(children, { error: error, ...rest })}
    <div className={"error-text"}>{error}</div>
  </>;
}