import { useCallback, useRef } from "react";
import { useLocation } from "react-router";
import Artalk from "artalk";
import "artalk/Artalk.css";

export function Comment({ postId }: { postId: number }) {
  const { pathname } = useLocation();
  const artalkRef = useRef<Artalk>();

  const initContainer = useCallback(
    (node: HTMLDivElement | null) => {
      if (!node) return;
      if (artalkRef.current) {
        artalkRef.current.destroy();
      }
      artalkRef.current = Artalk.init({
        el: node,
        pageKey: pathname || String(postId),
        pageTitle: document.title,
        server: "https://comment-api.fifsky.com",
        site: "FIFSKY",
        useBackendConf: true,
      });
    },
    [pathname, postId]
  );

  return <div className="comment" ref={initContainer}></div>;
}
