import SwiftUI

/// 登录视图
/// 分两步登录：第一步输入用户名密码，第二步（账号开启 TOTP 时）输入 6 位验证码
struct LoginView: View {

    @State private var viewModel = LoginViewModel()

    /// 登录成功回调
    var onLoginSuccess: (LoginViewModel) -> Void

    // MARK: - 焦点状态

    /// 用户名输入框焦点
    @FocusState private var userNameFocused: Bool

    /// 密码输入框焦点
    @FocusState private var passwordFocused: Bool

    /// OTP bridge 输入框焦点
    @FocusState private var otpFocused: Bool

    var body: some View {
        NavigationStack {
            ZStack {
                // 用 switch + transition 实现两步页面切换
                switch viewModel.step {
                case .credentials:
                    credentialsStep
                        .transition(.asymmetric(
                            insertion: .move(edge: .trailing).combined(with: .opacity),
                            removal: .move(edge: .leading).combined(with: .opacity)
                        ))
                case .totp:
                    totpStep
                        .transition(.asymmetric(
                            insertion: .move(edge: .trailing).combined(with: .opacity),
                            removal: .move(edge: .leading).combined(with: .opacity)
                        ))
                }
            }
            .animation(.easeInOut(duration: 0.25), value: viewModel.step)
            .alert("登录失败", isPresented: $viewModel.showError) {
                Button("确定", role: .cancel) {}
            } message: {
                Text(viewModel.errorMessage ?? "未知错误")
            }
            .onChange(of: viewModel.isLoginSuccess) {
                if viewModel.isLoginSuccess {
                    onLoginSuccess(viewModel)
                }
            }
            .onChange(of: viewModel.step) { _, newStep in
                // 切换到对应步骤时自动聚焦
                switch newStep {
                case .credentials:
                    userNameFocused = true
                case .totp:
                    otpFocused = true
                }
            }
        }
    }

    // MARK: - 第一步：用户名密码

