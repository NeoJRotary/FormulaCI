import React from 'react';
// import PropTypes from 'prop-types';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import ws from './api/websocket';
import './style.scss';

import sidemenu from './components/sidemenu';
import dashboard from './components/dashboard';
import repository from './components/repository';
import history from './components/history';
import config from './components/config';

class App extends React.Component {
  constructor() {
    super();
    this.wsChange = this.wsChange.bind(this);

    this.state = {
      wsopen: ws.open
    };

    ws.start(this.wsChange);
  }

  wsChange(status) {
    this.setState({ wsopen: status });
  }

  render() {
    if (!this.state.wsopen) return <div id="router-wrapper">Websocket Connecting</div>;
    return (
      <div id="router-wrapper">
        <Route component={sidemenu} />
        <Route exact path="/" component={dashboard} />
        <Route exact path="/repository" component={repository} />
        <Route exact path="/history/:id?" component={history} />
        <Route exact path="/config" component={config} />
      </div>
    );
  }
}

export default () => (
  <BrowserRouter>
    <Switch>
      <Route component={App} />
    </Switch>
  </BrowserRouter>
);
