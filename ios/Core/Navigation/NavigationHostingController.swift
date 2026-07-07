import SwiftUI

/// 页面底栏显示偏好
protocol BottomBarVisibilityProviding: AnyObject {
    /// 当前页面显示时是否隐藏底部 TabBar
    var prefersBottomBarHidden: Bool { get }
}

/// 带返回回调的 SwiftUI 承载控制器
final class NavigationHostingController<Content: View>: UIHostingController<Content>, BottomBarVisibilityProviding {

    /// 当前页面从导航栈返回后的回调
    var onPop: (() -> Void)?

    /// 当前页面显示时是否隐藏底部 TabBar
    var prefersBottomBarHidden = false

    private var didRunOnPop = false

    override func viewDidDisappear(_ animated: Bool) {
        super.viewDidDisappear(animated)

        if isMovingFromParent || navigationController?.isBeingDismissed == true {
            runOnPop()
        }
    }

    private func runOnPop() {
        guard !didRunOnPop else { return }
        didRunOnPop = true
        onPop?()
    }
}
