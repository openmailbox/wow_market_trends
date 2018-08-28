const path = require('path');

module.exports = {
  mode: 'production',
  entry: './web/static/src/scripts/index.js',
  devtool: 'inline-source-map',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'web/static/dist')
  },
  devServer: {
    contentBase: 'web/static/dist'
  },
  module: {
    rules: [{
      test: /\.css$/,
      use: [
        'style-loader',
        'css-loader'
      ]
    }]
  }
};
