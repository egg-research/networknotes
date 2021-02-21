import React, { useState, useContext } from 'react';
import { Typography, Space, Input, Button } from 'antd';
import { useHistory } from 'react-router-dom';
import UserContext from './context';
import { makeNewDoc } from './db';

export default function CreateDocumentForm() {
  const [title, setTitle] = useState('');
  const [submit, setSubmit] = useState(false);
  const userId = useContext(UserContext);
  const history = useHistory();

  const { Title } = Typography;

  return (
    <Space direction='vertical'>
      <Typography>
        <Title level={4}>Create Document</Title>
      </Typography>
      <Input
        placeholder='Document Title'
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <Button
        loading={submit}
        onClick={async (e) => {
          setSubmit(true);
          const docId = await makeNewDoc(userId, title);
          history.push(`/document/${docId}`);
        }}
      >
        Create
      </Button>
    </Space>
  );
}
