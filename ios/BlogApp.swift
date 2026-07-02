import SwiftUI

@main
struct BlogApp: App {

    @State private var isAuthenticated = false

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
            .task {
                // 启动时检查登录状态：Token 存于 Keychain（持久化），未过期则直接进入主页
                if AuthManager.shared.isLoggedIn {
                    isAuthenticated = true
                }
            }
        }
    }
}
