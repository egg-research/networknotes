import React, { useEffect, useContext } from 'react';

import { BrowserRouter, Switch, Route } from 'react-router-dom';

// import logo from './logo.svg';
import { Button } from 'antd';
import { useEffectOnce } from 'react-use';
import './Router.css';

import EditorPage from './EditorPage';
import GraphPage from './GraphPage';
import { getUserId, getDocGraph, getKeywordGraph, getAllDocs } from './db';
import UserContext from './context';

export default function Router() {
  // useEffectOnce(async () => {
  //   console.log(await getUserId('username'));
  //   console.log(await getDocGraph(31));
  //   console.log(await getKeywordGraph(31));
  // });

  return (
    <UserContext.Provider value={32}>
      <>
        <link
          rel='stylesheet'
          href='//maxcdn.bootstrapcdn.com/font-awesome/4.6.1/css/font-awesome.min.css'
        />
        <BrowserRouter>
          <Switch>
            <Route path='/document/:id'>
              <EditorPage />
            </Route>
            <Route path='/'>
              <GraphPage />
            </Route>
          </Switch>
        </BrowserRouter>
      </>
    </UserContext.Provider>
  );
}
