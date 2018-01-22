import React from 'react';
// import PropTypes from 'prop-types';

import gcloud from '../../api/gcloud';
import git from '../../api/git';
import './style.scss';

class config extends React.Component {
  constructor(props) {
    super(props);
    this.onChange = this.onChange.bind(this);
    this.setGCP = this.setGCP.bind(this);
    this.setGitEmail = this.setGitEmail.bind(this);
    this.setGitWebhookToken = this.setGitWebhookToken.bind(this);
    this.generateSSH = this.generateSSH.bind(this);

    this.state = {
      gcpAuthKey: '',
      gcpProject: '',
      gkeZone: '',
      gkeName: '',
      gitEmail: '',
      gitWebhookToken: '',
      gitPubkey: ''
    };
  }

  componentDidMount() {
    gcloud.getInfo().then(info => this.setState({
      gcpAuthKey: info.authKey,
      gcpProject: info.project,
      gkeZone: info.gkeZone,
      gkeName: info.gkeName
    }));
    git.getInfo().then(info => this.setState({
      gitEmail: info.email,
      gitWebhookToken: info.webhookToken,
      gitPubkey: info.pubkey
    }));
  }

  onChange(e, key) {
    const obj = {};
    obj[key] = e.target.value;
    this.setState(obj);
  }

  setGCP(key) {
    const {
      gcpAuthKey, gcpProject, gkeZone, gkeName
    } = this.state;
    switch (key) {
      case 'gcpAuthKey':
        gcloud.setAuthKey(gcpAuthKey).then(result => console.log(result));
        break;
      case 'gcpProject':
        gcloud.setProject(gcpProject).then(result => console.log(result));
        break;
      case 'gcpGKE':
        gcloud.setGKE(gkeZone, gkeName).then(result => console.log(result));
        break;
      default:
    }
  }

  setGitEmail() {
    git.setEmail(this.state.gitEmail);
    // .then(gitEmail => this.setState({ gitEmail }));
  }

  setGitWebhookToken() {
    git.setWebhookToken(this.state.gitWebhookToken);
    // .then(gitEmail => this.setState({ gitEmail }));
  }

  generateSSH() {
    git.generateSSH().then(gitPubkey => this.setState({ gitPubkey }));
  }

  render() {
    return (
      <section id="config">
        <h3>Google Cloud Platform</h3>
        <p>Service Account Key</p>
        <textarea value={this.state.gcpAuthKey} onChange={e => this.onChange(e, 'gcpAuthKey')} />
        <button onClick={() => this.setGCP('gcpAuthKey')}>update</button>
        <p>Project ID</p>
        <input value={this.state.gcpProject} onChange={e => this.onChange(e, 'gcpProject')} />
        <button onClick={() => this.setGCP('gcpProject')}>Update</button>
        <p>GKE Zone</p>
        <input value={this.state.gkeZone} onChange={e => this.onChange(e, 'gkeZone')} />
        <p>GKE Name</p>
        <input value={this.state.gkeName} onChange={e => this.onChange(e, 'gkeName')} />
        <button onClick={() => this.setGCP('gcpGKE')}>Connect GKE</button>
        {/* <textarea value={this.state.gcpKey} onChange={e => this.onChange(e, 'gcpKey')} />
        <button onClick={this.updateGCPKey}>update</button> */}
        <h3>Git</h3>
        <p>Email</p>
        <input value={this.state.gitEmail} onChange={e => this.onChange(e, 'gitEmail')} />
        <button onClick={this.setGitEmail}>Update</button>
        <p>Webhook Token</p>
        <input value={this.state.gitWebhookToken} onChange={e => this.onChange(e, 'gitWebhookToken')} />
        <button onClick={this.setGitWebhookToken}>Update</button>
        <p>SSH Public Key</p>
        <button onClick={this.generateSSH}>generate SSH Key</button>
        <textarea value={this.state.gitPubkey} readOnly />
      </section>
    );
  }
}

// config.propTypes = {
// };

export default config;
