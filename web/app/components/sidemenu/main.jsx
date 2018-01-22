import React from 'react';
// import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import './style.scss';

class sidemenu extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div id="sidemenu">
        <div className="title">Formula CI</div>
        <Link to="/">Dashboard</Link>
        <Link to="/repository">Repository</Link>
        <Link to="/history">History</Link>
        <Link to="/config">Config</Link>
      </div>
    )
  }
}

// header.propTypes = {
// };

export default sidemenu;
