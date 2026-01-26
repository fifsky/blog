import type { PropsWithChildren } from "react";

export default function Comp(props: PropsWithChildren) {
  return props.children as JSX.Element;
}

