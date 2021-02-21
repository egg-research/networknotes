import React, { useState } from 'react';
import { Typography, Space, Input, Button } from 'antd';

export default function CreateDocumentForm() {
  const [title, setTitle] = useState('');
  const [submit, setSubmit] = useState(false);

  const { Title } = Typography;

  return (
    <Space direction='vertical'>
      <Typography>
        <Title level={4}>Create New Document</Title>
      </Typography>
      <Input
        placeholder='Document Title'
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <Button
        loading={submit}
        onClick={(e) => {
          setSubmit(true);
          console.log(title);
          setTitle('');
        }}
      >
        Create
      </Button>
    </Space>
  );
}
