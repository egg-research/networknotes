import React, { useState } from 'react';

import { useMeasure } from 'react-use';
import { Radio } from 'antd';
import Graph from './Graph';
import Layout from './Layout';
import DocumentTable from './DocumentTable';
import SearchBar from './SearchBar';

import { genRandomTree } from './utils/Data';
import './GraphPage.css';

const data = genRandomTree(20);

const documentData = [
  {
    key: 1,
    name: 'Intro to ML',
    lastAccessed: '8:35pm Apr 20, 2020',
    keywords: ['CNN', 'RNN'],
  },
  {
    key: 2,
    name: 'Advanced Data Structures',
    lastAccessed: '8:35pm Apr 20, 2020',
    keywords: ['Trie', 'Skip List'],
  },
];

const keywords = ['CNN', 'RNN', 'Trie', 'Skip List'];

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
      <div className='search-bar'>
        <SearchBar
          documents={documentData.map((document) => document.name)}
          keywords={keywords}
        />
      </div>
      <DocumentTable data={tableData} className='table' />
    </div>
  );
}

export default function GraphPage() {
  const [view, setView] = useState('graph');
  const [ref, dimensions] = useMeasure();
  const { height, width } = dimensions;

  return (
    <Layout Sidebar={<h1>hi</h1>} contentRef={ref} contentPadding={false}>
      <ViewSwitch className='view-switch' setView={setView} />
      {view === 'graph' ? (
        <Graph data={data} height={height} width={width} />
      ) : (
        <TableContainer tableData={documentData} />
      )}
    </Layout>
  );
}
