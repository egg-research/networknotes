import React, { useState, useContext } from 'react';

import { useMeasure, useEffectOnce } from 'react-use';
import { Radio, Divider } from 'antd';
import Graph from './Graph';
import Layout from './Layout';
import DocumentTable from './DocumentTable';
import CreateDocumentForm from './CreateDocumentForm';
import SettingsCard from './SettingsCard';
import Bread from './Bread';
import UserContext from './context';
import { getDocGraph, getKeywordGraph, getAllKeywords, getAllDocs } from './db';

import { genRandomTree } from './utils/Data';
import { processGraph, applyGraphFilter } from './utils/graph';
import './GraphPage.css';

function ViewSwitch({ className, setView }) {
  const onChange = (e) => {
    setView(e.target.value);
  };

  return (
    <Radio.Group
      className={className}
      onChange={onChange}
      defaultValue='graph'
      buttonStyle='solid'
    >
      <Radio.Button value='graph'>Graph</Radio.Button>
      <Radio.Button value='table'>Table</Radio.Button>
    </Radio.Group>
  );
}

function TableContainer({ tableData, className }) {
  return (
    <div className={className}>
      <DocumentTable data={tableData} className='table' />
    </div>
  );
}

export default function GraphPage() {
  const [view, setView] = useState('graph');
  const [documentFilter, setDocumentFilter] = useState(new Set());
  const [keywordFilter, setKeywordFilter] = useState(new Set());
  const [allKeywords, setAllKeywords] = useState([]);
  const [allDocuments, setAllDocuments] = useState([]);
  const [ref, dimensions] = useMeasure();
  const [rawDocGraph, setRawDocGraph] = useState({ nodes: [], links: [] });
  const [rawKeywordGraph, setRawKeywordGraph] = useState({
    nodes: [],
    links: [],
  });
  const userId = useContext(UserContext);

  const { height, width } = dimensions;

  const Sidebar = () => (
    <>
      <Bread />
      <Divider dashed />
      <CreateDocumentForm />
    </>
  );

  useEffectOnce(async () => {
    Promise.all([
      getDocGraph(userId),
      getKeywordGraph(userId),
      getAllKeywords(userId),
      getAllDocs(userId),
    ]).then((values) => {
      setRawDocGraph(values[0]);
      setRawKeywordGraph(values[1]);
      setAllKeywords(values[2]);
      setAllDocuments(values[3]);
    });
  });

  const data = processGraph(rawDocGraph);
  const graphData = applyGraphFilter(data, documentFilter, keywordFilter);

  return (
    <Layout Sidebar={<Sidebar />} contentRef={ref} contentPadding={false}>
      <div className='elevated'>
        <ViewSwitch className='view-switch' setView={setView} />
        {view === 'graph' && (
          <SettingsCard
            className='settings-card'
            documents={allDocuments.filter(
              (doc) => !documentFilter.has(doc.id)
            )}
            keywords={allKeywords.filter((kw) => !keywordFilter.has(kw.name))}
            documentFilter={documentFilter}
            keywordFilter={keywordFilter}
            addDocument={(newDoc) => {
              documentFilter.add(newDoc);
              setDocumentFilter(new Set(documentFilter));
            }}
            addKeyword={(newKW) => {
              keywordFilter.add(newKW);
              setKeywordFilter(new Set(keywordFilter));
            }}
            removeDocument={(oldDoc) => {
              documentFilter.delete(oldDoc);
              setDocumentFilter(new Set(documentFilter));
            }}
            removeKeyword={(oldKW) => {
              keywordFilter.delete(oldKW);
              setKeywordFilter(new Set(keywordFilter));
            }}
            // {/* setOpacity={} */}
          />
        )}
      </div>
      {view === 'graph' ? (
        <Graph data={graphData} height={height} width={width} />
      ) : (
        <TableContainer
          className='content-container'
          // tableData={documentData}
        />
      )}
    </Layout>
  );
}
