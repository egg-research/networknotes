import React from 'react';

import { BrowserRouter, Switch, Route } from 'react-router-dom';

// import logo from './logo.svg';
import { Button } from 'antd';
import './Router.css';

import EditorPage from './EditorPage';
import GraphPage from './GraphPage';

export default function Router() {
  return (
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
  );
}