    /// 第一步登录页面（一体化分组输入框）
    private var credentialsStep: some View {
        VStack(spacing: 24) {
            Spacer()

            // 标题
            VStack(spacing: 8) {
                Text("博客")
                    .font(.largeTitle)
                    .bold()
                Text("登录以管理你的内容")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }

            // 一体化分组输入框：用户名 + 密码，中间用分割线分隔
            VStack(spacing: 0) {
                // 用户名
                HStack(spacing: 12) {
                    Image(systemName: "person")
                        .font(.body)
                        .foregroundStyle(.secondary)
                        .frame(width: 20)

                    TextField("用户名", text: $viewModel.userName)
                        .textContentType(.username)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .submitLabel(.next)
                        .focused($userNameFocused)
                        .onSubmit {
                            passwordFocused = true
                        }
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 14)

                // 中间分割线
                Divider()

                // 密码
                HStack(spacing: 12) {
                    Image(systemName: "lock")
                        .font(.body)
                        .foregroundStyle(.secondary)
                        .frame(width: 20)

                    SecureField("密码", text: $viewModel.password)
                        .textContentType(.password)
                        .submitLabel(.go)
                        .focused($passwordFocused)
                        .onSubmit {
                            submitCredentials()
                        }
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 14)
            }
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .strokeBorder(Color.secondary.opacity(0.3), lineWidth: 1)
            )
            .padding(.horizontal, 24)

            // 登录按钮
            Button {
                userNameFocused = false
                passwordFocused = false
                submitCredentials()
            } label: {
                HStack(spacing: 8) {
                    if viewModel.isLoading {
                        ProgressView()
                            .tint(.white)
                    }
                    Text("登录")
                        .font(.headline)
                }
                .frame(maxWidth: .infinity, minHeight: 44)
            }
            .buttonStyle(.borderedProminent)
            .disabled(viewModel.isCredentialsButtonDisabled || viewModel.isLoading)
            .padding(.horizontal, 24)

            Spacer()
            Spacer()
        }
        .navigationTitle("登录")
        .onAppear {
            // 首次进入自动聚焦用户名
            userNameFocused = true
        }
    }

    // MARK: - 第二步：TOTP 验证码

    /// 第二步验证码页面（6 个独立展示框 + 隐藏 bridge TextField）
    private var totpStep: some View {
        VStack(spacing: 24) {
            Spacer()

            // 图标 + 标题
            VStack(spacing: 12) {
                Image(systemName: "lock.shield")
                    .font(.system(size: 56))
                    .foregroundStyle(.blue)

                VStack(spacing: 8) {
                    Text("两步验证")
                        .font(.title2)
                        .bold()
                    Text("请输入身份验证器中的 6 位验证码")
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                        .multilineTextAlignment(.center)
                }
            }

            // 6 位验证码展示框（覆盖在隐藏的 bridge TextField 上）
            ZStack {
                // 隐藏的 bridge TextField：负责接收输入、键盘和自动填充
                TextField("", text: $viewModel.totpCode)
                    .keyboardType(.numberPad)
                    .textContentType(.oneTimeCode)
                    .focused($otpFocused)
                    .opacity(0.001)
                    .frame(width: 1, height: 1)
                    .onChange(of: viewModel.totpCode) { _, newValue in
                        // 仅保留数字，最多 6 位
                        let filtered = newValue.filter { $0.isNumber }
                        viewModel.totpCode = String(filtered.prefix(6))
                        // 输入满 6 位自动提交
                        if viewModel.totpCode.count == 6 {
                            submitTotp()
                        }
                    }

                // 6 个展示框
                HStack(spacing: 10) {
                    ForEach(0..<6, id: \.self) { index in
                        otpDigitBox(at: index)
                    }
                }
                .allowsHitTesting(false)
            }
            .padding(.horizontal, 24)
            .contentShape(Rectangle())
            .onTapGesture {
                // 点击展示框区域激活键盘
                otpFocused = true
            }

            // 提示输入位数
            Text("\(viewModel.totpCode.count) / 6")
                .font(.caption)
                .foregroundStyle(.secondary)

            // 验证按钮
            Button {
                otpFocused = false
                submitTotp()
            } label: {
                HStack(spacing: 8) {
                    if viewModel.isLoading {
                        ProgressView()
                            .tint(.white)
                    }
                    Text("验证")
                        .font(.headline)
                }
                .frame(maxWidth: .infinity, minHeight: 44)
            }
            .buttonStyle(.borderedProminent)
            .disabled(viewModel.isTotpButtonDisabled || viewModel.isLoading)
            .padding(.horizontal, 24)

            // 返回第一步
            Button {
                viewModel.backToCredentials()
            } label: {
                Text("返回重新输入账号")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }

            Spacer()
            Spacer()
        }
        .navigationTitle("两步验证")
        .onAppear {
            otpFocused = true
        }
    }

    // MARK: - 子视图

    /// 单个 OTP 数字展示框
    /// - Parameter index: 第几位（0...5）
    private func otpDigitBox(at index: Int) -> some View {
        let digits = Array(viewModel.totpCode)
        // 当前位是否已填入
        let isFilled = index < digits.count
        // 当前位数字
        let digit = isFilled ? String(digits[index]) : ""
        // 是否为下一个待输入位（高亮）
        let isCurrent = index == digits.count

        return Text(digit)
            .font(.system(size: 28, weight: .semibold, design: .monospaced))
            .frame(width: 40, height: 52)
            .background(
                RoundedRectangle(cornerRadius: 10)
                    .fill(Color(.systemGray6))
            )
            .overlay(
                RoundedRectangle(cornerRadius: 10)
                    .strokeBorder(
                        isCurrent ? Color.accentColor : Color.secondary.opacity(0.3),
                        lineWidth: isCurrent ? 2 : 1
                    )
            )
    }

    // MARK: - 辅助方法

    /// 提交第一步
    private func submitCredentials() {
        guard !viewModel.isCredentialsButtonDisabled else { return }
        Task {
            await viewModel.submitCredentials()
        }
    }

    /// 提交第二步验证码
    private func submitTotp() {
        guard !viewModel.isTotpButtonDisabled else { return }
        Task {
            await viewModel.submitTotp()
        }
    }
}
