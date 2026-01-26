# AI任务

以下记录后续需要AI完成的任务

## 微信小程序

采用Taro + TailwindCSS 4 开发微信小程序，实现以下功能：
1、小程序目录是miniapp，已经使用Taro初始化项目，你需要在miniapp目录下开发小程序
2、小程序appid和秘钥配置在config.MiniAPP，小程序appid为wxf55a6e31920e7294
3、请将你需要调用的小程序登录接口SDK封装到pkg/miniapp下面，登录使用users表中的openid进行匹配，登录后返回和web版本一样的JWT token
4、前端通过 Taro.request 调用服务端接口，接口的封装也参考web版本的openapi.ts定义和createApi函数，使用和web版本一样的JWT token进行认证。
5、登录后可以在小程序中发表心情、上传相册、提醒管理三大功能
6、支持上传相册，上传相册希望能够支持自动识别打卡地理位置或者主动选择省份城市
7、要求风格主色调和web版本保持一致，使用TailwindCSS 4的颜色类，使用lucide-react图标库，选择合适的图表作为tabbar的图标
8、本地调用http://127.0.0.1:8080，线上调用https://api.fifsky.com，接口的定义和web版本保持一致参见web/vite.config.ts

## 文章标签功能（已完成）

1、编写文章的时候，新增文章标签JSON字段，使用custom/input-tag.tsx组件
2、在组件后面增加一个按钮，点击按钮可以调用AI接口自动生成标签（将文章内容作为参数传递给AI接口），你需要实现AI接口，并使用openai sdk，项目中已经有类似的案例，这个接口不需要流式输出，直接输出既可
3、用户可以手动添加、删除标签
4、在文章的底部使用徽章组件展示标签，点击标签可以过滤文章列表
5、你需要实现所有功能，请不要中途退出，不要留任何优化给我，帮我直接优化
