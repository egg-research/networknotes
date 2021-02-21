import React from 'react';
import { Card, Divider, Typography, Tag, Slider, Space } from 'antd';
import Search from 'antd/lib/transfer/search';
import { HeavyText } from './utils/CustomText';
import SearchBar from './SearchBar';

import './SettingsCard.css';

export default function SettingsCard({
  className,
  documents,
  keywords,
  documentFilter,
  keywordFilter,
  addDocument,
  addKeyword,
  removeDocument,
  removeKeyword,
  setOpacity,
}) {
  const { Text } = Typography;

  const content = (
    <div style={{ padding: 16 }}>
      <SearchBar
        documents={documents}
        keywords={keywords}
        selectKeyword={addKeyword}
        selectDocument={addDocument}
        placeholder='Add Filter'
      />
      <div className='row'>
        <HeavyText>Documents</HeavyText>
      </div>
      <div className='row'>
        {documentFilter.size ? (
          Array.from(documentFilter).map((x) => (
            <Tag closable key={x} onClose={() => removeDocument(x)}>
              {x}
            </Tag>
          ))
        ) : (
          <Text type='secondary'>Showing all documents</Text>
        )}
      </div>
      <div className='row'>
        <HeavyText>Keywords</HeavyText>
      </div>
      <div className='row'>
        {keywordFilter.size ? (
          Array.from(keywordFilter).map((x) => (
            <Tag closable key={x} onClose={() => removeKeyword(x)}>
              {x}
            </Tag>
          ))
        ) : (
          <Text type='secondary'>Showing all keywords</Text>
        )}
      </div>
      <div className='row'>
        <HeavyText>Weight</HeavyText>
      </div>
      <Slider defaultValue={20} />
    </div>
  );

  return <Card cover={content} className={className} />;
}
