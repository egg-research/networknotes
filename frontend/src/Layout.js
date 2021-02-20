import React from 'react';

import { Row, Col } from 'antd';

import './Layout.css';

export default function Layout({
  Sidebar,
  contentRef,
  contentPadding = true,
  children,
}) {
  const contentClassName = contentPadding
    ? 'content content-padding'
    : 'content';

  return (
    <Row className='root'>
      <Col className='sidebar'>{Sidebar}</Col>
      <Col className={contentClassName} ref={contentRef}>
        {children}
      </Col>
    </Row>
  );
}
