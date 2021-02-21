import React, { useState, useRef, useCallback, useEffect } from 'react';
import { useHistory } from 'react-router-dom';

import { Empty, Spin } from 'antd';
import { ForceGraph2D } from 'react-force-graph';

const NODE_R = 8;

function Graph({ data, height, width, setNode, load }) {
  const [highlights, setHighlight] = useState({
    nodes: new Set(),
    edges: new Set(),
    hoverNode: null,
  });

  const ref = useRef(null);
  const history = useHistory();

  useEffect(() => {
    setTimeout(() => {
      if (ref && ref.current) {
        ref.current.d3Force('charge').strength(-50);
      }
    }, 0);
  }, []);

  if (load) {
    return (
      <div
        style={{
          display: 'flex',
          height: '100%',
          width: '100%',
          justifyContent: 'center',
          alignItems: 'center',
        }}
      >
        <Spin style={{ height: 200 }} size='large' />
      </div>
    );
  }

  if (data.nodes.length === 0) {
    return (
      <div
        style={{
          display: 'flex',
          height: '100%',
          width: '100%',
          justifyContent: 'center',
          alignItems: 'center',
        }}
      >
        <Empty style={{ height: 200 }} />
      </div>
    );
  }

  function paintNode(node, color, ctx) {
    const { id, name, x, y } = node;
    const fontSize = 8;
    ctx.font = `${fontSize}px Sans-Serif`;
    const textWidth = ctx.measureText(name).width + fontSize * 0.5;
    const bckgDimensions = [textWidth, fontSize].map((n) => n + fontSize * 0.2); // some padding

    ctx.fillStyle = 'rgba(255, 255, 255, 0.7)';
    ctx.fillRect(
      x - bckgDimensions[0] / 2,
      y - bckgDimensions[1] / 2,
      ...bckgDimensions
    );

    ctx.fillStyle = color;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(name, x, y);
    node._bckgDimensions = bckgDimensions;
  }

  function areaPaint(node, color, ctx) {
    ctx.fillStyle = color;
    const bckgDimensions = node._bckgDimensions;
    if (bckgDimensions) {
      ctx.fillRect(
        node.x - bckgDimensions[0] / 2,
        node.y - bckgDimensions[1] / 2,
        ...bckgDimensions
      );
    }
  }

  return (
    <ForceGraph2D
      ref={ref}
      graphData={data}
      height={height}
      width={width}
      nodeLabel={() => undefined}
      linkLabel={(link) => link.name}
      nodeCanvasObject={(node, ctx) => paintNode(node, '#63ace5', ctx)}
      nodePointerAreaPaint={areaPaint}
      onNodeClick={(node) => history.push(`/document/${node.id}`)}
    />
  );
}

export default React.memo(Graph);
