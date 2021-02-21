import { processGraph } from './graph';

export function genRandomTree(N = 3, reverse = false) {
  const data = {
    nodes: [...Array(N).keys()].map((i) => ({
      // id: [
      //   'Intro to Machine Learning',
      //   'Computer Vision',
      //   'Advanced Data Structures',
      // ][i],
      id: i,
      name: [
        'Intro to Machine Learning',
        'Computer Vision',
        'Advanced Data Structures',
      ][i],
    })),
    links: [
      {
        source: 0,
        target: 1,
        name: 'RNN',
        id: 'RNN',
      },
      {
        source: 0,
        target: 1,
        name: 'CNN',
        id: 'CNN',
      },
    ],
  };

  return processGraph(data);
}
