import { processGraph } from './graph';

export function genRandomTree(N = 3, reverse = false) {
  const data = {
    nodes: [...Array(N).keys()].map((i) => ({
      id: [
        'Intro to Machine Learning',
        'Computer Vision',
        'Advanced Data Structures',
      ][i],
      name: [
        'Intro to Machine Learning',
        'Computer Vision',
        'Advanced Data Structures',
      ][i],
    })),
    links: [
      {
        source: 'Intro to Machine Learning',
        target: 'Computer Vision',
        name: 'RNN',
        id: 'RNN',
      },
      {
        source: 'Intro to Machine Learning',
        target: 'Computer Vision',
        name: 'CNN',
        id: 'CNN',
      },
    ],
  };

  return processGraph(data);
}
