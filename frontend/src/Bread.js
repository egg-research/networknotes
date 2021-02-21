import React from 'react';
import { Breadcrumb } from 'antd';
import { useParams, Link } from 'react-router-dom';
import { HomeOutlined, FileOutlined } from '@ant-design/icons';

export default function Bread({ children }) {
  return (
    <Breadcrumb style={{ fontSize:"24px", marginBottom: 8 }}>
      <Breadcrumb.Item>
        <Link to='/'>
          <HomeOutlined style={{fontSize:"24px", paddingLeft:"10px", paddingTop:"10px"}} />
          <span style={{fontSize:"24px"}}> Home</span>
        </Link>
      </Breadcrumb.Item>
      {children && <Breadcrumb.Item style={{ fontSize: "24px" }} >{children}</Breadcrumb.Item>}
    </Breadcrumb>
  );
}
