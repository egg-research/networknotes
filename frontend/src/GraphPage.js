import React, { useState, useContext } from 'react';

import { useMeasure, useEffectOnce } from 'react-use';
import { Spin, Radio, Divider, Input } from 'antd';
import Graph from './Graph';
import Layout from './Layout';
import DocumentTable from './DocumentTable';
import CreateDocumentForm from './CreateDocumentForm';
import SettingsCard from './SettingsCard';
import Bread from './Bread';
import UserContext from './context';
import {
  makeNewDoc,
  getDocGraph,
  getKeywordGraph,
  getAllKeywords,
  getAllDocs,
} from './db';
import { getText } from './ml';

import InputFile from './InputFile';

import { genRandomTree } from './utils/Data';
import { processGraph, applyGraphFilter } from './utils/graph';
import './GraphPage.css';

const { Search } = Input;

const documentData = [
  {
    id: 0,
    key: 0,
    name: 'Intro to Machine Learning',
    lastAccessed: '8:35pm Apr 20, 2020',
    keywords: ['CNN', 'RNN'],
  },
  {
    id: 1,
    key: 1,
    name: 'Computer Vision',
    lastAccessed: '8:35pm Apr 20, 2020',
    keywords: ['CNN', 'RNN'],
  },
  {
    id: 2,
    key: 2,
    name: 'Advanced Data Structures',
    lastAccessed: '8:35pm Apr 20, 2020',
    keywords: ['Trie', 'Skip List'],
  },
];

const keywords = [
  { id: 0, name: 'CNN' },
  { id: 1, name: 'RNN' },
  { id: 2, name: 'Trie' },
  { id: 3, name: 'Skip List' },
];

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
  console.log('TABLE DATA', tableData);

  const docMap = new Map();
  const docKwMap = new Map();
  const links = tableData.links;
  const nodes = tableData.nodes;
  links.forEach((link) => {
    const source = link.source.id != null ? link.source.id : link.source;
    const target = link.target.id != null ? link.target.id : link.target;
    if (!docKwMap.has(source)) {
      docKwMap[source] = new Set();
    }
    if (!docKwMap.has(target)) {
      docKwMap[target] = new Set();
    }

    docKwMap[source].add(link.name);
    docKwMap[target].add(link.name);
  });

  console.log(docKwMap);

  nodes.forEach((node) => {
    docMap[node.id] = node.name;
  });

  const result = [];
  nodes.forEach((node) => {
    result.push({
      name: docMap[node.id],
      id: node.id,
      keywords: docKwMap[node.id] == null ? [] : Array.from(docKwMap[node.id]),
    });
  });
  console.log('result', result);
  return (
    <div className={className}>
      <DocumentTable data={result} className='table' />
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
  const [loading, setLoading] = useState(true);
  const userId = useContext(UserContext);

  const uploadFileHandler = (e) => {
    console.log(e.target.files[0]);
  };

  async function onUpload(e) {
    console.log(e);
    const docText = await getText(e);
    console.log(docText.text);
  }

  const { height, width } = dimensions;

  const Sidebar = () => (
    <>
      <Bread />
      <Divider dashed />
      <CreateDocumentForm />
      <InputFile uploadFileHandler={uploadFileHandler}>
        Select a PDF file
      </InputFile>
      <Search
        className='pdfURL'
        placeholder='Input URL'
        enterButton='Upload'
        size='small'
        onSearch={onUpload}
      />
    </>
  );

  useEffectOnce(async () => {
    Promise.all([
      getDocGraph(userId),
      getKeywordGraph(userId),
      getAllKeywords(userId),
      getAllDocs(userId),
    ]).then((values) => {
      // console.log('value-----------------');
      // console.log(values[0]);
      setRawDocGraph(values[0]);
      setRawKeywordGraph(values[1]);
      setAllKeywords(values[2]);
      setAllDocuments(values[3]);
      // console.log('hi');
      // setLoading(false);
    });
  });

  const graphCopy = { ...rawDocGraph };
  const data = processGraph(graphCopy);
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
        <Graph
          data={graphData}
          height={height}
          width={width}
          load={
            keywordFilter.size +
              documentFilter.size +
              graphData.nodes.length ===
            0
          }
        />
      ) : (
        <TableContainer className='content-container' tableData={graphData} />
      )}
    </Layout>
  );
}
