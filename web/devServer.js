const path = require('path');
const webpack = require('webpack');
const WebpackDevServer = require('webpack-dev-server');

function devServer() {
  const conf = {
    entry: {
      app: './app/index.jsx'
    },
    devtool: 'inline-source-map',
    output: {
      filename: 'bundle.js',
      path: path.resolve(__dirname, 'build'),
      publicPath: '/'
    },
    module: {
      rules: [
        {
          test: /\.js$|\.jsx?$/,
          exclude: /node_modules/,
          loader: 'babel-loader'
        },
        {
          test: /\.scss$|\.sass$/,
          use: [
            { loader: 'style-loader' },
            {
              loader: 'css-loader',
              options: {
                url: false
              }
            },
            {
              loader: 'sass-loader'
            }
          ]
        }
      ]
    },
    plugins: [
      new webpack.HotModuleReplacementPlugin(),
      new webpack.NamedModulesPlugin()
    ]
  };
  const options = {
    contentBase: './build',
    watchContentBase: true,
    hot: true,
    // hotOnly: true,
    // inline: true,
    host: 'localhost',
    historyApiFallback: true
  };

  WebpackDevServer.addDevServerEntrypoints(conf, options);
  const compiler = webpack(conf);
  const server = new WebpackDevServer(compiler, options);
  server.listen(8080, 'localhost');
}

devServer();
