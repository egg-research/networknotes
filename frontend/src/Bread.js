import React from 'react';
import { Breadcrumb } from 'antd';
import { useParams, Link } from 'react-router-dom';
import { HomeOutlined, FileOutlined } from '@ant-design/icons';

export default function Bread({ children }) {
  return (
    <Breadcrumb style={{ marginBottom: 8 }}>
      <Breadcrumb.Item>
        <Link to='/'>
          <HomeOutlined />
          <span> Home</span>
        </Link>
      </Breadcrumb.Item>
      {children && <Breadcrumb.Item>{children}</Breadcrumb.Item>}
    </Breadcrumb>
  );
}
