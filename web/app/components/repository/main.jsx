import React from 'react';
// import PropTypes from 'prop-types';
// import './style.scss';

import repo from '../../api/repo';

class repository extends React.Component {
  constructor() {
    super();
    this.onChange = this.onChange.bind(this);
    this.loadRepo = this.loadRepo.bind(this);
    this.addRepo = this.addRepo.bind(this);
    this.triggerRepo = this.triggerRepo.bind(this);
    this.rmRepo = this.rmRepo.bind(this);

    this.state = {
      newName: '',
      newSrc: '',
      newBranch: '',
      list: []
    };
  }

  componentDidMount() {
    this.loadRepo();
    // repo.getList().then(list => console.log(list));
  }

  onChange(e, key) {
    const obj = {};
    obj[key] = e.target.value;
    this.setState(obj);
  }

  loadRepo() {
    repo.getList().then(list => this.setState({ list }));
  }

  addRepo() {
    alert('add');
    const { newName, newSrc, newBranch } = this.state;
    repo.add(newName, newSrc, newBranch).then(() => this.loadRepo());
  }

  triggerRepo(name, branch) {
    alert('trigger');
    repo.trigger(name, branch).then(() => this.loadRepo());
  }

  rmRepo(name, branch) {
    alert('reomve');
    repo.remove(name, branch).then(() => this.loadRepo());
  }

  render() {
    return (
      <section id="repository">
        <p>Add</p>
        <div>Name <input value={this.state.newName} onChange={e => this.onChange(e, 'newName')} /></div>
        <div>Src <input value={this.state.newSrc} onChange={e => this.onChange(e, 'newSrc')} /></div>
        <div>Branch <input value={this.state.newBranch} onChange={e => this.onChange(e, 'newBranch')} /></div>
        <button onClick={this.addRepo}>addRepo</button>
        {
          this.state.list.map(r =>
            (
              <div>
                <span>{r.name} , </span>
                <span>{r.src} , </span>
                <span>{r.branch}</span>
                <button onClick={() => this.triggerRepo(r.name, r.branch)}>trigger</button>
                <button onClick={() => this.rmRepo(r.name, r.branch)}>remove</button>
              </div>
            ))
        }
      </section>
    );
  }
}

// header.propTypes = {
// };

export default repository;
