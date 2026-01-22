旅行相册功能开发
我新创建了两张表: regions和photos，用来存储中国省份城市区域和上传的照片，其中照片和区域通过region_id关联，现在我需要开发如下功能
1、在后台新增一个标签页，用于管理相册，参考AdminRemind.tsx提醒页面，列表页展示已经上传的相片，展示字段包含（标题，描述，照片查看，省份，城市，创建时间）支持编辑和删除（删除使用物理删除），编辑不支持重新上传，仅能编辑照片基本信息
2、管理页面新增相片按钮弹框表单，填写（标题，描述，选择省市两级：省份城市下拉框需要联动使用parent_id和level可以获取下一级，上传按钮上传相片）
2.1、上传采用阿里云OSS前端上传,你需要使用github.com/aliyun/aliyun-oss-go-sdk/oss sdk实现一个接口获取临时访问凭证给到前端上传组件，上传bucket为配置config.OSS.Bucket
2.2、前端上传保存路径为blog/photos/${YYYY}/${MM}/${DD}/${FILE}前端上传后将完整相片地址发送给后端接口保存，后端接口保存完整的图片路径，拼接CDN地址前缀保存https://static.fifsky.com/，src使用原图地址，thumbnail在地址后拼接oss图片处理后缀：!photothumb
3、提供一个开放接口给到TraveMap.tsx页面所需要的参数，前端根据返回的数据转换为所需的map_regions:对应provinces的name,map_footprints：对应citys
{"provinces":[{"region_id":"","name":"","longitude":"","latitude":""}]，"citys":[{"region_id":"","name":"","longitude":"","latitude":""}]}
当点击足迹上的城市，则请求获取照片接口，使用city的region_id
{"photos":[{"title":"","description":"","src":"","thumbnail"}]}
4、请合理规划代码，合理命名，实现接口和页面的所有功能，并完善后端单元测试