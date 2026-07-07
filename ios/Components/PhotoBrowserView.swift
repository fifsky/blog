import SwiftUI

/// 原生照片浏览器页面
/// 支持左右滑动翻页、双击/捏合缩放、放大后拖动查看细节
/// 通过 NavigationStack push 呈现，带返回按钮，符合系统原生交互
struct PhotoBrowserView: View {

    /// 待浏览的照片 URL 字符串列表
    let photoURLs: [String]

    /// 初始显示的照片索引
    let initialIndex: Int

    /// 地点名称（用作页面标题）
    let placeName: String

    /// 底部页码指示器是否可见
    @State private var indexVisible = true

    /// 页码指示器隐藏定时任务
    @State private var hideIndexTask: Task<Void, Never>?

    /// 当前显示的照片索引
    @State private var currentIndex: Int

    init(photoURLs: [String], initialIndex: Int, placeName: String) {
        self.photoURLs = photoURLs
        self.initialIndex = initialIndex
        self.placeName = placeName
        self._currentIndex = State(initialValue: initialIndex)
    }

    var body: some View {
        ZStack {
            // 浅灰背景
            Color(.systemGroupedBackground).ignoresSafeArea()

            // 照片翻页容器
            TabView(selection: $currentIndex) {
                ForEach(Array(photoURLs.enumerated()), id: \.offset) { index, url in
                    ZoomableImageView(urlString: url)
                        .tag(index)
                }
            }
            .tabViewStyle(.page(indexDisplayMode: photoURLs.count > 1 ? .automatic : .never))

            // 底部页码指示器
            VStack {
                Spacer()
                if photoURLs.count > 1 {
                    bottomBar
                }
            }
        }
        .navigationTitle(placeName)
        .navigationBarTitleDisplayMode(.inline)
        // 导航栏背景与内容区统一灰底，从顶到底一致
        .toolbarBackground(.visible, for: .navigationBar)
        .toolbarBackground(Color(.systemGroupedBackground), for: .navigationBar)
        .toolbar(.hidden, for: .tabBar)
        .statusBarHidden()
        .onAppear {
            scheduleHideIndex()
        }
        .onDisappear {
            hideIndexTask?.cancel()
        }
    }

    // MARK: - 底部页码指示器

    /// 底部页码指示器
    @ViewBuilder
    private var bottomBar: some View {
        if indexVisible {
            Text("\(currentIndex + 1) / \(photoURLs.count)")
                .font(.subheadline)
                .foregroundStyle(.secondary)
                .padding(.horizontal, 12)
                .padding(.vertical, 6)
                .background(.regularMaterial, in: Capsule())
                .padding(.bottom, 24)
                .transition(.opacity)
        }
    }

    // MARK: - 页码指示器显隐

    /// 延迟隐藏页码指示器（3 秒后自动隐藏，避免遮挡照片）
    private func scheduleHideIndex() {
        hideIndexTask?.cancel()
        hideIndexTask = Task {
            try? await Task.sleep(nanoseconds: 3_000_000_000)
            if !Task.isCancelled {
                withAnimation(.easeInOut) {
                    indexVisible = false
                }
            }
        }
    }
}

// MARK: - 可缩放图片视图

/// 基于 UIScrollView 的可缩放图片视图
///
/// 使用 UIScrollView 实现原生相册般的缩放与拖动体验：
/// - 双击放大/还原（点击位置为中心放大）
/// - 捏合缩放（1x ~ 4x）
/// - 放大后可自由拖动查看图片任意细节
/// - 缩放至 1x 时 TabView 恢复翻页手势
private struct ZoomableImageView: UIViewRepresentable {

    /// 图片 URL
    let urlString: String

