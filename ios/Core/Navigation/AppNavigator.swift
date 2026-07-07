import SwiftUI
import UIKit

/// SwiftUI 页面使用的统一导航入口
///
/// 主导航由 UIKit 的 `UINavigationController` 承担，SwiftUI 页面只负责声明目标页面。
/// 除根页面外，push 默认隐藏 TabBar；隐藏动作由 `AppNavigationController` 统一处理。
final class AppNavigator {

    weak var navigationController: UINavigationController?

    /// Push 到新页面
    /// - Parameters:
    ///   - view: SwiftUI 目标页面
    ///   - hidesBottomBar: 是否隐藏底部 TabBar，默认隐藏
    ///   - animated: 是否启用系统 push 动画
    ///   - onPop: 页面返回后的回调
    @MainActor
    func push<Content: View>(
        _ view: Content,
        hidesBottomBar: Bool = true,
        animated: Bool = true,
        onPop: (() -> Void)? = nil
    ) {
        guard let navigationController else { return }

        let controller = NavigationHostingController(
            rootView: view
                .environment(\.appNavigator, self)
                .tint(Color.themePrimary)
        )
        controller.prefersBottomBarHidden = hidesBottomBar
        controller.onPop = onPop
        navigationController.pushViewController(controller, animated: animated)
    }

    /// 返回上一页
    /// - Parameter animated: 是否启用系统 pop 动画
    /// - Returns: 是否成功从 UIKit 导航栈返回
    @MainActor
    @discardableResult
    func pop(animated: Bool = true) -> Bool {
        guard let navigationController,
              navigationController.viewControllers.count > 1 else {
            return false
        }

        navigationController.popViewController(animated: animated)
        return true
    }
}

private struct AppNavigatorKey: EnvironmentKey {
    static let defaultValue = AppNavigator()
}

extension EnvironmentValues {
    /// 当前 Tab 对应的 UIKit 导航器
    var appNavigator: AppNavigator {
        get { self[AppNavigatorKey.self] }
        set { self[AppNavigatorKey.self] = newValue }
    }
}
