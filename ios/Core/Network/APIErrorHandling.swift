import Foundation

/// 可展示 API 错误的视图模型
protocol APIErrorPresentable: AnyObject {
    /// 错误信息
    var errorMessage: String? { get set }

    /// 是否显示错误弹窗
    var showError: Bool { get set }
}

extension APIErrorPresentable {

    /// 统一处理 API 错误，客户端取消请求静默忽略
    func handleAPIError(_ error: Error, prefix: String? = nil) {
        let apiError = APIError.normalized(error)
        guard !apiError.isCancelled else {
            return
        }

        errorMessage = apiError.displayMessage(prefix: prefix)
        showError = true
    }
}
