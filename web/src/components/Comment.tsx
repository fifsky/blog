import Giscus from "@giscus/react";

export function Comment({ postId }: { postId: number }) {
  // const { pathname } = useLocation();
  console.log(postId)
  return <Giscus
      id="comments"
      repo="fifsky/blog"
      repoId="MDEwOlJlcG9zaXRvcnkyMDU5ODYyODc="
      category="Comment"
      categoryId="DIC_kwDODEcZ784C0rJu"
      mapping="pathname"
      term="Hello"
      reactionsEnabled="1"
      emitMetadata="0"
      inputPosition="top"
      theme="light_protanopia"
      lang="zh-CN"
      loading="lazy"
    />;
}
