# AI任务

以下记录后续需要AI完成的任务

## iOS开发

2026-07-04

### 文章详情页 Markdown 图片首次点击无响应

- **现象**：进入文章详情页后，在不滚动页面的情况下直接点击 Markdown 渲染的图片无响应，必须先滚动一次页面才能正常点击
- **根因**：Textual 库的 `TextLinkInteraction` 通过 `overlayPreferenceValue(Text.LayoutKey.self)` 读取 `Text` 发布的 layout 来建立点击区域。`TextFragment` 的 `textBuilder`（`@State`）初始为 nil，首帧 `text` 返回空，layout preference 为空，overlay 不创建手势。ScrollView 首次 layout 未完成时 preference 传播延迟，需滚动触发 re-layout 后 overlay 才收到 layout
- **涉及文件**：
  - `ios/Features/Article/ArticleDetailView.swift` — `StructuredText` 渲染及 `openURL` 图片点击处理
  - Textual 库 `TextLinkInteraction.swift` — `overlayPreferenceValue` 读取 `Text.LayoutKey`
  - Textual 库 `TextFragment.swift` — `textBuilder` 延迟初始化导致首帧 layout 为空
- **待选方案**：
  1. 在 `StructuredText` 上加 `.geometryGroup()` 强制同步 geometry/layout（iOS 17+）
  2. `onAppear` 后延迟触发 re-render，让 `overlayPreferenceValue` 重新评估 preference
  3. 触发 `.textContainer` 坐标空间 geometry 变化，让 `onGeometryChange` 重新调用 `textBuilder.sizeChanged`
  4. 如以上均无效，可能需 fork Textual 修改 `TextLinkInteraction` 或 `TextFragment` 的初始化时序
- **状态**：待处理
