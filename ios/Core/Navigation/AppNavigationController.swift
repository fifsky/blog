import SwiftUI
import UIKit

/// 单个 Tab 内部的 UIKit 导航栈
final class AppNavigationController: UINavigationController {

    let appNavigator: AppNavigator

    /// 初始化一个承载 SwiftUI 根页面的导航栈
    /// - Parameters:
    ///   - rootView: Tab 根页面
    ///   - title: TabBar 标题
    ///   - systemImage: TabBar 图标
    ///   - rootTitle: 根页面导航栏标题
    init<Root: View>(
        rootView: Root,
        title: String,
        systemImage: String,
        rootTitle: String? = nil
    ) {
        let navigator = AppNavigator()
        self.appNavigator = navigator

        let controller = NavigationHostingController(
            rootView: rootView
                .environment(\.appNavigator, navigator)
                .tint(Color.themePrimary)
        )
        controller.title = rootTitle ?? title
        controller.navigationItem.largeTitleDisplayMode = .always
        controller.prefersBottomBarHidden = false
        super.init(rootViewController: controller)

        navigator.navigationController = self
        delegate = self
        navigationBar.prefersLargeTitles = true
        tabBarItem = UITabBarItem(
            title: title,
            image: UIImage(systemName: systemImage),
            selectedImage: UIImage(systemName: systemImage)
        )
    }

    @available(*, unavailable)
    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
}

extension AppNavigationController: UINavigationControllerDelegate {

    func navigationController(
        _ navigationController: UINavigationController,
        willShow viewController: UIViewController,
        animated: Bool
    ) {
        let shouldHide = shouldHideBottomBar(for: viewController)
        guard let tabBar = tabBarController?.tabBar else { return }

        if shouldHide {
            // 隐藏 TabBar：瞬间隐藏（与原行为一致），不做淡出。
            // 瞬间释放 safeArea，详情页布局立即正确，评论框不会被顶上去。
            // 同时把 alpha 置 0，为下次 pop 淡入做准备。
            tabBar.alpha = 0
            tabBar.isHidden = true
            tabBar.isUserInteractionEnabled = false
            return
        }

        // 显示 TabBar：先恢复到视图层级（透明），再在转场中淡入。
        // isHidden=false 在 willShow 阶段恢复 safeArea，此时根页面正在滑入，
        // 布局调整被滑动过程掩盖；alpha 从 0 渐变到 1，实现淡入。
        tabBar.isHidden = false
        tabBar.isUserInteractionEnabled = true

        if let coordinator = transitionCoordinator, animated {
            tabBar.alpha = 0
            coordinator.animate(alongsideTransition: { _ in
                tabBar.alpha = 1
            }, completion: { [weak self] context in
                guard let self, let tabBar = self.tabBarController?.tabBar else { return }
                if context.isCancelled, let topVC = self.topViewController,
                   self.shouldHideBottomBar(for: topVC) {
                    // 交互式返回取消，恢复到栈顶页面（详情页）的隐藏状态
                    tabBar.alpha = 0
                    tabBar.isHidden = true
                    tabBar.isUserInteractionEnabled = false
                } else {
                    tabBar.alpha = 1
                }
            })
        } else {
            // 非动画转场（如程序化无动画切换），直接显示
            tabBar.alpha = 1
        }
    }

    /// 判断指定页面是否需要隐藏底部 TabBar
    private func shouldHideBottomBar(for viewController: UIViewController) -> Bool {
        guard viewController !== viewControllers.first else { return false }
        return (viewController as? BottomBarVisibilityProviding)?.prefersBottomBarHidden ?? true
    }
}
