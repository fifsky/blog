import SwiftUI

/// 原生照片浏览器页面
/// 支持左右滑动翻页、双击/捏合缩放
/// 通过 NavigationStack push 呈现，带返回按钮，符合系统原生交互
struct PhotoBrowserView: View {

    /// 待浏览的照片 URL 字符串列表
    let photoURLs: [String]

    /// 初始显示的照片索引
    let initialIndex: Int

    /// 底部页码指示器是否可见
    @State private var indexVisible = true

    /// 页码指示器隐藏定时任务
    @State private var hideIndexTask: Task<Void, Never>?

    /// 当前显示的照片索引
    @State private var currentIndex: Int

    init(photoURLs: [String], initialIndex: Int) {
        self.photoURLs = photoURLs
        self.initialIndex = initialIndex
        self._currentIndex = State(initialValue: initialIndex)
    }

    var body: some View {
        ZStack {
            // 纯黑背景，沉浸式浏览
            Color.black.ignoresSafeArea()

            // 照片翻页容器
            TabView(selection: $currentIndex) {
                ForEach(Array(photoURLs.enumerated()), id: \.offset) { index, url in
                    ZoomableImage(urlString: url)
                        .tag(index)
                }
            }
            .tabViewStyle(.page(indexDisplayMode: photoURLs.count > 1 ? .automatic : .never))
            .ignoresSafeArea()

            // 底部页码指示器
            VStack {
                Spacer()
                if photoURLs.count > 1 {
                    bottomBar
                }
            }
        }
        .navigationTitle("")
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(.regularMaterial, for: .navigationBar)
        .toolbarBackground(.visible, for: .navigationBar)
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
                .foregroundStyle(.white)
                .padding(.horizontal, 12)
                .padding(.vertical, 6)
                .background(.black.opacity(0.4), in: Capsule())
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

// MARK: - 可缩放单张图片

/// 支持双击缩放、捏合缩放的单张图片
private struct ZoomableImage: View {

    /// 图片 URL
    let urlString: String

    /// 当前缩放比例
    @State private var scale: CGFloat = 1.0

    /// 双击缩放后的目标比例
    @State private var targetScale: CGFloat = 1.0

    var body: some View {
        GeometryReader { geo in
            ZStack {
                AsyncImage(url: URL(string: urlString)) { phase in
                    switch phase {
                    case .success(let image):
                        image
                            .resizable()
                            .aspectRatio(contentMode: .fit)
                            .modifier(ZoomableModifier(
                                scale: $scale,
                                targetScale: $targetScale
                            ))
                            .frame(width: geo.size.width, height: geo.size.height)
                    case .failure:
                        VStack(spacing: 8) {
                            Image(systemName: "photo")
                                .font(.system(size: 40))
                            Text("加载失败")
                        }
                        .foregroundStyle(.gray)
                    default:
                        ProgressView()
                            .tint(.white)
                    }
                }
            }
            .frame(width: geo.size.width, height: geo.size.height)
        }
        // 重置缩放状态，避免翻页时残留上一张的缩放
        .id(urlString)
    }
}

// MARK: - 缩放手势修饰器

/// 封装双击与捏合缩放手势的 ViewModifier
private struct ZoomableModifier: ViewModifier {

    /// 当前缩放比例
    @Binding var scale: CGFloat

    /// 双击目标缩放比例
    @Binding var targetScale: CGFloat

    /// 双击放大后的目标比例
    private let maxScale: CGFloat = 3.0

    func body(content: Content) -> some View {
        content
            .scaleEffect(scale)
            // 双击切换原始尺寸与放大尺寸
            .onTapGesture(count: 2) {
                withAnimation(.spring()) {
                    if scale > 1 {
                        scale = 1
                        targetScale = 1
                    } else {
                        scale = maxScale
                        targetScale = maxScale
                    }
                }
            }
            // 捏合缩放
            .gesture(
                MagnifyGesture()
                    .onChanged { value in
                        let delta = value.magnification / targetScale
                        scale = min(max(scale * delta, 1), maxScale)
                    }
                    .onEnded { _ in
                        targetScale = scale
                        // 缩小到接近 1 时复位
                        if scale < 1.1 {
                            withAnimation(.spring()) {
                                scale = 1
                                targetScale = 1
                            }
                        }
                    }
            )
    }
}
