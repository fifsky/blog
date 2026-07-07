import SwiftUI

/// 通用加载占位视图
/// 以透明背景撑满父容器后居中展示旋转指示器与文案，供各列表页首次加载时统一使用
struct LoadingView: View {

    /// 加载提示文案，默认"加载中..."
    var text: String = "加载中..."

    var body: some View {
        // ScrollView 内拿不到全屏高度，用固定上边距把内容压到视觉中下区域，避开导航栏/搜索框
        VStack(spacing: 12) {
            ProgressView()
                .controlSize(.large)
                .tint(.secondary)
            Text(text)
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity)
        .padding(.top, 200)
    }
}

#Preview {
    LoadingView()
}
