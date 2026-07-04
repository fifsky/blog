import Foundation

extension Error {

    /// 是否为请求被取消的错误
    ///
    /// 覆盖两类来源：
    /// - `CancellationError`：Swift 协作式取消（Task 取消）
    /// - `URLError.cancelled`：URLSession 底层取消（用户离开页面、切 tab、刷新中断）
    ///
    /// ViewModel 的 `catch` 块应据此跳过 alert 提示，避免用户主动取消被误报为错误。
    var isCancellation: Bool {
        if self is CancellationError {
            return true
        }
        let urlError = self as? URLError
        return urlError?.code == .cancelled
    }
}
