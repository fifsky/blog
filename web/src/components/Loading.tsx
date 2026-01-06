// 页面加载中组件
export function Loading() {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-white/50 z-50">
      <div className="flex flex-col items-center justify-center">
        {/* 旋转加载动画 */}
        <div className="w-6 h-6 border-3 border-gray-200 border-t-blue-500 rounded-full animate-spin"></div>
        {/* 加载文字 */}
        <div className="mt-4 text-gray-600 text-sm">页面加载中</div>
      </div>
    </div>
  );
}

export default Loading;
