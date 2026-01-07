import React from 'react';

// 定义列配置类型
export interface Column<T> {
  title: React.ReactNode; // 列标题，支持字符串或React节点
  key: string; // 数据字段key，支持jsonpath格式
  render?: (value: any, record: T, index: number) => React.ReactNode; // 自定义渲染函数
  style?: React.CSSProperties; // 列样式
}

// 定义组件属性类型
interface MyTableProps<T> {
  data: T[]; // 表格数据
  columns: Column<T>[]; // 列配置
  className?: string; // 表格类名
}

// 解析jsonpath，获取对象中的值
const getValueByJsonPath = (obj: any, path: string): any => {
  if (!obj || !path) return undefined;
  
  // 分割路径，处理数组索引和对象属性
  const keys = path.split(/[.[\]]+/).filter(Boolean);
  
  let result = obj;
  for (const key of keys) {
    // 如果是数组索引，转换为数字
    const index = isNaN(Number(key)) ? key : Number(key);
    result = result[index];
    
    // 如果中间值为undefined，直接返回
    if (result === undefined) return undefined;
  }
  
  return result;
};

/**
 * 通用表格组件
 * @param data 表格数据
 * @param columns 列配置
 * @param className 表格类名
 * @returns 表格组件
 */
export const CTable = <T,>({ data, columns, className = '' }: MyTableProps<T>) => {
  return (
    <table className={`w-full text-[13px] ${className}`}>
      <tbody>
        <tr>
          {columns.map((column, index) => (
            <th key={index} style={column.style} className="border-b-2 border-[#d9e0ec] text-left pl-[10px]">
              {column.title}
            </th>
          ))}
        </tr>
        
        {data.length === 0 && (
          <tr>
            <td colSpan={columns.length} align="center">
              还没有数据！
            </td>
          </tr>
        )}
        
        {data.length > 0 && data.map((record, rowIndex) => (
          <tr key={rowIndex} className="hover:bg-[#eee]">
            {columns.map((column, colIndex) => {
              const value = getValueByJsonPath(record, column.key);
              
              return (
                <td key={colIndex} className="border-b border-dashed border-[#d9e0ec] border-t-0 py-[5px] px-[10px]">
                  {column.render ? column.render(value, record, rowIndex) : value}
                </td>
              );
            })}
          </tr>
        ))}
      </tbody>
    </table>
  );
};


