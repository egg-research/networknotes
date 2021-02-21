import React from 'react';

import { Table, Tag } from 'antd';
import { useHistory } from 'react-router-dom';

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
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
  const history = useHistory();
  return (
    <Table
      columns={columns}
      dataSource={data}
      className={className}
      onRow={(r) => ({ onClick: () => history.push(`/document/${r.id}`) })}
    />
  );
}
