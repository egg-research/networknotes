import React, { useState } from 'react';

import { useMeasure } from 'react-use';
import { Radio } from 'antd';
import Graph from './Graph';
import Layout from './Layout';

import { genRandomTree } from './utils/Data';
import './GraphPage.css';

const data = genRandomTree(20);

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

export default function GraphPage() {
  const [view, setView] = useState('graph');
  const [ref, dimensions] = useMeasure();
  const { height, width } = dimensions;

  return (
    <Layout Sidebar={<h1>hi</h1>} contentRef={ref} contentPadding={false}>
      <ViewSwitch className='view-switch' setView={setView} />
      <Graph data={data} height={height} width={width} />
    </Layout>
  );
}
