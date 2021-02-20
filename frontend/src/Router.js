import React from 'react';

import { BrowserRouter, Switch, Route } from 'react-router-dom';

// import logo from './logo.svg';
import { Button } from 'antd';
import './Router.css';

import EditorPage from './EditorPage';

export default function Router() {
  return (
    <>
      <link
        rel='stylesheet'
        href='//maxcdn.bootstrapcdn.com/font-awesome/4.6.1/css/font-awesome.min.css'
      />
      <BrowserRouter>
        <Switch>
          <Route path='/editor'>
            <EditorPage />
          </Route>
          <Route path='/'>
            <Content />
          </Route>
        </Switch>
      </BrowserRouter>
    </>
  );
}

function Content() {
  return (
    <div className='App'>
      <Button type='primary'>Button</Button>
      {/* <header className='App-header'> */}
      {/*   <img src={logo} className='App-logo' alt='logo' /> */}
      {/*   <p> */}
      {/*     Edit <code>src/App.js</code> and save to reload. */}
      {/*   </p> */}
      {/*   <a */}
      {/*     className='App-link' */}
      {/*     href='https://reactjs.org' */}
      {/*     target='_blank' */}
      {/*     rel='noopener noreferrer' */}
      {/*   > */}
      {/*     Learn React */}
      {/*   </a> */}
      {/* </header> */}
    </div>
  );
}
