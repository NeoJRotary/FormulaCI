import React from 'react';
// import PropTypes from 'prop-types';

import './style.scss';
import formula from '../../api/formula';

class history extends React.Component {
  constructor(props) {
    super(props);
    this.loadHistroy = this.loadHistroy.bind(this);

    this.state = {
      list: []
    };
  }

  componentDidMount() {
    this.loadHistroy();
  }

  loadHistroy() {
    formula.getHistory().then(list => this.setState({ list }));
  }

  render() {
    return (
      <section id="history">
        {/* <div className="record">
          <span className="col">result</span>
          <span className="col">flow</span>
          <span className="col">time</span>
          <span className="col">duration (s)</span>
        </div> */}
        {
          this.state.list.map(l => (
            <div className="record">
              <span className="col">{l.result}</span>
              <span className="col">{`${l.repo}/${l.branch}`}</span>
              <span className="col">{l.time}</span>
              <span className="col">{l.dur}</span>
              {
                l.flow.map(f => (
                  <div>
                    <p>{f}</p>
                    <p className="cliLog">
                      {
                        l.log[f].split('\n').map(s => (
                          <span>
                            {s}
                            <br />
                          </span>
                        ))
                      }
                    </p>
                  </div>
                ))
              }
            </div>
          ))
        }
      </section>
    );
  }
}

// dashboard.propTypes = {
// };

export default history;
