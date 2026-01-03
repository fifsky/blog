export function BatchHandle() {
  return (
    <div className="flex items-center gap-2">
      <a href="#" className="all-selected">全选</a><span className="line">|</span>
      <a href="#" className="inverse-selected">反选</a>&nbsp;&nbsp;
      <select name="batch_operation">
        <option value="" defaultChecked>批量操作</option>
        <option value="1">置顶</option>
        <option value="1">删除</option>
      </select>
    </div>
  )
}
