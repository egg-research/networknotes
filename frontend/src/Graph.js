import React, { useState, useRef, useCallback } from 'react';

import { ForceGraph2D } from 'react-force-graph';

const NODE_R = 8;

function Graph({ data, height, width }) {
  const [highlights, setHighlight] = useState({
    nodes: new Set(),
    edges: new Set(),
    hoverNode: null,
  });

  const ref = useRef(null);
  setTimeout(() => {
    if (ref && ref.current) {
      ref.current.d3Force('charge').strength(-50);
    }
  }, 0);

  const handleNodeHover = (node) => {
    const newHighlights = {
      nodes: new Set(),
      edges: new Set(),
      hoverNode: null,
    };

    console.log(node);

    if (node) {
      node.neighbors.forEach((neighbor) => newHighlights.nodes.add(neighbor));
      node.links.forEach((edge) => newHighlights.edges.add(edge));
      newHighlights.nodes.add(node);
      newHighlights.hoverNode = node;
    }

    setHighlight(newHighlights);
  };

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
      nodeCanvasObject={(node, ctx) => paintNode(node, 'red', ctx)}
      nodePointerAreaPaint={areaPaint}
      onNodeHover={handleNodeHover}
    />
  );
}

export default Graph;