    func makeUIView(context: Context) -> UIScrollView {
        let scrollView = UIScrollView()
        scrollView.delegate = context.coordinator
        scrollView.maximumZoomScale = 4.0
        scrollView.minimumZoomScale = 1.0
        scrollView.bouncesZoom = true
        scrollView.showsHorizontalScrollIndicator = false
        scrollView.showsVerticalScrollIndicator = false
        scrollView.backgroundColor = .clear
        scrollView.contentInsetAdjustmentBehavior = .never

        // 图片容器视图（缩放对象）
        let imageView = UIImageView()
        imageView.contentMode = .scaleAspectFit
        imageView.backgroundColor = .clear
        imageView.clipsToBounds = true
        scrollView.addSubview(imageView)
        context.coordinator.imageView = imageView

        // 双击缩放手势
        let doubleTap = UITapGestureRecognizer(target: context.coordinator, action: #selector(Coordinator.handleDoubleTap(_:)))
        doubleTap.numberOfTapsRequired = 2
        scrollView.addGestureRecognizer(doubleTap)

        // 加载图片
        context.coordinator.loadImage(urlString: urlString, scrollView: scrollView)

        return scrollView
    }

    func updateUIView(_ scrollView: UIScrollView, context: Context) {
        guard let imageView = context.coordinator.imageView else { return }
        // URL 变化时重新加载
        if context.coordinator.currentURLString != urlString {
            context.coordinator.currentURLString = urlString
            scrollView.zoomScale = 1.0
            context.coordinator.loadImage(urlString: urlString, scrollView: scrollView)
        }
        // 确保 imageView 尶寸与 scrollView 一致
        let bounds = scrollView.bounds
        if bounds.width > 0 && bounds.height > 0 {
            imageView.frame = bounds
        }
    }

    func makeCoordinator() -> Coordinator {
        Coordinator()
    }

    // MARK: - Coordinator

    final class Coordinator: NSObject, UIScrollViewDelegate {
        weak var imageView: UIImageView?
        var currentURLString: String?
        private var loadTask: Task<Void, Never>?

        /// 异步加载图片
        func loadImage(urlString: String, scrollView: UIScrollView) {
            currentURLString = urlString
            loadTask?.cancel()
            imageView?.image = nil
            loadTask = Task {
                guard let url = URL(string: urlString) else { return }
                do {
                    let (data, _) = try await URLSession.shared.data(from: url)
                    guard !Task.isCancelled, let image = UIImage(data: data) else { return }
                    await MainActor.run {
                        guard !Task.isCancelled, let imageView = self.imageView else { return }
                        imageView.image = image
                        // 设置 scrollView 的 contentSize 与缩放
                        let boundsSize = scrollView.bounds.size
                        if boundsSize.width > 0 {
                            imageView.frame = CGRect(origin: .zero, size: boundsSize)
                            scrollView.contentSize = boundsSize
                        }
                    }
                } catch {
                    // 加载失败，忽略
                }
            }
        }

        // MARK: - UIScrollViewDelegate

        /// 返回缩放视图
        func viewForZooming(in scrollView: UIScrollView) -> UIView? {
            imageView
        }

        /// 缩放时居中内容
        func scrollViewDidZoom(_ scrollView: UIScrollView) {
            guard let imageView = imageView, imageView.image != nil else { return }
            let boundsSize = scrollView.bounds.size
            let contentSize = imageView.frame.size

            let horizontalInset = max(0, (boundsSize.width - contentSize.width) / 2)
            let verticalInset = max(0, (boundsSize.height - contentSize.height) / 2)

            scrollView.contentInset = UIEdgeInsets(
                top: verticalInset,
                left: horizontalInset,
                bottom: verticalInset,
                right: horizontalInset
            )
        }

        // MARK: - 双击缩放

        /// 双击切换 1x 与 2.5x 缩放（以点击位置为中心）
        @objc func handleDoubleTap(_ recognizer: UITapGestureRecognizer) {
            guard let scrollView = recognizer.view as? UIScrollView else { return }
            if scrollView.zoomScale > 1.0 {
                scrollView.setZoomScale(1.0, animated: true)
            } else {
                let location = recognizer.location(in: scrollView)
                let targetScale: CGFloat = 2.5
                let zoomWidth = scrollView.bounds.width / targetScale
                let zoomHeight = scrollView.bounds.height / targetScale
                let zoomRect = CGRect(
                    x: location.x - zoomWidth / 2,
                    y: location.y - zoomHeight / 2,
                    width: zoomWidth,
                    height: zoomHeight
                )
                scrollView.zoom(to: zoomRect, animated: true)
            }
        }
    }
}
