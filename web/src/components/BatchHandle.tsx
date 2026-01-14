interface BatchHandleProps {
  selectedCount?: number;
  totalCount?: number;
  onSelectAll?: () => void;
  onInverseSelected?: () => void;
}

export function BatchHandle({ selectedCount = 0, totalCount = 0, onSelectAll, onInverseSelected }: BatchHandleProps) {
  return (
    <div className="flex items-center gap-2">
      <a
        href="#"
        className="all-selected"
        onClick={(e) => {
          e.preventDefault();
          onSelectAll?.();
        }}
      >
        全选
      </a>
      <span className="px-1.5 text-[#ccc]">|</span>
      <a
        href="#"
        className="inverse-selected"
        onClick={(e) => {
          e.preventDefault();
          onInverseSelected?.();
        }}
      >
        反选
      </a>
      &nbsp;&nbsp;
      {selectedCount > 0 && (
        <span className="text-[13px] text-[#666]">
          已选中 <span className="font-bold text-[#333]">{selectedCount}</span> / {totalCount}
        </span>
      )}
      &nbsp;&nbsp;
      <select name="batch_operation">
        <option value="" defaultChecked>
          批量操作
        </option>
        <option value="1">置顶</option>
        <option value="2">删除</option>
      </select>
    </div>
  );
}
