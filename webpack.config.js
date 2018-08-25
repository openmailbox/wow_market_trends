const path = require('path');

module.exports = {
  entry: './web/static/src/scripts/index.js',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'web/static/dist')
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