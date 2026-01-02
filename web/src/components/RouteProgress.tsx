import { useEffect } from 'react';
import { useNavigation } from 'react-router-dom'; // 核心钩子：获取导航状态
// @ts-ignore
import NProgress from 'nprogress';
import 'nprogress/nprogress.css'; // 必须导入样式

export function RouteProgress() {
    // 获取导航状态：idle（空闲）、loading（导航中/加载中）
    const navigation = useNavigation();

    useEffect(() => {
        // 当导航状态为 loading 时，启动进度条
        if (navigation.state === 'loading') {
            NProgress.start();
        } else {
            // 当导航状态为 idle 时，结束进度条
            NProgress.done();
        }

        // 组件卸载时确保进度条结束
        return () => {
            NProgress.done();
        };
    }, [navigation.state]); // 依赖导航状态变化

    return null; // 仅执行逻辑，不渲染DOM
}
