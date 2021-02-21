import React, { useState, useRef, useContext } from 'react';

import { useParams, Link } from 'react-router-dom';
import { Button, Space, Typography, Divider, Tag } from 'antd';
import { PlusOutlined, FileOutlined } from '@ant-design/icons';
import { useEffectOnce } from 'react-use';

import {
  Editor,
  createEditorState,
  BLOCK_BUTTONS,
  INLINE_BUTTONS,
} from 'medium-draft';
import SearchBar from './SearchBar';

import 'medium-draft/lib/index.css';

import './EditorPage.css';
import { HeavyText } from './utils/CustomText';
import Layout from './Layout';
import Bread from './Bread';

import UserContext from './context';
import { readDoc } from './db';

const { Text, Paragraph, Title } = Typography;

function EditorBlock({ ref, data, onChange }) {
  return (
    <Editor
      ref={ref}
      editorState={data}
      onChange={onChange}
      blockButtons={BLOCK_BUTTONS}
      inlineButtons={INLINE_BUTTONS}
    />
  );
}

function SideDisplay({
  document,
  allKeywords,
  suggestedKeywords,
  addKeyword,
  removeKeyword,
}) {
  const [showSearch, setShowSearch] = useState(false);

  const { title, keywords } = document;
  return (
    <>
      <Bread>
        <FileOutlined />
      </Bread>
      <Divider dashed />
      <Typography>
        <Title level={3}>{title}</Title>
      </Typography>
      <Typography>
        <Paragraph>
          <HeavyText>Last Edited:</HeavyText> <Text> Apr 25, 2021</Text>
        </Paragraph>
      </Typography>
      <Divider dashed />
      <div>
        <Typography>
          <div style={{ display: 'flex', justifyContent: 'flex-start' }}>
            <Title level={4} style={{ marginRight: 8 }}>
              Keywords
            </Title>
            <Button
              icon={<PlusOutlined />}
              size='small'
              shape='circle'
              type='dashed'
              onClick={() => setShowSearch(!showSearch)}
            />
          </div>
        </Typography>
        {showSearch && (
          <div style={{ marginBottom: 8 }}>
            <SearchBar
              keywords={allKeywords}
              selectKeyword={(keyword) => addKeyword(keyword)}
            />
          </div>
        )}
        <Space direction='vertical'>
          <div>
            {keywords.map((keyword) => (
              <Tag
                closable
                key={keyword.id}
                onClose={() => removeKeyword(keyword.id)}
              >
                {keyword.name}
              </Tag>
            ))}
          </div>
        </Space>
      </div>
      <Divider dashed />
      <div>
        <Typography>
          <Title level={4}>Suggested Keywords</Title>
        </Typography>
      </div>
    </>
  );
}

const KEYWORDS = [
  { id: 0, name: 'CNN' },
  { id: 1, name: 'RNN' },
];

export default function EditorPage() {
  const [data, setData] = useState(createEditorState());
  const [keywords, setKeywords] = useState([]);
  const [document, setDocument] = useState(null);
  const ref = useRef(null);
  const { id } = useParams();
  const docId = parseInt(id, 10);
  const userId = useContext(UserContext);

  useEffectOnce(async () => {
    const res = await readDoc(userId, docId);
    setDocument(res.document);
    setKeywords(res.keywords);
    const initData = res.document.rawText
      ? createEditorState(res.document.rawText)
      : createEditorState();
    setData(initData);
  });

  const onChange = (editorState) => {
    setData(editorState);
  };

  const allKeywords = [
    { id: 0, name: 'CNN' },
    { id: 1, name: 'RNN' },
    { id: 2, name: 'Trie' },
    { id: 3, name: 'LinkedList' },
  ];

  const d = {
    title: 'Introduction to Machine Learning',
    keywords,
  };

  const addKeyword = (keyword) => {
    for (const kw of keywords) {
      if (kw.id === keyword.id) {
        return;
      }
    }
    const newKeywords = Array.from(keywords.concat([keyword]));
    setKeywords(newKeywords);
  };

  const removeKeyword = (kwId) => {
    const newKeywords = Array.from(keywords.filter((y) => y.id !== kwId));
    setKeywords(newKeywords);
  };

  const Side = (
    <SideDisplay
      document={d}
      allKeywords={allKeywords}
      removeKeyword={removeKeyword}
      addKeyword={addKeyword}
    />
  );

  return (
    <Layout Sidebar={Side}>
      <EditorBlock ref={ref} data={data} onChange={onChange} />
    </Layout>
  );
}
