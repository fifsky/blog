# 开发规范

## 项目规范

项目采用 React + TypeScript + TailwindCSS 实现开发，技术栈如下

- React 19
- TypeScript 5.9
- TailwindCSS 4
- React Hook Form：表单统一使用 react-hook-form 处理
- Zod 4：用于表单验证
- React Router 7：请不要使用旧版本的 react-router-dom，统一采用新版的 react-router
- Shadcn UI：如果你需要安装响应的组件，请执行命令`pnpm dlx shadcn@latest add button`, 更多组件请参考[Shadcn UI 文档](https://ui.shadcn.com/docs/components)
- Zustand：状态管理统一使用 zustand 处理
- dayjs：日期处理统一使用 dayjs 处理

## 注意事项

- 请不要在每次任务执行完都执行`pnpm build`，在开发模式我可以实时看到编译结果，不需要执行打包命令
- 执行`pnpm tsc --noEmit`检查前端代码是否有类型错误的时候请忽略`node_modules`目录