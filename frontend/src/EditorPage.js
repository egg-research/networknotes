import React, { useState, useRef } from 'react';

import { Row, Col, Typography, Divider, Tag } from 'antd';

import {
  Editor,
  createEditorState,
  BLOCK_BUTTONS,
  INLINE_BUTTONS,
} from 'medium-draft';
import 'medium-draft/lib/index.css';

import './EditorPage.css';
import { HeavyText } from './utils/CustomText';

const { Text, Paragraph, Title } = Typography;

// const { useBreakpoint } = Grid;

export default function EditorPage() {
  const [data, setData] = useState(createEditorState());
  const ref = useRef(null);
  // const screens = useBreakpoint();

  const onChange = (editorState) => {
    // console.log(editorState.getCurrentContent().getPlainText());
    setData(editorState);
  };

  return (
    <Row gutter={64} className='root'>
      <Col xs={8} md={6} className='sidebar'>
        <Typography className='content first-content'>
          <Title level={3}>Introduction to Machine Learning</Title>
        </Typography>
        <Typography className='content'>
          <Paragraph>
            <HeavyText>Last Edited:</HeavyText> <Text> Apr 25, 2021</Text>
          </Paragraph>
        </Typography>
        <Divider dashed />
        <div className='content'>
          <Typography>
            <Title level={4}>Keywords</Title>
          </Typography>
          <Tag>CNN</Tag>
          <Tag>RNN</Tag>
        </div>
      </Col>
      <Col xs={16} md={18} className='first-content'>
        <Editor
          ref={ref}
          editorState={data}
          onChange={onChange}
          blockButtons={BLOCK_BUTTONS}
          inlineButtons={INLINE_BUTTONS}
        />
      </Col>
    </Row>
  );
}
