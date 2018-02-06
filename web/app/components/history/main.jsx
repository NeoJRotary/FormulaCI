import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import './style.scss';
import formula from '../../api/formula';
import { isEmpty } from '../../func';

const HisCard = ({ hs }) => {
  let result = '';
  switch (hs.result) {
    case -2:
      result = 'canceled';
      break;
    case -1:
      result = 'running';
      break;
    case 0:
      result = 'failed';
      break;
    case 1:
      result = 'success';
      break;
    default:
  }

  const dur = hs.result === 1 ? `( ${hs.dur} s )` : '';

  return (
    <Link to={`/history/${hs.rowid}`}>
      <div className="hisCard">
        <div className={`tag ${result}`}>{result}</div>
        <div className="repo">{`${hs.repo}/${hs.branch}`}</div>
        <div className="time">{hs.time}</div>
        <div className="flow">
          <div className="title">{`Flow ${dur}`}</div>
          {
            hs.flow.map((k, i) => {
              if (i === hs.flow.length - 1) return <span className={`tag ${result}`}>{k}</span>;
              return [
                <span className="tag success">{k}</span>,
                <span className="line" />
              ];
            })
          }
        </div>
      </div>
    </Link>
  );
};

HisCard.propTypes = {
  hs: PropTypes.objectOf(PropTypes.any).isRequired
};

class history extends React.Component {
  constructor(props) {
    super(props);
    this.loadHistroy = this.loadHistroy.bind(this);

    this.state = {
      list: [],
      id: props.match.params.id
    };
  }

  componentDidMount() {
    this.loadHistroy();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.match.params.id !== nextProps.match.params.id) {
      this.setState({ id: nextProps.match.params.id });
    }
  }

  loadHistroy() {
    formula.getHistory().then(list => this.setState({ list }));
  }

  render() {
    if (isEmpty(this.state.list)) return <section id="history">loading</section>;

    const isDetail = this.state.id !== undefined;
    let detail;
    if (isDetail) {
      detail = this.state.list.find(l => String(l.rowid) === this.state.id);
    }
    return (
      <section id="history">
        {
          isDetail &&
          <div className="detail">
            <Link to="/history"><button>Back</button></Link>
            <HisCard hs={detail} />
            {
              detail.flow.map(f => (
                <div className="log">
                  <div className="step">{` - ${f}`}</div>
                  <div className="cliLog">
                    {
                      detail.log[f].split('\n').map(s => (
                        <span>
                          {s}
                          <br />
                        </span>
                      ))
                    }
                  </div>
                </div>
              ))
            }
          </div>
        }
        {
          !isDetail &&
          this.state.list.map(l => <HisCard hs={l} />)
        }
      </section>
    );
  }
}

history.propTypes = {
  match: PropTypes.objectOf(PropTypes.any).isRequired
};

export default history;
