import React from 'react';

import { useMeasure } from 'react-use';
import Graph from './Graph';
import Layout from './Layout';

import { genRandomTree } from './utils/Data';

const data = genRandomTree(20);

export default function GraphPage() {
  const [ref, dimensions] = useMeasure();
  const { height, width } = dimensions;
  console.log(dimensions);

  return (
    <Layout Sidebar={<h1>hi</h1>} contentRef={ref} contentPadding={false}>
      <Graph data={data} height={height} width={width} />
    </Layout>
  );
}
