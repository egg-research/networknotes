import React, { useState, useRef, useContext } from 'react';

import { useParams, Link } from 'react-router-dom';
import { Button, Space, Typography, Divider, Tag, Spin } from 'antd';
import { PlusOutlined, FileOutlined } from '@ant-design/icons';
import { useEffectOnce } from 'react-use';


import {
  Editor,
  createEditorState,
  BLOCK_BUTTONS,
  INLINE_BUTTONS,
} from 'medium-draft';

import { convertToRaw } from 'draft-js';
import SearchBar from './SearchBar';

import 'medium-draft/lib/index.css';

import './EditorPage.css';
import { HeavyText } from './utils/CustomText';
import Layout from './Layout';
import Bread from './Bread';

import UserContext from './context';
import { readDoc, updateDoc, getAllKeywords, updateDocKeyword, getRelatedKwds} from './db';

const { Text, Paragraph, Title } = Typography;
const { CheckableTag } = Tag;

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
  addSuggestedKeywords,
  createKeyword,
  loading,
  loadingSuggested,
  addKeyword,
  removeKeyword,
}) {
  const [showSearch, setShowSearch] = useState(false);

  const { title, keywords, suggested } = document;
  return (
    <>
      <Bread>
        <FileOutlined style={{ fontSize: "24px" }}/>
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
        <Spin spinning={loading}>
          {showSearch && (
            <div style={{ marginBottom: 8 }}>
              <SearchBar
                keywords={allKeywords}
                newKeyword
                createKeyword={createKeyword}
                selectKeyword={(keyword) => addKeyword(keyword.name)}
              />
            </div>
          )}
          <Space direction='vertical'>
            <div>
              {keywords.map((keyword) => (
                <Tag
                  closable
                  key={keyword}
                  onClose={() => removeKeyword(keyword)}
                >
                  {keyword}
                </Tag>
              ))}
            </div>
          </Space>
        </Spin>
      </div>
      <Divider dashed />
      <div>
        <Typography>
          <Title level={4}>Suggested Keywords</Title>
          <Spin spinning={loadingSuggested}>
            <Space direction='vertical'>
              <div>
                {suggested.map((keyword) => (
                  <CheckableTag
                    key={keyword}
                    onChange={() => addSuggestedKeywords(keyword)}
                  >
                    {keyword}
                  </CheckableTag>
                ))}
              </div>
            </Space>
          </Spin>
        </Typography>
      </div>
    </>
  );
}

let tick = 0;

export default function EditorPage() {
  const [data, setData] = useState(createEditorState());
  const [keywords, setKeywords] = useState([]);
  const [document, setDocument] = useState(null);
  const [suggested, setSuggested] = useState([]);
  const [allKeywords, setAllKeywords] = useState([]);
  const [loading, setLoading] = useState(false);
  const [loadingSuggested, setLoadingSuggested] = useState(false);
  const ref = useRef(null);
  const { id } = useParams();
  const docId = parseInt(id, 10);
  const userId = useContext(UserContext);

  useEffectOnce(async () => {
    const res = await readDoc(userId, docId);
    const related_res = await getRelatedKwds(userId, docId);
    setDocument(res.document);
    setKeywords(res.keywords);
    let suggested_kws = related_res.generated.concat(related_res.kws);
    setSuggested(suggested_kws);
    const initData =
      res.document.rawText !== ''
        ? createEditorState(res.document.rawText)
        : createEditorState();
    setData(initData);
  });

  useEffectOnce(async () => {
    const res = await getAllKeywords(userId);
    setAllKeywords(res);
  });

  const onChange = (editorState) => {
    tick += 1;
    if (tick > 5) {
      const rawText = convertToRaw(editorState.getCurrentContent());
      const text = editorState.getCurrentContent().getPlainText();
      updateDoc(userId, docId, document.title, text, rawText);
      tick = 0;
    }

    setData(editorState);
  };

  const addKeyword = (keyword) => {
    for (const kw of keywords) {
      if (kw === keyword) {
        return;
      }
    }
    const newKeywords = Array.from(keywords.concat([keyword]));
    updateDocKeyword(userId, docId, newKeywords);

    setKeywords(newKeywords);
  };

  const removeKeyword = (kw) => {
    const newKeywords = Array.from(keywords.filter((y) => y !== kw));
    updateDocKeyword(userId, docId, newKeywords);
    setKeywords(newKeywords);
  };


  const createKeyword = async (kw) => {
    setLoading(true);
    const newKeywords = Array.from(keywords.concat([kw]));
    await updateDocKeyword(userId, docId, newKeywords);
    setKeywords(newKeywords);
    const res = await getAllKeywords(userId);
    setAllKeywords(res);
    setLoading(false);
  };

  const addSuggested = async (kw) => {
    setLoadingSuggested(true);
    await createKeyword(kw);
    const newKeywords = Array.from(suggested.filter((y) => y !== kw));
    setSuggested(newKeywords);
    setLoadingSuggested(false);
  };

  const Side = (
    <SideDisplay
      document={{ title: document ? document.title : '', keywords, suggested}}
      allKeywords={allKeywords}
      removeKeyword={removeKeyword}
      addKeyword={addKeyword}
      createKeyword={createKeyword}
      addSuggestedKeywords={addSuggested}
      loading={loading}
      loadingSuggested={loadingSuggested}
    />
  );

  return (
    <Layout Sidebar={Side}>
      <EditorBlock ref={ref} data={data} onChange={onChange} />
    </Layout>
  );
}
