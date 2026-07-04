import SwiftUI

/// 通用加载占位视图
/// 居中展示旋转指示器与文案，供各列表页首次加载时统一使用
struct LoadingView: View {

    /// 加载提示文案，默认"加载中..."
    var text: String = "加载中..."

    var body: some View {
        VStack(spacing: 12) {
            ProgressView()
                .controlSize(.large)
            Text(text)
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}

#Preview {
    LoadingView()
}
