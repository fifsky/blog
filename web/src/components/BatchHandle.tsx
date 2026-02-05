interface BatchHandleProps {
  selectedCount?: number;
  totalCount?: number;
  onSelectAll?: () => void;
  onInverseSelected?: () => void;
  onBatchOperation?: (operation: string) => void;
  disabled?: boolean;
}

export function BatchHandle({
  selectedCount = 0,
  totalCount = 0,
  onSelectAll,
  onInverseSelected,
  onBatchOperation,
  disabled = false,
}: BatchHandleProps) {
  const handleOperationChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    if (value && onBatchOperation) {
      onBatchOperation(value);
      e.target.value = ""; // 重置选择
    }
  };

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
      {selectedCount > 0 && (
        <span className="text-[13px] text-[#666]">
          已选中 <span className="font-bold text-[#333]">{selectedCount}</span> / {totalCount}
        </span>
      )}
      <select
        name="batch_operation"
        onChange={handleOperationChange}
        disabled={disabled || selectedCount === 0}
        defaultValue=""
      >
        <option value="">批量操作</option>
        <option value="2">删除</option>
      </select>
    </div>
  );
}
