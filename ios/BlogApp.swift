import SwiftUI

@main
struct BlogApp: App {

    @State private var isAuthenticated: Bool

    init() {
        // 启动时同步检查登录态：Keychain 读取是同步且快速的，
        // 让首帧即以正确状态渲染，避免冷启动闪过登录页（其 onAppear 会拉起键盘）
        _isAuthenticated = State(initialValue: AuthManager.shared.isLoggedIn)
    }

    var body: some Scene {
        WindowGroup {
            Group {
                if isAuthenticated {
                    ContentView()
                } else {
                    LoginView { viewModel in
                        // 登录成功后切换到主界面
                        isAuthenticated = true
                    }
                }
            }
            .onReceive(NotificationCenter.default.publisher(for: .didLogout)) { _ in
                // 收到退出登录通知后切回登录页
                isAuthenticated = false
            }
        }
    }
}
