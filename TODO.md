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

内容没有提取到，以下是提取规则
带有MOD参数的页面
昵称从 <img src="/web/20100716211623im_/http://www.windiness.com/guestbook/images/oicq.gif" alt="婷婷 的 QQ 号码：353534350" hspace="2" border="0"> 中提取,alt属性里面" 的"之前就是昵称
日期从 <font style="font-size:10px ; color:#000000;">Time: 2010-06-13 20:13:45</font> 中提取,Time: 后面就是日期
正文从 <td width="421" align="left" valign="top">随风而逝在那一年！</td> 中提取

不带MOD参数的页面
昵称从 <span class="name">莫一哲</span> 中提取
日期从<span class="input_time">2007-05-11 21:31:08</span>中提取
正文从 <div class="content" style="height:100px">
呵呵，突然想到，一只爬在窗玻璃上的蜗牛，从早晨的阳光爬到晚上的阳光，它和玻璃都有着晶莹剔透的美。呵，慢慢爬吧小蜗牛！</div> 中提取


带MOD参数的页面第一套留言提取之后
https://web.archive.org/web/20100716211639/http://www.windiness.com:80/guestbook/index.php?MOD=main&P=16
昵称：五月的雪
时间：2010-06-04 22:10:29
内容：真的 身处在这个网络时代 &nbsp; 每次面对屏幕总有一种无处藏身的感觉 &nbsp; 直到我来到这里 &nbsp;终于找到了一片净土 &nbsp;花火流年 &nbsp; 愿你越来越好

不带MOD参数的页面第一套留言提取之后
https://web.archive.org/web/20071012080740/http://windiness.com:80/guestbook/index.php?page=1
昵称：皇家澜澜
时间：2007-10-12 09:48:38
内容：<b>很佩服你，希望认识你．．．</b><br>　　　看了你的站感觉好美，很佩服有如此深的感受，我想认识你．．．

