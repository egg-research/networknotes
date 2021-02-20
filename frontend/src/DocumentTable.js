import React from 'react';

import { Table, Tag } from 'antd';

const columns = [
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: 'Last Accessed',
    dataIndex: 'lastAccessed',
    key: 'lastAccessed',
  },
  {
    title: 'Keywords',
    dataIndex: 'keywords',
    key: 'keywords',
    render: (keywords) => (
      <>
        {keywords.map((keyword) => (
          <Tag key={keyword}>{keyword}</Tag>
        ))}
      </>
    ),
  },
];

export default function DocumentTable({ data, className }) {
  return <Table columns={columns} dataSource={data} className={className} />;
}
